<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import Chart from 'chart.js/auto';
  import annotationPlugin from 'chartjs-plugin-annotation';
  import { auth } from '$lib/stores/auth';
  import html2canvas from 'html2canvas';

  // Registrasi plugin
  Chart.register(annotationPlugin);

  // Ambil ID & Tanggal dari URL
  let urlMachineId = $page.url.searchParams.get('no_mc') || $page.url.searchParams.get('id') || '11A';
  let urlDate = $page.url.searchParams.get('date') || new Date().toISOString().split('T')[0];

  // State Filter - shift default ke "1"
  let filters = {
    tanggal: urlDate,
    mesin: urlMachineId,
    shift: '1'
  };

  let chartData: any[] = [];
  let isLoading = false;
  let canvasTotal: HTMLCanvasElement;
  let canvasNG: HTMLCanvasElement;
  let chartTotal: Chart;
  let chartNG: Chart;

  // State untuk export
  let exportContainer: HTMLElement;
  let isExporting = false;

  // Shift labels for display
  const shiftLabels: Record<string, string> = {
    '1': 'Shift 1 (00:00 - 08:00)',
    '2': 'Shift 2 (08:00 - 16:00)',
    '3': 'Shift 3 (16:00 - 00:00)'
  };

  async function loadChartData() {
    isLoading = true;
    try {
      // IMPORTANT: Include shift parameter in API call
      const res = await fetch(`/api/chart/machine?tanggal=${filters.tanggal}&no_mc=${filters.mesin}&shift=${filters.shift}`, {
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
      chartData = [];
      renderCharts();
    } finally {
      isLoading = false;
    }
  }

  function renderCharts() {
    // Mapping Data
    const labels = chartData.map(d => d.label);
    const totalVals = chartData.map(d => d.actual);
    const ngVals = chartData.map(d => d.actual_ng);
    
    // Ambil target dari data pertama (jika ada) atau default 30
    const targetTotal = chartData.length > 0 ? (chartData[0].target || 30) : 30;
    const targetNG = Math.ceil(targetTotal * 0.05); // Target NG = 5% dari target total

    // Destroy existing charts
    if (chartTotal) chartTotal.destroy();
    if (chartNG) chartNG.destroy();

    // --- CHART 1: Total Output ---
    if (canvasTotal) {
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
                      text: `Total Output - Mesin ${filters.mesin} (${shiftLabels[filters.shift]})`,
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
    }

    // --- CHART 2: NG ---
    if (canvasNG) {
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
                      text: `Total NG - Mesin ${filters.mesin} (${shiftLabels[filters.shift]})`,
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
  }

  async function exportToImage() {
      if (!exportContainer) return;
      isExporting = true;
      
      try {
          // Pakai "as any" biar TypeScript nggak bawel soal properti scale
          const canvas = await html2canvas(exportContainer, { 
              scale: 2, 
              useCORS: true,
              backgroundColor: '#f8fafc' 
          } as any);
          
          const image = canvas.toDataURL('image/png');
          const link = document.createElement('a');
          
          link.download = `Laporan_Produksi_Mesin_${filters.mesin}_Shift${filters.shift}_${filters.tanggal}.png`;
          link.href = image;
          link.click();
      } catch (error) {
          console.error("Gagal mengekspor gambar:", error);
          alert("Gagal mengekspor laporan. Cek console log untuk detailnya.");
      } finally {
          isExporting = false;
      }
  }

  onMount(() => {
      loadChartData();
  });
</script>

<div bind:this={exportContainer} class="p-6 max-w-7xl mx-auto space-y-6 bg-slate-50 min-h-screen">
    <div class="flex justify-between items-center">
        <div>
            <h1 class="text-2xl font-bold text-slate-800">Laporan Produksi Per Jam</h1>
            <p class="text-sm text-slate-500 mt-1">Mesin {filters.mesin} - {shiftLabels[filters.shift]}</p>
        </div>
        
        <div class="flex gap-3 items-center">
            <button 
                on:click={exportToImage} 
                disabled={isExporting || chartData.length === 0}
                class="bg-emerald-600 hover:bg-emerald-700 disabled:bg-emerald-300 text-white px-4 py-2 rounded-lg text-sm font-bold transition-colors flex items-center gap-2 shadow-sm"
            >
                {#if isExporting}
                    <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Mengekspor...
                {:else}
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
                    </svg>
                    Simpan Laporan
                {/if}
            </button>
            <a href="/leader" class="text-sm text-indigo-600 hover:underline font-medium">← Kembali ke Dashboard</a>
        </div>
    </div>

    <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 flex flex-wrap gap-4 items-end data-ignore-export">
        <div>
            <label for="filter-tanggal" class="block text-xs font-bold text-slate-500 mb-1">Tanggal</label>
            <input 
                id="filter-tanggal"
                type="date" 
                bind:value={filters.tanggal} 
                class="px-3 py-2 border rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
        </div>
        <div>
            <label for="filter-shift" class="block text-xs font-bold text-slate-500 mb-1">Shift</label>
            <select 
                id="filter-shift"
                bind:value={filters.shift} 
                class="px-3 py-2 border rounded-lg text-sm bg-white focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
                <option value="1">Shift 1 (00:00 - 08:00)</option>
                <option value="2">Shift 2 (08:00 - 16:00)</option>
                <option value="3">Shift 3 (16:00 - 00:00)</option>
            </select>
        </div>
        <button 
            on:click={loadChartData} 
            disabled={isLoading}
            class="bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-300 text-white px-5 py-2 rounded-lg text-sm font-bold transition-colors flex items-center gap-2"
        >
            {#if isLoading}
                <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Loading...
            {:else}
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Tampilkan
            {/if}
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
        <div class="p-4 bg-slate-50 border-b border-slate-100 font-bold text-slate-700 flex items-center justify-between">
            <span>Tabel Data Detail</span>
            {#if chartData.length > 0}
                <span class="text-xs font-normal text-slate-500">
                    Total: {chartData.reduce((sum, row) => sum + (row.actual || 0), 0)} unit | 
                    NG: {chartData.reduce((sum, row) => sum + (row.actual_ng || 0), 0)} unit
                </span>
            {/if}
        </div>
        <div class="overflow-x-auto">
            <table class="w-full text-sm text-left">
                <thead class="bg-slate-100 text-slate-600 font-bold text-xs uppercase">
                    <tr>
                        <th class="px-6 py-3">Jam</th>
                        <th class="px-6 py-3">Nama Item</th>
                        <th class="px-6 py-3 text-right">Target</th>
                        <th class="px-6 py-3 text-right">Nilai Total</th>
                        <th class="px-6 py-3 text-right">Nilai NG</th>
                        <th class="px-6 py-3 text-right">Achievement</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                    {#each chartData as row}
                        {@const achievement = row.target > 0 ? ((row.actual / row.target) * 100).toFixed(1) : 0}
                        <tr class="hover:bg-slate-50 transition-colors">
                            <td class="px-6 py-3 font-mono font-bold text-indigo-600">{row.label}</td>
                            <td class="px-6 py-3 text-slate-700 font-medium max-w-xs truncate" title={row.item_name || '-'}>{row.item_name || '-'}</td>
                            <td class="px-6 py-3 text-right font-medium text-slate-600">{Math.round(row.target || 0)}</td>
                            <td class="px-6 py-3 text-right font-bold text-slate-800">{row.actual}</td>
                            <td class="px-6 py-3 text-right font-bold text-rose-600">{row.actual_ng}</td>
                            <td class="px-6 py-3 text-right">
                                <span class="inline-flex items-center px-2 py-1 rounded-md text-xs font-bold {
                                    parseFloat(achievement as string) >= 100 ? 'bg-emerald-100 text-emerald-700' :
                                    parseFloat(achievement as string) >= 80 ? 'bg-amber-100 text-amber-700' :
                                    'bg-rose-100 text-rose-700'
                                }">
                                    {achievement}%
                                </span>
                            </td>
                        </tr>
                    {/each}
                    {#if chartData.length === 0 && !isLoading}
                        <tr><td colspan="6" class="px-6 py-8 text-center text-slate-400">Tidak ada data untuk shift dan tanggal yang dipilih.</td></tr>
                    {/if}
                    {#if isLoading}
                        <tr><td colspan="6" class="px-6 py-8 text-center text-slate-400">
                            <svg class="animate-spin h-6 w-6 mx-auto mb-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Memuat data...
                        </td></tr>
                    {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>