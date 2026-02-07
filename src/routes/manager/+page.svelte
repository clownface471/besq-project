<script lang="ts">
    import { onMount } from 'svelte';
    import Chart from 'chart.js/auto';
    import { goto } from '$app/navigation';
    // 1. IMPORT STORE AUTH
    import { auth } from '$lib/stores/auth'; 

    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart;
    let selectedDate = new Date().toISOString().split('T')[0];
    let isLoading = false;
    
    // Pastikan URL sesuai dengan backend Anda (port 8080)
    const API_URL = 'http://localhost:8080';

    async function fetchData() {
        isLoading = true;
        try {
            // 2. TAMBAHKAN HEADER AUTHORIZATION DISINI
            const res = await fetch(`${API_URL}/api/chart/manager?tanggal=${selectedDate}`, {
                headers: {
                    'Authorization': `Bearer ${$auth.token}`, // Ambil token dari store
                    'Content-Type': 'application/json'
                }
            });

            if (res.status === 401) {
                alert("Sesi habis atau tidak ada izin. Silakan login ulang.");
                goto('/'); // Redirect ke login jika 401
                return;
            }

            const data = await res.json();
            updateChart(data);
        } catch (error) {
            console.error("Gagal ambil data:", error);
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
                        backgroundColor: 'rgba(200, 200, 200, 0.5)',
                        borderColor: 'gray',
                        borderWidth: 1
                    },
                    {
                        label: 'Actual',
                        data: data.map(d => d.actual),
                        backgroundColor: data.map(d => d.actual >= d.target ? 'rgba(75, 192, 192, 0.7)' : 'rgba(255, 99, 132, 0.7)'),
                        borderColor: data.map(d => d.actual >= d.target ? 'rgb(75, 192, 192)' : 'rgb(255, 99, 132)'),
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                onClick: (e, elements) => {
                    if (elements.length > 0) {
                        const index = elements[0].index;
                        const processName = data[index].label; 
                        goto(`/manager/prs-ldr?process=${processName}&tanggal=${selectedDate}`);
                    }
                },
                scales: {
                    y: { beginAtZero: true }
                }
            }
        });
    }

    onMount(() => {
        // Cek apakah user punya token sebelum fetch
        if ($auth.token) {
            fetchData();
        } else {
            goto('/');
        }
    });
</script>

<div class="p-6 bg-gray-50 min-h-screen">
    <div class="bg-white p-6 rounded-lg shadow-md">
        <h1 class="text-2xl font-bold mb-4">Manager Dashboard (Overview)</h1>
        
        <div class="flex gap-4 mb-6">
            <input 
                type="date" 
                class="border p-2 rounded" 
                bind:value={selectedDate} 
                on:change={fetchData}
            />
            <button on:click={fetchData} class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                Refresh
            </button>
        </div>

        <div class="relative h-96 w-full">
            {#if isLoading}
                <div class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
                    <p class="text-lg font-semibold text-gray-600">Loading data...</p>
                </div>
            {/if}
            <canvas bind:this={chartCanvas}></canvas>
        </div>
        <p class="text-sm text-gray-500 mt-2">*Klik pada batang grafik untuk melihat detail proses.</p>
    </div>
</div>