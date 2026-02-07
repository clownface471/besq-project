<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import Chart from 'chart.js/auto';
    import { goto } from '$app/navigation';
    // 1. IMPORT STORE AUTH
    import { auth } from '$lib/stores/auth';

    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart;
    let isLoading = false;
    
    // Ambil parameter dari URL
    $: processName = $page.url.searchParams.get('process') || 'PRESSING';
    $: selectedDate = $page.url.searchParams.get('tanggal') || new Date().toISOString().split('T')[0];

    // URL Backend
    const API_URL = 'http://localhost:8080';

    async function fetchData() {
        isLoading = true;
        try {
            // Mapping nama proses dari frontend ke kode database
            // Misal: "PRESSING" -> "PRS"
            let dbProcessCode = processName === 'PRESSING' ? 'PRS' : (processName === 'CUTTING' ? 'CUT' : processName);

            // 2. TAMBAHKAN HEADER AUTHORIZATION
            const res = await fetch(`${API_URL}/api/chart/process?tanggal=${selectedDate}&proses=${dbProcessCode}`, {
                headers: {
                    'Authorization': `Bearer ${$auth.token}`,
                    'Content-Type': 'application/json'
                }
            });

            // 3. HANDLING ERROR 401
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
                labels: data.map(d => d.label), // No Mesin (e.g., "04A", "05B")
                datasets: [
                    {
                        label: 'Target',
                        data: data.map(d => d.target),
                        type: 'line', // Target enak dilihat sebagai garis batas
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
                        // Drill Down: Pindah ke detail mesin per jam
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
        // Cek Token
        if ($auth.token) {
            fetchData();
        } else {
            goto('/');
        }
    });
</script>

<div class="p-6 bg-gray-50 min-h-screen">
    <button on:click={() => history.back()} class="mb-4 text-blue-600 hover:underline">← Kembali</button>
    
    <div class="bg-white p-6 rounded-lg shadow-md">
        <h1 class="text-2xl font-bold mb-2">Process Dashboard: {processName}</h1>
        <p class="text-gray-600 mb-6">Tanggal: {selectedDate}</p>

        <div class="relative h-96 w-full">
            {#if isLoading}
                <div class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
                    <p class="text-lg font-semibold text-gray-600">Loading data...</p>
                </div>
            {/if}
            <canvas bind:this={chartCanvas}></canvas>
        </div>
        <p class="text-sm text-gray-500 mt-2">*Klik pada batang mesin untuk melihat detail per jam.</p>
    </div>
</div>