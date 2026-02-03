<script context="module" lang="ts">
  declare const jsQR: (data: Uint8ClampedArray, width: number, height: number) => any;
</script>

<script lang="ts">
  import { onMount } from 'svelte';

  let videoElement: HTMLVideoElement;
  let canvasElement: HTMLCanvasElement;
  let cameraActive = false;
  let scannedData: string = '';
  let errorMessage = '';
  let successMessage = '';
  let scanningActive = false;

  async function startCamera() {
    try {
      errorMessage = '';
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { 
          facingMode: 'environment',
          width: { ideal: 1920 },
          height: { ideal: 1080 }
        }
      });
      if (videoElement) {
        videoElement.srcObject = stream;
        cameraActive = true;
        scanningActive = true;
        startQRCodeScanning();
      }
    } catch (err) {
      errorMessage = 'Gagal mengakses kamera. Pastikan Anda memberikan izin akses.';
      console.error('Camera error:', err);
    }
  }

  function stopCamera() {
    if (videoElement?.srcObject) {
      const tracks = (videoElement.srcObject as MediaStream).getTracks();
      tracks.forEach(track => track.stop());
      cameraActive = false;
      scanningActive = false;
    }
  }

  function startQRCodeScanning() {
    if (!scanningActive || !cameraActive) return;

    const canvas = canvasElement;
    const video = videoElement;
    const context = canvas.getContext('2d');

    if (!context) return;

    const scan = () => {
      if (video.readyState === video.HAVE_ENOUGH_DATA) {
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        context.drawImage(video, 0, 0, canvas.width, canvas.height);

        try {
          const imageData = context.getImageData(0, 0, canvas.width, canvas.height);
          const code = jsQR(imageData.data, imageData.width, imageData.height);

          if (code) {
            scannedData = code.data;
            handleScan();
            scanningActive = false;
            return;
          }
        } catch (err) {
          console.error('QR Code scanning error:', err);
        }
      }

      if (scanningActive) {
        requestAnimationFrame(scan);
      }
    };

    requestAnimationFrame(scan);
  }

  function handleScan() {
    if (scannedData.trim()) {
      const currentTime = new Date().toLocaleTimeString('id-ID');
      
      // Simpan data ke localStorage sebelum redirect
      localStorage.setItem('selectedLot', scannedData);
      localStorage.setItem('scanTime', currentTime);
      
      successMessage = `Lot ${scannedData} berhasil discan!`;
      
      // Redirect ke halaman detail setelah 1 detik
      setTimeout(() => {
        window.location.href = `/kpcp-detail?lot=${encodeURIComponent(scannedData)}`;
      }, 1000);
      
      scannedData = '';
    }
  }

  function handleBack() {
    stopCamera();
    window.location.href = '/pressing';
  }

  onMount(() => {
    // Load jsQR library
    const script = document.createElement('script');
    script.src = 'https://cdn.jsdelivr.net/npm/jsqr@1.4.0/dist/jsQR.js';
    script.onload = () => {
      startCamera();
    };
    document.head.appendChild(script);

    return () => {
      stopCamera();
    };
  });


</script>

