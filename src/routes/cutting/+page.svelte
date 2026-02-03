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

  // Form state for cutting data submission
  let formData = $state({
    partName: "",
    quantity: "",
    reject: "",
  });

  const employee = {
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

  // ===== LHC (Laporan Hasil Cutting) State =====
  let lhcData = $state({
    operatorName: "Tono Widiyanto",
    nik: "KRTP-2023-0456",
    tanggal: new Date().toLocaleDateString("id-ID"),
    shift: "I",
  });

  let cuttingDetails = $state<any[]>([
    {
      id: 1,
      kode: "PROD-001",
      compound: "A1",
      noLot: "KPCP-2309-A01",
      noPgc: "PGC-2309-001",
      waktuMulai: "10:30",
      waktuSelesai: "11:15",
      qty: 120,
      reject: 5,
    },
    {
      id: 2,
      kode: "PROD-002",
      compound: "B2",
      noLot: "KPCP-2309-A02",
      noPgc: "PGC-2309-002",
      waktuMulai: "10:15",
      waktuSelesai: "11:00",
      qty: 100,
      reject: 3,
    },
    {
      id: 3,
      kode: "PROD-001",
      compound: "A1",
      noLot: "KPCP-2309-B05",
      noPgc: "PGC-2309-003",
      waktuMulai: "09:45",
      waktuSelesai: "10:30",
      qty: 150,
      reject: 8,
    },
  ]);

  // Logic untuk Bar Chart Scale
  const maxChartValue = Math.max(
    ...dailyCompounds.map((d) => d.actual),
    ...dailyCompounds.map((d) => d.target),
    6
  );

  // ===== API Functions =====

  // Load today's cutting data from API
  async function loadData() {
    isLoading = true;
    errorMessage = "";

    try {
      const authToken = get(auth).token;
      if (!authToken) {
        errorMessage = "No authentication token found";
        return;
      }

      const response = await fetch(`${API_URL}/api/cutting/today`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to load cutting data");
      }

      const data = await response.json();
      recentScans = data.scans || recentScans;
      monthlyData = data.stats || monthlyData;
    } catch (error) {
      errorMessage =
        error instanceof Error ? error.message : "An error occurred";
      console.error("Error loading data:", error);
    } finally {
      isLoading = false;
    }
  }

  // Handle form submission for new cutting entry
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

      const response = await fetch(`${API_URL}/api/cutting`, {
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
        throw new Error(error.error || "Failed to submit cutting data");
      }

      successMessage = "Data cutting berhasil disimpan!";

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
      console.error("Error submitting cutting data:", error);
    }
  }

  function handleLogout() {
    if (confirm("Apakah Anda yakin ingin keluar?")) {
      window.location.href = "/";
    }
  }

  function handleScanKPCP() {
    window.location.href = "/barcode-ctng";
  }

  // Add new cutting entry from barcode scan
  function addCuttingEntry(scanData: any) {
    const newEntry = {
      id: cuttingDetails.length + 1,
      kode: scanData.kode || "",
      compound: scanData.compound || "",
      noLot: scanData.noLot || "",
      noPgc: scanData.noPgc || "",
      waktuMulai: scanData.waktuMulai || "",
      waktuSelesai: scanData.waktuSelesai || "",
      qty: scanData.qty || 0,
      reject: scanData.reject || 0,
    };
    cuttingDetails = [...cuttingDetails, newEntry];
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
      <div class="mb-8">
        <h3 class="font-bold text-slate-800 text-lg mb-6 flex items-center gap-2">
          <svg
            class="w-6 h-6 text-emerald-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            /></svg
            >
          LAPORAN HASIL CUTTING (LHC)
        </h3>

        <!-- Header Info Operator -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6 p-4 md:p-6 bg-gradient-to-br from-slate-50 to-slate-100 rounded-2xl border border-slate-200">
          <div>
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">
              Nama Operator
            </label>
            <input
              type="text"
              bind:value={lhcData.operatorName}
              class="w-full px-3 py-2 border border-slate-300 rounded-lg bg-white text-slate-800 font-medium focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
              disabled
            />
          </div>

          <div>
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">
              NIK
            </label>
            <input
              type="text"
              bind:value={lhcData.nik}
              class="w-full px-3 py-2 border border-slate-300 rounded-lg bg-white font-mono text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
              disabled
            />
          </div>

          <div>
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">
              Tanggal
            </label>
            <input
              type="text"
              bind:value={lhcData.tanggal}
              class="w-full px-3 py-2 border border-slate-300 rounded-lg bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
              disabled
            />
          </div>

          <div>
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-1">
              Shift
            </label>
            <input
              type="text"
              bind:value={lhcData.shift}
              class="w-full px-3 py-2 border border-slate-300 rounded-lg bg-white text-slate-800 font-bold focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
              disabled
            />
          </div>
        </div>

        <!-- Cutting Details Table -->
        <div class="overflow-x-auto rounded-xl border border-slate-200">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-gradient-to-r from-emerald-500 to-emerald-600 text-white">
                <th class="px-4 py-3 font-bold text-sm">No.</th>
                <th class="px-4 py-3 font-bold text-sm">Kode Produk</th>
                <th class="px-4 py-3 font-bold text-sm">Compound</th>
                <th class="px-4 py-3 font-bold text-sm">No. Lot</th>
                <th class="px-4 py-3 font-bold text-sm">No. PGC</th>
                <th class="px-4 py-3 font-bold text-sm">Waktu Mulai</th>
                <th class="px-4 py-3 font-bold text-sm">Waktu Selesai</th>
                <th class="px-4 py-3 font-bold text-sm text-right">Qty</th>
                <th class="px-4 py-3 font-bold text-sm text-right">Reject</th>
              </tr>
            </thead>
            <tbody class="text-sm divide-y divide-slate-200">
              {#each cuttingDetails as detail, index}
                <tr class="hover:bg-slate-50 transition-colors">
                  <td class="px-4 py-3 font-bold text-slate-800">
                    {index + 1}
                  </td>
                  <td class="px-4 py-3 font-mono font-bold text-slate-700">
                    {detail.kode}
                  </td>
                  <td class="px-4 py-3 font-bold text-slate-700">
                    {detail.compound}
                  </td>
                  <td class="px-4 py-3 font-mono text-slate-700 text-indigo-600 font-bold">
                    {detail.noLot}
                  </td>
                  <td class="px-4 py-3 font-mono text-slate-600">
                    {detail.noPgc}
                  </td>
                  <td class="px-4 py-3 text-slate-700 font-medium">
                    {detail.waktuMulai}
                  </td>
                  <td class="px-4 py-3 text-slate-700 font-medium">
                    {detail.waktuSelesai}
                  </td>
                  <td class="px-4 py-3 text-right font-bold text-slate-800">
                    {detail.qty}
                  </td>
                  <td class="px-4 py-3 text-right font-bold text-rose-600">
                    {detail.reject}
                  </td>
                </tr>
              {/each}

              {#if cuttingDetails.length === 0}
                <tr>
                  <td colspan="9" class="px-4 py-8 text-center text-slate-400">
                    Belum ada data cutting. Scan KPCP untuk menambahkan data.
                  </td>
                </tr>
              {/if}
            </tbody>
          </table>
        </div>

        <!-- Summary Footer -->
        <div class="mt-4 grid grid-cols-2 md:grid-cols-4 gap-3">
          <div class="p-3 bg-blue-50 rounded-lg border border-blue-200">
            <p class="text-xs text-blue-600 font-bold uppercase">Total Qty</p>
            <p class="text-xl font-bold text-blue-700">
              {cuttingDetails.reduce((sum, item) => sum + item.qty, 0)}
            </p>
          </div>
          <div class="p-3 bg-rose-50 rounded-lg border border-rose-200">
            <p class="text-xs text-rose-600 font-bold uppercase">Total Reject</p>
            <p class="text-xl font-bold text-rose-700">
              {cuttingDetails.reduce((sum, item) => sum + item.reject, 0)}
            </p>
          </div>
          <div class="p-3 bg-emerald-50 rounded-lg border border-emerald-200">
            <p class="text-xs text-emerald-600 font-bold uppercase">Total Good</p>
            <p class="text-xl font-bold text-emerald-700">
              {cuttingDetails.reduce((sum, item) => sum + (item.qty - item.reject), 0)}
            </p>
          </div>
          <div class="p-3 bg-amber-50 rounded-lg border border-amber-200">
            <p class="text-xs text-amber-600 font-bold uppercase">Efisiensi</p>
            <p class="text-xl font-bold text-amber-700">
              {cuttingDetails.length > 0
                ? (
                    ((cuttingDetails.reduce((sum, item) => sum + (item.qty - item.reject), 0) /
                      cuttingDetails.reduce((sum, item) => sum + item.qty, 0)) *
                      100) ||
                    0
                  ).toFixed(1)
                : 0}%
            </p>
          </div>
        </div>
      </div>
    </div>

    <div
      class="bg-white rounded-3xl shadow-sm border border-slate-100 p-6 lg:p-8"
    >
      <div
        class="flex flex-col md:flex-row md:items-center justify-between mb-6 gap-4"
      >
        <div>
          <h3 class="font-bold text-slate-800 text-lg flex items-center gap-2">
            <svg
              class="w-6 h-6 text-indigo-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><rect
                x="3"
                y="3"
                width="18"
                height="18"
                rx="2"
                stroke-width="2"
              /><path
                stroke-width="2"
                d="M7 7h.01M17 7h.01M7 17h.01M17 17h.01"
              /></svg
            >
            Antrian Scan KPCP
          </h3>
          <p class="text-sm text-slate-500">
            Daftar lot yang baru saja discan dan diproses. Klik tombol di bawah untuk scan lot baru.
          </p>
        </div>
        <button
          onclick={handleScanKPCP}
          class="flex items-center gap-2 px-4 md:px-6 py-3 md:py-2.5 bg-gradient-to-r from-indigo-600 to-indigo-700 text-white font-bold rounded-xl hover:from-indigo-700 hover:to-indigo-800 transition-all duration-300 shadow-lg hover:shadow-xl hover:-translate-y-0.5 active:translate-y-0 whitespace-nowrap text-sm md:text-base"
        >
          <svg
            class="w-5 h-5 md:w-5 md:h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 4v16m8-8H4"
            /></svg
          >
          Scan KPCP Baru
        </button>
      </div>

      <!-- Mobile View (Card Layout) -->
      <div class="md:hidden space-y-3">
        {#each recentScans as scan}
          <div
            class="p-4 border border-slate-100 rounded-lg bg-slate-50 transition-colors"
          >
            <div class="flex justify-between items-start mb-3">
              <span class="font-mono font-bold text-indigo-600">{scan.lot}</span
              >
              <span
                class={`px-2 py-1 rounded-full text-xs font-bold border ${getStatusColor(scan.status)}`}
              >
                {scan.status}
              </span>
            </div>
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-slate-500">Waktu Mulai:</span>
                <span class="font-medium text-slate-700"
                  >{scan.timeStart} WIB</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-slate-500">Waktu Selesai:</span>
                <span class="font-medium text-slate-700"
                  >{scan.timeEnd} WIB</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-slate-500">PIC:</span>
                <span class="font-medium text-slate-700">{scan.pic}</span>
              </div>
            </div>
          </div>
        {/each}
        {#if recentScans.length === 0}
          <div class="py-8 text-center text-slate-400">
            Belum ada data scan hari ini.
          </div>
        {/if}
      </div>

      <!-- Desktop View (Table Layout) -->
      <div class="hidden md:block overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr
              class="border-b border-slate-100 text-slate-400 text-xs uppercase tracking-wider"
            >
              <th class="pb-3 font-semibold">No. Lot</th>
              <th class="pb-3 font-semibold">Waktu Mulai</th>
              <th class="pb-3 font-semibold">Waktu Selesai</th>
              <th class="pb-3 font-semibold">PIC</th>
              <th class="pb-3 font-semibold text-right">Status</th>
            </tr>
          </thead>
          <tbody class="text-sm">
            {#each recentScans as scan}
              <tr
                class="group hover:bg-slate-50 transition-colors border-b last:border-0 border-slate-50"
              >
                <td
                  class="py-4 font-mono font-bold text-slate-700"
                >
                  {scan.lot}
                </td>
                <td class="py-4 text-slate-500">
                  {scan.timeStart} WIB
                </td>
                <td class="py-4 text-slate-500">
                  {scan.timeEnd} WIB
                </td>
                <td class="py-4 text-slate-600 font-medium">
                  <div class="flex items-center gap-2">
                    <div
                      class="w-6 h-6 rounded-full bg-slate-200 text-[10px] flex items-center justify-center font-bold text-slate-500"
                    >
                      {scan.pic.charAt(0)}
                    </div>
                    {scan.pic}
                  </div>
                </td>
                <td class="py-4 text-right">
                  <span
                    class={`px-3 py-1 rounded-full text-xs font-bold border ${getStatusColor(scan.status)}`}
                  >
                    {scan.status}
                  </span>
                </td>
              </tr>
            {/each}
            {#if recentScans.length === 0}
              <tr
                ><td colspan="5" class="py-8 text-center text-slate-400"
                  >Belum ada data scan hari ini.</td
                ></tr
              >
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  </main>
</div>

<style>
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>
