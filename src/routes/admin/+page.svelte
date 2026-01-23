<script>
  import Swal from "sweetalert2";

  // Data Admin (Tanpa Statistik Pribadi)
  const admin = {
    name: "Budi Santoso",
    position: "Production Manager",
    department: "Plant A - Management",
    nik: "ADM-2023-001",
    photo: "https://i.pravatar.cc/300?img=11",
  };

  let deptSummary = [
    {
      id: "pressing",
      name: "Pressing",
      icon: "M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10", // Ikon Mesin Press
      totalStaff: 12,
      activeNow: 10,
      avgEfficiency: 88,
      target: 2000,
      actual: 1850,
      color: "amber",
    },
    {
      id: "cutting",
      name: "Cutting",
      icon: "M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z", // Ikon Gunting/Cutting
      totalStaff: 8,
      activeNow: 8,
      avgEfficiency: 95,
      target: 2500,
      actual: 2400,
      color: "blue",
    },
    {
      id: "finishing",
      name: "Finishing",
      icon: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
      totalStaff: 15,
      activeNow: 14,
      avgEfficiency: 92,
      target: 1800,
      actual: 1750,
      color: "emerald",
    },
  ];

  let employeeStats = [
    {
      id: 101,
      name: "Alby Nizam P.",
      dept: "Pressing",
      target: 50,
      actual: 45,
      efficiency: 90,
      status: "Active",
    },
    {
      id: 102,
      name: "Siti Aminah",
      dept: "Finishing",
      target: 100,
      actual: 105,
      efficiency: 105,
      status: "Active",
    },
    {
      id: 103,
      name: "Rizky Ramadhan",
      dept: "Cutting",
      target: 80,
      actual: 40,
      efficiency: 50,
      status: "Break",
    },
    {
      id: 104,
      name: "Dewi Lestari",
      dept: "Finishing",
      target: 100,
      actual: 98,
      efficiency: 98,
      status: "Active",
    },
    {
      id: 105,
      name: "Joko Anwar",
      dept: "Pressing",
      target: 50,
      actual: 20,
      efficiency: 40,
      status: "Issue",
    },
    {
      id: 106,
      name: "Sarah W.",
      dept: "Cutting",
      target: 80,
      actual: 85,
      efficiency: 106,
      status: "Active",
    },
  ];

  function handleLogout() {
    Swal.fire({
        title: 'Logout',
        text: 'Apakah Anda yakin ingin logout?',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Ya, Logout',
        cancelButtonText: 'Batal'
        }).then((result) => {
            if (result.isConfirmed) {
                Swal.fire(
                    'Logged Out!',
                    'Anda telah berhasil logout.',
                    'success'
                );
                window.location.href = '../';
            }
    })
  }

  /**
   * @param {string} status
   */
  function getStatusColor(status) {
    if (status === "Active")
      return "bg-emerald-100 text-emerald-700 ring-emerald-600/20";
    if (status === "Break")
      return "bg-slate-100 text-slate-700 ring-slate-600/20";
    return "bg-rose-100 text-rose-700 ring-rose-600/20 animate-pulse";
  }

  /**
   * @param {number} efficiency
   */
  function getEfficiencyColor(efficiency) {
    if (efficiency >= 100) return "text-emerald-600";
    if (efficiency >= 80) return "text-blue-600";
    if (efficiency >= 60) return "text-amber-600";
    return "text-rose-600";
  }

  // Helper dinamis untuk warna background kartu
  /**
   * @param {string} color
   */
  function getCardColor(color) {
    const map = {
      amber:
        "bg-amber-50 border-amber-100 hover:border-amber-300 text-amber-900",
      blue: "bg-blue-50 border-blue-100 hover:border-blue-300 text-blue-900",
      emerald:
        "bg-emerald-50 border-emerald-100 hover:border-emerald-300 text-emerald-900",
    };
    // @ts-ignore
    return map[color];
  }

  /**
   * @param {string} color
   */
  function getIconBg(color) {
    const map = {
      amber: "bg-amber-100 text-amber-600",
      blue: "bg-blue-100 text-blue-600",
      emerald: "bg-emerald-100 text-emerald-600",
    };
    // @ts-ignore
    return map[color];
  }

  /**
   * @param {string} color
   */
  function getBarColor(color) {
    const map = {
      amber: "bg-amber-500",
      blue: "bg-blue-500",
      emerald: "bg-emerald-500",
    };
    // @ts-ignore
    return map[color];
  }
