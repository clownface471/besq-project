<script context="module" lang="ts">
  // UPDATE: Tambahkan 'options?: any' agar tidak merah saat dikasih argumen ke-4
  declare const jsQR: (data: Uint8ClampedArray, width: number, height: number, options?: any) => any;
</script>

<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';

  let videoElement: HTMLVideoElement;
  let canvasElement: HTMLCanvasElement;
  let cameraActive = false;
  let scannedData: string = '';
  let errorMessage = '';
  let successMessage = '';
  let scanningActive = false;
  let animationFrameId: any;
  let stream: MediaStream | null = null;

  async function startCamera() {
    try {
      errorMessage = '';
      stream = await navigator.mediaDevices.getUserMedia({
        video: { 
          facingMode: 'environment', // Kamera belakang
          width: { ideal: 1280 },
          height: { ideal: 720 }
        }
      });
      
      if (videoElement) {
        videoElement.srcObject = stream;
        videoElement.onloadedmetadata = () => {
             videoElement.play().then(() => {
                 cameraActive = true;
                 scanningActive = true;
                 startQRCodeScanning();
             });
        };
      }
    } catch (err) {
      errorMessage = 'Gagal mengakses kamera. Pastikan izin kamera diberikan.';
      console.error('Camera error:', err);
    }
  }

  function stopCamera() {
    if (stream) {
      stream.getTracks().forEach(track => track.stop());
      stream = null;
    }
    if (videoElement) {
        videoElement.srcObject = null;
    }
    if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
    }
    cameraActive = false;
    scanningActive = false;
  }

  function startQRCodeScanning() {
    if (!scanningActive || !cameraActive || !canvasElement || !videoElement) return;
    
    const canvas = canvasElement;
    const context = canvas.getContext('2d', { willReadFrequently: true });

    if (!context) return;

    const scan = () => {
      if (!scanningActive) return;

      if (videoElement.readyState === videoElement.HAVE_ENOUGH_DATA) {
        canvas.width = videoElement.videoWidth;
        canvas.height = videoElement.videoHeight;
        
        context.drawImage(videoElement, 0, 0, canvas.width, canvas.height);

        try {
          const imageData = context.getImageData(0, 0, canvas.width, canvas.height);
          // @ts-ignore - jsQR dimuat via CDN
          if (typeof jsQR !== 'undefined') {
              const code = jsQR(imageData.data, imageData.width, imageData.height, {
                  inversionAttempts: "dontInvert",
              });
    
              if (code && code.data) {
                drawQuad(code.location, context);
                scannedData = code.data;
                handleScan();
                return; 
              }
          }
        } catch (err) {
          console.error('QR Code scanning error:', err);
        }
      }

      animationFrameId = requestAnimationFrame(scan);
    };

    animationFrameId = requestAnimationFrame(scan);
  }
  
  function drawQuad(location: any, ctx: CanvasRenderingContext2D) {
      ctx.lineWidth = 4;
      ctx.strokeStyle = "#4f46e5"; 
      ctx.beginPath();
      ctx.moveTo(location.topLeftCorner.x, location.topLeftCorner.y);
      ctx.lineTo(location.topRightCorner.x, location.topRightCorner.y);
      ctx.lineTo(location.bottomRightCorner.x, location.bottomRightCorner.y);
      ctx.lineTo(location.bottomLeftCorner.x, location.bottomLeftCorner.y);
      ctx.lineTo(location.topLeftCorner.x, location.topLeftCorner.y);
      ctx.stroke();
  }

  function handleScan() {
    if (scannedData.trim()) {
      scanningActive = false; 

      let rawString = scannedData.trim();
      rawString = rawString.replace(/^\[|\]$/g, '');
      
      const parts = rawString.split(';');
      
      let productCode = 'KPCP Standard'; 
      let lotNo = rawString;             

      if (parts.length >= 2) {
          productCode = parts[0].trim();
          lotNo = parts[1].trim();
      }

      const currentTime = new Date().toLocaleTimeString('id-ID');
      
<<<<<<< HEAD
      // Simpan data ke localStorage sebelum redirect 
      localStorage.setItem('selectedLot', scannedData);
=======
      localStorage.setItem('selectedLot', lotNo);
      localStorage.setItem('selectedProduct', productCode); 
>>>>>>> ca9b029e240f91b16c92d31b59f4d4d035ea07cb
      localStorage.setItem('scanTime', currentTime);
      
      successMessage = `Scan Berhasil!`;
      
      setTimeout(() => {
        goto(`/kpcp-detail?lot=${encodeURIComponent(lotNo)}`);
      }, 1000);
    }
  }

  function handleBack() {
    stopCamera();
<<<<<<< HEAD
<<<<<<< HEAD:src/routes/barcode/+page.svelte
    window.location.href = '/cutting';
=======
    window.location.href = '/pressing';
>>>>>>> a1d43feec433a0361748632a0b5fce12d4acc101:src/routes/scan-barcode-prs/+page.svelte
=======
    goto('/pressing');
  }

  function handleManualSubmit() {
      if (scannedData) handleScan();
>>>>>>> ca9b029e240f91b16c92d31b59f4d4d035ea07cb
  }

  onMount(() => {
    if (!document.getElementById('jsqr-script')) {
        const script = document.createElement('script');
        script.id = 'jsqr-script';
        script.src = 'https://cdn.jsdelivr.net/npm/jsqr@1.4.0/dist/jsQR.js';
        script.onload = () => {
          startCamera();
        };
        script.onerror = () => {
            errorMessage = "Gagal memuat library QR Scanner. Cek koneksi internet.";
        };
        document.head.appendChild(script);
    } else {
        startCamera();
    }
  });

