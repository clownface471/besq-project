<script lang="ts">
  import { onMount } from 'svelte';
  import Swal from 'sweetalert2';

  let lotNumber = '';
  let scanTime = '';
  let cycleCount = 0;
  let lotDetail = {
    lot: '',
    product: 'KPCP Standard',
    quantity: 500,
    status: 'Proses',
    startTime: '',
    estimatedEnd: '',
    supervisor: 'Tono Widiyanto',
    machine: 'Machine A-01',
    noCompound: 'CMP-2309-001',
    noLotMixing: 'MIX-2309-A05'
  };

  function handleCompleteCycle() {
    cycleCount++;
    Swal.fire({
      title: 'Cycle Selesai!',
      html: `<p class="text-lg font-semibold">Cycle <span class="text-amber-600">${cycleCount}</span> telah berhasil diselesaikan</p>`,
      icon: 'success',
      confirmButtonText: 'Lanjutkan',
      confirmButtonColor: '#d97706',
      allowOutsideClick: false
    });
  }

  function handleComplete() {
    Swal.fire({
      title: 'Konfirmasi Penyelesaian',
      text: 'Apakah proses KPCP ini sudah selesai?',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#4f46e5',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'Ya, Selesaikan',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        Swal.fire({
          title: 'Berhasil!',
          text: 'Data berhasil disimpan!',
          icon: 'success',
          confirmButtonColor: '#4f46e5'
        }).then(() => {
          window.location.href = '/cutting';
        });
      }
    });
  }

  onMount(() => {
    const queryParams = new URLSearchParams(window.location.search);
    lotNumber = queryParams.get('lot') || '';
    scanTime = localStorage.getItem('scanTime') || new Date().toLocaleTimeString('id-ID');
    
    // Set data detail berdasarkan lot
    lotDetail.lot = lotNumber;
    lotDetail.startTime = scanTime;
    
    // Hitung estimated end (misal 1 jam dari sekarang)
    const endDate = new Date();
    endDate.setHours(endDate.getHours() + 1);
    lotDetail.estimatedEnd = endDate.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
  });

  function getStatusColor(status: string) {
    if (status === 'Selesai') return 'bg-emerald-50 text-emerald-700 border-emerald-200';
    if (status === 'Proses') return 'bg-blue-50 text-blue-700 border-blue-200';
    return 'bg-slate-50 text-slate-600 border-slate-200';
  }
</script>

