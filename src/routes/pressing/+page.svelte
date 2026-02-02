<script lang="ts">
  import { onMount } from "svelte";
  import { get } from "svelte/store";
  import { auth } from "$lib/stores/auth";

  const API_URL = "http://localhost:8080";

  // ===== Svelte 5 $state for reactive data =====
  let pressRequests = $state<any[]>([
    {
      id: 1,
      lotNo: "LOT-2024-001",
      qty: 50,
      status: "Urgent",
      note: "Material A1",
    },
    {
      id: 2,
      lotNo: "LOT-2024-005",
      qty: 120,
      status: "Normal",
      note: "Material B2",
    },
  ]);

  let isLoading = $state(false);
  let errorMessage = $state("");
  let successMessage = $state("");

  // Form state for pressing data submission
  let formData = $state({
    partName: "",
    quantity: "",
    reject: "",
  });

  let employee = {
    name: "Tono Widiyanto",
    position: "Senior Cutting Specialist",
    department: "Production - Cutting",
    nik: "KRTP-2023-0456",
    phone: "+62 812-3456-7890",
    photo:
      "https://i.pinimg.com/550x/26/38/08/2638086da29fccffa32c5666ea77ce09.jpg",
  };

  let dailyCompounds = $state([
    { day: "Senin", short: "Sen", target: 4, actual: 4, efficiency: 100 },
    { day: "Selasa", short: "Sel", target: 4, actual: 3.5, efficiency: 88 },
    { day: "Rabu", short: "Rab", target: 4, actual: 5, efficiency: 125 },
    { day: "Kamis", short: "Kam", target: 4, actual: 4, efficiency: 100 },
    { day: "Jumat", short: "Jum", target: 4, actual: 4, efficiency: 100 },
    { day: "Sabtu", short: "Sab", target: 3, actual: 2, efficiency: 67 },
    { day: "Minggu", short: "Min", target: 0, actual: 0, efficiency: 0 },
  ]);

  let monthlyData = $state({
    completed: 18,
    target: 20,
    efficiency: 90,
    todayCompleted: 3,
    todayTarget: 4,
  });

  let recentScans = $state<any[]>([
    {
      id: 1,
      lot: "KPCP-2309-A01",
      timeStart: "10:30",
      timeEnd: "11:15",
      status: "Proses",
      pic: "Tono W.",
    },
    {
      id: 2,
      lot: "KPCP-2309-A02",
      timeStart: "10:15",
      timeEnd: "11:00",
      status: "Selesai",
      pic: "Tono W.",
    },
    {
      id: 3,
      lot: "KPCP-2309-B05",
      timeStart: "09:45",
      timeEnd: "10:30",
      status: "Selesai",
      pic: "Tono W.",
    },
  ]);

  // ===== LWP (Laporan Waktu Produksi) Data =====
  let lwpRecords = $state<any[]>([
    {
      id: 1,
      noMesin: "PRESS-01",
      tanggal: "2024-01-15",
      shift: "I",
      nik: "KRTP-2023-0456",
      partName: "Part A1",
      kodePart: "PA-001",
      targetShift: 100,
      details: [
        {
          id: 1,
          noLot: "LOT-2024-001",
          jamMulai: "07:00",
          jamSelesai: "08:30",
          hasilOk: 95,
          ng: 5,
          klasifikasiReject: [
            { jenis: "Cacat Permukaan", qty: 3 },
            { jenis: "Dimensi", qty: 2 }
          ],
        },
        {
          id: 2,
          noLot: "LOT-2024-002",
          jamMulai: "08:45",
          jamSelesai: "10:15",
          hasilOk: 102,
          ng: 3,
          klasifikasiReject: [
            { jenis: "Goresan", qty: 3 }
          ],
        },
      ],
    },
  ]);

  // Form state untuk tambah LWP record
  let lwpFormData = $state({
    noMesin: "",
    tanggal: new Date().toISOString().split('T')[0],
    shift: "I",
    nik: employee.nik,
    partName: "",
    kodePart: "",
    jamMulai: "",
    jamSelesai: "",
    hasilOk: "",
    ng: "",
    klasifikasiReject: "Cacat Permukaan",
  });

  let isLwpFormOpen = $state(false);

  // Tambah state untuk status mesin (NoMesin, Item, NoLot, Operator)
  let machineStatuses = $state<any[]>([]);

  // Fungsi untuk membangun machineStatuses dari lwpRecords
  function updateMachineStatuses() {
    machineStatuses = lwpRecords.map((m: any) => {
      const latestDetail = m.details && m.details.length ? m.details[m.details.length - 1] : null;
      return {
        noMesin: m.noMesin || "-",
        item: m.partName || "-",
        noLot: latestDetail ? latestDetail.noLot : "-",
        operator: m.nik || "-",
      };
    });
  }

  // panggil setelah inisialisasi data awal
  updateMachineStatuses();

  // Logic untuk Bar Chart Scale
  const maxChartValue = Math.max(
    ...dailyCompounds.map((d) => d.actual),
    ...dailyCompounds.map((d) => d.target),
    6
  );

  // ===== API Functions =====

  // Load today's pressing data from API
  async function loadData() {
    isLoading = true;
    errorMessage = "";

    try {
      const authToken = get(auth).token;
      if (!authToken) {
        errorMessage = "No authentication token found";
        return;
      }

      const response = await fetch(`${API_URL}/api/pressing/today`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to load pressing data");
      }

      const data = await response.json();
      recentScans = data.scans || recentScans;
      monthlyData = data.stats || monthlyData;
      // jika API mengembalikan LWP, update lwpRecords
      if (data.lwpRecords) {
        lwpRecords = data.lwpRecords;
      }
      // update machineStatuses setiap kali data LWP/records berubah
      updateMachineStatuses();
    } catch (error) {
      errorMessage =
        error instanceof Error ? error.message : "An error occurred";
      console.error("Error loading data:", error);
    } finally {
      isLoading = false;
    }
  }

  // Handle form submission for new pressing entry
  async function handleSubmit(e?: Event) {
    if (e) {
      e.preventDefault();
    }

    errorMessage = "";
    successMessage = "";

    // Validation
    if (!formData.partName.trim() || !formData.quantity.trim()) {
      errorMessage = "Part Name and Quantity are required";
      return;
    }

    try {
      const authToken = get(auth).token;
      if (!authToken) {
        errorMessage = "No authentication token found";
        return;
      }

      const response = await fetch(`${API_URL}/api/pressing`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
        body: JSON.stringify({
          partName: formData.partName.trim(),
          quantity: parseFloat(formData.quantity),
          reject: formData.reject ? parseFloat(formData.reject) : 0,
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to submit pressing data");
      }

      successMessage = "Data pressing berhasil disimpan!";

      // Clear form
      formData.partName = "";
      formData.quantity = "";
      formData.reject = "";

      // Reload data
      await loadData();

      // Clear success message after 3 seconds
      setTimeout(() => {
        successMessage = "";
      }, 3000);
    } catch (error) {
      errorMessage =
        error instanceof Error ? error.message : "An error occurred";
      console.error("Error submitting pressing data:", error);
    }
  }

  // Handle LWP form submission
  async function handleLwpSubmit(e?: Event) {
    if (e) {
      e.preventDefault();
    }

    errorMessage = "";
    successMessage = "";

    // Validation
    if (!lwpFormData.noMesin.trim() || !lwpFormData.partName.trim() || !lwpFormData.hasilOk.trim()) {
      errorMessage = "No Mesin, Part Name, dan Hasil OK harus diisi";
      return;
    }

    try {
      const authToken = get(auth).token;
      if (!authToken) {
        errorMessage = "No authentication token found";
        return;
      }

      const response = await fetch(`${API_URL}/api/lwp`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
        body: JSON.stringify(lwpFormData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to submit LWP data");
      }

      successMessage = "Data LWP berhasil disimpan!";

      // Reset form
      lwpFormData = {
        noMesin: "",
        tanggal: new Date().toISOString().split('T')[0],
        shift: "I",
        nik: employee.nik,
        partName: "",
        kodePart: "",
        jamMulai: "",
        jamSelesai: "",
        hasilOk: "",
        ng: "",
        klasifikasiReject: "Cacat Permukaan",
      };
      isLwpFormOpen = false;

      // Reload data
      await loadData();

      setTimeout(() => {
        successMessage = "";
      }, 3000);
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : "An error occurred";
      console.error("Error submitting LWP data:", error);
    }
  }

  // Auto-fill LWP form dari barcode scan
  function autoFillLwpFromBarcode(barcodeData: any) {
    if (barcodeData.partName) lwpFormData.partName = barcodeData.partName;
    if (barcodeData.kodePart) lwpFormData.kodePart = barcodeData.kodePart;
    if (barcodeData.noMesin) lwpFormData.noMesin = barcodeData.noMesin;
    if (barcodeData.quantity) lwpFormData.hasilOk = barcodeData.quantity;
    isLwpFormOpen = true;
  }

  function getTotalOkForMachine(machine: any) {
    return machine.details.reduce((sum: number, d: any) => sum + d.hasilOk, 0);
  }

  function getTotalNgForMachine(machine: any) {
    return machine.details.reduce((sum: number, d: any) => sum + d.ng, 0);
  }

  function getEfficiencyForMachine(machine: any) {
    const totalOk = getTotalOkForMachine(machine);
    const totalNG = getTotalNgForMachine(machine);
    const total = totalOk + totalNG;
    return total > 0 ? ((totalOk / total) * 100).toFixed(1) : 0;
  }

  function handleLogout() {
    if (confirm("Apakah Anda yakin ingin keluar?")) {
      window.location.href = "/";
    }
  }

  function getEfficiencyColor(efficiency: number) {
    if (efficiency >= 100)
      return "text-emerald-600 bg-emerald-50 border-emerald-100";
    if (efficiency >= 80) return "text-amber-600 bg-amber-50 border-amber-100";
    return "text-rose-600 bg-rose-50 border-rose-100";
  }

  function getBarColor(efficiency: number) {
    if (efficiency >= 100)
      return "bg-emerald-500 from-emerald-500 to-emerald-400";
    if (efficiency >= 80) return "bg-amber-500 from-amber-500 to-amber-400";
    return "bg-rose-500 from-rose-500 to-rose-400";
  }

  function getStatusColor(status: string) {
    if (status === "Proses") return "bg-blue-50 text-blue-700 border-blue-100";
    if (status === "Selesai")
      return "bg-slate-100 text-slate-600 border-slate-200";
    return "bg-amber-50 text-amber-700 border-amber-100";
  }

  function handleRowClick() {
    window.location.href = "/barcode";
  }

  function handleScanMachine() {
    window.location.href = "/scan-mesin";
  }

  // Get color class for reject classification
  function getRejectColor(classify: string) {
    if (classify === "Cacat Permukaan")
      return "bg-red-50 text-red-700 border-red-200";
    if (classify === "Dimensi") return "bg-yellow-50 text-yellow-700 border-yellow-200";
    if (classify === "Goresan") return "bg-orange-50 text-orange-700 border-orange-200";
    return "bg-slate-50 text-slate-700 border-slate-200";
  }

  // Load data on mount
  onMount(() => {
    loadData();
  });
</script>

<div
  class="min-h-screen bg-slate-50 font-sans text-slate-800 pb-12 relative selection:bg-indigo-100 selection:text-indigo-800"
>
  <div
    class="absolute inset-0 z-0 opacity-[0.03] pointer-events-none"
    style="background-image: radial-gradient(#4f46e5 1px, transparent 1px); background-size: 24px 24px;"
  ></div>

  <main
    class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-4 md:mt-6 relative z-10 space-y-6"
  >
    <div
      class="bg-white rounded-2xl shadow-md border border-slate-100 overflow-visible group hover:shadow-lg transition-all duration-500 mx-auto w-full"
    >
      <div class="p-6 md:p-10">
        <div
          class="flex flex-col md:flex-row gap-6 md:gap-8 items-center md:items-start"
        >
          <div
            class="relative shrink-0 -mt-16 md:-mt-8 group-hover:-translate-y-2 transition-transform duration-500"
          >
            <div
              class="w-28 h-28 md:w-44 md:h-44 rounded-2xl overflow-hidden border-4 md:border-[6px] border-white shadow-lg ring-1 ring-slate-100 mt-10"
            >
              <img
                src={employee.photo}
                alt={employee.name}
                class="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-700"
              />
            </div>
            <div
              class="absolute bottom-2 right-2 md:bottom-4 md:right-4 bg-white p-1.5 rounded-full shadow-md"
            >
              <div
                class="w-4 h-4 md:w-5 md:h-5 bg-emerald-500 rounded-full border-2 md:border-[3px] border-white animate-pulse"
              ></div>
            </div>
          </div>

          <div class="flex-1 text-center md:text-left w-full pt-2">
            <div
              class="flex flex-col md:flex-row justify-between items-center md:items-start gap-4"
            >
              <div>
                <h2
                  class="text-2xl md:text-4xl font-bold text-slate-800 tracking-tight mb-2"
                >
                  {employee.name}
                </h2>
                <div
                  class="flex flex-wrap justify-center md:justify-start gap-2 mb-4"
                >
                  <span
                    class="px-3 py-1 bg-indigo-50 text-indigo-700 text-xs md:text-sm font-bold uppercase tracking-wide rounded-full border border-indigo-100"
                    >{employee.position}</span
                  >
                  <span
                    class="px-3 py-1 bg-slate-100 text-slate-600 text-xs md:text-sm font-semibold rounded-full border border-slate-200"
                    >{employee.department}</span
                  >
                </div>
              </div>
              <div class="hidden md:block text-right">
                <div class="text-sm font-medium text-slate-400">Hari ini</div>
                <div class="text-xl font-bold text-slate-700">
                  {new Date().toLocaleDateString("id-ID", {
                    weekday: "long",
                    day: "numeric",
                    month: "long",
                  })}
                </div>
              </div>
            </div>

            <div
              class="grid grid-cols-1 md:grid-cols-3 gap-3 border-t border-slate-100 pt-5 mt-2"
            >
              <div
                class="flex items-center gap-3 p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-default border border-transparent hover:border-slate-100"
              >
                <div
                  class="w-10 h-10 rounded-lg bg-slate-100 text-slate-500 flex items-center justify-center shrink-0"
                >
                  <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2"
                    /></svg
                  >
                </div>
                <div class="text-left overflow-hidden">
                  <p
                    class="text-xs text-slate-400 font-bold uppercase truncate"
                  >
                    NIK ID
                  </p>
                  <p class="font-mono font-medium text-slate-800 truncate">
                    {employee.nik}
                  </p>
                </div>
              </div>
              <div
                class="flex items-center gap-3 p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-default border border-transparent hover:border-slate-100"
              >
                <div
                  class="w-10 h-10 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0"
                >
                  <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                    /></svg
                  >
                </div>
                <div class="text-left overflow-hidden">
                  <p
                    class="text-xs text-slate-400 font-bold uppercase truncate"
                  >
                    Telepon
                  </p>
                  <p class="font-medium text-slate-800 truncate">
                    {employee.phone}
                  </p>
                </div>
              </div>
              <div
                class="flex items-center gap-3 p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-default border border-transparent hover:border-slate-100"
              >
                <div
                  class="w-10 h-10 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center shrink-0"
                >
                  <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    /></svg
                  >
                </div>
                <div class="text-left overflow-hidden">
                  <p
                    class="text-xs text-slate-400 font-bold uppercase truncate"
                  >
                    Shift
                  </p>
                  <p class="font-medium text-slate-800 truncate">
                    Shift 1 (07:00 - 15:00)
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-3xl shadow-sm border border-slate-100 p-6 lg:p-8 mb-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h3 class="font-bold text-slate-800 text-lg">Status Mesin</h3>
          <p class="text-sm text-slate-500">Ringkasan pekerjaan saat ini per mesin</p>
        </div>
        <button onclick={loadData} class="text-xs bg-slate-100 px-3 py-1 rounded-md text-slate-600">Refresh</button>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="bg-slate-100 border-y border-slate-200">
              <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">No. Mesin</th>
              <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">Item</th>
              <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">No. Lot</th>
              <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">Operator</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            {#each machineStatuses as ms}
              <tr class="hover:bg-slate-50 transition-colors">
                <td class="px-4 py-3 font-mono font-bold text-indigo-600">{ms.noMesin}</td>
                <td class="px-4 py-3">{ms.item}</td>
                <td class="px-4 py-3 font-mono text-slate-700">{ms.noLot}</td>
                <td class="px-4 py-3">{ms.operator}</td>
              </tr>
            {/each}

            {#if machineStatuses.length === 0}
              <tr>
                <td colspan="4" class="px-4 py-6 text-center text-slate-400 italic">Belum ada status mesin.</td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div
        class="bg-white rounded-3xl p-6 lg:p-8 shadow-sm border border-slate-100 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 flex flex-col h-full"
      >
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="font-bold text-slate-800 text-lg">Target Hari Ini</h3>
            <p class="text-sm text-slate-500">Progress pressing lot</p>
          </div>
          <div class="bg-indigo-50 p-2 rounded-lg text-indigo-600">
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              /></svg
            >
          </div>
        </div>

        <div class="flex-1 flex flex-col items-center justify-center py-4">
          <div class="relative w-48 h-48 lg:w-56 lg:h-56 group cursor-default">
            <svg class="w-full h-full transform -rotate-90">
              <circle
                cx="50%"
                cy="50%"
                r="45%"
                stroke="currentColor"
                stroke-width="15"
                fill="transparent"
                class="text-slate-100"
              />
              <circle
                cx="50%"
                cy="50%"
                r="45%"
                stroke="currentColor"
                stroke-width="15"
                fill="transparent"
                stroke-dasharray={440}
                stroke-dashoffset={440 -
                  440 * (monthlyData.todayCompleted / monthlyData.todayTarget)}
                class="text-indigo-600 transition-all duration-1000 ease-out drop-shadow-lg"
                stroke-linecap="round"
              />
            </svg>
            <div
              class="absolute inset-0 flex flex-col items-center justify-center transition-transform duration-300 group-hover:scale-110"
            >
              <span class="text-5xl lg:text-6xl font-bold text-slate-800"
                >{monthlyData.todayCompleted}</span
              >
              <span
                class="text-xs text-slate-400 font-bold uppercase tracking-wider mt-1"
                >/ {monthlyData.todayTarget} Lot</span
              >
            </div>
          </div>
        </div>

        <div class="mt-4 pt-4 border-t border-slate-100 text-center">
          <p class="text-sm text-slate-600">
            Kurang <strong class="text-rose-500"
              >{monthlyData.todayTarget - monthlyData.todayCompleted} lot</strong
            > lagi untuk mencapai target.
          </p>
        </div>
      </div>

      <div
        class="bg-white rounded-3xl shadow-sm border border-slate-100 p-6 lg:p-8 flex flex-col h-full"
      >
        <div class="flex justify-between items-center mb-6">
          <h3 class="font-bold text-slate-800 text-lg">
            Produktivitas Mingguan
          </h3>
          <span
            class="text-xs font-medium bg-slate-100 text-slate-500 px-2 py-1 rounded-md"
            >7 Hari Terakhir</span
          >
        </div>

        <!-- Mobile View (Horizontal Scroll) -->
        <div class="md:hidden -mx-6 px-6 mb-4">
          <div class="flex overflow-x-auto gap-3 pb-2 scrollbar-hide">
            {#each dailyCompounds as day}
              {#if day.short !== "Min"}
                <div
                  class="group flex flex-col items-center justify-end shrink-0 w-16 relative"
                >
                  <div
                    class="absolute -top-8 opacity-0 group-hover:opacity-100 transition-opacity bg-slate-800 text-white text-xs py-1 px-2 rounded pointer-events-none z-10 whitespace-nowrap"
                  >
                    {day.target}/{day.actual}
                  </div>
                  <div
                    class="relative w-10 h-40 bg-slate-50 rounded-t-lg overflow-hidden flex items-end shadow-sm border border-slate-100"
                  >
                    <div
                      class={`w-full rounded-t-lg bg-gradient-to-t transition-all duration-1000 ease-out ${getBarColor(day.efficiency)}`}
                      style={`height: ${(day.actual / maxChartValue) * 100}%`}
                    ></div>
                  </div>
                  <div class="mt-2 text-center">
                    <p class="text-xs font-bold text-slate-600">{day.short}</p>
                    <p class="text-[10px] text-slate-400 font-mono mt-0.5">
                      {day.actual}L
                    </p>
                  </div>
                </div>
              {/if}
            {/each}
          </div>
        </div>

        <div
          class="hidden md:flex flex-1 items-end justify-between gap-2 md:gap-4 h-52 md:h-64 pb-2 border-b border-slate-100"
        >
          {#each dailyCompounds as day}
            {#if day.short !== "Min"}
              <div
                class="group flex flex-col items-center justify-end h-full w-full relative"
              >
                <div
                  class="absolute -top-10 opacity-0 group-hover:opacity-100 transition-opacity bg-slate-800 text-white text-xs py-1 px-2 rounded pointer-events-none mb-2 z-10 whitespace-nowrap"
                >
                  Target: {day.target} | Actual: {day.actual}
                  <div
                    class="absolute -bottom-1 left-1/2 transform -translate-x-1/2 w-2 h-2 bg-slate-800 rotate-45"
                  ></div>
                </div>
                <div
                  class="relative w-full max-w-10 h-full bg-slate-50 rounded-t-lg overflow-hidden flex items-end"
                >
                  <div
                    class={`w-full rounded-t-lg bg-gradient-to-t transition-all duration-1000 ease-out ${getBarColor(day.efficiency)}`}
                    style={`height: ${(day.actual / maxChartValue) * 100}%`}
                  ></div>
                </div>

                <div class="mt-3 text-center">
                  <p
                    class="text-xs font-bold text-slate-500 group-hover:text-indigo-600 transition-colors"
                  >
                    {day.short}
                  </p>
                  <p class="text-[10px] text-slate-400 font-mono mt-0.5">
                    {day.actual} Lot
                  </p>
                </div>
              </div>
            {/if}
          {/each}
        </div>

        <div
          class="mt-4 md:mt-4 flex flex-wrap justify-center md:justify-start gap-3 md:gap-4 text-xs text-slate-400 border-t border-slate-100 pt-4 md:pt-0 md:border-t-0"
        >
          <span class="flex items-center gap-2"
            ><span class="w-3 h-3 rounded-full bg-emerald-500 shadow-sm"></span>
            <span class="hidden sm:inline">&gt;100%</span><span
              class="sm:hidden">Optimal</span
            ></span
          >
          <span class="flex items-center gap-2"
            ><span class="w-3 h-3 rounded-full bg-amber-500 shadow-sm"></span>
            <span class="hidden sm:inline">&gt;80%</span><span class="sm:hidden"
              >Baik</span
            ></span
          >
          <span class="flex items-center gap-2"
            ><span class="w-3 h-3 rounded-full bg-rose-500 shadow-sm"></span>
            <span class="hidden sm:inline">&lt;80%</span><span class="sm:hidden"
              >Perlu Tuning</span
            ></span
          >
        </div>
      </div>
    </div>

    <div
      class="bg-white rounded-3xl shadow-sm border border-slate-100 p-6 lg:p-8"
    >
      <a href="/scan-barcode-prs" class="block">
        <button
          class="w-full bg-indigo-600 rounded-3xl text-white font-bold py-3 active:scale-95 shadow-lg shadow-indigo-300/50 hover:bg-indigo-700 hover:shadow-indigo-400/50 transition-all duration-300 cursor-pointer"
        >
          Scan KPCP
        </button>
      </a>
    </div>

    <!-- Permintaan Tim Press Table -->
    

    <!-- LWP (Laporan Waktu Produksi) Section -->
    <div class="space-y-6">
      <div class="bg-white rounded-2xl border-2 border-indigo-100 shadow-sm overflow-hidden">
        <div class="bg-indigo-50 px-6 py-4 border-b border-indigo-100 flex justify-between items-center">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-indigo-600 rounded-lg text-white">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-indigo-900">LWP - Laporan Waktu Produksi</h2>
              <p class="text-xs text-indigo-700">Riwayat produksi harian per mesin</p>
            </div>
          </div>
        </div>

        <!-- LWP Records Display -->
        <div class="divide-y divide-slate-100">
          {#each lwpRecords as machine}
            <div class="p-6">
              <!-- Machine Header Info -->
              <div class="mb-4 p-4 bg-gradient-to-r from-indigo-50 to-blue-50 rounded-xl border border-indigo-100">
                <div class="grid grid-cols-2 md:grid-cols-6 gap-3">
                  <div>
                    <p class="text-xs text-slate-500 font-bold uppercase">No. Mesin</p>
                    <p class="font-bold text-slate-800 text-lg">{machine.noMesin}</p>
                  </div>
                  <div>
                    <p class="text-xs text-slate-500 font-bold uppercase">Tanggal</p>
                    <p class="font-bold text-slate-800">
                      {new Date(machine.tanggal).toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' })}
                    </p>
                  </div>
                  <div>
                    <p class="text-xs text-slate-500 font-bold uppercase">Shift</p>
                    <p class="font-bold text-slate-800">{machine.shift}</p>
                  </div>
                  <div>
                    <p class="text-xs text-slate-500 font-bold uppercase">NIK Operator</p>
                    <p class="font-mono font-bold text-indigo-600">{machine.nik}</p>
                  </div>
                  <div>
                    <p class="text-xs text-slate-500 font-bold uppercase">Part Name</p>
                    <p class="font-bold text-slate-800">{machine.partName}</p>
                  </div>
                  <div>
                    <p class="text-xs text-slate-500 font-bold uppercase">Kode Part</p>
                    <p class="font-mono font-bold text-blue-600 bg-blue-50 px-2 py-1 rounded">{machine.kodePart}</p>
                  </div>
                </div>
              </div>

              <!-- Detail Records Table -->
              <div class="overflow-x-auto -mx-6 px-6 mb-4">
                <table class="w-full text-left text-sm">
                  <thead>
                    <tr class="bg-slate-100 border-y border-slate-200">
                      <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">No. Lot</th>
                      <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">Jam Mulai</th>
                      <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">Jam Selesai</th>
                      <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase text-right">Hasil OK</th>
                      <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase text-right">NG</th>
                      <th class="px-4 py-3 font-semibold text-slate-600 text-xs uppercase">Klasifikasi Reject</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    {#each machine.details as detail}
                      <tr class="hover:bg-slate-50 transition-colors">
                        <td class="px-4 py-3">
                          <span class="font-mono font-bold text-indigo-600 bg-indigo-50 px-2 py-1 rounded text-xs">
                            {detail.noLot}
                          </span>
                        </td>
                        <td class="px-4 py-3 font-medium text-slate-700">{detail.jamMulai}</td>
                        <td class="px-4 py-3 font-medium text-slate-700">{detail.jamSelesai}</td>
                        <td class="px-4 py-3 text-right">
                          <span class="font-bold text-emerald-600 bg-emerald-50 px-3 py-1 rounded inline-block min-w-fit">
                            {detail.hasilOk}
                          </span>
                        </td>
                        <td class="px-4 py-3 text-right">
                          <span class="font-bold text-rose-600 bg-rose-50 px-3 py-1 rounded inline-block min-w-fit">
                            {detail.ng}
                          </span>
                        </td>
                        <td class="px-4 py-3">
                          <div class="flex flex-wrap gap-1">
                            {#each detail.klasifikasiReject as reject}
                              <span class={`px-2 py-1 rounded text-xs font-bold border ${getRejectColor(reject.jenis)}`}>
                                {reject.jenis} ({reject.qty})
                              </span>
                            {/each}
                          </div>
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>
              <!-- Summary Stats -->
              <div class="grid grid-cols-2 md:grid-cols-4 gap-3 p-4 bg-slate-50 rounded-lg border border-slate-200">
                <div class="text-center">
                  <p class="text-xs text-slate-500 font-bold uppercase mb-1">Total OK</p>
                  <p class="text-2xl font-bold text-emerald-600">{getTotalOkForMachine(machine)}</p>
                </div>
                <div class="text-center">
                  <p class="text-xs text-slate-500 font-bold uppercase mb-1">Total NG</p>
                  <p class="text-2xl font-bold text-rose-600">{getTotalNgForMachine(machine)}</p>
                </div>
                <div class="text-center">
                  <p class="text-xs text-slate-500 font-bold uppercase mb-1">Jumlah Records</p>
                  <p class="text-2xl font-bold text-blue-600">{machine.details.length}</p>
                </div>
                <div class="text-center">
                  <p class="text-xs text-slate-500 font-bold uppercase mb-1">Efisiensi</p>
                  <p class="text-2xl font-bold text-indigo-600">{getEfficiencyForMachine(machine)}%</p>
                </div>
              </div>
            </div>
          {/each}

          {#if lwpRecords.length === 0}
            <div class="p-8 text-center text-slate-400 italic">
              Belum ada data LWP untuk hari ini.
            </div>
          {/if}
        </div>
      </div>
    </div>
  </main>

  <!-- FAB Button for Scan Machine -->
  <div class="fixed bottom-6 right-6 z-50 md:bottom-8 md:right-8">
    <button
      onclick={handleScanMachine}
      class="group relative w-14 h-14 md:w-16 md:h-16 bg-gradient-to-br from-indigo-600 to-indigo-700 rounded-full shadow-xl hover:shadow-2xl hover:scale-110 active:scale-95 transition-all duration-300 flex items-center justify-center text-white hover:from-indigo-700 hover:to-indigo-800"
      title="Scan Mesin"
    >
      <svg class="w-7 h-7 md:w-8 md:h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      <span class="hidden md:block absolute right-full mr-3 bg-slate-800 text-white text-sm font-semibold px-3 py-2 rounded-lg whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">
        Scan Mesin
      </span>
    </button>
  </div>
</div>

<style>
  /* Hanya utility sederhana, layout responsive ditangani oleh Tailwind Classes */
  /* Hide Scrollbar for cleaner look */
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>
