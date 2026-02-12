<script lang="ts">
  import { onMount } from "svelte";
  import Swal from "sweetalert2";
  import { get } from "svelte/store";
  import { auth } from "$lib/stores/auth";

  const API_URL = "http://localhost:8080";

  let lotNumber = "";
  let activeMachine = "";
  let scanTime = "";
  let cycleCount = 0;

  // State untuk siklus produksi
  let isCycleActive = false; // false = Biru (Mulai), true = Oranye (Berjalan/Akhiri)
  let currentCycleStartTime = "";

  // Data detail
  let lotDetail = {
    lot: "",
    product: "KPCP Standard",
    quantity: 500,
    status: "Proses",
    startTime: "",
    estimatedEnd: "",
    supervisor: "Tono Widiyanto",
    machine: "",
    noCompound: "CMP-2309-001",
    noLotMixing: "MIX-2309-A05",
  };

  let isSubmitting = false;

  // --- LOGIKA BARU UNTUK TOMBOL PRODUKSI ---

  async function handleProductionAction() {
    if (!isCycleActive) {
      // START CYCLE (Biru -> Oranye)
      isCycleActive = true;
      currentCycleStartTime = new Date().toLocaleTimeString("id-ID");

      const Toast = Swal.mixin({
        toast: true,
        position: "top-end",
        showConfirmButton: false,
        timer: 2000,
        timerProgressBar: true,
      });
      Toast.fire({
        icon: "info",
        title: `Cycle dimulai pukul ${currentCycleStartTime}`,
      });
    } else {
      // END CYCLE (Oranye -> Muncul Form -> Kembali Biru)
      const endTime = new Date().toLocaleTimeString("id-ID");

      const { value: formValues } = await Swal.fire({
        title: "Input Hasil Cycle",
        html: `
    <div class="text-left space-y-4">
      
      <div class="p-3 bg-blue-50 rounded-lg border border-blue-100 mb-4">
        <p class="text-sm font-bold text-blue-800">OK Section (Time Monitoring)</p>
        <div class="grid grid-cols-2 gap-2 mt-1">
          <div>
             <label class="text-[10px] uppercase text-blue-600 font-bold">Mulai</label>
             <p class="text-sm font-mono">${currentCycleStartTime}</p>
          </div>
          <div>
             <label class="text-[10px] uppercase text-blue-600 font-bold">Selesai</label>
             <p class="text-sm font-mono">${endTime}</p>
          </div>
        </div>
        <div class="mt-3">
          <label class="block text-xs font-semibold text-slate-600 mb-1">Berapa OK?</label>
          <input id="swal-ok-qty" class="swal2-input !m-0 !w-full !text-sm" placeholder="Total OK Item">
          
        </div>
      </div>

      <div class="p-3 bg-red-50 rounded-lg border border-red-100">
        
        <div class="flex justify-between items-center mb-2">
          <p class="text-sm font-bold text-red-800">NG Section (Reject Data)</p>
          <button type="button" id="add-ng-row"
            class="text-xl bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600">
            + Tambah
          </button>
        </div>

        <label class="block text-xs font-semibold text-slate-600 mb-1">Berapa NG?</label>
        <input id="swal-ng-item" class="swal2-input !m-0 !w-full !text-sm mb-3" placeholder="Total NG Item">

        <div id="ng-container" class="space-y-3"></div>

        <div class="mt-3 text-right text-sm font-semibold text-red-700">
          Total NG: <span id="total-ng">0</span> pcs
        </div>

      </div>
    </div>
  `,
        focusConfirm: false,
        showCancelButton: true,
        confirmButtonText: "Simpan Cycle",
        cancelButtonText: "Batal",
        confirmButtonColor: "#f97316",

        didOpen: () => {
          const container = document.getElementById("ng-container")!;
          const addBtn = document.getElementById("add-ng-row")!;
          const totalDisplay = document.getElementById("total-ng")!;

          function createRow() {
            const row = document.createElement("div");
            row.className = "ng-row grid grid-cols-12 gap-2 items-end";

            row.innerHTML = `
        <div class="col-span-5">
          <label class="block text-xs font-semibold text-slate-600">Kategori</label>
          <select class="swal-ng-cat swal2-input !m-0 !w-full !text-sm">
            <option value="Flash">Flash</option>
            <option value="Dirty">Dirty</option>
            <option value="Short Shot">Short Shot</option>
            <option value="Other">Lainnya</option>
          </select>
        </div>

        <div class="col-span-5">
          <label class="block text-xs font-semibold text-slate-600">Jumlah (Pcs)</label>
          <input type="number" min="0" value="0"
            class="swal-ng-qty swal2-input !m-0 !w-full !text-sm">
        </div>

        <div class="col-span-2 flex justify-center">
          <button type="button"
            class="delete-row text-lg bg-gray-400 text-white px-2 py-1 rounded hover:bg-gray-500">
            ✕
          </button>
        </div>
      `;

            container.appendChild(row);
            updateTotal();
          }

          function updateTotal() {
            const qtyInputs = container.querySelectorAll(
              ".swal-ng-qty",
            ) as NodeListOf<HTMLInputElement>;
            let total = 0;

            qtyInputs.forEach((input) => {
              total += Number(input.value) || 0;
            });

            totalDisplay.textContent = total.toString();
          }

          // Event Delegation untuk qty & delete
          container.addEventListener("input", (e) => {
            if ((e.target as HTMLElement).classList.contains("swal-ng-qty")) {
              updateTotal();
            }
          });

          container.addEventListener("click", (e) => {
            if ((e.target as HTMLElement).classList.contains("delete-row")) {
              const rows = container.querySelectorAll(".ng-row");
              if (rows.length > 1) {
                (e.target as HTMLElement).closest(".ng-row")?.remove();
                updateTotal();
              }
            }
          });

          addBtn.addEventListener("click", createRow);

          // Buat 1 row default
          createRow();
        },

        preConfirm: () => {
          const item = (
            document.getElementById("swal-ng-item") as HTMLInputElement
          ).value;

          const rows = document.querySelectorAll(".ng-row");
          const ngData: any[] = [];
          let total = 0;

          rows.forEach((row) => {
            const category = (
              row.querySelector(".swal-ng-cat") as HTMLSelectElement
            ).value;
            const qty = Number(
              (row.querySelector(".swal-ng-qty") as HTMLInputElement).value,
            );

            if (qty > 0) {
              ngData.push({ category, qty });
              total += qty;
            }
          });

          return {
            startTime: currentCycleStartTime,
            endTime: endTime,
            item,
            ngData,
            totalNG: total,
          };
        },
      });

      if (formValues) {
        cycleCount++;
        isCycleActive = false; // Reset tombol ke biru
        Swal.fire({
          icon: "success",
          title: `Cycle ${cycleCount} Berhasil!`,
          text: `Data NG: ${formValues.qty} pcs (${formValues.category})`,
          timer: 1500,
        });
      }
    }
  }

  function handleBackToDashboard() {
    Swal.fire({
      title: "Kembali ke Dashboard?",
      text: "Perubahan yang belum disimpan akan hilang.",
      icon: "warning",
      showCancelButton: true,
      confirmButtonColor: "#4f46e5",
      confirmButtonText: "Ya, Kembali",
    }).then((result) => {
      if (result.isConfirmed) {
        window.location.href = "/pressing";
      }
    });
  }

