<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/stores/auth';
    import Chart from 'chart.js/auto';
    import annotationPlugin from 'chartjs-plugin-annotation';

    Chart.register(annotationPlugin);

    let selectedDate = new Date().toISOString().split('T')[0];
    let isLoading = true; 
    
    // Pastikan URL kosong jika proxy sudah diatur di vite.config.ts, atau isi jika perlu full URL
    const API_URL = '';

    // 1. UPDATE STRUKTUR DATA: Tambahkan field 'ng'
    let dashboardData = {
        mixing: { total: 0, pending: 0, completed: 0, ng: 0 },
        cutting: { total: 0, pending: 0, completed: 0, ng: 0 },
        pressing: { total: 0, pending: 0, completed: 0, ng: 0 },
        finishing: { total: 0, pending: 0, completed: 0, ng: 0 }
    };

    let charts: { [key: string]: { total?: Chart; status?: Chart } } = {};

    const divisions = [
        {
            id: 'mixing',
            name: 'Mixing',
            icon: 'fa-solid fa-blender',
            color: 'from-blue-500 to-cyan-500',
            bgColor: 'bg-blue-50',
            borderColor: 'border-blue-200',
            textColor: 'text-blue-700',
            description: 'Proses pencampuran bahan baku',
            chartColor: '#3b82f6'
        },
        {
            id: 'cutting',
            name: 'Cutting',
            icon: 'fa-solid fa-scissors',
            color: 'from-purple-500 to-pink-500',
            bgColor: 'bg-purple-50',
            borderColor: 'border-purple-200',
            textColor: 'text-purple-700',
            description: 'Proses pemotongan material',
            chartColor: '#a855f7'
        },
        {
            id: 'pressing',
            name: 'Pressing',
            icon: 'fa-solid fa-square',
            color: 'from-orange-500 to-red-500',
            bgColor: 'bg-orange-50',
            borderColor: 'border-orange-200',
            textColor: 'text-orange-700',
            description: 'Proses pengepresan produk',
            chartColor: '#f97316'
        },
        {
            id: 'finishing',
            name: 'Finishing',
            icon: 'fa-solid fa-star',
            color: 'from-green-500 to-emerald-500',
            bgColor: 'bg-green-50',
            borderColor: 'border-green-200',
            textColor: 'text-green-700',
            description: 'Proses penyelesaian akhir',
            chartColor: '#22c55e'
        }
    ];

    async function fetchDashboardData() {
        isLoading = true;
        try {
            const res = await fetch(`${API_URL}/api/chart/manager?tanggal=${selectedDate}`, {
                headers: { Authorization: `Bearer ${$auth.token}` }
            });

            if (res.ok) {
                const result = await res.json();
                
                // Reset data
                const newData = {
                    mixing: { total: 0, pending: 0, completed: 0, ng: 0 },
                    cutting: { total: 0, pending: 0, completed: 0, ng: 0 },
                    pressing: { total: 0, pending: 0, completed: 0, ng: 0 },
                    finishing: { total: 0, pending: 0, completed: 0, ng: 0 }
                };

                // Mapping Data
                result.forEach((item: any) => {
                    let key = '';
                    const label = item.label ? item.label.toUpperCase() : '';

                    if (label.includes('PRS') || label.includes('PRESS')) key = 'pressing';
                    else if (label.includes('CT') || label.includes('CUT')) key = 'cutting';
                    else if (label.includes('MIX')) key = 'mixing';
                    else if (label.includes('FIN')) key = 'finishing';

                    if (key && newData[key as keyof typeof newData]) {
                        newData[key as keyof typeof newData] = {
                            total: item.target, 
                            completed: item.actual,
                            pending: Math.max(0, item.target - item.actual),
                            ng: item.actual_ng || 0 // 2. UPDATE: Ambil data NG dari backend
                        };
                    }
                });
                
                dashboardData = newData;
            }
        } catch (error) {
            console.error('Error fetching dashboard data:', error);
        } finally {
            isLoading = false;
            initializeCharts();
        }
    }

    function initializeCharts() {
        setTimeout(() => {
            divisions.forEach((division) => {
                const data = dashboardData[division.id as keyof typeof dashboardData];
                
                if (charts[division.id]?.total) charts[division.id].total?.destroy();
                if (charts[division.id]?.status) charts[division.id].status?.destroy();

                // 3. LOGIKA TARGET DINAMIS
                const productionTarget = data.total || 0; 
                const notGoodTarget = Math.ceil(productionTarget * 0.02); // Max NG = 2% dari Target
                
                // Chart Total
                const totalCanvasId = `chart-total-${division.id}`;
                const totalCanvas = document.getElementById(totalCanvasId) as HTMLCanvasElement;
                if (totalCanvas) {
                    // Agar garis target tidak kepotong di atas
                    const maxValue = Math.max(productionTarget, (data.completed || 0) + (data.pending || 0)) * 1.2;

                    const newChart = new Chart(totalCanvas, {
                        type: 'bar',
                        data: {
                            labels: ['Total Produksi'],
                            datasets: [
                                {
                                    label: 'Actual',
                                    data: [data.completed || 0],
                                    backgroundColor: '#3b82f6',
                                    borderRadius: 4,
                                    borderSkipped: false,
                                    stack: 'Stack 0'
                                },
                                {
                                    label: 'Pending',
                                    data: [data.pending || 0],
                                    backgroundColor: '#e2e8f0', 
                                    borderRadius: 4,
                                    borderSkipped: false,
                                    stack: 'Stack 0'
                                }
                            ]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: true,
                            plugins: {
                                legend: { display: false },
                                tooltip: {
                                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                                    padding: 8
                                },
                                annotation: {
                                    annotations: {
                                        targetLine: {
                                            type: 'line',
                                            yMin: productionTarget,
                                            yMax: productionTarget,
                                            borderColor: '#10b981',
                                            borderWidth: 2,
                                            borderDash: [5, 5],
                                            label: {
                                                content: `Target: ${productionTarget}`,
                                                display: true,
                                                position: 'end',
                                                backgroundColor: '#10b981',
                                                color: 'white',
                                                font: { size: 10, weight: 'bold' },
                                                padding: 4,
                                                borderRadius: 4,
                                                yAdjust: -10
                                            }
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: {
                                    beginAtZero: true,
                                    max: maxValue, 
                                    ticks: { font: { size: 10 } },
                                    grid: { color: '#f0f0f0' }
                                },
                                x: { display: false }
                            }
                        }
                    });
                    if (!charts[division.id]) charts[division.id] = {};
                    charts[division.id].total = newChart;
                }

                // Chart NG
                const notGoodCanvasId = `chart-notgood-${division.id}`;
                const notGoodCanvas = document.getElementById(notGoodCanvasId) as HTMLCanvasElement;
                if (notGoodCanvas) {
                    const notGoodCount = data.ng || 0; // 4. DATA NG REAL
                    
                    const newChart = new Chart(notGoodCanvas, {
                        type: 'bar',
                        data: {
                            labels: ['NG'],
                            datasets: [
                                {
                                    label: 'Status',
                                    data: [notGoodCount], // Hapus dummy goodCount
                                    backgroundColor: '#ef4444',
                                    borderRadius: 4,
                                    borderSkipped: false
                                }
                            ]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: true,
                            plugins: {
                                legend: { display: false },
                                tooltip: {
                                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                                    padding: 8,
                                    displayColors: false,
                                    callbacks: {
                                        label: function(context) {
                                            return `NG: ${context.raw} pcs`;
                                        }
                                    }
                                },
                                annotation: {
                                    annotations: {
                                        limitLine: {
                                            type: 'line',
                                            yMin: notGoodTarget,
                                            yMax: notGoodTarget,
                                            borderColor: '#f59e0b',
                                            borderWidth: 2,
                                            borderDash: [5, 5],
                                            label: {
                                                content: `Max: ${notGoodTarget}`,
                                                display: true,
                                                position: 'end',
                                                backgroundColor: '#f59e0b',
                                                color: 'white',
                                                font: { size: 10, weight: 'bold' },
                                                padding: 4,
                                                borderRadius: 4,
                                                yAdjust: -10
                                            }
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: {
                                    beginAtZero: true,
                                    suggestedMax: Math.max(notGoodCount, notGoodTarget) * 1.5,
                                    ticks: { font: { size: 10 } },
                                    grid: { color: '#f0f0f0' }
                                },
                                x: { display: false }
                            }
                        }
                    });
                    if (!charts[division.id]) charts[division.id] = {};
                    charts[division.id].status = newChart;
                }
            });
        }, 100);
    }

    function getStatusColor(divisionId: string) {
        if (divisionId === 'cutting') return { badge: 'bg-red-500', label: 'Problem' };
        if (divisionId === 'mixing') return { badge: 'bg-green-500', label: 'Lancar' };
        if (divisionId === 'pressing') return { badge: 'bg-yellow-500', label: 'Sedang Berjalan' };
        if (divisionId === 'finishing') return { badge: 'bg-green-500', label: 'Lancar' };
        return { badge: 'bg-slate-500', label: 'Unknown' };
    }

    function handleDetailClick(divisionId: string) {
        if (divisionId === 'pressing') {
            goto('/leader'); 
        } else {
             goto(`/report-chart/line?division=${divisionId}`);
        }
    }

    onMount(() => {
        fetchDashboardData();
    });
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-slate-100 p-8">
    <div class="max-w-7xl mx-auto">
        <div class="mb-12">
            <h1 class="text-4xl font-bold text-slate-900 mb-3">Dashboard Manager</h1>
            <p class="text-slate-600 text-lg">Pantau performa setiap divisi produksi secara real-time</p>
        </div>

        {#if isLoading}
            <div class="flex justify-center items-center min-h-64">
                <div class="text-center">
                    <div class="inline-block animate-spin">
                        <i class="fa-solid fa-spinner text-indigo-600 text-4xl"></i>
                    </div>
                    <p class="text-slate-600 mt-4">Memuat data...</p>
                </div>
            </div>
        {:else}
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-10">
                {#each divisions as division (division.id)}
                    {@const status = getStatusColor(division.id)}
                    {@const data = dashboardData[division.id as keyof typeof dashboardData]}
                    <div class="group relative bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border border-slate-100 hover:border-slate-200">
                        <div class="absolute top-0 left-0 right-0 h-1.5 bg-gradient-to-r {division.color}"></div>
                        <div class="p-8 relative z-10">
                            <div class="flex items-center gap-4 mb-8">
                                <div class="{division.bgColor} p-4 rounded-xl group-hover:scale-110 transition-transform shadow-md">
                                    <i class="{division.icon} {division.textColor} text-2xl"></i>
                                </div>
                                <div class="flex-1">
                                    <h3 class="text-xl font-bold text-slate-800">{division.name}</h3>
                                    <p class="text-sm text-slate-500 mt-1">{division.description}</p>
                                </div>
                            </div>
                            <div class="mb-8 flex items-center gap-3 px-4 py-3 bg-slate-50 rounded-lg border border-slate-200">
                                <span class="{status.badge} w-3 h-3 rounded-full animate-pulse"></span>
                                <span class="text-slate-700 text-sm font-bold flex-1">{status.label}</span>
                                <span class="text-slate-600 text-sm font-semibold px-3 py-1 bg-white rounded-md">Total: <span class="text-slate-900 font-bold">{data.total}</span></span>
                            </div>
                            <div class="grid grid-cols-2 gap-6">
                                <div class="flex flex-col items-center p-4 bg-blue-50 rounded-lg border border-blue-200">
                                    <p class="text-xs text-blue-700 font-bold mb-3">Total Output</p>
                                    <canvas id="chart-total-{division.id}" style="max-height: 120px;"></canvas>
                                </div>
                                <div class="flex flex-col items-center p-4 bg-red-50 rounded-lg border border-red-200">
                                    <p class="text-xs text-red-700 font-bold mb-3">Total NG</p>
                                    <canvas id="chart-notgood-{division.id}" style="max-height: 120px;"></canvas>
                                </div>
                            </div>
                        </div>
                        <div class="px-8 py-4 bg-slate-50 border-t border-slate-100">
                            <button on:click={() => handleDetailClick(division.id)} class="w-full py-2 px-4 rounded-lg text-sm font-semibold text-white bg-slate-700 hover:bg-slate-800 transition-all duration-300">
                                Lihat Detail
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>