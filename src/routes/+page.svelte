<script lang="ts">
	// @ts-nocheck
	import { auth, login, logout, type User } from '$lib/stores/auth';
	import Swal from 'sweetalert2';
	import { goto } from '$app/navigation';

	// --- STATE (Svelte 5 Runes) ---
	let nik = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let errorMessage = $state('');

	let stats = $state({
		totalEmployees: 0,
		totalOutput: 0,
		rejectRate: '0%',
		activeShift: 'N/A',
		cuttingOutput: 0,
		pressingOutput: 0,
		finishingOutput: 0,
		activities: [] as { user: string; action: string; time: string; status: string }[]
	});

	const API_URL = 'http://localhost:8080';

	// --- LIFECYCLE & EFFECTS ---

	// Effect to fetch stats when the user logs in
	$effect(() => {
		if ($auth.isLoggedIn) {
			fetchStats();
		}
	});

	// --- ASYNC FUNCTIONS ---

	async function fetchStats() {
		console.log('User logged in, fetching dashboard stats...');
		const token = $auth.token;
		if (!token) return;

		try {
			const response = await fetch(`${API_URL}/api/dashboard/stats`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.error || 'Failed to fetch stats');
			}

			const data = await response.json();
			console.log('Fetching stats...', data); // Debug log
			stats = data;
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load dashboard data.';
			console.error('Stats fetch error:', error);
			Swal.fire({
				title: 'Error!',
				text: 'Gagal memuat data statistik.',
				icon: 'error'
			});
		}
	}

	async function handleLogin(e: Event) {
		e.preventDefault();
		isLoading = true;
		errorMessage = '';

		try {
			const response = await fetch(`${API_URL}/login`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ nik, password })
			});

			const data = await response.json();
			if (!response.ok) {
				throw new Error(data.error || 'Login failed');
			}

			const token = data.token;
			const profileResponse = await fetch(`${API_URL}/api/users/profile`, {
				headers: { Authorization: `Bearer ${token}` }
			});

			if (!profileResponse.ok) {
				throw new Error('Failed to fetch user profile');
			}

			const profileData = await profileResponse.json();
			login(token, profileData.user as User);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'An error occurred during login';
			console.error('Login error:', error);
		} finally {
			isLoading = false;
		}
	}

	async function handleLogout() {
		const result = await Swal.fire({
			title: 'Yakin ingin keluar?',
			text: 'Anda harus login ulang untuk mengakses aplikasi.',
			icon: 'warning',
			showCancelButton: true,
			confirmButtonText: 'Ya, Keluar',
			cancelButtonText: 'Batal'
		});

		if (result.isConfirmed) {
			logout();
			goto('/');
		}
	}

	// --- STATIC DATA (for UI rendering as per admin/+page.svelte) ---
	let deptSummary = [
		{
			id: 'pressing',
			name: 'Pressing',
			icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10',
			totalStaff: 12,
			activeNow: 10,
			avgEfficiency: 88,
			target: 2000,
			color: 'amber'
		},
		{
			id: 'cutting',
			name: 'Cutting',
			icon: 'M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z',
			totalStaff: 8,
			activeNow: 8,
			avgEfficiency: 95,
			target: 2500,
			color: 'blue'
		},
		{
			id: 'finishing',
			name: 'Finishing',
			icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
			totalStaff: 15,
			activeNow: 14,
			avgEfficiency: 92,
			target: 1800,
			color: 'emerald'
		}
	];

	// --- UI HELPER FUNCTIONS ---
	function getStatusColor(status: string) {
		if (status === 'Completed') return 'bg-emerald-100 text-emerald-700 ring-emerald-600/20';
		if (status === 'Pending') return 'bg-slate-100 text-slate-700 ring-slate-600/20';
		return 'bg-amber-100 text-amber-700 ring-amber-600/20';
	}

	function getDeptOutput(deptId: string) {
		if (deptId === 'pressing') return stats.pressingOutput;
		if (deptId === 'cutting') return stats.cuttingOutput;
		if (deptId === 'finishing') return stats.finishingOutput;
		return 0;
	}

	function getCardColor(color: string) {
		const map = {
			amber: 'bg-amber-50 border-amber-100 hover:border-amber-300 text-amber-900',
			blue: 'bg-blue-50 border-blue-100 hover:border-blue-300 text-blue-900',
			emerald: 'bg-emerald-50 border-emerald-100 hover:border-emerald-300 text-emerald-900'
		};
		return map[color];
	}

	function getIconBg(color: string) {
		const map = {
			amber: 'bg-amber-100 text-amber-600',
			blue: 'bg-blue-100 text-blue-600',
			emerald: 'bg-emerald-100 text-emerald-600'
		};
		return map[color];
	}

	function getBarColor(color: string) {
		const map = {
			amber: 'bg-amber-500',
			blue: 'bg-blue-500',
			emerald: 'bg-emerald-500'
		};
		return map[color];
	}
