<script>
  import Swal from "sweetalert2";

  // Data & Logic Script (Tidak berubah dari versi sebelumnya)
  const employee = {
    name: "Alby Nizam P.",
    position: "Senior Pressing Specialist",
    department: "Production - Pressing",
    nik: "KRTP-2023-0456",
    phone: "+62 812-3456-7890",
    photo: "https://i.pinimg.com/550x/26/38/08/2638086da29fccffa32c5666ea77ce09.jpg",
  };

  let dailyCompounds = [
    { day: 'Senin', short: 'Sen', target: 4, actual: 4, efficiency: 100 },
    { day: 'Selasa', short: 'Sel', target: 4, actual: 4, efficiency: 88 },
    { day: 'Rabu', short: 'Rab', target: 4, actual: 5, efficiency: 125 },
    { day: 'Kamis', short: 'Kam', target: 4, actual: 4, efficiency: 100 },
    { day: 'Jumat', short: 'Jum', target: 4, actual: 4, efficiency: 100 },
    { day: 'Sabtu', short: 'Sab', target: 3, actual: 2, efficiency: 67 },
    { day: 'Minggu', short: 'Min', target: 0, actual: 0, efficiency: 0 }
  ];

  let monthlyData = {
    completed: 18,
    target: 20,
    efficiency: 90,
    todayCompleted: 3,
    todayTarget: 4
  };

  let scannedLot = '';
  let scanMessage = '';

  function handleLogout() {
    Swal.fire({
      title: 'logout',
      text: 'Apakah Anda yakin ingin logout?',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Ya, Logout',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        window.location.href = '../';
      }
    })
  }

  /**
   * @param {{ target: { value: string; }; }} e
   */
  function handleScanInput(e) {
    scannedLot = e.target.value;
    if (scannedLot.trim()) {
      scanMessage = `Lot ${scannedLot} berhasil discan!`;
      setTimeout(() => scanMessage = '', 2000);
      scannedLot = '';
    }
  }

  /**
   * @param {number} efficiency
   */
  function getEfficiencyColor(efficiency) {
    if (efficiency >= 100) return 'text-emerald-600 bg-emerald-50 border-emerald-100';
    if (efficiency >= 80) return 'text-amber-600 bg-amber-50 border-amber-100';
    return 'text-rose-600 bg-rose-50 border-rose-100';
  }

  /**
   * @param {number} efficiency
   */
  function getBarColor(efficiency) {
    if (efficiency >= 100) return 'bg-emerald-500';
    if (efficiency >= 80) return 'bg-amber-500';
    return 'bg-rose-500';
  }
</script>

