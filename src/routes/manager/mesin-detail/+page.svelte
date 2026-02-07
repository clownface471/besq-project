<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import Chart from 'chart.js/auto';
    import { auth } from '$lib/stores/auth';
    import { goto } from '$app/navigation';

    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart;
    
    $: noMC = $page.url.searchParams.get('no_mc') || '';
    
    // Ambil tanggal dari URL (hasil operan dari halaman sebelumnya) atau default hari ini
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
            
            // Reset info item
            itemsProduced = "-";
            const validItem = data.find((d: any) => d.extra_info && d.extra_info !== '- (-)');
            if (validItem) itemsProduced = validItem.extra_info;

            updateChart(data);
        } catch (error) {
            console.error("Error fetching machine detail:", error);
        } finally {
            isLoading = false;
        }
    }

    function updateChart(data: any[]) {
        if (chartInstance) chartInstance.destroy();

        chartInstance = new Chart(chartCanvas, {
            type: 'bar',
            data: {
                labels: data.map(d => d.label),
                datasets: [
                    {
                        label: 'Target Speed/Jam',
                        data: data.map(d => d.target),
                        type: 'line',
                        borderColor: 'red',
                        borderDash: [5, 5], 
                        borderWidth: 2,
                        pointRadius: 0,
                        order: 0
                    },
                    {
                        label: 'Total Output',
                        data: data.map(d => d.actual),
                        backgroundColor: 'rgba(54, 162, 235, 0.8)',
                        order: 1
                    },
                    {
                        label: 'NG (Reject)',
                        data: data.map(d => d.actual_ng),
                        backgroundColor: 'rgba(255, 99, 132, 0.8)',
                        order: 2
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false, // Agar grafik menyesuaikan container
                plugins: {
                    tooltip: {
                        callbacks: {
                            afterLabel: function(context) {
                                const index = context.dataIndex;
                                return 'Info: ' + (data[index].extra_info || '-');
                            }
                        }
                    }
                },
                scales: {
                    x: { title: { display: true, text: 'Jam Produksi' } },
                    y: { title: { display: true, text: 'Qty Pcs' }, beginAtZero: true }
                }
            }
        });
    }

    onMount(() => {
        if ($auth.token) {
            fetchData();
        } else {
            goto('/');
        }
    });
</script>

<div class="p-6 bg-gray-50 min-h-screen">
    <div class="flex justify-between items-center mb-4">
        <button on:click={() => history.back()} class="text-blue-600 hover:underline flex items-center gap-2">
            <i class="fa-solid fa-arrow-left"></i> Kembali
        </button>

        <div class="flex items-center gap-2">
            <label for="date-detail" class="text-sm font-semibold text-gray-600">Tanggal:</label>
            <input 
                type="date" 
                id="date-detail"
                class="border border-gray-300 p-2 rounded-lg shadow-sm focus:ring focus:ring-blue-200"
                bind:value={selectedDate}
                on:change={fetchData}
            />
        </div>
    </div>

    <div class="bg-white p-6 rounded-lg shadow-md h-full">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
            <div>
                <h1 class="text-2xl font-bold text-gray-800">Detail Mesin: <span class="text-blue-600">{noMC}</span></h1>
                <p class="text-gray-500 text-sm mt-1">Laporan Produksi Per Jam</p>
            </div>
            
            <div class="bg-blue-50 px-4 py-3 rounded-lg border border-blue-100 flex items-center gap-3">
                <div class="bg-blue-500 text-white rounded-full w-8 h-8 flex items-center justify-center">
                    <i class="fa-solid fa-box-open"></i>
                </div>
                <div>
                    <span class="block text-xs font-bold text-blue-800 uppercase tracking-wider">Sedang Diproduksi</span>
                    <p class="text-gray-900 font-medium truncate max-w-[200px]" title={itemsProduced}>{itemsProduced || "-"}</p>
                </div>
            </div>
        </div>

        <div class="relative w-full h-[500px] border rounded-lg p-2">
            {#if isLoading}
                <div class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10 rounded-lg">
                    <div class="text-center">
                        <i class="fa-solid fa-spinner animate-spin text-3xl text-blue-600"></i>
                        <p class="mt-2 text-gray-600 font-semibold">Mengambil data mesin...</p>
                    </div>
                </div>
            {/if}
            <canvas bind:this={chartCanvas}></canvas>
        </div>
    </div>
</div>