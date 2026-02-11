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
    mesin: '11A',
    shift: '1'
  };

  let chartData: any[] = [];
  let isLoading = false;

  function getDummyData() {
    // Dummy data matching screenshot pattern
    return [
      { jam_label: '07:00', nilai_total: 27, nilai_ng: 6 },
      { jam_label: '08:00', nilai_total: 0, nilai_ng: 0 },
      { jam_label: '09:00', nilai_total: 36, nilai_ng: 0 },
      { jam_label: '10:00', nilai_total: 36, nilai_ng: 0 },
      { jam_label: '11:00', nilai_total: 18, nilai_ng: 0 },
      { jam_label: '12:00', nilai_total: 0, nilai_ng: 0 }
    ];
  }

  async function loadChartData() {
    isLoading = true;
    try {
      // Using dummy data for now
      chartData = getDummyData();
      renderCharts();
    } catch (err) {
      console.error("Error:", err);
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
    
    // Tentukan target untuk Total Output (contoh: target 30 unit per jam)
    const targetTotal = 30;
    
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
                    type: 'line', // Garis target
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
                    labels: {
                        usePointStyle: true,
                        padding: 20
                    }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    callbacks: {
                        label: function(context) {
                            let label = context.dataset.label || '';
                            if (label) {
                                label += ': ';
                            }
                            if (context.datasetIndex === 0) {
                                label += context.raw + ' unit';
                                // Tambahkan indikator target
                                const diff = (context.raw as number) - targetTotal;
                                if (diff >= 0) {
                                    label += ` (${diff} di atas target)`;
                                } else {
                                    label += ` (${Math.abs(diff)} di bawah target)`;
                                }
                            } else {
                                label += 'Target: ' + context.raw + ' unit';
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
                                font: {
                                    size: 12,
                                    weight: 'bold'
                                },
                                padding: 6
                            }
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: { 
                        display: true, 
                        text: 'Jumlah Output',
                        font: { weight: 'bold' }
                    },
                    grid: {
                        color: 'rgba(0, 0, 0, 0.1)'
                    }
                },
                x: {
                    title: { 
                        display: true, 
                        text: 'Jam Produksi',
                        font: { weight: 'bold' }
                    },
                    grid: {
                        display: false
                    }
                }
            },
            interaction: {
                intersect: false,
                mode: 'index'
            }
        }
    });

    // --- CHART 2: NG (Vertical Bar) ---
    if (chartNG) chartNG.destroy();
    
    // Tentukan target untuk NG (contoh: maksimal 5 unit per jam)
    const targetNG = 5;
    
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
                    type: 'line', // Garis target
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
                    labels: {
                        usePointStyle: true,
                        padding: 20
                    }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    callbacks: {
                        label: function(context) {
                            let label = context.dataset.label || '';
                            if (label) {
                                label += ': ';
                            }
                            if (context.datasetIndex === 0) {
                                label += context.raw + ' unit';
                                // Tambahkan indikator target
                                const diff = (context.raw as number) - targetNG;
                                if (diff <= 0) {
                                    label += ` (${Math.abs(diff)} di bawah target)`;
                                } else {
                                    label += ` (${diff} di atas target)`;
                                }
                            } else {
                                label += 'Target maksimal: ' + context.raw + ' unit';
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
                                content: `Target: ≤${targetNG}`,
                                position: 'end',
                                backgroundColor: '#f59e0b',
                                color: 'white',
                                font: {
                                    size: 12,
                                    weight: 'bold'
                                },
                                padding: 6
                            }
                        }
                    }
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: { 
                        display: true, 
                        text: 'Jumlah NG',
                        font: { weight: 'bold' }
                    },
                    grid: {
                        color: 'rgba(0, 0, 0, 0.1)'
                    }
                },
                x: {
                    title: { 
                        display: true, 
                        text: 'Jam Produksi',
                        font: { weight: 'bold' }
                    },
                    grid: {
                        display: false
                    }
                }
            },
            interaction: {
                intersect: false,
                mode: 'index'
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
        <a href="/leader" class="text-sm text-indigo-600 hover:underline">Kembali ke Dashboard</a>
    </div>
    <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 flex flex-wrap gap-4 items-end">
        <div>
            <!-- svelte-ignore a11y_label_has_associated_control -->
            <label class="block text-xs font-bold text-slate-500 mb-1">Tanggal</label>
            <input type="date" bind:value={filters.tanggal} class="px-3 py-2 border rounded-lg text-sm">
        </div>
        <div>
            <!-- svelte-ignore a11y_label_has_associated_control -->
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