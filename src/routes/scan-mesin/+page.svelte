<script lang="ts">
  import { onMount } from "svelte";
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
    const scanInterval = setInterval(() => {
      if (!cameraActive || !videoElement || !canvasElement) {
        clearInterval(scanInterval);
        return;
      }

      const context = canvasElement.getContext("2d");
      if (context) {
        context.drawImage(videoElement, 0, 0, canvasElement.width, canvasElement.height);
        const imageData = context.getImageData(0, 0, canvasElement.width, canvasElement.height);
        
        // Simple barcode detection (detects dark patterns)
        if (detectBarcode(imageData)) {
          const code = extractBarcodeData(imageData);
          if (code) {
            scanResult = code;
            handleScanResult(code);
            clearInterval(scanInterval);
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
    
    return darkPixels > data.length * 0.1;
  }

  function extractBarcodeData(imageData: ImageData): string {
    // Simplified barcode extraction
    // For production, use a library like jsQR or quagga.js
    const canvas = document.createElement("canvas");
    canvas.width = imageData.width;
    canvas.height = imageData.height;
    const ctx = canvas.getContext("2d");
    if (ctx) {
      ctx.putImageData(imageData, 0, 0);
      // Return mock data for demonstration
      return "MESIN-" + Math.random().toString(36).substr(2, 9).toUpperCase();
    }
    return "";
  }

  async function handleScanResult(code: string) {
    isProcessing = true;
    errorMessage = "";

    try {
      const authToken = get(auth).token;
      if (!authToken) {
        errorMessage = "Tidak ada token autentikasi";
        isProcessing = false;
        return;
      }

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
        throw new Error("Gagal memproses scan mesin");
      }

      successMessage = `Mesin ${code} berhasil discan!`;
      
      // Redirect ke halaman pressing setelah 2 detik
      setTimeout(() => {
        window.location.href = "/pressing";
      }, 2000);
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : "Terjadi kesalahan";
      console.error("Error:", error);
      isProcessing = false;
      cameraActive = false;
      
      // Restart scanning
      setTimeout(() => {
        cameraActive = true;
        startScanning();
      }, 1500);
    }
  }

  function stopCamera() {
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
    return () => {
      stopCamera();
    };
  });
</script>

<div class="min-h-screen bg-slate-900 text-white flex items-center justify-center p-4">
  <div class="w-full max-w-lg">
    <!-- Header -->
    <div class="text-center mb-8">
      <h1 class="text-3xl md:text-4xl font-bold mb-2">Scan Mesin</h1>
      <p class="text-slate-400">Arahkan kamera ke barcode mesin untuk scan</p>
    </div>

    <!-- Camera Preview Container -->
    <div class="relative bg-black rounded-2xl overflow-hidden shadow-2xl mb-6 border-4 border-indigo-500">
      <video
        bind:this={videoElement}
        autoplay
        playsinline
        class="w-full aspect-video object-cover"
      />
      <canvas bind:this={canvasElement} class="hidden" width="1280" height="720" />
      
      <!-- Scanning Overlay -->
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

      <!-- Loading Indicator -->
      {#if isProcessing}
        <div class="absolute inset-0 bg-black/50 flex items-center justify-center backdrop-blur-sm">
          <div class="flex flex-col items-center gap-4">
            <div class="w-12 h-12 border-4 border-indigo-400 border-t-indigo-600 rounded-full animate-spin"></div>
            <p class="text-sm font-semibold">Memproses scan...</p>
          </div>
        </div>
      {/if}
    </div>

    <!-- Messages -->
    {#if errorMessage}
      <div class="mb-4 p-4 bg-red-900/30 border border-red-500 rounded-lg text-red-200 text-sm">
        {errorMessage}
      </div>
    {/if}

    {#if successMessage}
      <div class="mb-4 p-4 bg-emerald-900/30 border border-emerald-500 rounded-lg text-emerald-200 text-sm">
        {successMessage}
      </div>
    {/if}

    {#if scanResult}
      <div class="mb-4 p-4 bg-indigo-900/30 border border-indigo-500 rounded-lg">
        <p class="text-xs text-slate-400 mb-1">Hasil Scan:</p>
        <p class="text-lg font-mono font-bold text-indigo-300">{scanResult}</p>
      </div>
    {/if}

    <!-- Manual Input -->
    <div class="mb-6">
      <label class="block text-sm font-semibold mb-2 text-slate-300">
        Atau masukkan kode mesin secara manual:
      </label>
      <input
        type="text"
        placeholder="Masukkan kode mesin..."
        onchange={manualInput}
        class="w-full px-4 py-3 bg-slate-800 border border-slate-600 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/30"
      />
    </div>

    <!-- Actions -->
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
        Mulai Ulang
      </button>
    </div>

    <!-- Info -->
    <div class="mt-6 p-4 bg-slate-800/50 rounded-lg border border-slate-700">
      <p class="text-xs text-slate-400 text-center">
        💡 Pastikan kamera memiliki izin akses dan cahaya yang cukup untuk hasil scan optimal.
      </p>
    </div>
  </div>
</div>

<style>
  :global(body) {
    background-color: #0f172a;
  }
</style>
