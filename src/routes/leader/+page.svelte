<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/stores/auth';
    import { ROLES } from '$lib/roles';
    import Chart from 'chart.js/auto';
    import annotationPlugin from 'chartjs-plugin-annotation';

    Chart.register(annotationPlugin);

    let machineData: any[] = [];
    let isLoading = true; 
    let charts: Record<string, any> = {};

    function goBack() {
        goto('/');
    }

    function selectMachine(machine: any) {
        goto(`/leader/mc-detail?no_mc=${machine.id}`);
    }

    function getMachineStatusColor(machine: any) {
        if (machine.isProblem) {
            return {
                bg: 'bg-red-500 hover:bg-red-600',
                badge: 'bg-red-700',
                text: 'text-white',
                indicator: '🔴'
            };
        }
        return {
            bg: 'bg-green-500 hover:bg-green-600',
            badge: 'bg-green-700',
            text: 'text-white',
            indicator: '🟢'
        };
    }

    async function fetchMachineStatus() {
        isLoading = true;
        try {
            const today = new Date().toISOString().split('T')[0];
            const res = await fetch(`/api/chart/process?tanggal=${today}&proses=PRS`, {
                headers: { Authorization: `Bearer ${$auth.token}` }
            });

            if (res.ok) {
                const data = await res.json();
                
                machineData = data.map((item: any) => {
                    // PERBAIKAN: Gunakan nilai RAW (Angka Asli)
                    const target = Math.round(item.target || 0);
                    const actual = Math.round(item.actual || 0);
                    const ng = Math.round(item.actual_ng || 0); // Ambil data NG dari backend

                    // Hitung % hanya untuk logika status warna (isProblem)
                    const achievement = target > 0 ? (actual / target) * 100 : 0;
                    const ngRate = (actual + ng) > 0 ? (ng / (actual + ng)) * 100 : 0;

                    return {
                        id: item.label,
                        name: `Mesin Press ${item.label}`,
                        target: target,      // Simpan target per mesin
                        completed: actual,   // Angka aktual
                        notGood: ng,         // Angka aktual NG
                        // Logic Status: Problem jika pencapaian < 80% ATAU NG > 5% (Sesuaikan standar pabrik)
                        isProblem: achievement < 80 || ngRate > 5 
                    };
                });
            }
        } catch (e) {
            console.error("Gagal ambil data mesin:", e);
        } finally {
            isLoading = false;
            setTimeout(initializeCharts, 100);
        }
    }

    function initializeCharts() {
        Object.keys(charts).forEach(machineId => {
            if (charts[machineId]?.total) charts[machineId].total.destroy();
            if (charts[machineId]?.notGood) charts[machineId].notGood.destroy();
        });
        charts = {};

        machineData.forEach((machine) => {
            // Target per mesin (bukan global 90% lagi)
            const machineTarget = machine.target; 
            const maxScale = Math.max(machineTarget, machine.completed) * 1.2; // Skala grafik sedikit lebih tinggi dari nilai terbesar

            // --- Chart Total Produksi (Angka Aktual) ---
            const totalCanvasId = `chart-total-${machine.id}`;
            const totalCanvas = document.getElementById(totalCanvasId) as HTMLCanvasElement;
            
            if (totalCanvas) {
                const ctx = totalCanvas.getContext('2d');
                if (ctx) {
                    const gradient = ctx.createLinearGradient(0, 0, 0, totalCanvas.height);
                    gradient.addColorStop(0, '#3b82f6');
                    gradient.addColorStop(1, '#1d4ed8');

                    const newChart = new Chart(totalCanvas, {
                        type: 'bar',
                        data: {
                            labels: ['Produksi'],
                            datasets: [{
                                label: 'Produksi',
                                data: [machine.completed],
                                backgroundColor: gradient,
                                borderColor: '#1e40af',
                                borderWidth: 1,
                                borderRadius: 6,
                                barPercentage: 0.6,
                                categoryPercentage: 0.8
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            animation: { duration: 1000, easing: 'easeOutQuart' },
                            plugins: {
                                legend: { display: false },
                                tooltip: {
                                    backgroundColor: 'rgba(15, 23, 42, 0.9)',
                                    padding: 10,
                                    displayColors: false,
                                    titleFont: { size: 12 },
                                    bodyFont: { size: 14, weight: 'bold' },
                                    callbacks: {
                                        label: function(context) {
                                            const value = context.raw as number;
                                            const diff = value - machineTarget;
                                            const status = diff >= 0 ? '✓' : '✗';
                                            return [
                                                `${status} Aktual: ${value} pcs`, 
                                                `Target: ${machineTarget} pcs`, 
                                                `Selisih: ${diff >= 0 ? '+' : ''}${diff} pcs`
                                            ];
                                        }
                                    }
                                },
                                annotation: {
                                    annotations: {
                                        targetLine: {
                                            type: 'line',
                                            yMin: machineTarget,
                                            yMax: machineTarget,
                                            borderColor: '#10b981',
                                            borderWidth: 2,
                                            borderDash: [5, 5],
                                            label: {
                                                content: `Target: ${machineTarget}`,
                                                position: 'end',
                                                backgroundColor: '#10b981',
                                                color: 'white',
                                                font: { size: 9, weight: 'bold' },
                                                padding: 4,
                                                borderRadius: 4
                                            }
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: {
                                    beginAtZero: true,
                                    max: maxScale, // Skala dinamis
                                    ticks: { font: { size: 9 }, color: '#64748b', padding: 4 },
                                    grid: { color: 'rgba(100, 116, 139, 0.1)', drawOnChartArea: true, drawTicks: false },
                                    border: { display: false }
                                },
                                x: { display: false }
                            }
                        }
                    });
                    if (!charts[machine.id]) charts[machine.id] = {};
                    charts[machine.id].total = newChart;
                }
            }

            // --- Chart Not Good (Angka Aktual) ---
            const notGoodCanvasId = `chart-notgood-${machine.id}`;
            const notGoodCanvas = document.getElementById(notGoodCanvasId) as HTMLCanvasElement;
            
            if (notGoodCanvas) {
                const ctx = notGoodCanvas.getContext('2d');
                if (ctx) {
                    const gradient = ctx.createLinearGradient(0, 0, 0, notGoodCanvas.height);
                    gradient.addColorStop(0, '#ef4444');
                    gradient.addColorStop(1, '#dc2626');

                    // Target NG (Misal toleransi 5% dari target produksi)
                    const maxNG = Math.ceil(machineTarget * 0.05); 

                    const newChart = new Chart(notGoodCanvas, {
                        type: 'bar',
                        data: {
                            labels: ['NG'],
                            datasets: [{
                                label: 'NG',
                                data: [machine.notGood],
                                backgroundColor: gradient,
                                borderColor: '#b91c1c',
                                borderWidth: 1,
                                borderRadius: 6,
                                barPercentage: 0.6,
                                categoryPercentage: 0.8
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            animation: { duration: 1000, easing: 'easeOutQuart' },
                            plugins: {
                                legend: { display: false },
                                tooltip: {
                                    backgroundColor: 'rgba(15, 23, 42, 0.9)',
                                    padding: 10,
                                    displayColors: false,
                                    titleFont: { size: 12 },
                                    bodyFont: { size: 14, weight: 'bold' },
                                    callbacks: {
                                        label: function(context) {
                                            const value = context.raw as number;
                                            return [`Jumlah NG: ${value} pcs`];
                                        }
                                    }
                                },
                                // Opsional: Garis batas maksimal NG
                                annotation: {
                                    annotations: {
                                        limitLine: {
                                            type: 'line',
                                            yMin: maxNG,
                                            yMax: maxNG,
                                            borderColor: '#f59e0b',
                                            borderWidth: 1,
                                            borderDash: [2, 2],
                                            label: {
                                                content: `Max: ${maxNG}`,
                                                position: 'start',
                                                backgroundColor: '#f59e0b',
                                                color: 'white',
                                                font: { size: 8 },
                                                padding: 2,
                                                borderRadius: 2
                                            }
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: {
                                    beginAtZero: true,
                                    // Berikan sedikit ruang di atas
                                    suggestedMax: Math.max(machine.notGood, maxNG) * 1.5, 
                                    ticks: { font: { size: 9 }, color: '#64748b', padding: 4 },
                                    grid: { color: 'rgba(100, 116, 139, 0.1)', drawOnChartArea: true, drawTicks: false },
                                    border: { display: false }
                                },
                                x: { display: false }
                            }
                        }
                    });
                    if (!charts[machine.id]) charts[machine.id] = {};
                    charts[machine.id].notGood = newChart;
                }
            }
        });
    }

    // Reactive variables
    $: pressMachines = machineData.filter(m => ['07B', '16B', '16A', '12B', '12A', '14B', '14A', '09A', '09B', '11A', '11B', '15A', '15B', '02A', '04A', '04B', '13B', '13A', '10B', '10A'].includes(m.id));
    $: injectMachines = machineData.filter(m => ['20A', '21A', '22A', '19A', '18A', '17A'].includes(m.id));

    onMount(() => {
        fetchMachineStatus();
        window.addEventListener('resize', handleResize);
    });

    function handleResize() {
        setTimeout(initializeCharts, 300);
    }

    onDestroy(() => {
        Object.keys(charts).forEach(machineId => {
            if (charts[machineId]?.total) charts[machineId].total.destroy();
            if (charts[machineId]?.notGood) charts[machineId].notGood.destroy();
        });
        window.removeEventListener('resize', handleResize);
    });
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-4 md:p-6">
    <div class="max-w-7xl mx-auto">
        <div class="mb-6 flex items-center justify-between">
            <div>
                <button
                    on:click={goBack}
                    class="mb-3 inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-white text-slate-700 hover:bg-slate-50 transition-all duration-200 border border-slate-200 shadow-sm hover:shadow"
                >
                    <i class="fa-solid fa-arrow-left"></i>
                    Kembali
                </button>
                <h1 class="text-2xl md:text-3xl font-bold text-slate-800">Pressing Line Dashboard</h1>
                <p class="text-slate-600 mt-1 text-sm md:text-base">Pantau status semua mesin press dengan diagram real-time</p>
            </div>
        </div>

        {#if isLoading}
            <div class="flex justify-center items-center min-h-64 bg-white rounded-2xl shadow-lg border border-slate-200">
                <div class="text-center">
                    <div class="inline-block animate-spin">
                        <i class="fa-solid fa-spinner text-indigo-600 text-4xl"></i>
                    </div>
                    <p class="text-slate-600 mt-4">Memuat data mesin...</p>
                </div>
            </div>
        {:else}
            <div class="bg-white rounded-2xl shadow-lg p-4 md:p-6 border border-slate-200">
                <h2 class="text-lg font-bold text-slate-800 mb-6">Layout Pabrik Pressing</h2>

                <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <div>
                        <h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider flex items-center gap-2">
                            <span class="w-2 h-2 bg-purple-500 rounded-full"></span>
                            Inject Machines
                        </h3>
                        <div class="space-y-4">
                            {#each injectMachines as machine}
                                {@const status = getMachineStatusColor(machine)}
                                <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
                                    <button
                                        on:click={() => selectMachine(machine)}
                                        class="w-full flex items-center gap-3 p-4 hover:bg-slate-50/50 transition-colors"
                                    >
                                        <div class="relative">
                                            <div class="w-10 h-10 flex items-center justify-center rounded-full {status.badge} text-white font-bold text-sm flex-shrink-0 shadow">
                                                {machine.id}
                                            </div>
                                            <div class="absolute -top-1 -right-1">
                                                {status.indicator}
                                            </div>
                                        </div>
                                        <div class="flex-1 text-left min-w-0">
                                            <div class="font-semibold text-slate-800 truncate">{machine.name}</div>
                                            <div class="text-xs text-slate-600 flex items-center gap-1">
                                                <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
                                                {machine.completed} / {machine.target} pcs
                                            </div>
                                        </div>
                                        <i class="fa-solid fa-arrow-right text-slate-400"></i>
                                    </button>

                                    <div class="px-4 pb-4 pt-2 border-t border-slate-200/50 space-y-4">
                                        <div>
                                            <div class="flex justify-between items-center mb-2">
                                                <span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
                                                    <i class="fa-solid fa-chart-bar text-blue-500"></i>
                                                    Total Produksi
                                                </span>
                                                <span class="text-sm font-bold {machine.completed >= machine.target ? 'text-green-600' : 'text-orange-600'}">
                                                    {machine.completed}
                                                    {#if machine.completed >= machine.target}
                                                        <i class="fa-solid fa-check ml-1 text-xs"></i>
                                                    {:else}
                                                        <i class="fa-solid fa-exclamation ml-1 text-xs"></i>
                                                    {/if}
                                                </span>
                                            </div>
                                            <div class="h-32">
                                                <canvas id="chart-total-{machine.id}"></canvas>
                                            </div>
                                        </div>
                                        
                                        <div>
                                            <div class="flex justify-between items-center mb-2">
                                                <span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
                                                    <i class="fa-solid fa-triangle-exclamation text-red-500"></i>
                                                    Not Good
                                                </span>
                                                <span class="text-sm font-bold text-red-600">
                                                    {machine.notGood}
                                                </span>
                                            </div>
                                            <div class="h-32">
                                                <canvas id="chart-notgood-{machine.id}"></canvas>
                                            </div>
                                        </div>
                                        
                                        <div class="grid grid-cols-2 gap-3 pt-3 border-t border-slate-200/50">
                                            <div class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-lg p-3 text-center border border-blue-200">
                                                <p class="text-xs text-blue-700 font-semibold mb-1">PRODUKSI</p>
                                                <div class="flex items-center justify-center gap-1">
                                                    <p class="text-xl font-bold text-blue-800">{machine.completed}</p>
                                                </div>
                                                <p class="text-[10px] text-blue-600 mt-1">Target: {machine.target}</p>
                                            </div>
                                            <div class="bg-gradient-to-r from-red-50 to-red-100 rounded-lg p-3 text-center border border-red-200">
                                                <p class="text-xs text-red-700 font-semibold mb-1">NOT GOOD</p>
                                                <div class="flex items-center justify-center gap-1">
                                                    <p class="text-xl font-bold text-red-800">{machine.notGood}</p>
                                                </div>
                                                <p class="text-[10px] text-red-600 mt-1">Total NG</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>

                    <div>
                        <h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider flex items-center gap-2">
                            <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
                            Press Machines
                        </h3>
                        <div class="space-y-4">
                            {#each pressMachines as machine}
                                {@const status = getMachineStatusColor(machine)}
                                <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
                                    <button
                                        on:click={() => selectMachine(machine)}
                                        class="w-full flex items-center gap-3 p-4 hover:bg-slate-50/50 transition-colors"
                                    >
                                        <div class="relative">
                                            <div class="w-10 h-10 flex items-center justify-center rounded-full {status.badge} text-white font-bold text-sm flex-shrink-0 shadow">
                                                {machine.id}
                                            </div>
                                            <div class="absolute -top-1 -right-1">
                                                {status.indicator}
                                            </div>
                                        </div>
                                        <div class="flex-1 text-left min-w-0">
                                            <div class="font-semibold text-slate-800 truncate">{machine.name}</div>
                                            <div class="text-xs text-slate-600 flex items-center gap-1">
                                                <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
                                                {machine.completed} / {machine.target} pcs
                                            </div>
                                        </div>
                                        <i class="fa-solid fa-arrow-right text-slate-400"></i>
                                    </button>

                                    <div class="px-4 pb-4 pt-2 border-t border-slate-200/50 space-y-4">
                                        <div>
                                            <div class="flex justify-between items-center mb-2">
                                                <span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
                                                    <i class="fa-solid fa-chart-bar text-blue-500"></i>
                                                    Total Produksi
                                                </span>
                                                <span class="text-sm font-bold {machine.completed >= machine.target ? 'text-green-600' : 'text-orange-600'}">
                                                    {machine.completed}
                                                    {#if machine.completed >= machine.target}
                                                        <i class="fa-solid fa-check ml-1 text-xs"></i>
                                                    {:else}
                                                        <i class="fa-solid fa-exclamation ml-1 text-xs"></i>
                                                    {/if}
                                                </span>
                                            </div>
                                            <div class="h-32">
                                                <canvas id="chart-total-{machine.id}"></canvas>
                                            </div>
                                        </div>
                                        
                                        <div>
                                            <div class="flex justify-between items-center mb-2">
                                                <span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
                                                    <i class="fa-solid fa-triangle-exclamation text-red-500"></i>
                                                    Not Good
                                                </span>
                                                <span class="text-sm font-bold text-red-600">
                                                    {machine.notGood}
                                                </span>
                                            </div>
                                            <div class="h-32">
                                                <canvas id="chart-notgood-{machine.id}"></canvas>
                                            </div>
                                        </div>
                                        
                                        <div class="grid grid-cols-2 gap-3 pt-3 border-t border-slate-200/50">
                                            <div class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-lg p-3 text-center border border-blue-200">
                                                <p class="text-xs text-blue-700 font-semibold mb-1">PRODUKSI</p>
                                                <div class="flex items-center justify-center gap-1">
                                                    <p class="text-xl font-bold text-blue-800">{machine.completed}</p>
                                                </div>
                                                <p class="text-[10px] text-blue-600 mt-1">Target: {machine.target}</p>
                                            </div>
                                            <div class="bg-gradient-to-r from-red-50 to-red-100 rounded-lg p-3 text-center border border-red-200">
                                                <p class="text-xs text-red-700 font-semibold mb-1">NOT GOOD</p>
                                                <div class="flex items-center justify-center gap-1">
                                                    <p class="text-xl font-bold text-red-800">{machine.notGood}</p>
                                                </div>
                                                <p class="text-[10px] text-red-600 mt-1">Total NG</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>
            </div>

            <div class="mt-6 bg-white rounded-xl p-4 border border-slate-200 shadow-sm">
                <h3 class="text-sm font-semibold text-slate-800 mb-3">Keterangan Chart</h3>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div class="flex items-center gap-2">
                        <div class="w-4 h-4 bg-gradient-to-b from-blue-500 to-blue-600 rounded"></div>
                        <span class="text-sm text-slate-600">Total Produksi (Pcs)</span>
                    </div>
                    <div class="flex items-center gap-2">
                        <div class="w-4 h-4 bg-gradient-to-b from-red-500 to-red-600 rounded"></div>
                        <span class="text-sm text-slate-600">Not Good / NG (Pcs)</span>
                    </div>  
                    <div class="flex items-center gap-2">
                        <div class="w-4 h-4 border-2 border-green-500 border-dashed"></div>
                        <span class="text-sm text-slate-600">Garis Target Produksi</span>
                    </div>
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    :global(body) {
        background-color: #f8fafc;
        font-family: 'Inter', sans-serif;
    }

    canvas {
        width: 100% !important;
        height: 100% !important;
    }
    
    ::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
    
    ::-webkit-scrollbar-track {
        background: #f1f5f9;
        border-radius: 4px;
    }
    
    ::-webkit-scrollbar-thumb {
        background: #cbd5e1;
        border-radius: 4px;
    }
    
    ::-webkit-scrollbar-thumb:hover {
        background: #94a3b8;
    }
</style>