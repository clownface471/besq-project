<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import Chart from 'chart.js/auto';
  import annotationPlugin from 'chartjs-plugin-annotation';
  import { auth } from '$lib/stores/auth';

  // Registrasi plugin
  Chart.register(annotationPlugin);

  let canvasTotal: HTMLCanvasElement;
  let canvasNG: HTMLCanvasElement;
  let chartTotal: Chart;
  let chartNG: Chart;

  // 1. Ambil ID & Tanggal dari URL
  // Karena folder sudah diubah jadi 'mc-detail', pastikan URL param terbaca
  let urlMachineId = $page.url.searchParams.get('no_mc') || $page.url.searchParams.get('id') || '11A';
  let urlDate = $page.url.searchParams.get('date') || new Date().toISOString().split('T')[0];

  // State Filter
  let filters = {
    tanggal: urlDate,
    mesin: urlMachineId,
    shift: '1'
  };

  let chartData: any[] = [];
  let isLoading = false;

  async function loadChartData() {
    isLoading = true;
    try {
      // PERBAIKAN: Ubah 'machine-detail' menjadi 'machine' sesuai main.go
      const res = await fetch(`/api/chart/machine?tanggal=${filters.tanggal}&no_mc=${filters.mesin}`, {
          headers: { Authorization: `Bearer ${$auth.token}` }
      });

      if (res.ok) {
          chartData = await res.json();
          if (!chartData) chartData = [];
          renderCharts();
      } else {
          console.error("Gagal load data, status:", res.status);
          chartData = [];
          renderCharts();
      }
    } catch (err) {
      console.error("Error:", err);
    } finally {
      isLoading = false;
    }
  }

  function renderCharts() {
    // 3. Mapping Data (Perbaikan Kunci Data)
    // Backend Go: label, actual, actual_ng
    const labels = chartData.map(d => d.label);
    const totalVals = chartData.map(d => d.actual);
    const ngVals = chartData.map(d => d.actual_ng);
    
    // Ambil target dari data pertama (jika ada) atau default 30
    const targetTotal = chartData.length > 0 ? (chartData[0].target || 30) : 30;
    const targetNG = 5; // Target NG static atau hitung % (misal Math.ceil(targetTotal * 0.05))

    // --- CHART 1: Total Output ---
    if (chartTotal) chartTotal.destroy();
    
    chartTotal = new Chart(canvasTotal, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Total Output',
                    data: totalVals,
                    backgroundColor: '#4f46e5',
                    borderColor: '#4338ca',
                    borderWidth: 1,
                    borderRadius: 4,
                    barPercentage: 0.7
                },
                {
                    label: 'Target Output',
                    type: 'line',
                    data: Array(labels.length).fill(targetTotal),
                    borderColor: '#10b981',
                    borderWidth: 2,
                    borderDash: [5, 5],
                    pointRadius: 0,
                    fill: false,
                    tension: 0
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: { 
                    display: true, 
                    text: `Grafik Total Output - Mesin ${filters.mesin}`,
                    font: { size: 16, weight: 'bold' }
                },
                legend: {
                    display: true,
                    position: 'top',
                    labels: { usePointStyle: true, padding: 20 }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    callbacks: {
                        label: function(context) {
                            let label = context.dataset.label || '';
                            if (label) label += ': ';
                            if (context.datasetIndex === 0) {
                                label += context.raw + ' unit';
                                const diff = (context.raw as number) - targetTotal;
                                label += diff >= 0 ? ` (+${diff})` : ` (${diff})`;
                            } else {
                                label += 'Target: ' + context.raw;
                            }
                            return label;
                        }
                    }
                },
                annotation: {
                    annotations: {
                        targetLine: {
                            type: 'line',
                            yMin: targetTotal,
                            yMax: targetTotal,
                            borderColor: '#10b981',
                            borderWidth: 2,
                            borderDash: [5, 5],
                            label: {
                                content: `Target: ${targetTotal}`,
                                position: 'end',
                                backgroundColor: '#10b981',
                                color: 'white',
                                font: { size: 12, weight: 'bold' },
                                padding: 6
                            }
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: { display: true, text: 'Jumlah Output', font: { weight: 'bold' } },
                    grid: { color: 'rgba(0, 0, 0, 0.1)' }
                },
                x: {
                    title: { display: true, text: 'Jam Produksi', font: { weight: 'bold' } },
                    grid: { display: false }
                }
            },
            interaction: { intersect: false, mode: 'index' }
        }
    });

    // --- CHART 2: NG ---
    if (chartNG) chartNG.destroy();
    
    chartNG = new Chart(canvasNG, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Total NG',
                    data: ngVals,
                    backgroundColor: '#e11d48',
                    borderColor: '#be123c',
                    borderWidth: 1,
                    borderRadius: 4,
                    barPercentage: 0.7
                },
                {
                    label: 'Target Maksimal NG',
                    type: 'line',
                    data: Array(labels.length).fill(targetNG),
                    borderColor: '#f59e0b',
                    borderWidth: 2,
                    borderDash: [5, 5],
                    pointRadius: 0,
                    fill: false,
                    tension: 0
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: { 
                    display: true, 
                    text: `Grafik NG - Mesin ${filters.mesin}`,
                    font: { size: 16, weight: 'bold' }
                },
                legend: {
                    display: true,
                    position: 'top',
                    labels: { usePointStyle: true, padding: 20 }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    callbacks: {
                        label: function(context) {
                            let label = context.dataset.label || '';
                            if (label) label += ': ';
                            if (context.datasetIndex === 0) {
                                label += context.raw + ' unit';
                            } else {
                                label += 'Max: ' + context.raw;
                            }
                            return label;
                        }
                    }
                },
                annotation: {
                    annotations: {
                        targetLine: {
                            type: 'line',
                            yMin: targetNG,
                            yMax: targetNG,
                            borderColor: '#f59e0b',
                            borderWidth: 2,
                            borderDash: [5, 5],
                            label: {
                                content: `Max: ${targetNG}`,
                                position: 'end',
                                backgroundColor: '#f59e0b',
                                color: 'white',
                                font: { size: 12, weight: 'bold' },
                                padding: 6
                            }
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: { display: true, text: 'Jumlah NG', font: { weight: 'bold' } },
                    grid: { color: 'rgba(0, 0, 0, 0.1)' }
                },
                x: {
                    title: { display: true, text: 'Jam Produksi', font: { weight: 'bold' } },
                    grid: { display: false }
                }
            },
            interaction: { intersect: false, mode: 'index' }
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
        <a href="/leader" class="text-sm text-indigo-600 hover:underline">Kembali ke Dashboard</a>
    </div>
    <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 flex flex-wrap gap-4 items-end">
        <div>
            <label class="block text-xs font-bold text-slate-500 mb-1">Tanggal</label>
            <input type="date" bind:value={filters.tanggal} class="px-3 py-2 border rounded-lg text-sm">
        </div>
        <div>
            <label class="block text-xs font-bold text-slate-500 mb-1">Shift</label>
            <select bind:value={filters.shift} class="px-3 py-2 border rounded-lg text-sm">
                <option value="1">Shift 1 (07-15)</option>
                <option value="2">Shift 2 (15-23)</option>
                <option value="3">Shift 3 (23-07)</option>
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
                            <td class="px-6 py-3 font-mono font-bold text-indigo-600">{row.label}</td>
                            <td class="px-6 py-3 text-right font-medium">{row.actual}</td>
                            <td class="px-6 py-3 text-right font-bold text-rose-600">{row.actual_ng}</td>
                        </tr>
                    {/each}
                    {#if chartData.length === 0 && !isLoading}
                        <tr><td colspan="3" class="px-6 py-8 text-center text-slate-400">Tidak ada data.</td></tr>
                    {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>