<div class="min-h-screen bg-slate-50 font-sans text-slate-800 pb-12 relative selection:bg-indigo-100 selection:text-indigo-800">
  
  <div class="absolute inset-0 z-0 opacity-[0.03] pointer-events-none" 
       style="background-image: radial-gradient(#4f46e5 1px, transparent 1px); background-size: 24px 24px;">
  </div>

  <header class="sticky top-0 md:static z-50 transition-all duration-300">
    

    <div class="bg-white/80 backdrop-blur-md md:bg-transparent border-b border-slate-200 md:border-none shadow-sm md:shadow-none">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 md:py-8 flex justify-between items-center">
            <div class="flex items-center gap-4">
                <div class="md:hidden w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center text-white font-bold shadow-lg shadow-indigo-500/30">D</div>
                <div>
                    <h1 class="text-xl md:text-3xl font-bold text-slate-800 md:text-white tracking-tight">Dashboard Karyawan</h1>
                    <div class="flex items-center gap-2 mt-1">
                        <span class="relative flex h-2 w-2 md:hidden">
                          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                          <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                        </span>
                        <p class="text-xs text-slate-500 md:text-indigo-200 font-medium">Selamat datang kembali, Semangat Bekerja!</p>
                    </div>
                </div>
            </div>

            <button 
                on:click={handleLogout}
                class="group relative overflow-hidden px-5 py-2.5 rounded-xl transition-all duration-300
                       bg-white border border-slate-200 text-slate-600 hover:text-rose-600 hover:border-rose-200 hover:shadow-lg
                       md:bg-white/10 md:border-white/20 md:text-white md:hover:bg-white/20"
            >
                <div class="flex items-center gap-2 relative z-10">
                    <span class="text-sm font-semibold">Logout</span>
                    <svg class="w-4 h-4 transition-transform group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                </div>
            </button>
        </div>
    </div>
  </header>

  <main class="max-w-7xl mx-auto px-2 sm:px-4 lg:px-8 mt-4 md:-mt-24 relative z-10 space-y-8">
  <div class="bg-white rounded-2xl shadow-md border border-slate-100 overflow-visible group hover:shadow-lg transition-all duration-500 mx-auto w-full md:w-auto mt-4 md:mt-0">
    <div class="p-4 md:p-10">
      <div class="flex flex-col md:flex-row gap-4 md:gap-8 items-center md:items-start">
        <div class="relative shrink-0 -mt-12 md:-mt-16 group-hover:-translate-y-2 transition-transform duration-500">
          <div class="w-24 h-24 md:w-44 md:h-44 rounded-2xl overflow-hidden border-4 md:border-[6px] border-white shadow-lg ring-1 ring-slate-100">
            <img src={employee.photo} alt={employee.name} class="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-700" />
          </div>
          <div class="absolute bottom-2 right-2 md:bottom-4 md:right-4 bg-white p-1.5 rounded-full shadow-md">
            <div class="w-4 h-4 md:w-5 md:h-5 bg-emerald-500 rounded-full border-2 md:border-[3px] border-white animate-pulse"></div>
          </div>
        </div>
        <div class="flex-1 text-center md:text-left w-full pt-2">
          <div class="flex flex-col md:flex-row justify-between items-center md:items-start gap-2 md:gap-4">
            <div>
              <h2 class="text-xl md:text-4xl font-bold text-slate-800 tracking-tight mb-2">{employee.name}</h2>
              <div class="flex flex-wrap justify-center md:justify-start gap-2 md:gap-3 mb-4 md:mb-6">
                <span class="px-3 py-1 bg-indigo-50 text-indigo-700 text-xs md:text-sm font-bold uppercase tracking-wide rounded-full border border-indigo-100">{employee.position}</span>
                <span class="px-3 py-1 bg-slate-100 text-slate-600 text-xs md:text-sm font-semibold rounded-full border border-slate-200">{employee.department}</span>
              </div>
            </div>
            <div class="hidden md:block text-right">
              <div class="text-sm font-medium text-slate-400">Hari ini</div>
              <div class="text-xl font-bold text-slate-700">{new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long' })}</div>
            </div>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-2 md:gap-4 border-t border-slate-100 pt-4 md:pt-6 mt-2">
            <div class="flex items-center gap-2 md:gap-3 p-2 md:p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-default">
              <div class="w-8 h-8 md:w-10 md:h-10 rounded-lg bg-slate-100 text-slate-500 flex items-center justify-center">
                <svg class="w-4 h-4 md:w-5 md:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" /></svg>
              </div>
              <div class="text-left">
                <p class="text-xs text-slate-400 font-bold uppercase">NIK ID</p>
                <p class="font-mono font-medium text-slate-800 text-xs md:text-base">{employee.nik}</p>
              </div>
            </div>
            <div class="flex items-center gap-2 md:gap-3 p-2 md:p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-default">
              <div class="w-8 h-8 md:w-10 md:h-10 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center">
                <svg class="w-4 h-4 md:w-5 md:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
              </div>
              <div class="text-left">
                <p class="text-xs text-slate-400 font-bold uppercase">Telepon</p>
                <p class="font-medium text-slate-800 text-xs md:text-base">{employee.phone}</p>
              </div>
            </div>
            <div class="flex items-center gap-2 md:gap-3 p-2 md:p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-default">
              <div class="w-8 h-8 md:w-10 md:h-10 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
                <svg class="w-4 h-4 md:w-5 md:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              </div>
              <div class="text-left">
                <p class="text-xs text-slate-400 font-bold uppercase">Jam Kerja</p>
                <p class="font-medium text-slate-800 text-xs md:text-base">Shift 1 (07:00 - 15:00)</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6 lg:gap-8 mt-4">
        
        <div class="bg-white rounded-3xl p-6 lg:p-8 shadow-sm border border-slate-100 hover:shadow-lg hover:-translate-y-1 transition-all duration-300 flex flex-col h-full">
            <div class="flex justify-between items-start mb-6">
                <div>
                     <h3 class="font-bold text-slate-800 text-lg">Target Hari Ini</h3>
                     <p class="text-sm text-slate-500">Progress pressing lot</p>
                </div>
                <div class="bg-indigo-50 p-2 rounded-lg text-indigo-600">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
                </div>
            </div>
            
            <div class="flex-1 flex flex-col items-center justify-center py-2">
                <div class="relative w-48 h-48 lg:w-56 lg:h-56 group cursor-default">
                    <svg class="w-full h-full transform -rotate-90">
                        <circle cx="50%" cy="50%" r="45%" stroke="currentColor" stroke-width="15" fill="transparent" class="text-slate-100" />
                        <circle cx="50%" cy="50%" r="45%" stroke="currentColor" stroke-width="15" fill="transparent" 
                                stroke-dasharray={440} 
                                stroke-dashoffset={440 - (440 * (monthlyData.todayCompleted / monthlyData.todayTarget))} 
                                class="text-indigo-600 transition-all duration-1000 ease-out drop-shadow-lg" stroke-linecap="round" />
                    </svg>
                    <div class="absolute inset-0 flex flex-col items-center justify-center transition-transform duration-300 group-hover:scale-110">
                        <span class="text-5xl lg:text-6xl font-bold text-slate-800">{monthlyData.todayCompleted}</span>
                        <span class="text-sm text-slate-400 font-bold uppercase tracking-wider mt-1">/ {monthlyData.todayTarget} Lot</span>
                    </div>
                </div>
            </div>
            
            <div class="mt-6 pt-6 border-t border-slate-100 text-center">
                 <p class="text-sm text-slate-600">Kurang <strong class="text-rose-500">{monthlyData.todayTarget - monthlyData.todayCompleted} lot</strong> lagi untuk mencapai target.</p>
            </div>
        </div>

        <div class="bg-gradient-to-br from-indigo-600 to-violet-700 rounded-3xl p-6 lg:p-8 shadow-xl shadow-indigo-500/30 text-white relative overflow-hidden flex flex-col justify-between hover:shadow-2xl hover:scale-[1.02] transition-all duration-300 h-full">
            <div class="absolute top-0 right-0 -mr-8 -mt-8 w-40 h-40 bg-white opacity-10 rounded-full blur-2xl animate-pulse"></div>
            <div class="absolute bottom-0 left-0 -ml-8 -mb-8 w-40 h-40 bg-pink-500 opacity-20 rounded-full blur-2xl"></div>
            
            <div class="relative z-10">
                <h3 class="text-indigo-100 font-medium mb-8 flex items-center gap-2">
                    <span class="bg-white/20 p-1.5 rounded-lg">📊</span> Performance Bulan Ini
                </h3>
                
                <div class="space-y-8">
                    <div>
                        <p class="text-indigo-200 text-xs font-bold uppercase tracking-wider mb-1">Total Output</p>
                        <div class="flex items-baseline gap-2">
                            <h3 class="text-5xl font-bold">156</h3>
                            <span class="text-lg text-indigo-200">Lot Selesai</span>
                        </div>
                    </div>
                    
                    <div class="grid grid-cols-2 gap-4">
                         <div class="bg-white/10 rounded-2xl p-4 backdrop-blur-sm">
                            <p class="text-indigo-200 text-xs font-bold uppercase mb-1">Quality</p>
                            <h3 class="text-2xl font-bold">9.2<span class="text-sm opacity-60">/10</span></h3>
                         </div>
                         <div class="bg-white/10 rounded-2xl p-4 backdrop-blur-sm">
                            <p class="text-indigo-200 text-xs font-bold uppercase mb-1">Ranking</p>
                            <h3 class="text-2xl font-bold">#3<span class="text-sm opacity-60">/12</span></h3>
                         </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="bg-white rounded-xl shadow border border-slate-100 flex flex-col h-full mt-4 md:mt-0">
            <h3 class="font-bold text-slate-800 text-lg mb-6 flex items-center justify-between">
                Rekap Mingguan
                <span class="text-xs font-normal bg-slate-100 px-2 py-1 rounded text-slate-500">7 Hari Terakhir</span>
            </h3>
            
            <div class="flex-1 overflow-y-auto pr-2 scrollbar-hide space-y-3">
                {#each dailyCompounds as day}
                  {#if day.target > 0} 
                    <div class="group flex items-center justify-between p-3 rounded-2xl hover:bg-slate-50 transition-all border border-transparent hover:border-slate-100 cursor-default">
                        <div class="flex items-center gap-4">
                            <div class="h-10 w-1.5 bg-slate-100 rounded-full overflow-hidden">
                                <div class={`w-full rounded-full ${getBarColor(day.efficiency)}`} style={`height: ${day.efficiency}%`}></div>
                            </div>
                            <div>
                                <p class="font-bold text-slate-700 text-sm group-hover:text-indigo-700 transition-colors">{day.day}</p>
                                <p class="text-xs text-slate-400">{day.actual}/{day.target} Lot</p>
                            </div>
                        </div>
                        <span class={`text-xs font-bold px-3 py-1 rounded-lg border ${getEfficiencyColor(day.efficiency)}`}>
                            {day.efficiency}%
                        </span>
                    </div>
                  {/if}
                {/each}
            </div>
            <button class="w-full mt-4 py-3 text-sm font-semibold text-slate-500 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl transition-colors">
                Lihat Laporan Lengkap →
            </button>
        </div>
    </div>

    <div class="bg-white rounded-3xl shadow-sm border border-slate-100 overflow-hidden hover:shadow-md transition-shadow duration-300">
        <div class="p-6 lg:p-8 border-b border-slate-100 bg-slate-50/50 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <div>
                 <h3 class="font-bold text-slate-800 text-xl">Aktivitas Terkini</h3>
                 <p class="text-sm text-slate-500">Log aktivitas mesin dan pressing hari ini</p>
            </div>
            <div class="flex gap-2">
                <span class="w-3 h-3 bg-emerald-500 rounded-full animate-pulse"></span>
                <span class="text-xs font-bold text-emerald-600 uppercase tracking-wide">Live Updates</span>
            </div>
        </div>
        
        <div class="divide-y divide-slate-100">
            <div class="group p-6 lg:p-8 hover:bg-slate-50 transition-colors flex flex-col md:flex-row md:items-center gap-6">
                <div class="flex gap-4 md:w-1/3">
                    <div class="shrink-0">
                        <div class="w-12 h-12 rounded-2xl bg-emerald-100 text-emerald-600 flex items-center justify-center text-xl group-hover:scale-110 transition-transform">✅</div>
                    </div>
                    <div>
                        <h4 class="font-bold text-slate-800 text-base">Pressing NoLot HA2007</h4>
                        <span class="inline-flex mt-1 items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800">
                          Selesai
                        </span>
                    </div>
                </div>
                <div class="md:w-1/3">
                     <p class="text-sm text-slate-500 flex items-center gap-2">
                        <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" /></svg>
                        Mesin No. 17A
                     </p>
                </div>
                <div class="md:w-1/3 md:text-right">
                    <span class="text-sm font-bold text-slate-600">09:30 WIB</span>
                    <p class="text-xs text-slate-400 mt-1">Output Quality: OK</p>
                </div>
            </div>

            <div class="group p-6 lg:p-8 hover:bg-slate-50 transition-colors flex flex-col md:flex-row md:items-center gap-6">
                <div class="flex gap-4 md:w-1/3">
                    <div class="shrink-0">
                        <div class="w-12 h-12 rounded-2xl bg-blue-100 text-blue-600 flex items-center justify-center text-xl group-hover:rotate-90 transition-transform duration-700">⚙️</div>
                    </div>
                    <div>
                        <h4 class="font-bold text-slate-800 text-base">Pressing NoLot HA2008</h4>
                         <span class="inline-flex mt-1 items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                          Dalam Proses
                        </span>
                    </div>
                </div>
                <div class="md:w-1/3 flex flex-col justify-center">
                     <p class="text-sm text-slate-500 mb-2 flex items-center gap-2">
                        <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        Estimasi sisa waktu: 45 menit
                     </p>
                     <div class="w-full h-2 bg-slate-100 rounded-full overflow-hidden">
                        <div class="h-full bg-blue-500 rounded-full w-[75%] animate-[shimmer_2s_infinite] relative overflow-hidden">
                            <div class="absolute inset-0 bg-white/30 skew-x-12 -translate-x-full animate-[shimmer_1.5s_infinite]"></div>
                        </div>
                     </div>
                </div>
                <div class="md:w-1/3 md:text-right">
                    <span class="text-sm font-bold text-slate-600">Mulai 10:15 WIB</span>
                    <p class="text-xs text-slate-400 mt-1">Operator: Alby Nizam</p>
                </div>
            </div>
        </div>
        
        <div class="bg-slate-50 p-4 text-center border-t border-slate-100">
            <button class="text-sm font-semibold text-indigo-600 hover:text-indigo-800 transition-colors">Tampilkan Semua Aktivitas ↓</button>
        </div>
    </div>

    <!-- Scan Barcode Section -->
    <div class="bg-white rounded-3xl shadow-sm border border-slate-100 mt-4 p-6 flex flex-col items-center gap-4">
      <h3 class="font-bold text-slate-800 text-lg mb-2 flex items-center gap-2">
        <svg class="w-6 h-6 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><rect x="3" y="3" width="18" height="18" rx="2" stroke-width="2"/><path stroke-width="2" d="M7 7h.01M17 7h.01M7 17h.01M17 17h.01"/></svg>
        Scan Barcode Lot KPCP
      </h3>
      <button
        class="w-full max-w-xs px-4 py-3 rounded-xl bg-indigo-600 text-white font-bold text-lg shadow hover:bg-indigo-700 transition-all"
        on:click={() => window.location.href = '/scan-barcode'}
      >
        Mulai Scan Barcode
      </button>
      <p class="text-xs text-slate-400 mt-1">Klik tombol di atas untuk scan barcode nomor lot KPCP.</p>
    </div>
    <!-- End Scan Barcode Section -->

  </main>
</div>

<style>
  /* Utility Scrollbar */
  .scrollbar-hide::-webkit-scrollbar {
      display: none;
  }
  .scrollbar-hide {
      -ms-overflow-style: none;
      scrollbar-width: none;
  }
  
  /* Shimmer Animation for Progress Bar */
  @keyframes shimmer {
    100% { transform: translateX(100%); }
  }

  /* Responsive fix for mobile */
  @media (max-width: 640px) {
    .rounded-2xl, .rounded-3xl, .rounded-xl {
      border-radius: 1rem !important;
    }
    .shadow-md, .shadow-xl, .shadow-2xl, .shadow-sm, .shadow {
      box-shadow: 0 2px 8px 0 rgb(30 41 59 / 0.08) !important;
    }
    .mt-4 {
      margin-top: 1rem !important;
    }
    .p-4, .md\:p-10 {
      padding: 1rem !important;
    }
    .space-y-8 > :not([hidden]) ~ :not([hidden]) {
      margin-top: 1rem !important;
    }
    .grid-cols-3, .md\:grid-cols-3, .lg\:grid-cols-3 {
      grid-template-columns: 1fr !important;
    }
  }
</style>