function openLineStop() {
  Swal.fire({
    title: 'Input Kode Line Stop',
    input: 'text',
    inputLabel: 'Kode Line Stop',
    inputPlaceholder: 'Masukkan kode line stop...',
    showCancelButton: true,
    confirmButtonText: 'Simpan',
    cancelButtonText: 'Batal',
    confirmButtonColor: '#dc2626',
    inputValidator: (value) => {
      if (!value) {
        return 'Kode Line Stop tidak boleh kosong!';
      }
      if (value.length < 3) {
        return 'Kode minimal 3 karakter';
      }
    }
  }).then((result) => {

    if (result.isConfirmed) {

      Swal.fire({
        icon: 'success',
        title: 'Line Stop Berhasil Dicatat!',
        text: `Kode: ${result.value}`,
        timer: 1500,
        showConfirmButton: false
      });

    }

  });
}

  // --- LOGIKA BACKEND ORIGINAL (TIDAK DIUBAH) ---

  async function handleComplete() {
    const result = await Swal.fire({
      title: "Konfirmasi Simpan Data?",
      text: "Data produksi akhir akan dikirim ke server.",
      icon: "question",
      showCancelButton: true,
      confirmButtonColor: "#4f46e5",
      confirmButtonText: "Ya, Kirim",
    });

    if (result.isConfirmed) {
      isSubmitting = true;
      try {
        const authToken = get(auth).token;
        if (!authToken) throw new Error("Token tidak ditemukan.");

        const payload = {
          noMesin: lotDetail.machine,
          tanggal: new Date().toISOString(),
          shift: "I",
          partName: lotDetail.product,
          kodePart: "KPCP-STD",
          jamMulai: lotDetail.startTime,
          jamSelesai: new Date().toLocaleTimeString("id-ID"),
          hasilOk: cycleCount,
          ng: 0,
          klasifikasiReject: "",
        };

        const response = await fetch(`${API_URL}/api/lwp`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authToken}`,
          },
          body: JSON.stringify(payload),
        });

        if (!response.ok) throw new Error("Gagal menyimpan data LWP");

        await Swal.fire("Berhasil!", "Data tersimpan di Dashboard.", "success");
        window.location.href = "/pressing";
      } catch (error: any) {
        Swal.fire("Gagal", error.message, "error");
      } finally {
        isSubmitting = false;
      }
    }
  }

  onMount(() => {
    const queryParams = new URLSearchParams(window.location.search);
    lotNumber = queryParams.get("lot") || "";
    activeMachine = localStorage.getItem("activeMachine") || "UNKNOWN-MC";
    scanTime =
      localStorage.getItem("scanTime") ||
      new Date().toLocaleTimeString("id-ID");

    const scannedProduct = localStorage.getItem("selectedProduct");
    if (scannedProduct) lotDetail.product = scannedProduct;

    lotDetail.lot = lotNumber;
    lotDetail.machine = activeMachine;
    lotDetail.startTime = scanTime;

    const endDate = new Date();
    endDate.setHours(endDate.getHours() + 1);
    lotDetail.estimatedEnd = endDate.toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
    });
  });

  function getStatusColor(status: string) {
    if (status === "Selesai")
      return "bg-emerald-50 text-emerald-700 border-emerald-200";
    if (status === "Proses") return "bg-blue-50 text-blue-700 border-blue-200";
    return "bg-slate-50 text-slate-600 border-slate-200";
  }
</script>

<div class="min-h-screen bg-slate-50 text-slate-800 pb-12">
  <header
    class="sticky top-0 z-50 bg-white border-b border-slate-200 shadow-sm"
  >
    <div class="max-w-6xl mx-auto px-4 py-4">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-slate-800">Monitoring Cycle</h1>
          <p class="text-sm text-slate-500 mt-1">
            Lot: <span class="font-bold text-indigo-600">{lotNumber}</span>
          </p>
        </div>
        <div class="text-right">
          <p class="text-xs text-slate-500">
            {new Date().toLocaleDateString("id-ID")}
          </p>
          <div class="flex items-center gap-2 mt-1">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"
            ></span>
            <span class="text-sm font-semibold uppercase"
              >{lotDetail.machine}</span
            >
          </div>
        </div>
      </div>
    </div>
  </header>

  <main class="max-w-6xl mx-auto px-4 mt-8 space-y-6">
    <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
      <div
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-center md:text-left"
      >
        <div class="p-3 border-r border-slate-100 last:border-0">
          <p class="text-[10px] font-bold text-slate-400 uppercase">Produk</p>
          <p class="text-lg font-bold text-slate-700">{lotDetail.product}</p>
        </div>
        <div class="p-3 border-r border-slate-100 last:border-0">
          <p class="text-[10px] font-bold text-slate-400 uppercase">
            Cycle Counter
          </p>
          <p class="text-2xl font-black text-indigo-600">{cycleCount}</p>
        </div>
        <div class="p-3 border-r border-slate-100 last:border-0">
          <p class="text-[10px] font-bold text-slate-400 uppercase">
            No. Compound
          </p>
          <p class="text-sm font-semibold text-slate-700">
            {lotDetail.noCompound}
          </p>
        </div>
        <div class="p-3">
          <p class="text-[10px] font-bold text-slate-400 uppercase">
            Status Mesin
          </p>
          <span
            class={`px-3 py-1 rounded-full text-xs font-bold inline-block mt-1 ${getStatusColor(lotDetail.status)}`}
          >
            {isCycleActive ? "Sedang Running" : "Idle / Siap"}
          </span>
        </div>
      </div>
      <div>
        <button
          on:click={openLineStop}
          class="bg-red-600 hover:bg-red-700 mt-8 text-white w-full text-lg font-bold px-5 py-3 rounded-lg shadow"
        >
          LineStop
        </button>
      </div>
    </div>

    <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <button
          on:click={handleBackToDashboard}
          class="flex items-center justify-center gap-3 px-6 py-4 bg-slate-500 hover:bg-slate-600 text-white font-bold rounded-xl transition-all shadow-lg shadow-slate-200 text-lg"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M10 19l-7-7m0 0l7-7m-7 7h18"
            />
          </svg>
          Kembali ke Dashboard
        </button>

        <button
          on:click={handleProductionAction}
          class="flex items-center justify-center gap-3 px-6 py-4 font-bold rounded-xl transition-all shadow-lg text-lg text-white
          {isCycleActive
            ? 'bg-orange-500 hover:bg-orange-600 animate-pulse shadow-orange-200'
            : 'bg-blue-600 hover:bg-blue-700 shadow-blue-200'}"
        >
          {#if isCycleActive}
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z"
              />
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            Akhiri Satu Cycle ({cycleCount + 1})
          {:else}
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
              />
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            Mulai Produksi
          {/if}
        </button>
      </div>

      <div
        class="mt-8 pt-6 border-t border-slate-100 flex flex-col items-center"
      >
        <button
          on:click={handleComplete}
          disabled={isSubmitting || isCycleActive}
          class="w-full md:w-auto px-12 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-300 text-white font-bold rounded-lg transition-all shadow-md"
        >
          {isSubmitting ? "Menyimpan..." : "Finalisasi & Kirim LWP"}
        </button>
        <p class="text-[10px] text-slate-400 mt-2">
          Pastikan semua cycle selesai sebelum menekan Finalisasi
        </p>
      </div>
    </div>
  </main>
</div>

<style>
  :global(body) {
    background-color: #f8fafc;
  }

  /* Styling custom untuk SweetAlert agar input lebih rapi */
  :global(.swal2-input) {
    height: 2.5rem !important;
    font-size: 0.875rem !important;
  }
</style>