<<<<<<< HEAD
=======
  onDestroy(() => {
    stopCamera();
  });
>>>>>>> ca9b029e240f91b16c92d31b59f4d4d035ea07cb
</script>

<div class="min-h-screen bg-slate-900 text-white pb-6 relative overflow-hidden">
  
  <header class="absolute top-0 left-0 right-0 z-20 bg-gradient-to-b from-black/80 to-transparent p-4">
    <div class="max-w-md mx-auto flex justify-between items-center">
      <button 
        on:click={handleBack} 
        aria-label="Kembali ke Dashboard"
        class="p-2 bg-white/10 backdrop-blur-md rounded-full hover:bg-white/20 transition-all">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
      </button>
      <div class="text-center">
        <h1 class="text-lg font-bold tracking-wide">Scan QR KPCP</h1>
        <p class="text-xs text-slate-300">Format: [Produk;NoLot]</p>
      </div>
      <div class="w-10"></div> 
    </div>
  </header>

  <main class="h-screen flex flex-col items-center justify-center relative">
    
    <div class="absolute inset-0 z-0 bg-black">
         <video
            bind:this={videoElement}
            autoplay
            playsinline
            class="w-full h-full object-cover"
          ></video>
         <canvas bind:this={canvasElement} class="hidden"></canvas>
    </div>

    <div class="absolute inset-0 z-10 pointer-events-none flex items-center justify-center">
        <div class="absolute inset-0 bg-black/50 mask-scan"></div>
        
        <div class="relative w-72 h-72 border-2 border-white/50 rounded-3xl overflow-hidden">
            <div class="absolute top-0 left-0 w-8 h-8 border-t-4 border-l-4 border-indigo-500 rounded-tl-lg"></div>
            <div class="absolute top-0 right-0 w-8 h-8 border-t-4 border-r-4 border-indigo-500 rounded-tr-lg"></div>
            <div class="absolute bottom-0 left-0 w-8 h-8 border-b-4 border-l-4 border-indigo-500 rounded-bl-lg"></div>
            <div class="absolute bottom-0 right-0 w-8 h-8 border-b-4 border-r-4 border-indigo-500 rounded-br-lg"></div>
            
            {#if scanningActive}
                <div class="absolute top-0 left-0 right-0 h-1 bg-indigo-400 shadow-[0_0_15px_rgba(99,102,241,0.8)] animate-scan-line"></div>
            {/if}

            {#if !cameraActive && !errorMessage}
               <div class="absolute inset-0 flex items-center justify-center">
                   <svg class="animate-spin h-10 w-10 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                   </svg>
               </div>
            {/if}
        </div>
        
        <p class="absolute mt-80 text-white/80 text-sm font-medium animate-pulse">
            Arahkan kamera ke QR Code
        </p>
    </div>

    <div class="absolute bottom-0 left-0 right-0 z-20 p-6 bg-gradient-to-t from-slate-900 via-slate-900/90 to-transparent">
        
        {#if successMessage}
          <div class="mb-4 p-4 bg-emerald-500/20 backdrop-blur-md border border-emerald-500/50 rounded-xl flex items-center gap-3 animate-bounce-in">
            <div class="w-10 h-10 rounded-full bg-emerald-500 flex items-center justify-center shrink-0">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
            </div>
            <div>
                <h3 class="font-bold text-emerald-100">Scan Berhasil!</h3>
                <p class="text-xs text-emerald-200">{scannedData.replace(';', ' | ')}</p>
            </div>
          </div>
        {/if}

        {#if errorMessage}
          <div class="mb-4 p-4 bg-rose-500/20 backdrop-blur-md border border-rose-500/50 rounded-xl flex items-center gap-3">
            <svg class="w-6 h-6 text-rose-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4v.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            <span class="text-sm text-rose-100 font-medium">{errorMessage}</span>
          </div>
        {/if}

        <div class="bg-slate-800/80 backdrop-blur-md rounded-2xl p-2 border border-slate-700">
             <div class="flex gap-2">
                 <input 
                    type="text" 
                    bind:value={scannedData}
                    placeholder="Input manual (Kode;Lot)..." 
                    class="flex-1 bg-transparent border-none text-white placeholder-slate-500 focus:ring-0 px-4 py-2"
                    on:keydown={(e) => e.key === 'Enter' && handleManualSubmit()}
                 />
                 <button 
                    on:click={handleManualSubmit}
                    disabled={!scannedData}
                    aria-label="Submit Manual Input"
                    class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-600 text-white p-2 rounded-xl transition-colors">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                 </button>
             </div>
        </div>
    </div>

  </main>
</div>

<style>
  @keyframes scan-line {
      0% { top: 0%; opacity: 0; }
      10% { opacity: 1; }
      90% { opacity: 1; }
      100% { top: 100%; opacity: 0; }
  }
  .animate-scan-line {
      animation: scan-line 2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
  }
  .animate-bounce-in {
      animation: bounce-in 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275) both;
  }
  @keyframes bounce-in {
      0% { transform: translateY(20px); opacity: 0; }
      100% { transform: translateY(0); opacity: 1; }
  }
</style>