</script>

<div class="min-h-screen bg-slate-50 font-sans text-slate-800 pb-12 relative">
  <div
    class="absolute inset-0 z-0 opacity-[0.03] pointer-events-none"
    style="background-image: radial-gradient(#64748b 1px, transparent 1px); background-size: 24px 24px;"
  ></div>

  <header class="sticky top-0 z-50 transition-all duration-300">
    <div
      class="bg-white/90 backdrop-blur-md border-b border-slate-200 shadow-sm"
    >
      <div
        class="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center"
      >
        <div class="flex items-center gap-4">
          <div
            class="w-10 h-10 bg-slate-800 rounded-xl flex items-center justify-center text-white font-bold shadow-lg shadow-slate-500/30"
          >
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
              /></svg
            >
          </div>
          <div>
            <h1 class="text-xl font-bold text-slate-800 tracking-tight">
              Production Admin
            </h1>
            <p class="text-xs text-slate-500 font-medium">
              Monitoring Dashboard
            </p>
          </div>
        </div>

        <button
          on:click={handleLogout}
          class="hidden md:flex items-center gap-2 px-4 py-2 rounded-xl bg-white border border-slate-200 text-slate-600 hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all cursor-pointer active:scale-95"
        >
          <span class="text-sm font-semibold">Logout</span>
          <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            ><path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
            /></svg
          >
        </button>
      </div>
    </div>
  </header>

  <main class="max-w-7xl mx-auto px-4 mt-6 space-y-6 relative z-10">
    <div
      class="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 flex flex-col md:flex-row items-center justify-between gap-6"
    >
      <div class="flex items-center gap-5 w-full md:w-auto">
        <img
          src={admin.photo}
          alt={admin.name}
          class="w-16 h-16 rounded-full object-cover border-2 border-white shadow-md ring-2 ring-slate-100"
        />
        <div>
          <h2 class="text-xl font-bold text-slate-800">Halo, {admin.name}</h2>
          <div class="flex items-center gap-2 text-sm text-slate-500 mt-1">
            <span
              class="font-medium bg-slate-100 px-2 py-0.5 rounded text-slate-600"
              >{admin.nik}</span
            >
            <span>•</span>
            <span>{admin.position}</span>
          </div>
        </div>
      </div>
      <div
        class="w-full md:w-auto flex items-center justify-between md:justify-end gap-8 border-t md:border-t-0 border-slate-100 pt-4 md:pt-0"
      >
        <div class="text-center md:text-right">
          <p class="text-xs text-slate-400 font-bold uppercase tracking-wider">
            Total Output Hari Ini
          </p>
          <p class="text-2xl font-bold text-slate-800">
            6,000 <span class="text-sm font-normal text-slate-500">Lot</span>
          </p>
        </div>
        <div class="text-center md:text-right">
          <p class="text-xs text-slate-400 font-bold uppercase tracking-wider">
            Active Staff
          </p>
          <p class="text-2xl font-bold text-emerald-600">
            32 <span class="text-sm font-normal text-slate-500">/ 35</span>
          </p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      {#each deptSummary as dept}
        <div
          class={`rounded-2xl p-6 border transition-all duration-300 hover:shadow-lg hover:-translate-y-1 ${getCardColor(dept.color)}`}
        >
          <div class="flex justify-between items-start mb-4">
            <div class={`p-3 rounded-xl ${getIconBg(dept.color)}`}>
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
                  d={dept.icon}
                />
              </svg>
            </div>
            <div class="text-right">
              <span class="block text-2xl font-bold">{dept.avgEfficiency}%</span
              >
              <span class="text-xs font-semibold opacity-70"
                >Efisiensi Rata-rata</span
              >
            </div>
          </div>

          <h3 class="text-lg font-bold mb-1">{dept.name} Department</h3>
          <p class="text-sm opacity-80 mb-4">
            {dept.activeNow} dari {dept.totalStaff} staff aktif
          </p>

          <div class="space-y-3">
            <div
              class="flex justify-between text-xs font-bold uppercase tracking-wide opacity-60"
            >
              <span>Progress Lot</span>
              <span>{dept.actual} / {dept.target}</span>
            </div>
            <div class="w-full h-2 bg-white/50 rounded-full overflow-hidden">
              <div
                class={`h-full rounded-full transition-all duration-1000 ${getBarColor(dept.color)}`}
                style={`width: ${(dept.actual / dept.target) * 100}%`}
              ></div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <div
      class="bg-white rounded-2xl shadow-sm border border-slate-200 flex flex-col"
    >
      <div
        class="p-6 border-b border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"
      >
        <div>
          <h3 class="font-bold text-slate-800 text-lg">Live Floor Activity</h3>
          <p class="text-sm text-slate-500">
            Monitoring real-time aktivitas semua divisi
          </p>
        </div>

        <div class="flex gap-2">
          <select
            class="text-sm border-slate-200 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option>Semua Divisi</option>
            <option>Pressing</option>
            <option>Cutting</option>
            <option>Finishing</option>
          </select>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/50 border-b border-slate-100">
              <th
                class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider"
                >Karyawan</th
              >
              <th
                class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider"
                >Divisi</th
              >
              <th
                class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider text-center"
                >Output Lot</th
              >
              <th
                class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider text-center"
                >Efisiensi</th
              >
              <th
                class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider text-center"
                >Status</th
              >
              <th
                class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider text-right"
                >Aksi</th
              >
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            {#each employeeStats as emp}
              <tr class="group hover:bg-slate-50 transition-colors">
                <td class="px-6 py-4">
                  <div class="flex items-center">
                    <div
                      class="h-8 w-8 rounded-full bg-slate-200 flex items-center justify-center text-xs font-bold text-slate-600 mr-3"
                    >
                      {emp.name.charAt(0)}
                    </div>
                    <div>
                      <p class="font-bold text-slate-700 text-sm">{emp.name}</p>
                      <p class="text-xs text-slate-400">ID: {emp.id}</p>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <span
                    class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border
                                ${emp.dept === "Pressing" ? "bg-amber-50 text-amber-700 border-amber-100" : ""}
                                ${emp.dept === "Cutting" ? "bg-blue-50 text-blue-700 border-blue-100" : ""}
                                ${emp.dept === "Finishing" ? "bg-emerald-50 text-emerald-700 border-emerald-100" : ""}
                            `}
                  >
                    {emp.dept}
                  </span>
                </td>
                <td class="px-6 py-4 text-center">
                  <div class="flex flex-col items-center">
                    <span class="font-bold text-slate-700">{emp.actual}</span>
                    <span class="text-[10px] text-slate-400"
                      >Target: {emp.target}</span
                    >
                  </div>
                </td>
                <td class="px-6 py-4 text-center">
                  <span
                    class={`font-bold ${getEfficiencyColor(emp.efficiency)}`}
                  >
                    {emp.efficiency}%
                  </span>
                </td>
                <td class="px-6 py-4 text-center">
                  <span
                    class={`inline-flex items-center px-2 py-1 rounded-md text-xs font-bold uppercase ring-1 ring-inset ${getStatusColor(emp.status)}`}
                  >
                    {emp.status === "Issue" ? "⚠️ Issue" : emp.status}
                  </span>
                </td>
                <td class="px-6 py-4 text-right">
                  <!-- svelte-ignore a11y_consider_explicit_label -->
                  <button
                    class="text-slate-400 hover:text-indigo-600 transition-colors"
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
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      /><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      /></svg
                    >
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <div
        class="p-4 border-t border-slate-100 bg-slate-50/50 rounded-b-2xl text-center"
      >
        <button
          class="text-sm font-semibold text-slate-500 hover:text-slate-800 transition-colors"
          >Lihat Semua Log Aktivitas</button
        >
      </div>
    </div>
  </main>
</div>

<style>
  /* Custom Scrollbar for Table if needed */
  .overflow-x-auto::-webkit-scrollbar {
    height: 8px;
  }
  .overflow-x-auto::-webkit-scrollbar-track {
    background: transparent;
  }
  .overflow-x-auto::-webkit-scrollbar-thumb {
    background-color: #cbd5e1;
    border-radius: 20px;
    border: 3px solid transparent;
    background-clip: content-box;
  }
</style>
