<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

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

	// tambah reactive kelompok mesin
	$: pressMachines = machineData.filter(m => ['07B', '16B', '16A', '12B', '12A', '14B'].includes(m.id));
	$: injectMachines = machineData.filter(m => ['20A', '21A', '22A', '17A'].includes(m.id));

	onMount(() => {
		// Data sudah di-load
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-6">
	<div class="max-w-7xl mx-auto">
		<div class="mb-6 flex items-center justify-between">
			<div>
				<h1 class="text-3xl font-bold text-slate-800">Pressing Line Dashboard</h1>
				<p class="text-slate-600 mt-1">Pantau status semua mesin press</p>
			</div>
		</div>

		<div class="bg-white rounded-2xl shadow-md p-6 border border-slate-200">
			<h2 class="text-lg font-bold text-slate-800 mb-6">Layout Pabrik Pressing</h2>

			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<!-- Inject Machines Column -->
				<div>
					<h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider">Inject Machines</h3>
					<div class="space-y-4">
						{#each injectMachines as machine}
							{@const status = getMachineStatusColor(machine)}
							<div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-lg border border-slate-200 overflow-hidden hover:shadow-lg transition">
								<!-- Machine Header -->
								<button
									on:click={() => selectMachine(machine)}
									class="w-full flex items-center gap-3 p-4 hover:bg-slate-200 transition"
								>
									<div class="w-10 h-10 flex items-center justify-center rounded-full {status.badge} text-white font-bold text-sm flex-shrink-0">
										{machine.id}
									</div>
									<div class="flex-1 text-left min-w-0">
										<div class="font-semibold text-slate-800 truncate">{machine.name}</div>
										<div class="text-xs text-slate-600">{machine.completed}% selesai</div>
									</div>
									<i class="fa-solid fa-arrow-right text-slate-400"></i>
								</button>

								<!-- Machine Details -->
								<div class="px-4 pb-4 pt-2 border-t border-slate-200 space-y-3">
									<div>
										<div class="flex justify-between items-center mb-1">
											<span class="text-xs font-medium text-slate-600">Selesai</span>
											<span class="text-sm font-bold text-blue-600">{machine.completed}%</span>
										</div>
										<div class="w-full bg-slate-300 rounded-full h-2 overflow-hidden">
											<div class="h-full bg-blue-500" style="width: {machine.completed}%"></div>
										</div>
									</div>
									<div>
										<div class="flex justify-between items-center mb-1">
											<span class="text-xs font-medium text-slate-600">Not Good</span>
											<span class="text-sm font-bold text-red-600">{machine.notGood}%</span>
										</div>
										<div class="w-full bg-slate-300 rounded-full h-2 overflow-hidden">
											<div class="h-full bg-red-500" style="width: {machine.notGood}%"></div>
										</div>
									</div>
									<div class="grid grid-cols-2 gap-2 pt-2">
										<div class="bg-blue-50 rounded p-2 text-center">
											<p class="text-xs text-slate-600">Selesai</p>
											<p class="text-lg font-bold text-blue-600">{machine.completed}</p>
										</div>
										<div class="bg-red-50 rounded p-2 text-center">
											<p class="text-xs text-slate-600">Not Good</p>
											<p class="text-lg font-bold text-red-600">{machine.notGood}</p>
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>

				<!-- Press Machines Column -->
				<div>
					<h3 class="text-xs font-bold text-slate-600 mb-4 uppercase tracking-wider">Press Machines</h3>
					<div class="space-y-4">
						{#each pressMachines as machine}
							{@const status = getMachineStatusColor(machine)}
							<div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-lg border border-slate-200 overflow-hidden hover:shadow-lg transition">
								<!-- Machine Header -->
								<button
									on:click={() => selectMachine(machine)}
									class="w-full flex items-center gap-3 p-4 hover:bg-slate-200 transition"
								>
									<div class="w-10 h-10 flex items-center justify-center rounded-full {status.badge} text-white font-bold text-sm flex-shrink-0">
										{machine.id}
									</div>
									<div class="flex-1 text-left min-w-0">
										<div class="font-semibold text-slate-800 truncate">{machine.name}</div>
										<div class="text-xs text-slate-600">{machine.completed}% selesai</div>
									</div>
									<i class="fa-solid fa-arrow-right text-slate-400"></i>
								</button>

								<!-- Machine Details -->
								<div class="px-4 pb-4 pt-2 border-t border-slate-200 space-y-3">
									<div>
										<div class="flex justify-between items-center mb-1">
											<span class="text-xs font-medium text-slate-600">Selesai</span>
											<span class="text-sm font-bold text-blue-600">{machine.completed}%</span>
										</div>
										<div class="w-full bg-slate-300 rounded-full h-2 overflow-hidden">
											<div class="h-full bg-blue-500" style="width: {machine.completed}%"></div>
										</div>
									</div>
									<div>
										<div class="flex justify-between items-center mb-1">
											<span class="text-xs font-medium text-slate-600">Not Good</span>
											<span class="text-sm font-bold text-red-600">{machine.notGood}%</span>
										</div>
										<div class="w-full bg-slate-300 rounded-full h-2 overflow-hidden">
											<div class="h-full bg-red-500" style="width: {machine.notGood}%"></div>
										</div>
									</div>
									<div class="grid grid-cols-2 gap-2 pt-2">
										<div class="bg-blue-50 rounded p-2 text-center">
											<p class="text-xs text-slate-600">Selesai</p>
											<p class="text-lg font-bold text-blue-600">{machine.completed}</p>
										</div>
										<div class="bg-red-50 rounded p-2 text-center">
											<p class="text-xs text-slate-600">Not Good</p>
											<p class="text-lg font-bold text-red-600">{machine.notGood}</p>
										</div>
									</div>
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
	:global(body) {
		background-color: #f8fafc;
	}
</style>
