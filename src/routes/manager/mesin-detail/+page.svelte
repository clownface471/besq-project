<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';
    import Chart from 'chart.js/auto';
    import { auth } from '$lib/stores/auth';
    import { goto } from '$app/navigation';

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
    
    let itemsProduced = "";
    let isLoading = false;
    
    const API_URL = 'http://localhost:8080';

    async function fetchData() {
        if (!noMC) return;
        isLoading = true;
        try {
            const res = await fetch(`${API_URL}/api/chart/machine?tanggal=${selectedDate}&no_mc=${noMC}`, {
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
            const validItem = data.find((d: any) => d.extra_info && d.extra_info !== '- (-)');
            if (validItem) itemsProduced = validItem.extra_info;

            renderCharts(data);
        } catch (error) {
            console.error("Error fetching machine detail:", error);
        } finally {
            isLoading = false;
        }
    }

    function renderCharts(data: any[]) {
        // Hapus chart lama jika ada sebelum render baru (untuk refresh/ganti tanggal)
        if (chartTotal) chartTotal.destroy();
        if (chartOK) chartOK.destroy();
        if (chartNG) chartNG.destroy();

        const labels = data.map(d => d.label); // Jam (07:00, 08:00...)
        const targets = data.map(d => d.target); // Target dari DB

        // --- 1. GRAFIK TOTAL OUTPUT ---
        chartTotal = new Chart(canvasTotal, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Target Speed (DB)',
                        data: targets,
                        type: 'line',
                        borderColor: '#2563eb', // Biru Tua
                        borderDash: [5, 5], 
                        borderWidth: 2,
                        pointRadius: 0,
                        tension: 0.1
                    },
                    {
                        label: 'Total Output',
                        data: data.map(d => d.actual),
                        backgroundColor: '#60a5fa', // Biru Muda
                        barPercentage: 0.6
                    }
                ]
            },
            options: getChartOptions('Total Produksi vs Target')
        });

        // --- 2. GRAFIK OK OUTPUT ---
        chartOK = new Chart(canvasOK, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Minimal OK (Target)',
                        data: targets, // Asumsi Target DB adalah target barang OK
                        type: 'line',
                        borderColor: '#16a34a', // Hijau Tua
                        borderDash: [5, 5],
                        borderWidth: 2,
                        pointRadius: 0,
                        tension: 0.1
                    },
                    {
                        label: 'Actual OK',
                        data: data.map(d => d.actual_ok),
                        backgroundColor: '#4ade80', // Hijau Muda
                        barPercentage: 0.6
                    }
                ]
            },
            options: getChartOptions('Barang OK vs Target')
        });

        // --- 3. GRAFIK NG (REJECT) ---
        chartNG = new Chart(canvasNG, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Batas Max NG (Visual)',
                        // Karena belum ada kolom max_ng di DB, kita buat garis visual statis (misal 5 pcs)
                        // Nanti bisa diganti: data.map(d => d.max_ng_limit)
                        data: Array(data.length).fill(5), 
                        type: 'line',
                        borderColor: '#dc2626', // Merah Tua
                        borderWidth: 2,
                        pointRadius: 0
                    },
                    {
                        label: 'Actual NG',
                        data: data.map(d => d.actual_ng),
                        backgroundColor: '#f87171', // Merah Muda
                        barPercentage: 0.6
                    }
                ]
            },
            options: getChartOptions('Jumlah Reject (NG)')
        });
    }

    // Helper untuk opsi grafik yang seragam
    function getChartOptions(title: string): any {
        return {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { position: 'top' },
                title: { display: true, text: title, font: { size: 14 } }
            },
            scales: {
                y: { beginAtZero: true, grid: { color: '#f3f4f6' } },
                x: { grid: { display: false } }
            }
        };
    }

    onMount(() => {
        if ($auth.token) {
            fetchData();
        } else {
            goto('/');
        }
    });

    onDestroy(() => {
        if (chartTotal) chartTotal.destroy();
        if (chartOK) chartOK.destroy();
        if (chartNG) chartNG.destroy();
    });
</script>

<div class="p-6 bg-gray-50 min-h-screen">
    <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
        <div class="flex items-center gap-4 w-full md:w-auto">
            <button on:click={() => history.back()} class="text-blue-600 hover:text-blue-800 flex items-center gap-2 font-medium">
                <i class="fa-solid fa-arrow-left"></i> Kembali
            </button>
            <div class="h-6 w-px bg-gray-300"></div>
            <div class="flex items-center gap-2">
                <label for="date-detail" class="text-sm font-semibold text-gray-600">Tanggal:</label>
                <input 
                    type="date" 
                    id="date-detail"
                    class="border border-gray-300 p-2 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-200 text-sm"
                    bind:value={selectedDate}
                    on:change={fetchData}
                />
            </div>
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