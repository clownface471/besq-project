<script lang="ts">
  import { onMount } from 'svelte';
  import Swal from 'sweetalert2';
  import { get } from "svelte/store";
  import { auth } from "$lib/stores/auth";

  const API_URL = "http://localhost:8080";

  let lotNumber = '';
  let activeMachine = ''; // Tambahan: Ambil dari localStorage
  let scanTime = '';
  let cycleCount = 0;
  
  // Data detail yang akan dikirim
  let lotDetail = {
    lot: '',
    product: 'KPCP Standard',
    quantity: 500,
    status: 'Proses',
    startTime: '',
    estimatedEnd: '',
    supervisor: 'Tono Widiyanto', // Bisa diganti dinamis nanti
    machine: '',
    noCompound: 'CMP-2309-001',
    noLotMixing: 'MIX-2309-A05'
  };

  // State untuk button loading
  let isSubmitting = false;

  function handleCompleteCycle() {
    cycleCount++;
    Swal.fire({
      title: 'Cycle Selesai!',
      html: `<p class="text-lg font-semibold">Cycle <span class="text-amber-600">${cycleCount}</span> telah berhasil diselesaikan</p>`,
      icon: 'success',
      confirmButtonText: 'Lanjutkan',
      confirmButtonColor: '#d97706',
      allowOutsideClick: false,
      timer: 1500
    });
  }

  async function handleComplete() {
    const result = await Swal.fire({
      title: 'Konfirmasi Mulai Produksi',
      text: 'Apakah data sudah benar dan siap memulai produksi?',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#4f46e5',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'Ya, Mulai',
      cancelButtonText: 'Batal'
    });

    if (result.isConfirmed) {
        isSubmitting = true;
        try {
            // 1. Ambil Token
            const authToken = get(auth).token;
            if (!authToken) throw new Error("Token tidak ditemukan. Silakan login ulang.");

            // 2. Siapkan Payload untuk Backend (sesuai struct LWP / PerCycle)
            // Kita pakai endpoint LWP agar masuk ke tabel dashboard
            const payload = {
                noMesin: lotDetail.machine,
                tanggal: new Date().toISOString(),
                shift: "I",
                // nik: $auth.user?.username, // Backend otomatis ambil dari token
                partName: lotDetail.product,
                kodePart: "KPCP-STD",
                jamMulai: lotDetail.startTime,
                jamSelesai: "-",
                hasilOk: 0, // Awal mulai 0
                ng: 0,
                klasifikasiReject: ""
            };

            // 3. Kirim ke Backend
            const response = await fetch(`${API_URL}/api/lwp`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${authToken}`
                },
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                const errData = await response.json();
                throw new Error(errData.error || "Gagal menyimpan data LWP");
            }

            // 4. Sukses
            await Swal.fire({
                title: 'Berhasil!',
                text: 'Data produksi berhasil dibuat. Silakan cek Dashboard.',
                icon: 'success',
                confirmButtonColor: '#4f46e5'
            });

            // Redirect ke dashboard
            window.location.href = '/pressing';

        } catch (error: any) {
            Swal.fire({
                title: 'Gagal',
                text: error.message,
                icon: 'error'
            });
        } finally {
            isSubmitting = false;
        }
    }
  }

  onMount(() => {
    // Ambil data dari URL & LocalStorage
    const queryParams = new URLSearchParams(window.location.search);
    lotNumber = queryParams.get('lot') || '';
    
    // Ambil kode mesin yang di-scan sebelumnya
    activeMachine = localStorage.getItem('activeMachine') || 'UNKNOWN-MC';
    scanTime = localStorage.getItem('scanTime') || new Date().toLocaleTimeString('id-ID');

    const scannedProduct = localStorage.getItem('selectedProduct');
    if (scannedProduct) lotDetail.product = scannedProduct;
    
    lotDetail.lot = lotNumber;
    lotDetail.machine = activeMachine; // Update mesin sesuai scan
    lotDetail.startTime = scanTime;
    
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
          <p class="text-xs text-blue-600 font-bold uppercase mb-1">Mesin (Scan)</p>
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

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      
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

    <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
      <button 
        on:click={handleComplete}
        disabled={isSubmitting}
        class="w-full px-6 py-4 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-300 text-white font-bold rounded-lg transition-colors flex items-center justify-center gap-3 shadow-lg shadow-indigo-200 text-lg">
        {#if isSubmitting}
             <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
               <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
               <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
             </svg>
             Menyimpan Data...
        {:else}
             <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
             </svg>
             Mulai Produksi & Kembali ke Dashboard
        {/if}
      </button>
      <p class="text-xs text-slate-500 text-center mt-3">Tekan tombol untuk menyimpan data dan mengaktifkan status mesin di Dashboard</p>
    </div>

  </main>
</div>

<style>
  :global(body) {
    background-color: #f8fafc;
  }
</style>