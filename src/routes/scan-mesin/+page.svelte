<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { get } from "svelte/store";
  import { auth } from "$lib/stores/auth";

  const API_URL = "http://localhost:8080";

  let videoElement: HTMLVideoElement;
  let canvasElement: HTMLCanvasElement;
  let cameraActive = $state(false);
  let scanResult = $state("");
  let errorMessage = $state("");
  let successMessage = $state("");
  let isProcessing = $state(false);
  
  // Interval ID disimpan agar bisa dibersihkan
  let scanIntervalId: any;

  async function startCamera() {
    try {
      errorMessage = "";
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { 
          facingMode: "environment",
          width: { ideal: 1280 },
          height: { ideal: 720 }
        },
      });
      if (videoElement) {
        videoElement.srcObject = stream;
        cameraActive = true;
        startScanning();
      }
    } catch (error) {
      errorMessage = "Tidak dapat mengakses kamera. Periksa izin kamera Anda.";
      console.error("Camera error:", error);
    }
  }

  function startScanning() {
    // Bersihkan interval lama jika ada
    if (scanIntervalId) clearInterval(scanIntervalId);

    scanIntervalId = setInterval(() => {
      if (!cameraActive || !videoElement || !canvasElement) {
        clearInterval(scanIntervalId);
        return;
      }

      const context = canvasElement.getContext("2d");
      if (context && videoElement.readyState === videoElement.HAVE_ENOUGH_DATA) {
        canvasElement.width = videoElement.videoWidth;
        canvasElement.height = videoElement.videoHeight;
        
        context.drawImage(videoElement, 0, 0, canvasElement.width, canvasElement.height);
        const imageData = context.getImageData(0, 0, canvasElement.width, canvasElement.height);
        
        // Deteksi sederhana (Dark pixels) -> Bisa diganti library jsQR di masa depan
        if (detectBarcode(imageData)) {
          // Simulasi pembacaan sukses
          // Di production, gunakan library 'jsqr' atau 'quagga' di sini
          const code = "PRESS-" + Math.floor(Math.random() * 10) + 1; // Simulasi: PRESS-1 s/d PRESS-10
          
          // Jika ingin test manual input barcode asli, gunakan logic extractBarcodeData yg asli
          // const code = extractBarcodeData(imageData);

          if (code && !isProcessing) {
            scanResult = code;
            handleScanResult(code);
            clearInterval(scanIntervalId);
          }
        }
      }
    }, 500);
  }

  function detectBarcode(imageData: ImageData): boolean {
    const data = imageData.data;
    let darkPixels = 0;
    for (let i = 0; i < data.length; i += 4) {
      const gray = (data[i] + data[i + 1] + data[i + 2]) / 3;
      if (gray < 100) darkPixels++;
    }
    // Ambang batas sederhana
    return darkPixels > data.length * 0.15;
  }

  async function handleScanResult(code: string) {
    if (isProcessing) return; // Mencegah double submit
    isProcessing = true;
    errorMessage = "";

    try {
      const authToken = get(auth).token;
      if (!authToken) {
        errorMessage = "Tidak ada token autentikasi";
        isProcessing = false;
        return;
      }

      // 1. Validasi ke Backend
      const response = await fetch(`${API_URL}/api/scan-machine`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
        body: JSON.stringify({
          machineCode: code,
          timestamp: new Date().toISOString(),
        }),
      });

      if (!response.ok) {
        throw new Error("Gagal memvalidasi mesin");
      }

      // 2. SUKSES SCAN
      successMessage = `Mesin ${code} Valid! Lanjut ke Scan KPCP...`;
      
      // 3. SIMPAN DATA PENTING KE STORAGE
      // Ini yang akan dibaca oleh halaman 'kpcp-detail' nanti
      localStorage.setItem('activeMachine', code);

      // 4. REDIRECT KE ALUR SELANJUTNYA (SCAN KPCP)
      setTimeout(() => {
        window.location.href = "/scan-barcode-prs"; 
      }, 1500);

    } catch (error) {
      errorMessage = error instanceof Error ? error.message : "Terjadi kesalahan";
      console.error("Error:", error);
      isProcessing = false;
      
      // Restart scanning setelah error
      setTimeout(() => {
        if(cameraActive) startScanning();
      }, 2000);
    }
  }

  function stopCamera() {
    if (scanIntervalId) clearInterval(scanIntervalId);
    if (videoElement && videoElement.srcObject) {
      const tracks = (videoElement.srcObject as MediaStream).getTracks();
      tracks.forEach((track) => track.stop());
      cameraActive = false;
    }
  }

  function handleCancel() {
    stopCamera();
    window.location.href = "/pressing";
  }

  function manualInput(e: Event) {
    const input = (e.target as HTMLInputElement).value.trim();
    if (input) {
      scanResult = input;
      handleScanResult(input);
    }
  }

  onMount(() => {
    startCamera();
  });

  onDestroy(() => {
    stopCamera();
  });
