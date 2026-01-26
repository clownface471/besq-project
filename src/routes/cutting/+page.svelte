<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	interface FormData {
		partName: string;
		quantity: number;
		reject: number;
	}

	interface DailyLog {
		id: number;
		time: string;
		partName: string;
		quantity: number;
		reject: number;
		total: number;
	}

	// Reactive state
	let formData = $state<FormData>({
		partName: '',
		quantity: 0,
		reject: 0,
	});

	let dailyLogs = $state<DailyLog[]>([]);
	let isLoading = $state(false);
	let error = $state<string | null>(null);
	let successMessage = $state<string | null>(null);

	// Derived state for total calculation
	let total = $derived(Math.max(0, formData.quantity - formData.reject));

	// Role guard on mount
	onMount(() => {
		if (!$auth.user || ($auth.user.role !== 'cutting' && $auth.user.role !== 'admin')) {
			alert('⚠️ Unauthorized: Only Cutting or Admin users can access this page');
			goto('/');
			return;
		}
		
		// Fetch daily logs
		fetchDailyLogs();
	});

	// Fetch daily logs from API
	async function fetchDailyLogs() {
		try {
			isLoading = true;
			const response = await fetch('/api/cutting/today', {
				headers: {
					'Authorization': `Bearer ${$auth.token}`,
				},
			});

			if (!response.ok) {
				throw new Error('Failed to fetch daily logs');
			}

			const data = await response.json();
			dailyLogs = data.cutting_logs || [];
		} catch (err) {
			console.error('Error fetching daily logs:', err);
			error = 'Failed to load daily logs';
		} finally {
			isLoading = false;
		}
	}

	// Handle form submission
	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!formData.partName.trim()) {
			error = 'Part name is required';
			return;
		}

		if (formData.quantity <= 0) {
			error = 'Quantity must be greater than 0';
			return;
		}

		if (formData.reject < 0) {
			error = 'Reject cannot be negative';
			return;
		}

		try {
			isLoading = true;
			error = null;
			successMessage = null;

			const response = await fetch('/api/cutting', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${$auth.token}`,
				},
				body: JSON.stringify({
					partName: formData.partName,
					quantity: formData.quantity,
					reject: formData.reject,
				}),
			});

			if (!response.ok) {
				throw new Error('Failed to create cutting entry');
			}

			successMessage = '✓ Cutting entry recorded successfully';
			
			// Reset form
			formData = {
				partName: '',
				quantity: 0,
				reject: 0,
			};

			// Refresh daily logs
			await fetchDailyLogs();
		} catch (err) {
			error = err instanceof Error ? err.message : 'An error occurred';
			console.error('Error submitting form:', err);
		} finally {
			isLoading = false;
		}
	}

	// Format time from ISO string
	function formatTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
	}

	// Calculate total for the day
	let totalQuantity = $derived(dailyLogs.reduce((sum, log) => sum + log.quantity, 0));
	let totalReject = $derived(dailyLogs.reduce((sum, log) => sum + log.reject, 0));
	let totalProduction = $derived(dailyLogs.reduce((sum, log) => sum + log.total, 0));

	// Clear error message
	function clearError() {
		error = null;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 pb-12">
	<!-- Header -->
	<header class="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200 shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
			<div>
				<h1 class="text-2xl md:text-3xl font-bold text-slate-800">
					✂️ Cutting Department
				</h1>
				<p class="text-sm text-slate-500 mt-1">Production Input & Daily Log</p>
			</div>
			<div class="text-right">
				<p class="text-sm text-slate-600 font-medium">
					{$auth.user?.name || 'User'} • <span class="capitalize bg-blue-50 text-blue-700 px-2 py-1 rounded inline-block">{$auth.user?.role}</span>
				</p>
				<p class="text-xs text-slate-400 mt-1">
					{new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
				</p>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-8 space-y-6">
		<!-- Alert Messages -->
		{#if error}
			<div class="bg-rose-50 border border-rose-200 rounded-lg p-4 flex justify-between items-center animate-pulse">
				<span class="text-rose-700 font-medium">⚠️ {error}</span>
				<button onclick={clearError} class="text-rose-600 hover:text-rose-800">✕</button>
			</div>
		{/if}

		{#if successMessage}
			<div class="bg-emerald-50 border border-emerald-200 rounded-lg p-4 animate-pulse">
				<span class="text-emerald-700 font-medium">{successMessage}</span>
			</div>
		{/if}

		<!-- Grid Layout -->
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Input Form Card -->
			<div class="lg:col-span-1">
				<div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6 sticky top-24">
					<div class="flex items-center gap-3 mb-6">
						<div class="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
							<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 1112 2.944a11.954 11.954 0 018.618 3.04A12.02 12.02 0 0121 12" />
							</svg>
						</div>
						<h2 class="text-xl font-bold text-slate-800">Input Form</h2>
					</div>

					<form onsubmit={handleSubmit} class="space-y-4">
						<!-- Part Name -->
						<div>
							<label for="partName" class="block text-sm font-medium text-slate-700 mb-2">
								Part Name <span class="text-rose-500">*</span>
							</label>
							<input
								type="text"
								id="partName"
								bind:value={formData.partName}
								placeholder="e.g., PART-001"
								class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
								disabled={isLoading}
							/>
						</div>

						<!-- Quantity -->
						<div>
							<label for="quantity" class="block text-sm font-medium text-slate-700 mb-2">
								Quantity <span class="text-rose-500">*</span>
							</label>
							<input
								type="number"
								id="quantity"
								bind:value={formData.quantity}
								min="1"
								placeholder="0"
								class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
								disabled={isLoading}
							/>
						</div>

						<!-- Reject -->
						<div>
							<label for="reject" class="block text-sm font-medium text-slate-700 mb-2">
								Reject <span class="text-slate-400 text-xs">(Optional)</span>
							</label>
							<input
								type="number"
								id="reject"
								bind:value={formData.reject}
								min="0"
								placeholder="0"
								class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
								disabled={isLoading}
							/>
						</div>

						<!-- Total Display -->
						<div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
							<div class="text-xs font-medium text-slate-600 mb-1">TOTAL PRODUCTION</div>
							<div class="text-3xl font-bold text-blue-600">{total}</div>
							<div class="text-xs text-slate-500 mt-2">
								= Quantity ({formData.quantity}) - Reject ({formData.reject})
							</div>
						</div>

						<!-- Submit Button -->
						<button
							type="submit"
							disabled={isLoading}
							class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-slate-400 text-white font-semibold py-2 px-4 rounded-lg transition flex items-center justify-center gap-2"
						>
							{#if isLoading}
								<svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Submitting...
							{:else}
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
								</svg>
								Record Entry
							{/if}
						</button>
					</form>
				</div>
			</div>

			<!-- Daily Logs Table -->
			<div class="lg:col-span-2">
				<div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
					<!-- Header -->
					<div class="bg-gradient-to-r from-slate-50 to-slate-100 border-b border-slate-200 p-6 flex items-center justify-between">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
								<svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
								</svg>
							</div>
							<div>
								<h3 class="text-lg font-bold text-slate-800">Daily Logs</h3>
								<p class="text-xs text-slate-500">Today's production entries</p>
							</div>
						</div>
						<span class="text-sm font-medium text-slate-500 bg-white px-3 py-1 rounded-lg border border-slate-200">
							{dailyLogs.length} entries
						</span>
					</div>

					<!-- Summary Badges -->
					<div class="grid grid-cols-3 gap-4 p-6 border-b border-slate-200 bg-slate-50">
						<div class="text-center">
							<div class="text-2xl font-bold text-blue-600">{totalQuantity}</div>
							<div class="text-xs font-medium text-slate-600 mt-1">Total Qty</div>
						</div>
						<div class="text-center">
							<div class="text-2xl font-bold text-rose-600">{totalReject}</div>
							<div class="text-xs font-medium text-slate-600 mt-1">Total Reject</div>
						</div>
						<div class="text-center">
							<div class="text-2xl font-bold text-emerald-600">{totalProduction}</div>
							<div class="text-xs font-medium text-slate-600 mt-1">Total Hari Ini</div>
						</div>
					</div>

					<!-- Table Content -->
					<div class="overflow-x-auto">
						{#if isLoading}
							<div class="flex items-center justify-center py-12">
								<div class="text-slate-500">Loading logs...</div>
							</div>
						{:else if dailyLogs.length === 0}
							<div class="flex flex-col items-center justify-center py-12 text-slate-500">
								<svg class="w-12 h-12 mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
								</svg>
								<p class="text-sm font-medium">No entries yet</p>
								<p class="text-xs">Start recording production data to see logs here</p>
							</div>
						{:else}
							<table class="w-full">
								<thead>
									<tr class="border-b border-slate-200 bg-slate-50">
										<th class="px-6 py-3 text-left text-xs font-bold text-slate-700 uppercase tracking-wider">Time</th>
										<th class="px-6 py-3 text-left text-xs font-bold text-slate-700 uppercase tracking-wider">Part Name</th>
										<th class="px-6 py-3 text-center text-xs font-bold text-slate-700 uppercase tracking-wider">Qty</th>
										<th class="px-6 py-3 text-center text-xs font-bold text-slate-700 uppercase tracking-wider">Reject</th>
										<th class="px-6 py-3 text-center text-xs font-bold text-slate-700 uppercase tracking-wider">Total</th>
									</tr>
								</thead>
								<tbody>
									{#each dailyLogs as log (log.id)}
										<tr class="border-b border-slate-100 hover:bg-slate-50 transition">
											<td class="px-6 py-4 text-sm font-mono text-slate-600">
												{formatTime(log.time)}
											</td>
											<td class="px-6 py-4 text-sm font-semibold text-slate-700">
												{log.partName}
											</td>
											<td class="px-6 py-4 text-sm text-center font-medium text-slate-700">
												{log.quantity}
											</td>
											<td class="px-6 py-4 text-sm text-center font-medium text-rose-600">
												{log.reject}
											</td>
											<td class="px-6 py-4 text-sm text-center font-bold text-emerald-600">
												{log.total}
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</main>
</div>