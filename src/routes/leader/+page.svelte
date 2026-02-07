<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import Chart from 'chart.js/auto';
	import annotationPlugin from 'chartjs-plugin-annotation';
	
	// Register annotation plugin
	Chart.register(annotationPlugin);

	let machineData = [
		{ id: '07B', name: 'Mesin Press 07B', completed: 85, notGood: 15, isProblem: false },
		{ id: '16B', name: 'Mesin Press 16B', completed: 72, notGood: 28, isProblem: false },
		{ id: '16A', name: 'Mesin Press 16A', completed: 90, notGood: 10, isProblem: false },
		{ id: '12B', name: 'Mesin Press 12B', completed: 78, notGood: 22, isProblem: false },
		{ id: '12A', name: 'Mesin Press 12A', completed: 88, notGood: 12, isProblem: false },
		{ id: '14B', name: 'Mesin Press 14B', completed: 65, notGood: 35, isProblem: true },
		{ id: '14A', name: 'Mesin Press 14A', completed: 81, notGood: 19, isProblem: false },
		{ id: '09A', name: 'Mesin Press 09A', completed: 92, notGood: 8, isProblem: false },
		{ id: '09B', name: 'Mesin Press 09B', completed: 76, notGood: 24, isProblem: false },
		{ id: '11A', name: 'Mesin Press 11A', completed: 84, notGood: 16, isProblem: false },
		{ id: '11B', name: 'Mesin Press 11B', completed: 79, notGood: 21, isProblem: false },
		{ id: '15A', name: 'Mesin Press 15A', completed: 87, notGood: 13, isProblem: true },
		{ id: '15B', name: 'Mesin Press 15B', completed: 73, notGood: 27, isProblem: false },
		{ id: '02A', name: 'Mesin Press 02A', completed: 94, notGood: 6, isProblem: false },
		{ id: '04A', name: 'Mesin Press 04A', completed: 80, notGood: 20, isProblem: false },
		{ id: '04B', name: 'Mesin Press 04B', completed: 77, notGood: 23, isProblem: false },
		{ id: '20A', name: 'Mesin Press (Inject) 20A', completed: 86, notGood: 14, isProblem: false },
		{ id: '21A', name: 'Mesin Press (Inject) 21A', completed: 75, notGood: 25, isProblem: true },
		{ id: '22A', name: 'Mesin Press (Inject) 22A', completed: 89, notGood: 11, isProblem: false },
		{ id: '19A', name: 'Mesin Press (Inject) 19A', completed: 82, notGood: 18, isProblem: false },
		{ id: '18A', name: 'Mesin Press (Inject) 18A', completed: 91, notGood: 9, isProblem: false },
		{ id: '17A', name: 'Mesin Press (Inject) 17A', completed: 74, notGood: 26, isProblem: false },
		{ id: '13B', name: 'Mesin Press 13B', completed: 83, notGood: 17, isProblem: false },
		{ id: '13A', name: 'Mesin Press 13A', completed: 79, notGood: 21, isProblem: false },
		{ id: '10B', name: 'Mesin Press 10B', completed: 85, notGood: 15, isProblem: false },
		{ id: '10A', name: 'Mesin Press 10A', completed: 88, notGood: 12, isProblem: false }
	];

	let charts: Record<string, any> = {};

	function goBack() {
		goto('/leader');
	}

	function selectMachine(machine: any) {
		goto(`/leader/mc_detail`);
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

	function initializeCharts() {
		// Hancurkan chart yang ada sebelumnya
		Object.keys(charts).forEach(machineId => {
			if (charts[machineId]?.total) {
				charts[machineId].total.destroy();
			}
			if (charts[machineId]?.notGood) {
				charts[machineId].notGood.destroy();
			}
		});
		charts = {};

		// Target untuk semua mesin
		const productionTarget = 90; // Target produksi 90%
		const notGoodTarget = 10;   // Target NG maksimal 10%

		machineData.forEach((machine) => {
			// Chart untuk Total Produksi
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
							datasets: [
								{
									label: 'Produksi',
									data: [machine.completed],
									backgroundColor: gradient,
									borderColor: '#1e40af',
									borderWidth: 1,
									borderRadius: 6,
									borderSkipped: false,
									barPercentage: 0.6,
									categoryPercentage: 0.8
								}
							]
						},
						options: {
							responsive: true,
							maintainAspectRatio: false,
							animation: {
								duration: 1000,
								easing: 'easeOutQuart'
							},
							plugins: {
								legend: {
									display: false
								},
								tooltip: {
									backgroundColor: 'rgba(15, 23, 42, 0.9)',
									padding: 10,
									displayColors: false,
									titleFont: {
										size: 12
									},
									bodyFont: {
										size: 14,
										weight: 'bold'
									},
									callbacks: {
										label: function(context) {
											const value = context.raw as number;
											const target = productionTarget;
											const diff = value - target;
											const status = diff >= 0 ? '✓' : '✗';
											return [`${status} Produksi: ${value}%`, `Target: ${target}%`, `Selisih: ${diff >= 0 ? '+' : ''}${diff}%`];
										}
									}
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
												content: `Target: ${productionTarget}%`,
												position: 'end',
												backgroundColor: '#10b981',
												color: 'white',
												font: {
													size: 9,
													weight: 'bold'
												},
												padding: {
													top: 4,
													bottom: 4,
													left: 8,
													right: 8
												},
												borderRadius: 4
											}
										}
									}
								}
							},
							scales: {
								y: {
									beginAtZero: true,
									max: 100,
									ticks: {
										font: {
											size: 9
										},
										color: '#64748b',
										padding: 4,
										callback: function(value) {
											return value + '%';
										}
									},
									grid: {
										color: 'rgba(100, 116, 139, 0.1)'
									},
									border: {
										display: false
									}
								},
								x: {
									ticks: {
										font: {
											size: 11,
											weight: 'bold'
										},
										color: '#1e293b'
									},
									grid: {
										display: false
									},
									border: {
										display: false
									}
								}
							}
						}
					});
					
					if (!charts[machine.id]) charts[machine.id] = {};
					charts[machine.id].total = newChart;
				}
			}

			// Chart untuk Not Good
			const notGoodCanvasId = `chart-notgood-${machine.id}`;
			const notGoodCanvas = document.getElementById(notGoodCanvasId) as HTMLCanvasElement;
			
			if (notGoodCanvas) {
				const ctx = notGoodCanvas.getContext('2d');
				if (ctx) {
					const gradient = ctx.createLinearGradient(0, 0, 0, notGoodCanvas.height);
					gradient.addColorStop(0, '#ef4444');
					gradient.addColorStop(1, '#dc2626');

					const newChart = new Chart(notGoodCanvas, {
						type: 'bar',
						data: {
							labels: ['Not Good'],
							datasets: [
								{
									label: 'NG',
									data: [machine.notGood],
									backgroundColor: gradient,
									borderColor: '#b91c1c',
									borderWidth: 1,
									borderRadius: 6,
									borderSkipped: false,
									barPercentage: 0.6,
									categoryPercentage: 0.8
								}
							]
						},
						options: {
							responsive: true,
							maintainAspectRatio: false,
							animation: {
								duration: 1000,
								easing: 'easeOutQuart'
							},
							plugins: {
								legend: {
									display: false
								},
								tooltip: {
									backgroundColor: 'rgba(15, 23, 42, 0.9)',
									padding: 10,
									displayColors: false,
									titleFont: {
										size: 12
									},
									bodyFont: {
										size: 14,
										weight: 'bold'
									},
									callbacks: {
										label: function(context) {
											const value = context.raw as number;
											const target = notGoodTarget;
											const diff = target - value;
											const status = value <= target ? '✓' : '✗';
											return [
												`${status} Not Good: ${value}%`,
												`Target: ≤${target}%`,
												value <= target ? `Baik: ${diff}% di bawah target` : `Perlu perbaikan: ${-diff}% di atas target`
											];
										}
									}
								},
								annotation: {
									annotations: {
										targetLine: {
											type: 'line',
											yMin: notGoodTarget,
											yMax: notGoodTarget,
											borderColor: '#f59e0b',
											borderWidth: 2,
											borderDash: [5, 5],
											label: {
												content: `Target: ≤${notGoodTarget}%`,
												position: 'end',
												backgroundColor: '#f59e0b',
												color: 'white',
												font: {
													size: 9,
													weight: 'bold'
												},
												padding: {
													top: 4,
													bottom: 4,
													left: 8,
													right: 8
												},
												borderRadius: 4
											}
										}
									}
								}
							},
							scales: {
								y: {
									beginAtZero: true,
									max: 40,
									ticks: {
										font: {
											size: 9
										},
										color: '#64748b',
										padding: 4,
										callback: function(value) {
											return value + '%';
										}
									},
									grid: {
										color: 'rgba(100, 116, 139, 0.1)'
									},
									border: {
										display: false
									}
								},
								x: {
									ticks: {
										font: {
											size: 11,
											weight: 'bold'
										},
										color: '#1e293b'
									},
									grid: {
										display: false
									},
									border: {
										display: false
									}
								}
							}
						}
					});
					
					if (!charts[machine.id]) charts[machine.id] = {};
					charts[machine.id].notGood = newChart;
				}
			}
		});
	}

	// Reactive variables for filtering
	$: pressMachines = machineData.filter(m => ['07B', '16B', '16A', '12B', '12A', '14B', '14A', '09A', '09B', '11A', '11B', '15A', '15B', '02A', '04A', '04B', '13B', '13A', '10B', '10A'].includes(m.id));
	$: injectMachines = machineData.filter(m => ['20A', '21A', '22A', '19A', '18A', '17A'].includes(m.id));

	onMount(() => {
		setTimeout(() => {
			initializeCharts();
		}, 100);
		
		// Reinitialize charts on window resize
		window.addEventListener('resize', handleResize);
	});

	function handleResize() {
		setTimeout(() => {
			initializeCharts();
		}, 300);
	}

	onDestroy(() => {
		// Hancurkan semua chart saat komponen dihancurkan
		Object.keys(charts).forEach(machineId => {
			if (charts[machineId]?.total) {
				charts[machineId].total.destroy();
			}
			if (charts[machineId]?.notGood) {
				charts[machineId].notGood.destroy();
			}
		});
		window.removeEventListener('resize', handleResize);
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-4 md:p-6">
	<div class="max-w-7xl mx-auto">
		<!-- Header -->
		<div class="mb-6 flex items-center justify-between">
			<div>
				<h1 class="text-2xl md:text-3xl font-bold text-slate-800">Pressing Line Dashboard</h1>
				<p class="text-slate-600 mt-1 text-sm md:text-base">Pantau status semua mesin press dengan diagram real-time</p>
			</div>
		</div>

		<!-- Dashboard Content -->
		<div class="bg-white rounded-2xl shadow-lg p-4 md:p-6 border border-slate-200">
			<h2 class="text-lg font-bold text-slate-800 mb-6">Layout Pabrik Pressing</h2>

			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<!-- Inject Machines Column -->
				<div>
					<h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider flex items-center gap-2">
						<span class="w-2 h-2 bg-purple-500 rounded-full"></span>
						Inject Machines
					</h3>
					<div class="space-y-4">
						{#each injectMachines as machine}
							{@const status = getMachineStatusColor(machine)}
							<div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
								<!-- Machine Header -->
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
											{machine.completed}% selesai
										</div>
									</div>
									<i class="fa-solid fa-arrow-right text-slate-400"></i>
								</button>

								<!-- Machine Charts -->
								<div class="px-4 pb-4 pt-2 border-t border-slate-200/50 space-y-4">
									<!-- Total Production Chart -->
									<div>
										<div class="flex justify-between items-center mb-2">
											<span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
												<i class="fa-solid fa-chart-bar text-blue-500"></i>
												Total Produksi
											</span>
											<span class="text-sm font-bold {machine.completed >= 90 ? 'text-green-600' : 'text-orange-600'}">
												{machine.completed}%
												{#if machine.completed >= 90}
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
									
									<!-- Not Good Chart -->
									<div>
										<div class="flex justify-between items-center mb-2">
											<span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
												<i class="fa-solid fa-triangle-exclamation text-red-500"></i>
												Not Good
											</span>
											<span class="text-sm font-bold {machine.notGood <= 10 ? 'text-green-600' : 'text-red-600'}">
												{machine.notGood}%
												{#if machine.notGood <= 10}
													<i class="fa-solid fa-check ml-1 text-xs"></i>
												{:else}
													<i class="fa-solid fa-xmark ml-1 text-xs"></i>
												{/if}
											</span>
										</div>
										<div class="h-32">
											<canvas id="chart-notgood-{machine.id}"></canvas>
										</div>
									</div>
									
									<!-- Stats Summary -->
									<div class="grid grid-cols-2 gap-3 pt-3 border-t border-slate-200/50">
										<div class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-lg p-3 text-center border border-blue-200">
											<p class="text-xs text-blue-700 font-semibold mb-1">PRODUKSI</p>
											<div class="flex items-center justify-center gap-1">
												<p class="text-xl font-bold text-blue-800">{machine.completed}%</p>
												{#if machine.completed >= 90}
													<span class="text-green-600 text-xs">
														<i class="fa-solid fa-arrow-up"></i>
													</span>
												{:else}
													<span class="text-orange-600 text-xs">
														<i class="fa-solid fa-arrow-down"></i>
													</span>
												{/if}
											</div>
											<p class="text-[10px] text-blue-600 mt-1">Target: 90%</p>
										</div>
										<div class="bg-gradient-to-r from-red-50 to-red-100 rounded-lg p-3 text-center border border-red-200">
											<p class="text-xs text-red-700 font-semibold mb-1">NOT GOOD</p>
											<div class="flex items-center justify-center gap-1">
												<p class="text-xl font-bold text-red-800">{machine.notGood}%</p>
												{#if machine.notGood <= 10}
													<span class="text-green-600 text-xs">
														<i class="fa-solid fa-arrow-down"></i>
													</span>
												{:else}
													<span class="text-red-600 text-xs">
														<i class="fa-solid fa-arrow-up"></i>
													</span>
												{/if}
											</div>
											<p class="text-[10px] text-red-600 mt-1">Target: ≤10%</p>
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>

				<!-- Press Machines Column -->
				<div>
					<h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider flex items-center gap-2">
						<span class="w-2 h-2 bg-blue-500 rounded-full"></span>
						Press Machines
					</h3>
					<div class="space-y-4">
						{#each pressMachines as machine}
							{@const status = getMachineStatusColor(machine)}
							<div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border border-slate-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
								<!-- Machine Header -->
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
											{machine.completed}% selesai
										</div>
									</div>
									<i class="fa-solid fa-arrow-right text-slate-400"></i>
								</button>

								<!-- Machine Charts -->
								<div class="px-4 pb-4 pt-2 border-t border-slate-200/50 space-y-4">
									<!-- Total Produksi Chart -->
									<div>
										<div class="flex justify-between items-center mb-2">
											<span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
												<i class="fa-solid fa-chart-bar text-blue-500"></i>
												Total Produksi
											</span>
											<span class="text-sm font-bold {machine.completed >= 90 ? 'text-green-600' : 'text-orange-600'}">
												{machine.completed}%
												{#if machine.completed >= 90}
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
									
									<!-- Not Good Chart -->
									<div>
										<div class="flex justify-between items-center mb-2">
											<span class="text-xs font-semibold text-slate-700 flex items-center gap-1">
												<i class="fa-solid fa-triangle-exclamation text-red-500"></i>
												Not Good
											</span>
											<span class="text-sm font-bold {machine.notGood <= 10 ? 'text-green-600' : 'text-red-600'}">
												{machine.notGood}%
												{#if machine.notGood <= 10}
													<i class="fa-solid fa-check ml-1 text-xs"></i>
												{:else}
													<i class="fa-solid fa-xmark ml-1 text-xs"></i>
												{/if}
											</span>
										</div>
										<div class="h-32">
											<canvas id="chart-notgood-{machine.id}"></canvas>
										</div>
									</div>
									
									<!-- Stats Summary -->
									<div class="grid grid-cols-2 gap-3 pt-3 border-t border-slate-200/50">
										<div class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-lg p-3 text-center border border-blue-200">
											<p class="text-xs text-blue-700 font-semibold mb-1">PRODUKSI</p>
											<div class="flex items-center justify-center gap-1">
												<p class="text-xl font-bold text-blue-800">{machine.completed}%</p>
												{#if machine.completed >= 90}
													<span class="text-green-600 text-xs">
														<i class="fa-solid fa-arrow-up"></i>
													</span>
												{:else}
													<span class="text-orange-600 text-xs">
														<i class="fa-solid fa-arrow-down"></i>
													</span>
												{/if}
											</div>
											<p class="text-[10px] text-blue-600 mt-1">Target: 90%</p>
										</div>
										<div class="bg-gradient-to-r from-red-50 to-red-100 rounded-lg p-3 text-center border border-red-200">
											<p class="text-xs text-red-700 font-semibold mb-1">NOT GOOD</p>
											<div class="flex items-center justify-center gap-1">
												<p class="text-xl font-bold text-red-800">{machine.notGood}%</p>
												{#if machine.notGood <= 10}
													<span class="text-green-600 text-xs">
														<i class="fa-solid fa-arrow-down"></i>
													</span>
												{:else}
													<span class="text-red-600 text-xs">
														<i class="fa-solid fa-arrow-up"></i>
													</span>
												{/if}
											</div>
											<p class="text-[10px] text-red-600 mt-1">Target: ≤10%</p>
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<!-- Legend -->
		<div class="mt-6 bg-white rounded-xl p-4 border border-slate-200 shadow-sm">
			<h3 class="text-sm font-semibold text-slate-800 mb-3">Keterangan Chart</h3>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				<div class="flex items-center gap-2">
					<div class="w-4 h-4 bg-gradient-to-b from-blue-500 to-blue-600 rounded"></div>
					<span class="text-sm text-slate-600">Total Produksi</span>
				</div>
				<div class="flex items-center gap-2">
					<div class="w-4 h-4 bg-gradient-to-b from-red-500 to-red-600 rounded"></div>
					<span class="text-sm text-slate-600">Not Good (NG)</span>
				</div>	
				<div class="flex items-center gap-2">
					<div class="w-4 h-4 border-2 border-green-500 border-dashed"></div>
					<span class="text-sm text-slate-600">Garis Target (90% Produksi, ≤10% NG)</span>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	:global(body) {
		background-color: #f8fafc;
		font-family: 'Inter', sans-serif;
	}

	/* Styling untuk canvas charts */
	canvas {
		width: 100% !important;
		height: 100% !important;
	}
	
	/* Custom scrollbar */
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