</script>

<div class="min-h-screen bg-slate-900 text-white flex items-center justify-center p-4">
  <div class="w-full max-w-lg">
    <div class="text-center mb-8">
      <h1 class="text-3xl md:text-4xl font-bold mb-2">Scan Mesin</h1>
      <p class="text-slate-400">Scan Barcode Mesin Pressing</p>
    </div>

    <div class="relative bg-black rounded-2xl overflow-hidden shadow-2xl mb-6 border-4 border-indigo-500">
      <video
        bind:this={videoElement}
        autoplay
        playsinline
        class="w-full aspect-video object-cover"
      />
      <canvas bind:this={canvasElement} class="hidden" />
      
      {#if cameraActive}
        <div class="absolute inset-0 border-2 border-dashed border-indigo-400 rounded-lg flex items-center justify-center">
          <div class="absolute w-56 h-56 border-2 border-indigo-500 rounded-lg animate-pulse"></div>
          <div class="relative">
             <div class="absolute top-0 left-0 w-4 h-4 border-t-2 border-l-2 border-green-400"></div>
            <div class="absolute top-0 right-0 w-4 h-4 border-t-2 border-r-2 border-green-400"></div>
            <div class="absolute bottom-0 left-0 w-4 h-4 border-b-2 border-l-2 border-green-400"></div>
            <div class="absolute bottom-0 right-0 w-4 h-4 border-b-2 border-r-2 border-green-400"></div>
          </div>
        </div>
      {/if}

      {#if isProcessing}
        <div class="absolute inset-0 bg-black/50 flex items-center justify-center backdrop-blur-sm z-10">
          <div class="flex flex-col items-center gap-4">
            <div class="w-12 h-12 border-4 border-indigo-400 border-t-indigo-600 rounded-full animate-spin"></div>
            <p class="text-sm font-semibold">Memvalidasi Mesin...</p>
          </div>
        </div>
      {/if}
    </div>

    {#if errorMessage}
      <div class="mb-4 p-4 bg-red-900/30 border border-red-500 rounded-lg text-red-200 text-sm">
        {errorMessage}
      </div>
    {/if}

    {#if successMessage}
      <div class="mb-4 p-4 bg-emerald-900/30 border border-emerald-500 rounded-lg text-emerald-200 text-sm animate-pulse">
        {successMessage}
      </div>
    {/if}

    <div class="mb-6">
      <label class="block text-sm font-semibold mb-2 text-slate-300" for="manual-input">
        Atau masukkan kode mesin manual:
      </label>
      <input
        id="manual-input"
        type="text"
        placeholder="Contoh: PRESS-01"
        onchange={manualInput}
        class="w-full px-4 py-3 bg-slate-800 border border-slate-600 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/30"
      />
    </div>

    <div class="flex gap-3">
      <button
        onclick={handleCancel}
        class="flex-1 px-4 py-3 bg-slate-700 hover:bg-slate-600 text-white font-semibold rounded-lg transition-colors duration-200 active:scale-95"
      >
        Batal
      </button>
      <button
        onclick={startCamera}
        disabled={cameraActive}
        class="flex-1 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-600 text-white font-semibold rounded-lg transition-colors duration-200 active:scale-95"
      >
        Scan Ulang
      </button>
    </div>
  </div>
</div>

<style>
  :global(body) {
    background-color: #0f172a;
  }
</style>