<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { auth, logout, type User } from '$lib/stores/auth';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import Swal from 'sweetalert2';

	let { children } = $props();
	
	let sidebarOpen = $state(false);
	let isMobile = $state(false);

	const routePermissions: Record<string, string[]> = {
		"/karyawan": ["admin"],
		"/cutting": ["admin", "cutting"],
		"/pressing": ["admin", "pressing"],
	};
	function checkScreenSize() {
		if (typeof window !== 'undefined') {
			isMobile = window.innerWidth < 768;
			if (!isMobile) {
				sidebarOpen = false;
			}
		}
	}

	$effect(() => {
		if (!browser) return;
		
		checkScreenSize();
		
		window.addEventListener('resize', checkScreenSize);
		
		const handleClickOutside = (e: MouseEvent) => {
			const sidebar = document.querySelector('aside');
			const menuButton = document.querySelector('#mobile-menu-button');
			if (sidebarOpen && sidebar && !sidebar.contains(e.target as Node) && 
				menuButton && !menuButton.contains(e.target as Node)) {
				sidebarOpen = false;
			}
		};
		
		document.addEventListener('click', handleClickOutside);
		
		return () => {
			window.removeEventListener('resize', checkScreenSize);
			document.removeEventListener('click', handleClickOutside);
		};
	});

	$effect(() => {
		if (!browser) return;

		const pathname = $page.url.pathname;
		const user = $auth.user;
		const isLoggedIn = $auth.isLoggedIn;

		if (!isLoggedIn && pathname !== '/') {
			goto('/');
			return;
		}

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
				goto('/karyawan');
			}
		}
		
		if (isMobile) {
			sidebarOpen = false;
		}
	});

	async function handleLogout() {
		const result = await Swal.fire({
			title: "Yakin ingin keluar?",
			text: "Anda harus login ulang untuk mengakses aplikasi.",
			icon: "warning",
			showCancelButton: true,
			confirmButtonText: "Ya, Keluar",
			cancelButtonText: "Batal",
		});

		if (result.isConfirmed) {
			logout();
			goto("/");
		}
	}

	function canAccess(path: string): boolean {
		if (!$auth.isLoggedIn || !$auth.user) return false;
		const requiredRoles = routePermissions[path];
		if (!requiredRoles) return true;
		return requiredRoles.includes($auth.user.role);
	}
	
	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	function getInitials(name: string): string {
		return name
			.split(' ')
			.slice(0, 2)
			.map(n => n.charAt(0).toUpperCase())
			.join('');
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
</svelte:head>

{#if $auth.isLoggedIn} 
	<div class="md:hidden fixed top-0 left-0 right-0 bg-white/95 backdrop-blur-xl border-b border-slate-100 shadow-md z-50 px-4 py-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<button
					id="mobile-menu-button"
					onclick={toggleSidebar}
					class="p-2.5 rounded-lg hover:bg-slate-100 active:bg-slate-200 transition-all duration-200"
					aria-label="Toggle menu"
				>
					{#if sidebarOpen}
						<svg class="w-6 h-6 text-slate-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					{:else}
						<svg class="w-6 h-6 text-slate-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
						</svg>
					{/if}
				</button>
				<div>
					<h2 class="text-base font-bold text-slate-800">BESQ</h2>
					<p class="text-xs text-slate-500 font-medium">{$auth.user?.name}</p>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<div class="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center border border-white shadow-sm">
					<span class="text-xs font-bold text-white">
						{getInitials($auth.user?.name || '')}
					</span>
				</div>
			</div>
		</div>
	</div>

	<!-- Sidebar Overlay -->
	{#if sidebarOpen && isMobile}
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div
			class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40 transition-opacity duration-300 md:hidden"
			onclick={() => sidebarOpen = false}
		></div>
	{/if}

	<!-- Sidebar -->
	<aside
		class="fixed left-0 top-0 h-screen w-64 bg-gradient-to-b from-white to-slate-50 border-r border-slate-200 shadow-xl z-40 transition-transform duration-300 ease-in-out md:translate-x-0 flex flex-col"
		class:translate-x-0={sidebarOpen}
		class:-translate-x-full={!sidebarOpen && isMobile}
	>
		<!-- Sidebar Header -->
		<div class="p-6 border-b border-slate-100/50">
			<div class="flex items-center gap-3 mb-1">
				<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center shadow-lg">
					<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
					</svg>
				</div>
				<div>
					<h2 class="text-lg font-bold bg-gradient-to-r from-indigo-600 to-blue-600 bg-clip-text text-transparent">BESQ Portal</h2>
				</div>
			</div>
		</div>

		<!-- User Info Card (Desktop) -->
		<div class="hidden md:block p-4 mx-4 mt-4 bg-gradient-to-br from-indigo-50 to-blue-50 rounded-2xl border border-indigo-100/50">
			<div class="flex items-center gap-3">
				<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center shadow-md border border-white flex-shrink-0">
					<span class="text-sm font-bold text-white">
						{getInitials($auth.user?.name || '')}
					</span>
				</div>
				<div class="flex-1 min-w-0">
					<h3 class="font-bold text-slate-800 truncate text-sm">{$auth.user?.name}</h3>
					<p class="text-xs text-slate-500 truncate">{$auth.user?.nik}</p>
					<p class="text-xs text-indigo-600 font-medium capitalize">{$auth.user?.role}</p>
				</div>
			</div>
		</div>

		<!-- Navigation -->
		<nav class="flex-1 p-4 space-y-1.5 overflow-y-auto">
			{#if $auth.isLoggedIn}
				{#if $auth.user && $auth.user.role === "admin"}
					<a
						href="/admin"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-all duration-200 group relative overflow-hidden"
						class:text-indigo-600={$page.url.pathname === "/admin"}
						class:bg-indigo-50={$page.url.pathname === "/admin"}
					>
						<div class="relative z-10 w-5 h-5">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
	d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
</svg>
						</div>
						<span class="font-semibold text-sm relative z-10">Admin</span>
						{#if $page.url.pathname === "/admin"}
							<span class="absolute inset-0 bg-gradient-to-r from-indigo-500/10 to-transparent -z-0"></span>
							<span class="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-600 shadow-sm"></span>
						{/if}
					</a>
				{/if}
				{#if canAccess("/karyawan")}
					<a
						href="/karyawan"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-all duration-200 group relative overflow-hidden"
						class:text-indigo-600={$page.url.pathname === "/karyawan"}
						class:bg-indigo-50={$page.url.pathname === "/karyawan"}
					>
						<div class="relative z-10 w-5 h-5">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-3m0 0l7-4 7 4M5 9v10a1 1 0 001 1h12a1 1 0 001-1V9m-9 16l-7-4m0 0V5m7 4l7-4m0 0v10a1 1 0 01-1 1H4a1 1 0 01-1-1V9m9 16l7 4m0 0V5" />
							</svg>
						</div>
						<span class="font-semibold text-sm relative z-10">Management karyawan</span>
						{#if $page.url.pathname === "/karyawan"}
							<span class="absolute inset-0 bg-gradient-to-r from-indigo-500/10 to-transparent -z-0"></span>
							<span class="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-600 shadow-sm"></span>
						{/if}
					</a>
				{/if}
				{#if canAccess("/pressing")}
					<a
						href="/pressing"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-all duration-200 group relative overflow-hidden"
						class:text-indigo-600={$page.url.pathname === "/pressing"}
						class:bg-indigo-50={$page.url.pathname === "/pressing"}
					>
						<div class="relative z-10 w-5 h-5">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 20h16M7 20V4c0-.5.5-1 1-1h8c.5 0 1 .5 1 1v16" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 8h6m-3-5v5m-5 4h10v3H7v-3z" />
							</svg>
						</div>
						<span class="font-semibold text-sm relative z-10">Pressing</span>
						{#if $page.url.pathname === "/pressing"}
							<span class="absolute inset-0 bg-gradient-to-r from-indigo-500/10 to-transparent -z-0"></span>
							<span class="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-600 shadow-sm"></span>
						{/if}
					</a>
				{/if}

				{#if canAccess("/cutting")}
					<a
						href="/cutting"
						class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-all duration-200 group relative overflow-hidden"
						class:text-indigo-600={$page.url.pathname === "/cutting"}
						class:bg-indigo-50={$page.url.pathname === "/cutting"}
					>
						<div class="relative z-10 w-5 h-5">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z" />
							</svg>
						</div>
						<span class="font-semibold text-sm relative z-10">Cutting</span>
						{#if $page.url.pathname === "/cutting"}
							<span class="absolute inset-0 bg-gradient-to-r from-indigo-500/10 to-transparent -z-0"></span>
							<span class="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-600 shadow-sm"></span>
						{/if}
					</a>
				{/if}

				<!-- Divider -->
				<div class="my-2 h-px bg-gradient-to-r from-transparent via-slate-200 to-transparent"></div>

				<!-- User Info Card (Mobile) -->
				<div class="md:hidden p-3 bg-gradient-to-br from-slate-100 to-slate-50 rounded-xl border border-slate-200">
					<div class="flex items-center gap-2">
						<div class="w-10 h-10 rounded-lg bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center shadow-md flex-shrink-0 border border-white">
							<span class="text-xs font-bold text-white">
								{getInitials($auth.user?.name || '')}
							</span>
						</div>
						<div class="flex-1 min-w-0">
							<h3 class="font-bold text-slate-800 truncate text-xs">{$auth.user?.name}</h3>
							<p class="text-[10px] text-slate-500 truncate">{$auth.user?.nik}</p>
						</div>
					</div>
				</div>
			{/if}
		</nav>
		<div class="p-4 border-t border-slate-100/50 bg-gradient-to-t from-slate-50 to-transparent">
			<button
				onclick={handleLogout}
				class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-rose-600 hover:text-rose-700 hover:bg-rose-50/80 transition-all duration-200 font-semibold text-sm group"
			>
				<div class="w-5 h-5 flex items-center justify-center rounded-lg">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
					</svg>
				</div>
				<span>Keluar</span>
			</button>
		</div>
	</aside>

	<!-- Main Content -->
	<div class="min-h-screen bg-slate-50 transition-all duration-300"
		class:md:ml-64={!isMobile}
		class:pt-16={isMobile}
	>
		<div class="p-4 sm:p-5 md:p-6">
			{@render children()}
		</div>
	</div>
{:else}
	{@render children()}
{/if}

<style>
	aside nav::-webkit-scrollbar {
		width: 6px;
	}
	aside nav::-webkit-scrollbar-track {
		background: transparent;
	}
	aside nav::-webkit-scrollbar-thumb {
		background: #e2e8f0;
		border-radius: 10px;
	}
	aside nav::-webkit-scrollbar-thumb:hover {
		background: #cbd5e1;
	}

	@media (max-width: 768px) {
		a, button {
			min-height: 44px;
		}
	}
</style>