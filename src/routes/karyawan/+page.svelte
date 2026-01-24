<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { get } from 'svelte/store';

	interface User {
		id: number;
		nik: string;
		name: string;
		role: string;
		created_at: string;
		updated_at: string;
	}

	const API_URL = 'http://localhost:8080';

	let users = $state<User[]>([]);
	let isLoading = $state(false);
	let errorMessage = $state('');
	let showModal = $state(false);
	let isEditing = $state(false);
	let editingUser = $state<User | null>(null);

	// Form state using Svelte 5 $state rune for reactivity
	let formData = $state({
		id: 0,
		nik: '',
		name: '',
		password: '',
		role: 'cutting' as 'admin' | 'cutting' | 'pressing'
	});

	// Check admin role on mount
	onMount(() => {
		if (browser) {
			const authState = get(auth);
			if (!authState.user || authState.user.role !== 'admin') {
				goto('/');
			} else {
				fetchUsers();
			}
		}
	});

	// Fetch users from API
	async function fetchUsers() {
		const authToken = get(auth).token;
		if (!authToken) return;

		isLoading = true;
		errorMessage = '';

		try {
			const response = await fetch(`${API_URL}/api/users/`, {
				headers: {
					Authorization: `Bearer ${authToken}`
				}
			});

			if (!response.ok) {
				throw new Error('Failed to fetch users');
			}

			const data = await response.json();
			users = data.users || [];
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'An error occurred';
			console.error('Error fetching users:', error);
		} finally {
			isLoading = false;
		}
	}

	// Open modal for adding new user
	function openAddModal() {
		isEditing = false;
		editingUser = null;
		
		formData.id = 0;
		formData.nik = '';
		formData.name = '';
		formData.password = '';
		formData.role = 'cutting';

		showModal = true;
		errorMessage = '';
	}

	// Open modal for editing user
	function openEditModal(user: User) {
		isEditing = true;
		editingUser = user;

		formData.id = user.id;
		formData.nik = user.nik;
		formData.name = user.name;
		formData.password = ''; // Password is optional when editing
		formData.role = user.role as 'admin' | 'cutting' | 'pressing';

		showModal = true;
		errorMessage = '';
	}

	// Close modal
	function closeModal() {
		showModal = false;
		isEditing = false;
		editingUser = null;

		formData.id = 0;
		formData.nik = '';
		formData.name = '';
		formData.password = '';
		formData.role = 'cutting';
		
		errorMessage = '';
	}

	// Handle form submission (Add or Edit)
	async function handleSubmit() {
		const authToken = get(auth).token;
		if (!authToken) return;

		errorMessage = '';

		// Validation - read from reactive $state object
		if (!formData.nik.trim() || !formData.name.trim()) {
			errorMessage = 'NIK and Name are required';
			return;
		}

		if (!isEditing && !formData.password.trim()) {
			errorMessage = 'Password is required for new users';
			return;
		}

		try {
			const url = isEditing
				? `${API_URL}/api/users/${editingUser?.id}`
				: `${API_URL}/api/users/`;

			const method = isEditing ? 'PUT' : 'POST';

			const body: any = {
				nik: formData.nik.trim(),
				name: formData.name.trim(),
				role: formData.role
			};

			// Only include password if provided (for edit) or always (for add)
			if (isEditing && formData.password.trim()) {
				body.password = formData.password.trim();
			} else if (!isEditing) {
				body.password = formData.password.trim();
			}

			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${authToken}`
				},
				body: JSON.stringify(body)
			});

			console.log('Sending Payload:', body);

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Operation failed');
			}

			// Success - close modal and refetch
			closeModal();
			await fetchUsers();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'An error occurred';
			console.error('Error saving user:', error);
		}
	}

	// Handle delete user
	async function handleDelete(user: User) {
		const authToken = get(auth).token;
		if (!authToken) return;

		const confirmed = confirm(`Are you sure you want to delete user "${user.name}" (NIK: ${user.nik})?`);
		if (!confirmed) return;

		try {
			const response = await fetch(`${API_URL}/api/users/${user.id}`, {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${authToken}`
				}
			});

			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error || 'Failed to delete user');
			}

			// Success - refetch users
			await fetchUsers();
		} catch (error) {
			alert(error instanceof Error ? error.message : 'Failed to delete user');
			console.error('Error deleting user:', error);
		}
	}

	// Get role badge color
	function getRoleBadgeColor(role: string): string {
		switch (role) {
			case 'admin':
				return 'bg-purple-100 text-purple-700 border-purple-200';
			case 'cutting':
				return 'bg-blue-100 text-blue-700 border-blue-200';
			case 'pressing':
				return 'bg-amber-100 text-amber-700 border-amber-200';
			default:
				return 'bg-slate-100 text-slate-700 border-slate-200';
		}
	}
