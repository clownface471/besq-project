<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';

	const API_URL = 'http://localhost:8080';

	let dashboardData = {
		mixing: { total: 0, pending: 0, completed: 0 },
		cutting: { total: 0, pending: 0, completed: 0 },
		pressing: { total: 0, pending: 0, completed: 0 },
		finishing: { total: 0, pending: 0, completed: 0 }
	};

	let isLoading = true;

	const divisions = [
		{
			id: 'mixing',
			name: 'Mixing',
			icon: 'fa-solid fa-blender',
			color: 'from-blue-500 to-cyan-500',
			bgColor: 'bg-blue-50',
			borderColor: 'border-blue-200',
			textColor: 'text-blue-700',
			description: 'Proses pencampuran bahan baku'
		},
		{
			id: 'cutting',
			name: 'Cutting',
			icon: 'fa-solid fa-scissors',
			color: 'from-purple-500 to-pink-500',
			bgColor: 'bg-purple-50',
			borderColor: 'border-purple-200',
			textColor: 'text-purple-700',
			description: 'Proses pemotongan material'
		},
		{
			id: 'pressing',
			name: 'Pressing',
			icon: 'fa-solid fa-square',
			color: 'from-orange-500 to-red-500',
			bgColor: 'bg-orange-50',
			borderColor: 'border-orange-200',
			textColor: 'text-orange-700',
			description: 'Proses pengepresan produk'
		},
		{
			id: 'finishing',
			name: 'Finishing',
			icon: 'fa-solid fa-star',
			color: 'from-green-500 to-emerald-500',
			bgColor: 'bg-green-50',
			borderColor: 'border-green-200',
			textColor: 'text-green-700',
			description: 'Proses penyelesaian akhir'
		}
	];

	async function fetchDashboardData() {
		try {
			const res = await fetch(`${API_URL}/manager/dashboard`, {
				headers: { Authorization: `Bearer ${$auth.token}` }
			});

			if (res.ok) {
				const result = await res.json();
				dashboardData = result.data;
			}
		} catch (error) {
			console.error('Error fetching dashboard data:', error);
			// Fallback data untuk demonstration
			dashboardData = {
				mixing: { total: 45, pending: 12, completed: 33 },
				cutting: { total: 38, pending: 8, completed: 30 },
				pressing: { total: 52, pending: 15, completed: 37 },
				finishing: { total: 41, pending: 5, completed: 36 }
			};
		} finally {
			isLoading = false;
		}
	}

	function getProgressPercentage(division: string) {
		if (dashboardData[division as keyof typeof dashboardData].total === 0) return 0;
		return Math.round(
			((dashboardData[division as keyof typeof dashboardData].completed /
				dashboardData[division as keyof typeof dashboardData].total) *
				100)
		);
	}

	function getStatusColor(percentage: number, divisionId: string) {
		if (divisionId === 'cutting') {
			return { bg: 'bg-red-100', bgFill: '#ef4444', text: 'text-red-700', badge: 'bg-red-500', label: 'Problem' };
		}
		if (divisionId === 'mixing') {
			return { bg: 'bg-green-100', bgFill: '#22c55e', text: 'text-green-700', badge: 'bg-green-500', label: 'Lancar' };
		}
		if (divisionId === 'pressing') {
			return { bg: 'bg-yellow-100', bgFill: '#eab308', text: 'text-yellow-700', badge: 'bg-yellow-500', label: 'Sedang Berjalan' };
		}
		if (divisionId === 'finishing') {
			return { bg: 'bg-green-100', bgFill: '#22c55e', text: 'text-green-700', badge: 'bg-green-500', label: 'Lancar' };
		}
		return { bg: 'bg-slate-100', bgFill: '#64748b', text: 'text-slate-700', badge: 'bg-slate-500', label: 'Unknown' };
	}

	function handleDetailClick(divisionId: string) {
		if (divisionId === 'pressing') {
			goto('/manager/prs-ldr');
		}
	}

	onMount(() => {
		fetchDashboardData();
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-6">
	<div class="max-w-7xl mx-auto">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-slate-800 mb-2">Dashboard Manager</h1>
			<p class="text-slate-600">Kelola dan pantau semua divisi produksi Anda</p>
		</div>

		<!-- Loading State -->
		{#if isLoading}
			<div class="flex justify-center items-center min-h-64">
				<div class="text-center">
					<div class="inline-block animate-spin">
						<i class="fa-solid fa-spinner text-indigo-600 text-4xl"></i>
					</div>
					<p class="text-slate-600 mt-4">Memuat data...</p>
				</div>
			</div>
		{:else}
			<!-- Cards Grid -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
				{#each divisions as division (division.id)}
					{@const percentage = getProgressPercentage(division.id)}
					{@const status = getStatusColor(percentage, division.id)}
					<div
						class="group relative bg-white rounded-2xl shadow-md hover:shadow-2xl transition-all duration-300 overflow-hidden border border-slate-200 hover:border-slate-300"
					>
						<!-- Background Gradient Line -->
						<div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r {division.color}"></div>

						<!-- Card Content -->
						<div class="p-6 relative z-10">
							<!-- Icon & Title -->
							<div class="flex items-center justify-between mb-4">
								<div class="flex items-center gap-3">
									<div class="{division.bgColor} p-3 rounded-xl group-hover:scale-110 transition-transform">
										<i class="{division.icon} {division.textColor} text-xl"></i>
									</div>
									<div>
										<h3 class="text-lg font-bold text-slate-800">{division.name}</h3>
										<p class="text-xs text-slate-500">{division.description}</p>
									</div>
								</div>
							</div>

							<!-- Status Badge -->
							<div class="mb-4 flex items-center gap-2">
								<span class="{status.badge} w-3 h-3 rounded-full"></span>
								<span class="{status.text} text-sm font-semibold">{status.label}</span>
							</div>

							<!-- Progress Bar with Percentage -->
							<div class="space-y-2">
								<div class="flex justify-between items-center">
									<span class="text-xs text-slate-500">Progress</span>
									<span class="{status.text} font-bold text-lg">{percentage}%</span>
								</div>
								<div class="w-full {status.bg} rounded-full h-3 overflow-hidden border border-slate-200">
									<div
										class="h-full rounded-full transition-all duration-500"
										style="width: {percentage}%; background-color: {status.bgFill};"
									></div>
								</div>
							</div>

							<!-- Task Count -->
							<p class="text-xs text-slate-500 mt-3 text-center">
								{dashboardData[division.id as keyof typeof dashboardData].completed} dari {dashboardData[division.id as keyof typeof dashboardData].total}
							</p>
						</div>

						<!-- Hover Action Button -->
						<div
							class="px-6 py-3 bg-gradient-to-r {division.color} bg-opacity-0 group-hover:bg-opacity-5 transition-all duration-300 border-t border-slate-200 group-hover:border-slate-300"
						>
							<button
								on:click={() => handleDetailClick(division.id)}
								class="w-full py-2 px-3 rounded-lg text-sm font-medium {division.textColor} hover:{division.bgColor} transition-all duration-300 flex items-center justify-center gap-2"
							>
								<i class="fa-solid fa-arrow-right"></i>
								Lihat Detail
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		background-color: #f8fafc;
	}
</style>