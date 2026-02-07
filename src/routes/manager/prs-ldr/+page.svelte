<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';
    import Chart from 'chart.js/auto';
    import annotationPlugin from 'chartjs-plugin-annotation';
    import { auth } from '$lib/stores/auth';

    Chart.register(annotationPlugin);

    let machineData: any[] = [];
    let isLoading = false;
    const API_URL = ''; // Diisi URL backend
    
    let charts: Record<string, any> = {};

    function goBack() {
        goto('/manager');
    }

    function selectMachine(machine: any) {
        // Asumsi machine.id adalah "07B", "16A", dst.
        goto(`/manager/mesin-detail?no_mc=${machine.id}`);
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
            const res = await fetch(`${API_URL}/api/chart/process?tanggal=${today}&proses=PRS`, {
                headers: { Authorization: `Bearer ${$auth.token}` }
            });

            if (res.ok) {
                const data = await res.json();
                
                machineData = data.map((item: any) => {
                    const target = item.target || 1;
                    const actual = item.actual || 0;
                    const percentage = Math.round((actual / target) * 100);
                    const ngPercentage = 0; // Data NG belum ada di endpoint ini

                    return {
                        id: item.label,
                        name: `Mesin Press ${item.label}`,
                        completed: percentage,
                        notGood: ngPercentage,
                        isProblem: percentage < 50 || ngPercentage > 10
                    };
                });
            }
        } catch (e) {
            console.error("Gagal ambil data mesin:", e);
        } finally {
            isLoading = false;
            // Beri jeda sedikit agar DOM render loop selesai sebelum inisialisasi chart
            setTimeout(initializeCharts, 100);
        }
    }

    function initializeCharts() {
        // Hancurkan chart sebelumnya
        Object.keys(charts).forEach(machineId => {
            if (charts[machineId]?.total) charts[machineId].total.destroy();
            if (charts[machineId]?.notGood) charts[machineId].notGood.destroy();
        });
        charts = {};

        const productionTarget = 90;
        const notGoodTarget = 10;

        machineData.forEach((machine) => {
            // Chart Total
            const totalCanvasId = `chart-total-${machine.id}`;
            const totalCanvas = document.getElementById(totalCanvasId) as HTMLCanvasElement;
            
            if (totalCanvas) {
                const ctx = totalCanvas.getContext('2d');
                if (ctx) {
                    const gradient = ctx.createLinearGradient(0, 0, 0, totalCanvas.height);
                    gradient.addColorStop(0, '#3b82f6');
                    gradient.addColorStop(1, '#1d4ed8');

                    charts[machine.id] = charts[machine.id] || {};
                    charts[machine.id].total = new Chart(totalCanvas, {
                        type: 'bar',
                        data: {
                            labels: ['Produksi'],
                            datasets: [{
                                label: 'Produksi',
                                data: [machine.completed],
                                backgroundColor: gradient,
                                borderColor: '#1e40af',
                                borderWidth: 1,
                                borderRadius: 6
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: {
                                legend: { display: false },
                                annotation: {
                                    annotations: {
                                        targetLine: {
                                            type: 'line',
                                            yMin: productionTarget,
                                            yMax: productionTarget,
                                            borderColor: '#10b981',
                                            borderWidth: 2,
                                            borderDash: [5, 5],
                                            label: { content: `Target: ${productionTarget}%`, position: 'end', backgroundColor: '#10b981', color: 'white', font: { size: 9 } }
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: { beginAtZero: true, max: 100, display: false },
                                x: { display: false }
                            }
                        }
                    });
                }
            }

            // Chart NG
            const notGoodCanvasId = `chart-notgood-${machine.id}`;
            const notGoodCanvas = document.getElementById(notGoodCanvasId) as HTMLCanvasElement;
            
            if (notGoodCanvas) {
                const ctx = notGoodCanvas.getContext('2d');
                if (ctx) {
                    const gradient = ctx.createLinearGradient(0, 0, 0, notGoodCanvas.height);
                    gradient.addColorStop(0, '#ef4444');
                    gradient.addColorStop(1, '#dc2626');

                    charts[machine.id] = charts[machine.id] || {};
                    charts[machine.id].notGood = new Chart(notGoodCanvas, {
                        type: 'bar',
                        data: {
                            labels: ['NG'],
                            datasets: [{
                                label: 'NG',
                                data: [machine.notGood],
                                backgroundColor: gradient,
                                borderColor: '#b91c1c',
                                borderWidth: 1,
                                borderRadius: 6
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: {
                                legend: { display: false },
                                annotation: {
                                    annotations: {
                                        targetLine: {
                                            type: 'line',
                                            yMin: notGoodTarget,
                                            yMax: notGoodTarget,
                                            borderColor: '#f59e0b',
                                            borderWidth: 2,
                                            borderDash: [5, 5],
                                            label: { content: `Target: ≤${notGoodTarget}%`, position: 'end', backgroundColor: '#f59e0b', color: 'white', font: { size: 9 } }
                                        }
                                    }
                                }
                            },
                            scales: {
                                y: { beginAtZero: true, max: 40, display: false },
                                x: { display: false }
                            }
                        }
                    });
                }
            }
        });
    }

    // Reactive variables (Filter jika ada data)
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
                <button on:click={goBack} class="mb-3 inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-white text-slate-700 hover:bg-slate-50 transition-all duration-200 border border-slate-200 shadow-sm hover:shadow">
                    <i class="fa-solid fa-arrow-left"></i> Kembali
                </button>
                <h1 class="text-2xl md:text-3xl font-bold text-slate-800">Pressing Line Dashboard</h1>
                <p class="text-slate-600 mt-1 text-sm md:text-base">Pantau status semua mesin press dengan diagram real-time</p>
            </div>
        </div>

        <div class="bg-white rounded-2xl shadow-lg p-4 md:p-6 border border-slate-200">
            <h2 class="text-lg font-bold text-slate-800 mb-6">Layout Pabrik Pressing</h2>
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div>
                    <h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider flex items-center gap-2">
                        <span class="w-2 h-2 bg-purple-500 rounded-full"></span> Inject Machines
                    </h3>
                    <div class="space-y-4">
                        {#each injectMachines as machine}
                            {@const status = getMachineStatusColor(machine)}
                            <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
                                <button on:click={() => selectMachine(machine)} class="w-full flex items-center gap-3 p-4 hover:bg-slate-50/50 transition-colors">
                                    <div class="relative">
                                        <div class="w-10 h-10 flex items-center justify-center rounded-full {status.badge} text-white font-bold text-sm flex-shrink-0 shadow">
                                            {machine.id}
                                        </div>
                                        <div class="absolute -top-1 -right-1">{status.indicator}</div>
                                    </div>
                                    <div class="flex-1 text-left min-w-0">
                                        <div class="font-semibold text-slate-800 truncate">{machine.name}</div>
                                        <div class="text-xs text-slate-600 flex items-center gap-1"><span class="w-2 h-2 bg-blue-500 rounded-full"></span> {machine.completed}% selesai</div>
                                    </div>
                                    <i class="fa-solid fa-arrow-right text-slate-400"></i>
                                </button>
                                <div class="px-4 pb-4 pt-2 border-t border-slate-200/50 space-y-4">
                                    <div class="h-32"><canvas id="chart-total-{machine.id}"></canvas></div>
                                    <div class="h-32"><canvas id="chart-notgood-{machine.id}"></canvas></div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>

                <div>
                    <h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider flex items-center gap-2">
                        <span class="w-2 h-2 bg-blue-500 rounded-full"></span> Press Machines
                    </h3>
                    <div class="space-y-4">
                        {#each pressMachines as machine}
                             {@const status = getMachineStatusColor(machine)}
                            <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
                                <button on:click={() => selectMachine(machine)} class="w-full flex items-center gap-3 p-4 hover:bg-slate-50/50 transition-colors">
                                    <div class="relative">
                                        <div class="w-10 h-10 flex items-center justify-center rounded-full {status.badge} text-white font-bold text-sm flex-shrink-0 shadow">
                                            {machine.id}
                                        </div>
                                        <div class="absolute -top-1 -right-1">{status.indicator}</div>
                                    </div>
                                    <div class="flex-1 text-left min-w-0">
                                        <div class="font-semibold text-slate-800 truncate">{machine.name}</div>
                                        <div class="text-xs text-slate-600 flex items-center gap-1"><span class="w-2 h-2 bg-blue-500 rounded-full"></span> {machine.completed}% selesai</div>
                                    </div>
                                    <i class="fa-solid fa-arrow-right text-slate-400"></i>
                                </button>
                                <div class="px-4 pb-4 pt-2 border-t border-slate-200/50 space-y-4">
                                    <div class="h-32"><canvas id="chart-total-{machine.id}"></canvas></div>
                                    <div class="h-32"><canvas id="chart-notgood-{machine.id}"></canvas></div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
        </div>
        </div>
</div>
<style>
    :global(body) { background-color: #f8fafc; font-family: 'Inter', sans-serif; }
    canvas { width: 100% !important; height: 100% !important; }
</style>