</script>

<div class="min-h-screen bg-slate-50 font-sans text-slate-800 pb-12 relative">
	<div
		class="absolute inset-0 z-0 opacity-[0.03] pointer-events-none"
		style="background-image: radial-gradient(#4f46e5 1px, transparent 1px); background-size: 24px 24px;"
	></div>

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative z-10">
		<!-- Header -->
		<div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6 mb-6">
			<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
				<div>
					<h1 class="text-2xl md:text-3xl font-bold text-slate-800 tracking-tight">
						Manajemen Karyawan
					</h1>
					<p class="text-sm text-slate-500 mt-1">Kelola data karyawan dan pengguna sistem</p>
				</div>
				<button
					onclick={openAddModal}
					class="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-3 rounded-xl shadow-lg shadow-indigo-500/30 hover:shadow-indigo-500/50 transition-all duration-300 hover:-translate-y-0.5 active:scale-95 flex items-center gap-2"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
					<span>Tambah Karyawan</span>
				</button>
			</div>
		</div>

		<!-- Error Message -->
		{#if errorMessage && !showModal}
			<div
				class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl mb-6 flex items-center justify-between"
			>
				<span>{errorMessage}</span>
				<button
					onclick={() => (errorMessage = '')}
					class="text-red-500 hover:text-red-700"
					aria-label="Close"
				>
					<svg
						class="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>
		{/if}

		<!-- Loading Skeleton -->
		{#if isLoading}
			<div class="bg-white rounded-2xl shadow-sm border border-slate-100 p-6">
				<div class="animate-pulse space-y-4">
					<div class="h-4 bg-slate-200 rounded w-3/4"></div>
					<div class="h-4 bg-slate-200 rounded w-1/2"></div>
					<div class="h-4 bg-slate-200 rounded w-5/6"></div>
				</div>
			</div>
		{:else}
			<!-- Users Table -->
			<div class="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead class="bg-slate-50 border-b border-slate-200">
							<tr>
								<th
									class="px-6 py-4 text-left text-xs font-bold text-slate-600 uppercase tracking-wider"
								>
									NIK
								</th>
								<th
									class="px-6 py-4 text-left text-xs font-bold text-slate-600 uppercase tracking-wider"
								>
									Name
								</th>
								<th
									class="px-6 py-4 text-left text-xs font-bold text-slate-600 uppercase tracking-wider"
								>
									Role
								</th>
								<th
									class="px-6 py-4 text-right text-xs font-bold text-slate-600 uppercase tracking-wider"
								>
									Actions
								</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#if users.length === 0}
								<tr>
									<td
										colspan="4"
										class="px-6 py-12 text-center text-slate-500"
									>
										No users found. Click "Tambah Karyawan" to add one.
									</td>
								</tr>
							{:else}
								{#each users as user (user.id)}
									<tr class="hover:bg-slate-50 transition-colors">
										<td class="px-6 py-4 whitespace-nowrap">
											<span class="font-mono font-medium text-slate-800">{user.nik}</span>
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<span class="text-slate-800 font-medium">{user.name}</span>
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<span
												class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold border {getRoleBadgeColor(user.role)}"
											>
												{user.role}
											</span>
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-right">
											<div class="flex items-center justify-end gap-2">
												<button
													onclick={() => openEditModal(user)}
													class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
													title="Edit"
													aria-label="Edit"
												>
													<svg
														class="w-5 h-5"
														fill="none"
														stroke="currentColor"
														viewBox="0 0 24 24"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
														/>
													</svg>
												</button>
												<button
													onclick={() => handleDelete(user)}
													class="p-2 text-rose-600 hover:bg-rose-50 rounded-lg transition-colors"
													title="Delete"
													aria-label="Delete"
												>
													<svg
														class="w-5 h-5"
														fill="none"
														stroke="currentColor"
														viewBox="0 0 24 24"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
														/>
													</svg>
												</button>
											</div>
										</td>
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</div>
		{/if}
	</main>
</div>

<!-- Add/Edit Modal -->
{#if showModal}
	<div
		class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
		onclick={closeModal}
		role="button"
		tabindex="0"
		onkeydown={(e) => {
			if (e.key === 'Escape') closeModal();
		}}
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<div
			class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
			aria-modal="true"
			tabindex="-1"
		>
			<div class="flex justify-between items-center mb-6">
				<h2 class="text-2xl font-bold text-slate-800">
					{isEditing ? 'Edit Karyawan' : 'Tambah Karyawan'}
				</h2>
				<button
					onclick={closeModal}
					class="text-slate-400 hover:text-slate-600 transition-colors"
					aria-label="Close"
				>
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
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>

			{#if errorMessage}
				<div
					class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl mb-4 text-sm"
				>
					{errorMessage}
				</div>
			{/if}

			<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
				<div>
					<label
						for="nik"
						class="block text-sm font-semibold text-slate-700 mb-2"
					>
						NIK <span class="text-rose-500">*</span>
					</label>
					<input
						id="nik"
						type="text"
						bind:value={formData.nik}
						required
						disabled={isEditing}
						class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all disabled:bg-slate-50 disabled:text-slate-500"
						placeholder="Enter NIK"
					/>
				</div>

				<div>
					<label
						for="name"
						class="block text-sm font-semibold text-slate-700 mb-2"
					>
						Name <span class="text-rose-500">*</span>
					</label>
					<input
						id="name"
						type="text"
						bind:value={formData.name}
						required
						class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
						placeholder="Enter name"
					/>
				</div>

				<div>
					<label
						for="role"
						class="block text-sm font-semibold text-slate-700 mb-2"
					>
						Role <span class="text-rose-500">*</span>
					</label>
					<select
						id="role"
						bind:value={formData.role}
						required
						class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all bg-white"
					>
						<option value="admin">Admin</option>
						<option value="cutting">Cutting</option>
						<option value="pressing">Pressing</option>
					</select>
				</div>

				<div>
					<label
						for="password"
						class="block text-sm font-semibold text-slate-700 mb-2"
					>
						Password {#if isEditing}
							<span class="text-slate-400 text-xs font-normal">(leave blank to keep current)</span>
						{:else}
							<span class="text-rose-500">*</span>
						{/if}
					</label>
					<input
						id="password"
						type="password"
						bind:value={formData.password}
						required={!isEditing}
						class="w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
						placeholder={isEditing ? 'Enter new password (optional)' : 'Enter password'}
					/>
				</div>

				<div class="flex gap-3 pt-4">
					<button
						type="button"
						onclick={closeModal}
						class="flex-1 px-4 py-3 border border-slate-200 text-slate-700 font-semibold rounded-xl hover:bg-slate-50 transition-colors"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="flex-1 px-4 py-3 bg-indigo-600 text-white font-semibold rounded-xl hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-500/30"
					>
						{isEditing ? 'Update' : 'Create'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}