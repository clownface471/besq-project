<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import Chart from 'chart.js/auto';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/stores/auth';

    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart;
    let isLoading = false;
    
    // Default ambil dari URL, atau hari ini
    $: processName = $page.url.searchParams.get('process') || 'PRESSING';
    
    // Kita gunakan let biasa agar bisa diubah oleh input date
    let selectedDate = $page.url.searchParams.get('tanggal') || new Date().toISOString().split('T')[0];

    const API_URL = 'http://localhost:8080';

    async function fetchData() {
        isLoading = true;
        try {
            let dbProcessCode = processName === 'PRESSING' ? 'PRS' : (processName === 'CUTTING' ? 'CUT' : processName);

            // Gunakan selectedDate variable yang reaktif
            const res = await fetch(`${API_URL}/api/chart/process?tanggal=${selectedDate}&proses=${dbProcessCode}`, {
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
            updateChart(data);
        } catch (error) {
            console.error("Error fetching process data:", error);
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
                        label: 'Target',
                        data: data.map(d => d.target),
                        type: 'line',
                        borderColor: 'black',
                        borderWidth: 2,
                        pointRadius: 0,
                        order: 0
                    },
                    {
                        label: 'Actual Output',
                        data: data.map(d => d.actual),
                        backgroundColor: 'rgba(54, 162, 235, 0.7)',
                        order: 1
                    }
                ]
            },
            options: {
                responsive: true,
                onClick: (e, elements) => {
                    if (elements.length > 0) {
                        const index = elements[0].index;
                        const machineNo = data[index].label;
                        
                        // FITUR UTAMA: Mengirim tanggal yang sedang dipilih ke halaman detail
                        goto(`/manager/mesin-detail?no_mc=${machineNo}&tanggal=${selectedDate}`);
                    }
                },
                scales: {
                     y: { beginAtZero: true }
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
            <label for="date" class="text-sm font-semibold text-gray-600">Pilih Tanggal:</label>
            <input 
                type="date" 
                id="date"
                class="border border-gray-300 p-2 rounded-lg shadow-sm focus:ring focus:ring-blue-200"
                bind:value={selectedDate}
                on:change={fetchData}
            />
        </div>
    </div>
    
    <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="flex justify-between items-center mb-6">
            <div>
                <h1 class="text-2xl font-bold">Process Dashboard: {processName}</h1>
                <p class="text-gray-600">Menampilkan data tanggal: <span class="font-semibold text-blue-600">{selectedDate}</span></p>
            </div>
            <button on:click={fetchData} class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition">
                Refresh Data
            </button>
        </div>

        <div class="relative h-96 w-full">
            {#if isLoading}
                <div class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
                    <div class="text-center">
                        <i class="fa-solid fa-spinner animate-spin text-3xl text-blue-600"></i>
                        <p class="mt-2 text-gray-600 font-semibold">Memuat data...</p>
                    </div>
                </div>
            {/if}
            <canvas bind:this={chartCanvas}></canvas>
        </div>
        <p class="text-sm text-gray-500 mt-4 text-center border-t pt-2">*Klik pada batang mesin untuk melihat detail per jam pada tanggal ini.</p>
    </div>
</div>