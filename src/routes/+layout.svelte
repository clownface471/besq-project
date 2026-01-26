<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';


	import { auth, logout, type User } from '$lib/stores/auth';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import Swal from 'sweetalert2';

	export let children;

	const routePermissions: Record<string, string[]> = {
		'/admin': ['admin'],
		'/karyawan': ['admin'],
		'/cutting': ['admin', 'cutting'],
		'/pressing': ['admin', 'pressing']
	};

	$: if (browser) {
		const pathname = $page.url.pathname;
		const user = $auth.user;
		const isLoggedIn = $auth.isLoggedIn;

		// 1. If not logged in and not on the public root page, redirect to root
		if (!isLoggedIn && pathname !== '/') {
			goto('/');
		}

		// 2. If logged in, check role-based permissions
		if (isLoggedIn && user) {
			const requiredRoles = routePermissions[pathname];
			if (requiredRoles && !requiredRoles.includes(user.role)) {
				Swal.fire({
					toast: true,
					position: 'top-end',
					icon: 'error',
					title: 'Anda tidak memiliki akses ke halaman ini',
					showConfirmButton: false,
					timer: 3000,
					timerProgressBar: true
				});
				goto('/'); // Redirect to a safe default page
			}
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

	function canAccess(path: string): boolean {
		if (!$auth.isLoggedIn || !$auth.user) return false;
		const requiredRoles = routePermissions[path];
		if (!requiredRoles) return true; // No specific role required, so allow
		return requiredRoles.includes($auth.user.role);
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

{#if $auth.isLoggedIn}
	<aside class="fixed left-0 top-0 h-full w-64 bg-white border-r border-slate-200 shadow-lg z-40">
		<div class="flex flex-col h-full">
			<!-- Sidebar Header -->
			<div class="p-6 border-b border-slate-200">
				<h2 class="text-xl font-bold text-slate-800">BESQ Portal</h2>
				<p class="text-xs text-slate-500 mt-1">Welcome, {$auth.user?.name}</p>
			</div>

			<!-- Sidebar Navigation -->
			<nav class="flex-1 p-4 space-y-2">
				{#if canAccess('/')}
					<a
						href="/"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors"
						class:bg-slate-100={$page.url.pathname === '/'}
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
							/>
						</svg>
						<span class="font-medium">Dashboard</span>
					</a>
				{/if}
				{#if canAccess('/karyawan')}
					<a
						href="/karyawan"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
							/>
						</svg>
						<span class="font-medium">Karyawan</span>
					</a>
				{/if}
				{#if canAccess('/cutting')}
					<a
						href="/cutting"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z"
							/>
						</svg>
						<span class="font-medium">Cutting</span>
					</a>
				{/if}
			</nav>

			<!-- Sidebar Footer with Logout Button -->
			<div class="p-4 border-t border-slate-200">
				<button
					onclick={handleLogout}
					class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-rose-600 hover:bg-rose-50 transition-colors font-medium"
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
							d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
						/>
					</svg>
					<span>Keluar / Logout</span>
				</button>
			</div>
		</div>
	</aside>

	<!-- Main Content with Sidebar Offset -->
	<div class="ml-64">
		{@render children()}
	</div>
{:else}
	<!-- No Sidebar when not logged in -->
	{@render children()}
{/if}