<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';
    import Chart from 'chart.js/auto';
    import { auth } from '$lib/stores/auth';
    import { goto } from '$app/navigation';
    import annotationPlugin from 'chartjs-plugin-annotation';

    Chart.register(annotationPlugin);

    // Definisi Canvas untuk 3 Grafik
    let canvasTotal: HTMLCanvasElement;
    let canvasOK: HTMLCanvasElement;
    let canvasNG: HTMLCanvasElement;

    // Instance Chart
    let chartTotal: Chart;
    let chartOK: Chart;
    let chartNG: Chart;
    
    $: noMC = $page.url.searchParams.get('no_mc') || '';
    let selectedDate = $page.url.searchParams.get('tanggal') || new Date().toISOString().split('T')[0];
    
    // Perbaikan: Definisi filters object
    let filters = {
        tanggal: selectedDate,
        mesin: noMC
    };

    // Reactivity: Update filters jika URL berubah
    $: filters.mesin = noMC;

    let itemsProduced = "";
    let isLoading = false;
    
    const API_URL = '';

    // Perbaikan: Rename fetchData -> loadChartData sesuai onMount
    async function loadChartData() {
        if (!noMC) return;
        isLoading = true;
        try {
            const res = await fetch(`${API_URL}/api/chart/machine?tanggal=${filters.tanggal}&no_mc=${filters.mesin}`, {
                 headers: {
                    'Authorization': `Bearer ${$auth.token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (res.status === 401) {
                alert("Sesi habis. Silakan login ulang.");
                goto('/');
                return;
            }

            const data = await res.json();
            
            itemsProduced = "-";
            const validItem = data.find((d: any) => d.extra_info && d.extra_info !== '- (-)' && d.extra_info !== '-');
            if (validItem) itemsProduced = validItem.extra_info;

            renderCharts(data);
        } catch (error) {
            console.error("Error fetching machine detail:", error);
        } finally {
            isLoading = false;
        }
    }

    function renderCharts(data: any[]) {
        const labels = data.map((d: any) => d.label); 
        const totalVals = data.map((d: any) => d.actual); 
        const okVals = data.map((d: any) => d.actual_ok); // Tambahan untuk Chart OK
        const ngVals = data.map((d: any) => d.actual_ng); 

        // --- CHART 1: Total Output (Vertical Bar) ---
        if (chartTotal) chartTotal.destroy();
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
                    legend: { display: true, position: 'top' },
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
                                    label += context.raw + ' unit';
                                }
                                return label;
                            }
                        }
                    }
                },
                scales: {
                    y: { beginAtZero: true, title: { display: true, text: 'Output' } },
                    x: { display: false }
                }
            }
        });

        // --- CHART 2: OK Output (Green) ---
        // Perbaikan: Implementasi Chart OK
        if (chartOK) chartOK.destroy();
        
        chartOK = new Chart(canvasOK, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Total OK',
                        data: okVals,
                        backgroundColor: '#10b981',
                        borderColor: '#059669',
                        borderWidth: 1,
                        borderRadius: 4,
                        barPercentage: 0.7
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: { 
                        display: true, 
                        text: `Grafik OK - Mesin ${filters.mesin}`,
                        font: { size: 16, weight: 'bold' }
                    }
                },
                scales: {
                    y: { beginAtZero: true, title: { display: true, text: 'OK' } },
                    x: { display: false }
                }
            }
        });

        // --- CHART 3: NG (Red) ---
        if (chartNG) chartNG.destroy();
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
                        label: 'Max NG',
                        type: 'line',
                        data: Array(labels.length).fill(targetNG),
                        borderColor: '#f59e0b',
                        borderWidth: 2,
                        borderDash: [5, 5],
                        pointRadius: 0
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
                    }
                },
                scales: {
                    y: { beginAtZero: true, title: { display: true, text: 'NG' } },
                    x: { display: false }
                }
            }
        });
    }

    // Perbaikan: Jalankan ulang jika filter tanggal berubah
    $: if (filters.tanggal) {
        loadChartData();
    }

    onMount(() => {
        loadChartData();
    });
</script>

<div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold text-slate-800">Laporan Produksi Per Jam</h1>
        <a href="/manager/prs-ldr" class="text-sm text-indigo-600 hover:underline">Kembali ke Dashboard</a>
    </div>
    <div class="bg-white p-4 rounded-xl shadow-sm border border-slate-200 flex flex-wrap gap-4 items-end">
        <div>
            <label class="block text-xs font-bold text-slate-500 mb-1">Tanggal</label>
            <input type="date" bind:value={filters.tanggal} class="px-3 py-2 border rounded-lg text-sm">
        </div>

        <div class="bg-white px-4 py-2 rounded-lg shadow-sm border border-gray-200 flex items-center gap-4 w-full md:w-auto">
            <div>
                <span class="text-xs text-gray-500 font-bold uppercase">Mesin</span>
                <p class="text-xl font-bold text-blue-700">{noMC}</p>
            </div>
            <div class="h-8 w-px bg-gray-200"></div>
            <div>
                <span class="text-xs text-gray-500 font-bold uppercase">Item Aktif</span>
                <p class="text-sm font-medium text-gray-800 truncate max-w-[200px]" title={itemsProduced}>{itemsProduced}</p>
            </div>
        </div>
    </div>

    {#if isLoading}
        <div class="flex justify-center items-center h-64 bg-white rounded-lg shadow">
            <div class="text-center">
                <i class="fa-solid fa-spinner animate-spin text-4xl text-blue-500"></i>
                <p class="mt-4 text-gray-500 font-medium">Memuat data grafik...</p>
            </div>
        </div>
    {:else}
        <div class="space-y-6">
            <div class="bg-white p-4 rounded-xl shadow border border-blue-100">
                <div class="h-[300px]">
                    <canvas bind:this={canvasTotal}></canvas>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="bg-white p-4 rounded-xl shadow border border-green-100">
                    <div class="h-[300px]">
                        <canvas bind:this={canvasOK}></canvas>
                    </div>
                </div>

                <div class="bg-white p-4 rounded-xl shadow border border-red-100">
                    <div class="h-[300px]">
                        <canvas bind:this={canvasNG}></canvas>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>