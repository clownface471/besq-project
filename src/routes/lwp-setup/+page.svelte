<script>
  import { goto } from '$app/navigation';
  // State untuk menampung data input
  let npk = "KRT-08348";
  let noMc = "136A";
  let shift = "1";

  function handleBack() {
    goto('/pressing'); 
  }

  function handleShiftStart() {
    alert(`Shift dimulai untuk NPK: ${npk}, NoMc: ${noMc}, Shift: ${shift}`);
    goto('/lwp-input');
  }
</script>

<div class="min-h-screen bg-[#f3f3f3] flex flex-col items-center p-6 font-sans">
  
  <div onclick={handleBack} class="w-full max-w-md mt-4">
    <button class="bg-[#facc15] hover:bg-yellow-500 p-3 rounded-xl shadow-md transition-all active:scale-90">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-black" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
  </div>

  <div class="w-full max-w-[340px] bg-[#111827] rounded-[24px] p-8 mt-24 shadow-2xl border border-gray-800">
    <div class="text-center mb-8">
      <h1 class="text-white text-2xl font-black tracking-tight uppercase">Shift Setup</h1>
      <p class="text-gray-400 text-[12px] font-medium mt-1">Laporan Waktu Produksi</p>
    </div>

    <div class="space-y-3">
      <div class="relative">
        <input 
          type="text" 
          bind:value={npk}
          placeholder="NPK: KRT-08348"
          class="w-full bg-[#eeeeee] border-none rounded-xl py-3 px-5 text-gray-500 placeholder-gray-400 font-medium focus:ring-2 focus:ring-yellow-400 outline-none"
        />
      </div>

      <div class="relative">
        <input 
          type="text" 
          bind:value={noMc}
          placeholder="NoMc: 136A"
          class="w-full bg-[#eeeeee] border-none rounded-xl py-3 px-5 text-gray-500 placeholder-gray-400 font-medium focus:ring-2 focus:ring-yellow-400 outline-none"
        />
      </div>

      <div class="relative">
        <input 
          type="text" 
          bind:value={shift}
          placeholder="Shift: 1"
          class="w-full bg-[#eeeeee] border-none rounded-xl py-3 px-5 text-gray-500 placeholder-gray-400 font-medium focus:ring-2 focus:ring-yellow-400 outline-none"
        />
      </div>

      <button 
        onclick={handleShiftStart}
        class="w-full bg-[#facc15] hover:bg-yellow-500 text-black font-bold py-3.5 rounded-full mt-6 transition-all active:scale-95 shadow-lg text-sm"
      >
        Mulai Shift
      </button>
    </div>
  </div>
</div>

<style>
  /* Menghilangkan ring default browser */
  :global(input:focus) {
    outline: none;
  }
</style>