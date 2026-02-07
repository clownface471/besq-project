<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import Chart from 'chart.js/auto';
    import { auth } from '$lib/stores/auth'; // IMPORT AUTH
    import { goto } from '$app/navigation';

    let chartCanvas: HTMLCanvasElement;
    let chartInstance: Chart;
    
    $: noMC = $page.url.searchParams.get('no_mc') || '';
    $: selectedDate = $page.url.searchParams.get('tanggal') || new Date().toISOString().split('T')[0];
    let itemsProduced = "";
    
    const API_URL = 'http://localhost:8080';

    async function fetchData() {
        if (!noMC) return;
        try {
            // TAMBAHKAN HEADER AUTH
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
            
            const validItem = data.find((d: any) => d.extra_info && d.extra_info !== '- (-)');
            if (validItem) itemsProduced = validItem.extra_info;

            updateChart(data);
        } catch (error) {
            console.error("Error fetching machine detail:", error);
        }
    }

    function updateChart(data: any[]) {
        if (chartInstance) chartInstance.destroy();

        chartInstance = new Chart(chartCanvas, {
            type: 'bar',
            data: {
                labels: data.map(d => d.label), // Jam (07:00, 08:00...)
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
    <button on:click={() => history.back()} class="mb-4 text-blue-600 hover:underline">← Kembali</button>

    <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="flex justify-between items-start mb-6">
            <div>
                <h1 class="text-2xl font-bold">Detail Mesin: {noMC}</h1>
                <p class="text-gray-600">Tanggal: {selectedDate}</p>
            </div>
            {#if itemsProduced}
                <div class="bg-blue-50 p-3 rounded border border-blue-200 text-sm">
                    <span class="font-bold text-blue-800">Produksi Aktif:</span>
                    <p>{itemsProduced}</p>
                </div>
            {/if}
        </div>

        <div class="relative h-[500px] w-full">
            <canvas bind:this={chartCanvas}></canvas>
        </div>
    </div>
</div>