</script>

<<<<<<< HEAD
<body class="bg-gray-100 min-h-screen flex items-center justify-center p-4">
  <div
    class="border w-full max-w-sm overflow-hidden shadow-2xl bg-white rounded-3xl p-8 transition-all duration-500 hover:shadow-[0_20px_50px_rgba(0,0,0,0.2)]"
  >
    <div class="mb-8">
      <div class="flex justify-center mb-2">
        <span
          class="bg-indigo-100 text-indigo-500 text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-widest"
          >Official Portal</span
        >
      </div>
      <h1
        class="font-bold text-2xl text-center text-gray-800 font-[Poppins] tracking-tight"
      >
        Besq User Login
      </h1>
      <div
        class="text-[#0065F8] text-7xl text-center flex items-center justify-center mt-6 drop-shadow-md"
      >
        <i class="fas fa-circle-user text-indigo-600"></i>
      </div>
    </div>
=======
<div class="font-sans text-slate-800">
	{#if !$auth.isLoggedIn}
		<!-- LOGIN FORM -->
		<div class="bg-gray-100 min-h-screen flex items-center justify-center p-4">
			<div
				class="border w-full max-w-sm overflow-hidden shadow-2xl bg-white rounded-3xl p-8 transition-all duration-500 hover:shadow-[0_20px_50px_rgba(0,0,0,0.2)]"
			>
				<div class="mb-8">
					<div class="flex justify-center mb-2">
						<span
							class="bg-blue-100 text-[#0065F8] text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-widest"
						>
							Official Portal
						</span>
					</div>
					<h1
						class="font-bold text-2xl text-center text-gray-800 font-[Poppins] tracking-tight"
					>
						Besq User Login
					</h1>
					<div class="text-[#0065F8] text-7xl text-center mt-6 drop-shadow-md">
						<i class="fas fa-circle-user text-blue-600"></i>
					</div>
				</div>
>>>>>>> f2fb86746a74cfc4be4db64df563914a12820552

				<form onsubmit={handleLogin} class="space-y-5">
					{#if errorMessage}
						<div
							class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl text-sm"
						>
							{errorMessage}
						</div>
					{/if}

					<div>
						<label
							for="nik"
							class="text-[11px] font-bold text-gray-400 uppercase tracking-widest ml-1"
						>
							NIK
						</label>
						<input
							id="nik"
							type="text"
							placeholder="Masukkan NIK"
							bind:value={nik}
							required
							disabled={isLoading}
							class="mt-1 border border-gray-100 bg-gray-50 rounded-xl w-full py-3 px-4 font-[inter] outline-none focus:bg-white focus:border-[#0065F8] focus:ring-4 focus:ring-blue-100 transition-all placeholder:text-gray-300 disabled:opacity-50"
						/>
					</div>

<<<<<<< HEAD
    <a href="/karyawan"
      class="bg-indigo-600 shadow-lg shadow-blue-300/50 transition-all duration-300 hover:bg-indigo-700 hover:shadow-blue-400/50 hover:-translate-y-1 text-white font-bold py-3.5 px-4 rounded-xl w-full mt-8 cursor-pointer active:scale-95 flex justify-center items-center gap-2"
    >
      <span>Submit</span>
    </a>
=======
					<div>
						<label
							for="password"
							class="text-[11px] font-bold text-gray-400 uppercase tracking-widest ml-1"
						>
							Password
						</label>
						<input
							id="password"
							type="password"
							placeholder="••••••••"
							bind:value={password}
							required
							disabled={isLoading}
							class="mt-1 border border-gray-100 bg-gray-50 rounded-xl w-full py-3 px-4 font-[inter] outline-none focus:bg-white focus:border-[#0065F8] focus:ring-4 focus:ring-blue-100 transition-all placeholder:text-gray-300 disabled:opacity-50"
						/>
					</div>
>>>>>>> f2fb86746a74cfc4be4db64df563914a12820552

					<button
						type="submit"
						disabled={isLoading}
						class="bg-[#0065F8] shadow-lg shadow-blue-300/50 transition-all duration-300 hover:bg-[#004bbd] hover:shadow-blue-400/50 hover:-translate-y-1 text-white font-bold py-3.5 px-4 rounded-xl w-full mt-8 flex justify-center items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if isLoading}
							<svg
								class="animate-spin h-5 w-5 text-white"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
							>
								<circle
									class="opacity-25"
									cx="12"
									cy="12"
									r="10"
									stroke="currentColor"
									stroke-width="4"
								></circle>
								<path
									class="opacity-75"
									fill="currentColor"
									d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
								></path>
							</svg>
							<span>Logging in...</span>
						{:else}
							<span>Login</span>
						{/if}
					</button>
				</form>
			</div>
		</div>
	{:else}
		<!-- CONSOLIDATED DYNAMIC DASHBOARD -->
		<div class="min-h-screen bg-slate-50 pb-12 relative">
			<div
				class="absolute inset-0 z-0 opacity-[0.03] pointer-events-none"
				style="background-image: radial-gradient(#64748b 1px, transparent 1px); background-size: 24px 24px;"
			></div>

			<main class="max-w-7xl mx-auto px-4 mt-6 space-y-6 relative z-10">
				<!-- Header section from admin page -->
				<div
					class="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 flex flex-col md:flex-row items-center justify-between gap-6"
				>
					<div class="flex items-center gap-5 w-full md:w-auto">
						<img
							src="https://i.pravatar.cc/300?u={$auth.user?.nik}"
							alt={$auth.user?.name}
							class="w-16 h-16 rounded-full object-cover border-2 border-white shadow-md ring-2 ring-slate-100"
						/>
						<div>
							<h2 class="text-xl font-bold text-slate-800">
								Halo, {$auth.user?.name}!
							</h2>
							<div class="flex items-center gap-2 text-sm text-slate-500 mt-1">
								<span
									class="font-medium bg-slate-100 px-2 py-0.5 rounded text-slate-600"
								>
									{$auth.user?.nik}
								</span>
								<span>•</span>
								<span class="capitalize">{$auth.user?.role}</span>
							</div>
						</div>
					</div>
					<div
						class="w-full md:w-auto flex flex-wrap items-center justify-center md:justify-end gap-6 border-t md:border-t-0 border-slate-100 pt-4 md:pt-0"
					>
						<div class="text-center">
							<p
								class="text-xs text-slate-400 font-bold uppercase tracking-wider"
							>
								Total Output
							</p>
							<p class="text-2xl font-bold text-slate-800">
								{stats.totalOutput.toLocaleString('id-ID')}
								<span class="text-sm font-normal text-slate-500">Lot</span>
							</p>
						</div>
						<div class="text-center">
							<p
								class="text-xs text-slate-400 font-bold uppercase tracking-wider"
							>
								Tingkat Reject
							</p>
							<p class="text-2xl font-bold text-rose-600">{stats.rejectRate}</p>
						</div>
						<div class="text-center">
							<p
								class="text-xs text-slate-400 font-bold uppercase tracking-wider"
							>
								Shift Aktif
							</p>
							<p class="text-2xl font-bold text-blue-600">{stats.activeShift}</p>
						</div>
						<div class="text-center">
							<p
								class="text-xs text-slate-400 font-bold uppercase tracking-wider"
							>
								Total Karyawan
							</p>
							<p class="text-2xl font-bold text-emerald-600">
								{stats.totalEmployees}
								<span class="text-sm font-normal text-slate-500">Orang</span>
							</p>
						</div>
					</div>
				</div>

				<!-- Department Summary Cards -->
				<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
					{#each deptSummary as dept}
						{@const actual = getDeptOutput(dept.id)}
						<div
							class={`rounded-2xl p-6 border transition-all duration-300 hover:shadow-lg hover:-translate-y-1 ${getCardColor(
								dept.color
							)}`}
						>
							<div class="flex justify-between items-start mb-4">
								<div class={`p-3 rounded-xl ${getIconBg(dept.color)}`}>
									<svg
										class="w-6 h-6"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d={dept.icon}
										/>
									</svg>
								</div>
								<div class="text-right">
									<span class="block text-2xl font-bold"
										>{dept.avgEfficiency}%</span
									>
									<span class="text-xs font-semibold opacity-70"
										>Efisiensi Rata-rata</span
									>
								</div>
							</div>

							<h3 class="text-lg font-bold mb-1">{dept.name} Department</h3>
							<p class="text-sm opacity-80 mb-4">
								{dept.activeNow} dari {dept.totalStaff} staff aktif
							</p>

							<div class="space-y-3">
								<div
									class="flex justify-between text-xs font-bold uppercase tracking-wide opacity-60"
								>
									<span>Progress Lot</span>
									<span>{actual} / {dept.target}</span>
								</div>
								<div
									class="w-full h-2 bg-white/50 rounded-full overflow-hidden"
								>
									<div
										class={`h-full rounded-full transition-all duration-1000 ${getBarColor(
											dept.color
										)}`}
										style={`width: ${(actual / dept.target) * 100}%`}
									></div>
								</div>
							</div>
						</div>
					{/each}
				</div>

				<!-- Live Floor Activity Table -->
				<div
					class="bg-white rounded-2xl shadow-sm border border-slate-200 flex flex-col"
				>
					<div
						class="p-6 border-b border-slate-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"
					>
						<div>
							<h3 class="font-bold text-slate-800 text-lg">
								Live Floor Activity
							</h3>
							<p class="text-sm text-slate-500">
								Monitoring real-time aktivitas semua divisi
							</p>
						</div>
						<div class="flex gap-2">
							<select
								class="text-sm border-slate-200 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
							>
								<option>Semua Divisi</option>
								<option>Pressing</option>
								<option>Cutting</option>
								<option>Finishing</option>
							</select>
						</div>
					</div>

					<div class="overflow-x-auto">
						<table class="w-full text-left border-collapse">
							<thead>
								<tr class="bg-slate-50/50 border-b border-slate-100">
									<th
										class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider"
										>User</th
									>
									<th
										class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider"
										>Action</th
									>
									<th
										class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider text-center"
										>Time</th
									>
									<th
										class="px-6 py-4 text-xs font-bold text-slate-400 uppercase tracking-wider text-center"
										>Status</th
									>
								</tr>
							</thead>
							<tbody class="divide-y divide-slate-100">
								{#if stats.activities.length === 0}
									<tr>
										<td
											colspan="4"
											class="text-center py-12 text-slate-500"
										>
											Belum ada aktivitas.
										</td>
									</tr>
								{:else}
									{#each stats.activities as activity}
										<tr class="group hover:bg-slate-50 transition-colors">
											<td class="px-6 py-4">
												<div class="font-bold text-slate-700 text-sm">
													{activity.user}
												</div>
											</td>
											<td class="px-6 py-4">
												<span class="text-slate-700">{activity.action}</span>
											</td>
											<td class="px-6 py-4 text-center">
												<span class="text-slate-500 text-sm">{activity.time}</span>
											</td>
											<td class="px-6 py-4 text-center">
												<span
													class={`inline-flex items-center px-2 py-1 rounded-md text-xs font-bold uppercase ring-1 ring-inset ${getStatusColor(
														activity.status
													)}`}
												>
													{activity.status}
												</span>
											</td>
										</tr>
									{/each}
								{/if}
							</tbody>
						</table>
					</div>

					<div
						class="p-4 border-t border-slate-100 bg-slate-50/50 rounded-b-2xl text-center"
					>
						<button
							class="text-sm font-semibold text-slate-500 hover:text-slate-800 transition-colors"
						>
							Lihat Semua Log Aktivitas
						</button>
					</div>
				</div>
			</main>
		</div>
	{/if}
</div>

<style>
	/* Custom Scrollbar for Table if needed */
	.overflow-x-auto::-webkit-scrollbar {
		height: 8px;
	}
	.overflow-x-auto::-webkit-scrollbar-track {
		background: transparent;
	}
	.overflow-x-auto::-webkit-scrollbar-thumb {
		background-color: #cbd5e1;
		border-radius: 20px;
		border: 3px solid transparent;
		background-clip: content-box;
	}
</style>