<div class="min-h-screen bg-slate-50 text-slate-800 pb-12">
  
  <!-- Header -->
  <header class="sticky top-0 z-50 bg-white border-b border-slate-200 shadow-sm">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
      <div class="flex justify-between items-start md:items-center">
        <div>
          <h1 class="text-2xl font-bold text-slate-800">Detail KPCP Lot</h1>
          <p class="text-sm text-slate-500 mt-1">Monitoring proses produksi</p>
        </div>
        <div class="text-right">
          <p class="text-sm font-semibold text-indigo-600">{lotNumber}</p>
          <p class="text-xs text-slate-500">{new Date().toLocaleDateString('id-ID')}</p>
        </div>
      </div>
    </div>
  </header>

  <main class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 mt-8 space-y-6">
    
    <!-- Lot Info Card -->
    <div class="bg-white rounded-2xl p-6 lg:p-8 shadow-sm border border-slate-100">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <div class="p-4 bg-indigo-50 rounded-xl border border-indigo-100">
          <p class="text-xs text-indigo-600 font-bold uppercase mb-1">No. Lot</p>
          <p class="text-xl font-mono font-bold text-indigo-700">{lotDetail.lot}</p>
        </div>
        <div class="p-4 bg-emerald-50 rounded-xl border border-emerald-100">
          <p class="text-xs text-emerald-600 font-bold uppercase mb-1">Produk</p>
          <p class="text-lg font-semibold text-emerald-700">{lotDetail.product}</p>
        </div>
        <div class="p-4 bg-blue-50 rounded-xl border border-blue-100">
          <p class="text-xs text-blue-600 font-bold uppercase mb-1">Mesin</p>
          <p class="text-lg font-semibold text-blue-700">{lotDetail.machine}</p>
        </div>
        <div class="p-4 bg-purple-50 rounded-xl border border-purple-100">
          <p class="text-xs text-purple-600 font-bold uppercase mb-1">No Compound</p>
          <p class="text-lg font-semibold text-purple-700">{lotDetail.noCompound}</p>
        </div>
        <div class="p-4 bg-slate-100 rounded-xl border border-slate-200">
          <p class="text-xs text-slate-600 font-bold uppercase mb-1">Status</p>
          <p class={`text-lg font-bold px-3 py-1 rounded-full inline-block ${getStatusColor(lotDetail.status)}`}>
            {lotDetail.status}
          </p>
        </div>
      </div>
    </div>

    <!-- Detail Info -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      
      <!-- Waktu & Supervisor -->
      <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
        <h3 class="font-bold text-lg text-slate-800 mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Informasi Waktu
        </h3>
        <div class="space-y-4">
          <div>
            <p class="text-sm text-slate-600 font-medium">Waktu Mulai</p>
            <p class="text-lg font-semibold text-slate-800 mt-1">{lotDetail.startTime} WIB</p>
          </div>
          <div>
            <p class="text-sm text-slate-600 font-medium">Estimasi Selesai</p>
            <p class="text-lg font-semibold text-slate-800 mt-1">{lotDetail.estimatedEnd} WIB</p>
          </div>
        </div>
      </div>

      <!-- Supervisor & Lot Mixing -->
      <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
        <h3 class="font-bold text-lg text-slate-800 mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Informasi Tambahan
        </h3>
        <div class="space-y-4">
          <div>
            <p class="text-sm text-slate-600 font-medium">Supervisor</p>
            <p class="text-lg font-semibold text-slate-800 mt-1">{lotDetail.supervisor}</p>
          </div>
          <div>
            <p class="text-sm text-slate-600 font-medium">No. Lot Mixing</p>
            <p class="text-lg font-semibold text-slate-800 mt-1">{lotDetail.noLotMixing}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Status Card -->
    <div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-2xl p-6 border border-blue-200">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-blue-600 flex items-center justify-center shrink-0 animate-pulse">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <div>
          <p class="text-sm text-blue-600 font-semibold uppercase">Status Proses</p>
          <p class="text-2xl font-bold text-blue-900 mt-1">Sedang Berlangsung</p>
          <p class="text-sm text-blue-700 mt-2">Harap menyelesaikan proses untuk melanjutkan</p>
        </div>
      </div>
    </div>

    <!-- Cycle Counter Card -->
    <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-bold text-lg text-slate-800 flex items-center gap-2">
          <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Progress Cycle
        </h3>
        <p class="text-3xl font-bold text-amber-600">{cycleCount}</p>
      </div>
      <button on:click={handleCompleteCycle}
        class="w-full px-6 py-3 bg-amber-500 hover:bg-amber-600 text-white font-bold rounded-lg transition-colors flex items-center justify-center gap-3 shadow-lg shadow-amber-200">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        Selesaikan 1 Cycle
      </button>
    </div>

    <!-- Action Button -->
    <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
      <button on:click={handleComplete}
        class="w-full px-6 py-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-lg transition-colors flex items-center justify-center gap-3 shadow-lg shadow-indigo-200 text-lg">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Selesaikan Proses & Kembali ke Dashboard
      </button>
      <p class="text-xs text-slate-500 text-center mt-3">Tekan tombol di atas untuk menyelesaikan proses dan kembali ke dashboard</p>
    </div>

  </main>
</div>

<style>
  :global(body) {
    background-color: #f8fafc;
  }
</style>
