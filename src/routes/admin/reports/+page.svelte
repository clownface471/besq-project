<script lang="ts">
  import { onMount } from 'svelte';
  import Chart from 'chart.js/auto';
  import { auth } from '$lib/stores/auth';
  
  const API_URL = "http://localhost:8080";

  let canvasTotal: HTMLCanvasElement;
  let canvasNG: HTMLCanvasElement;
  let chartTotal: Chart;
  let chartNG: Chart;

  // Filter State
  let filters = {
    tanggal: new Date().toISOString().split('T')[0],
    mesin: '11A', // Default contoh
    shift: 1
  };

  let chartData: any[] = [];
  let isLoading = false;

  async function loadChartData() {
    isLoading = true;
    try {
        const res = await fetch(`${API_URL}/admin/chart-production`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${$auth.token}`
            },
            body: JSON.stringify(filters)
        });
        
        const json = await res.json();
        if (res.ok) {
            chartData = json.data || [];
            renderCharts();
        } else {
            console.error("API Error:", json.error);
        }
    } catch (err) {
        console.error("Fetch error:", err);
    } finally {
        isLoading = false;
    }
  }

function renderCharts() {
    const labels = chartData.map(d => d.jam_label);
    const totalVals = chartData.map(d => d.nilai_total);
    const ngVals = chartData.map(d => d.nilai_ng);

    // --- CHART 1: Total Output (Vertical Bar) ---
    if (chartTotal) chartTotal.destroy();
    chartTotal = new Chart(canvasTotal, {
        type: 'bar',
        data: {
            labels: labels, // Label Jam otomatis masuk ke Sumbu X
            datasets: [{
                label: 'Total Output',
                data: totalVals, // Nilai otomatis masuk ke Sumbu Y
                backgroundColor: '#4f46e5',
                borderColor: '#4338ca',
                borderWidth: 1
            }]
        },
        options: {
            // HAPUS baris 'indexAxis: 'y'' agar kembali vertikal (default)
            responsive: true,
            maintainAspectRatio: false,
            plugins: { title: { display: true, text: `Grafik Total Output - Mesin ${filters.mesin}` } },
            scales: {
                y: {
                    beginAtZero: true,
                    title: { display: true, text: 'Jumlah Output' }
                },
                x: {
                    title: { display: true, text: 'Jam Produksi' }
                }
            }
        }
    });

    // --- CHART 2: NG (Vertical Bar) ---
    if (chartNG) chartNG.destroy();
    chartNG = new Chart(canvasNG, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Total NG',
                data: ngVals,
                backgroundColor: '#e11d48',
                borderColor: '#be123c',
                borderWidth: 1
            }]
        },
        options: {
            // HAPUS baris 'indexAxis: 'y''
            responsive: true,
            maintainAspectRatio: false,
            plugins: { title: { display: true, text: `Grafik NG - Mesin ${filters.mesin}` } },
            scales: {
                y: {
                    beginAtZero: true,
                    title: { display: true, text: 'Jumlah NG' }
                },
                x: {
                    title: { display: true, text: 'Jam Produksi' }
                }
            }
        }
    });
  }

  onMount(() => {
      loadChartData();
  });
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold text-slate-800">Laporan Produksi Per Jam</h1>
        <a href="/admin" class="text-sm text-indigo-600 hover:underline">Kembali ke Dashboard</a>
    </div>

    <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 flex flex-wrap gap-4 items-end">
        <div>
            <label class="block text-xs font-bold text-slate-500 mb-1">Tanggal</label>
            <input type="date" bind:value={filters.tanggal} class="px-3 py-2 border rounded-lg text-sm">
        </div>
        <div>
            <label class="block text-xs font-bold text-slate-500 mb-1">No. Mesin</label>
            <input type="text" bind:value={filters.mesin} class="px-3 py-2 border rounded-lg text-sm w-24">
        </div>
        <div>
            <label class="block text-xs font-bold text-slate-500 mb-1">Shift</label>
            <select bind:value={filters.shift} class="px-3 py-2 border rounded-lg text-sm">
                <option value={1}>Shift 1 (07-15)</option>
                <option value={2}>Shift 2 (15-23)</option>
                <option value={3}>Shift 3 (23-07)</option>
            </select>
        </div>
        <button on:click={loadChartData} class="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2 rounded-lg text-sm font-bold transition-colors">
            {isLoading ? 'Loading...' : 'Tampilkan'}
        </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 h-96">
            <canvas bind:this={canvasTotal}></canvas>
        </div>
        <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 h-96">
            <canvas bind:this={canvasNG}></canvas>
        </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <div class="p-4 bg-slate-50 border-b border-slate-100 font-bold text-slate-700">Tabel Data Detail</div>
        <div class="overflow-x-auto">
            <table class="w-full text-sm text-left">
                <thead class="bg-slate-100 text-slate-600 font-bold text-xs uppercase">
                    <tr>
                        <th class="px-6 py-3">Jam</th>
                        <th class="px-6 py-3 text-right">Nilai Total</th>
                        <th class="px-6 py-3 text-right">Nilai NG</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                    {#each chartData as row}
                        <tr class="hover:bg-slate-50">
                            <td class="px-6 py-3 font-mono font-bold text-indigo-600">{row.jam_label}</td>
                            <td class="px-6 py-3 text-right font-medium">{row.nilai_total}</td>
                            <td class="px-6 py-3 text-right font-bold text-rose-600">{row.nilai_ng}</td>
                        </tr>
                    {/each}
                    {#if chartData.length === 0}
                        <tr><td colspan="3" class="px-6 py-8 text-center text-slate-400">Tidak ada data.</td></tr>
                    {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>