<div class="min-h-screen bg-white text-slate-800 pb-6 relative">
  
  <!-- Header -->
  <header class="sticky top-0 z-50 bg-white border-b border-slate-200 shadow-sm">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
      <div class="flex items-center gap-3">
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button on:click={handleBack}
          class="group p-2 hover:bg-slate-100 rounded-lg transition-colors">
          <svg class="w-6 h-6 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div>
          <h1 class="text-xl md:text-2xl font-bold text-slate-800 tracking-tight">Scan Barcode KPCP</h1>
          <p class="text-xs md:text-sm text-slate-500">Arahkan kamera ke barcode atau QR code lot</p>
        </div>
      </div>
      <div class="text-right">
        <p class="text-sm text-slate-500">{new Date().toLocaleTimeString('id-ID')}</p>
      </div>
    </div>
  </header>

  <main class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 mt-6 space-y-6">
    
    <!-- Camera Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Video Feed -->
      <div class="lg:col-span-2">
        <div class="relative bg-black rounded-2xl overflow-hidden shadow-xl border-2 border-slate-200">
          <!-- svelte-ignore element_invalid_self_closing_tag -->
          <video
            bind:this={videoElement}
            autoplay
            playsinline
            class="w-full h-auto aspect-video object-cover"
          />
          <!-- svelte-ignore element_invalid_self_closing_tag -->
          <canvas bind:this={canvasElement} class="hidden" />
          
          <!-- Scanner Frame Overlay -->
          <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div class="relative w-72 h-72">
              <!-- Corner marks -->
              <div class="absolute top-0 left-0 w-8 h-8 border-t-2 border-l-2 border-indigo-500"></div>
              <div class="absolute top-0 right-0 w-8 h-8 border-t-2 border-r-2 border-indigo-500"></div>
              <div class="absolute bottom-0 left-0 w-8 h-8 border-b-2 border-l-2 border-indigo-500"></div>
              <div class="absolute bottom-0 right-0 w-8 h-8 border-b-2 border-r-2 border-indigo-500"></div>
              
              <!-- Scanning line animation -->
              <div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-transparent via-red-500 to-transparent animate-pulse"></div>
              
              <!-- Center circle for QR Code -->
              <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-16 h-16 border-2 border-indigo-400 rounded-full opacity-50"></div>
            </div>
          </div>

          <!-- Status indicator -->
          <div class="absolute top-4 right-4 flex items-center gap-2 bg-black/50 px-3 py-2 rounded-full">
            <div class="w-2 h-2 bg-red-500 rounded-full animate-pulse"></div>
            <span class="text-xs font-medium text-white">LIVE</span>
          </div>

          <!-- Scanning indicator -->
          {#if scanningActive}
            <div class="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-full">
              <div class="w-2 h-2 bg-white rounded-full animate-pulse"></div>
              <span class="text-xs font-semibold">Scanning...</span>
            </div>
          {/if}
        </div>

        {#if successMessage}
          <div class="mt-4 p-4 bg-emerald-50 border border-emerald-200 rounded-lg flex items-center gap-3 animate-pulse">
            <svg class="w-5 h-5 text-emerald-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span class="text-sm text-emerald-800 font-medium">{successMessage}</span>
            <span class="text-xs text-emerald-600 ml-auto">Mengalihkan...</span>
          </div>
        {/if}

        {#if errorMessage}
          <div class="mt-4 p-4 bg-rose-50 border border-rose-200 rounded-lg flex items-center gap-3">
            <svg class="w-5 h-5 text-rose-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4v.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-sm text-rose-800 font-medium">{errorMessage}</span>
          </div>
        {/if}
      </div>

      <!-- Input & Controls Panel -->
      <div class="lg:col-span-1 space-y-4">
        
        <!-- QR Code Info Card -->
        <div class="bg-blue-50 rounded-2xl p-6 border border-blue-200">
          <h3 class="font-bold text-lg mb-3 flex items-center gap-2 text-blue-900">
            <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Cara Menggunakan
          </h3>
          <ul class="text-sm text-blue-800 space-y-2">
            <li class="flex gap-2">
              <span class="font-bold">1.</span>
              <span>Arahkan kamera ke QR code lot</span>
            </li>
            <li class="flex gap-2">
              <span class="font-bold">2.</span>
              <span>Tunggu scanning otomatis selesai</span>
            </li>
            <li class="flex gap-2">
              <span class="font-bold">3.</span>
              <span>Atau input manual jika diperlukan</span>
            </li>
          </ul>
        </div>

        <!-- Manual Input Card -->
        <div class="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
          <h3 class="font-bold text-lg mb-4 flex items-center gap-2 text-slate-800">
            <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Input Manual
          </h3>
          
          <div class="space-y-3">
            <input
              type="text"
              bind:value={scannedData}
              placeholder="Masukkan No. Lot"
              class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 placeholder-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 transition-all"
              on:keydown={(e) => {
                if (e.key === 'Enter') handleScan();
              }}
            />
            
            <button
              on:click={handleScan}
              disabled={successMessage !== ''}
              class="w-full px-4 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-300 disabled:cursor-not-allowed rounded-lg font-semibold transition-colors flex items-center justify-center gap-2 text-white"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Tambah Scan
            </button>
          </div>
        </div>

        <!-- Camera Controls Card -->
        <div class="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
          <h3 class="font-bold text-lg mb-4 flex items-center gap-2 text-slate-800">
            <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            Kontrol Kamera
          </h3>
          
          <div class="space-y-2">
            <button
              on:click={startCamera}
              disabled={cameraActive}
              class="w-full px-4 py-2 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-300 disabled:cursor-not-allowed rounded-lg font-medium transition-colors text-sm text-white"
            >
              Mulai Kamera
            </button>
            <button
              on:click={stopCamera}
              disabled={!cameraActive}
              class="w-full px-4 py-2 bg-rose-600 hover:bg-rose-700 disabled:bg-slate-300 disabled:cursor-not-allowed rounded-lg font-medium transition-colors text-sm text-white"
            >
              Hentikan Kamera
            </button>
          </div>
        </div>

        <!-- Status Card -->
        <div class="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm">
          <div class="space-y-3">
            <div class="flex items-baseline gap-2">
              <span class="text-sm text-slate-600 font-medium">Kamera:</span>
              <span class={`px-3 py-1 rounded-full text-sm font-semibold ${cameraActive ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-600 border border-slate-200'}`}>
                {cameraActive ? 'Aktif' : 'Tidak Aktif'}
              </span>
            </div>
            <div class="flex items-baseline gap-2">
              <span class="text-sm text-slate-600 font-medium">Scanning:</span>
              <span class={`px-3 py-1 rounded-full text-sm font-semibold ${scanningActive ? 'bg-blue-50 text-blue-700 border border-blue-200' : 'bg-slate-100 text-slate-600 border border-slate-200'}`}>
                {scanningActive ? 'Aktif' : 'Siaga'}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</div>

<style>
  :global(body) {
    background-color: #ffffff;
  }

  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>
