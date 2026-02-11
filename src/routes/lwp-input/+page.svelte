<script lang="ts">
  // State untuk data profil (dari halaman sebelumnya)
  let user = {
    nama: "TONO WIDIYANTO",
    npk: "KRT-84938",
    noMc: "136A",
    shift: "1"
  };

  // State untuk form input
  let item: string = "";
  let lot: string = "";
  let waktu: string = "";
  let qty: string = "";

  // Tipe data untuk history
  interface HistoryData {
    id: number;
    item: string;
    lot: string;
    waktu: string;
    qty: string;
    timestamp: string;
  }

  // State untuk History Lokal (Live Table)
  let history: HistoryData[] = [];

  import { goto } from '$app/navigation';

  function handleBack() {
    goto('/lwp-setup'); // Redirect to home or another route as needed
  }

  function handleSubmit() {
    if (item && lot && waktu && qty) {
      // Tambah data ke awal array (terbaru di atas)
      const newData = {
        id: Date.now(),
        item,
        lot,
        waktu,
        qty,
        timestamp: new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
      };
      
      history = [newData, ...history];

      // Reset form setelah input
      item = ""; lot = ""; waktu = ""; qty = "";
    } else {
      alert("Mohon isi semua field!");
    }
  }
</script>

<div class="min-h-screen bg-[#f3f3f3] flex flex-col items-center p-4 font-sans pb-10">
  
  <div class="w-full max-w-md bg-[#0f172a] rounded-xl p-5 mb-4 shadow-lg text-white">
    <div class="flex justify-between items-start mb-2">
      <h2 class="font-black text-lg leading-tight uppercase">{user.nama}</h2>
      <span class="text-xs font-bold text-gray-400">NPK: {user.npk}</span>
    </div>
    <div class="text-sm space-y-1">
      <p class="font-semibold">NoMc: <span class="font-normal">{user.noMc}</span></p>
      <p class="font-semibold">Shift: <span class="font-normal">{user.shift}</span></p>
    </div>
  </div>

  <div class="w-full max-w-md bg-[#0f172a] rounded-xl p-6 mb-6 shadow-lg">
    <h3 class="text-white text-center font-black text-xl mb-6 tracking-wide">INPUT FORM LWP</h3>
    
    <div class="space-y-4">
      {#each [
        {label: 'Item', val: 'item', bind: item},
        {label: 'Lot', val: 'lot', bind: lot},
        {label: 'Waktu', val: 'waktu', bind: waktu},
        {label: 'Qty', val: 'qty', bind: qty}
      ] as field}
        <div class="flex flex-col">
          <label class="text-white text-[10px] font-bold mb-1 ml-1 uppercase" for={field.val}>
            {field.label}
          </label>
          <input 
            type="text" 
            id={field.val}
            bind:value={field.bind}
            class="w-full bg-[#eeeeee] border-none rounded-xl py-3 px-4 text-gray-700 font-medium focus:ring-2 focus:ring-yellow-400 outline-none transition-all"
          />
        </div>
      {/each}

      <button 
        onclick={handleSubmit}
        class="w-full bg-yellow-400 hover:bg-yellow-500 text-black font-black py-3 rounded-full mt-4 transition-all active:scale-95 shadow-md"
      >
        SUBMIT DATA
      </button>
    </div>
  </div>

  <div class="w-full max-w-md bg-white rounded-xl shadow-md overflow-hidden">
    <div class="bg-gray-200 px-4 py-2 border-b border-gray-300">
      <h4 class="text-[11px] font-black text-gray-600 uppercase tracking-wider">History Input Hari Ini</h4>
    </div>
    
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead class="bg-gray-50 text-[10px] uppercase text-gray-500 font-bold">
          <tr>
            <th class="px-4 py-2 border-b">Item</th>
            <th class="px-4 py-2 border-b">Lot</th>
            <th class="px-4 py-2 border-b text-center">Waktu</th>
            <th class="px-4 py-2 border-b text-right">Qty</th>
          </tr>
        </thead>
        <tbody class="text-xs">
          {#if history.length === 0}
            <tr>
              <td colspan="4" class="px-4 py-8 text-center text-gray-400 italic">Belum ada data yang dimasukkan.</td>
            </tr>
          {:else}
            {#each history as data (data.id)}
              <tr class="hover:bg-gray-50 transition-colors border-b">
                <td class="px-4 py-3 font-semibold text-gray-800">{data.item}</td>
                <td class="px-4 py-3 text-gray-600">{data.lot}</td>
                <td class="px-4 py-3 text-center text-gray-500">{data.waktu}</td>
                <td class="px-4 py-3 text-right font-bold text-blue-600">{data.qty}</td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>

  <div class="mt-6">
    <button class="text-gray-500 text-sm underline hover:text-gray-700" onclick={handleBack}>
      Selesaikan Shift
    </button>
  </div>

</div>

<style>
  /* Menghaluskan scrollbar jika data banyak */
  ::-webkit-scrollbar {
    width: 4px;
  }
  ::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 10px;
  }
</style>