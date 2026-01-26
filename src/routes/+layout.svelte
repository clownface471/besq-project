<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { auth, logout, type User } from '$lib/stores/auth';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import Swal from 'sweetalert2';

	let { children } = $props();
	
	// State untuk toggle sidebar di mobile
	let sidebarOpen = $state(false);
	let isMobile = $state(false);

	const routePermissions: Record<string, string[]> = {
		"/admin": ["admin"],
		"/karyawan": ["admin"],
		"/cutting": ["admin", "cutting"],
		"/pressing": ["admin", "pressing"],
	};

	// Cek ukuran layar
	function checkScreenSize() {
		if (typeof window !== 'undefined') {
			isMobile = window.innerWidth < 768;
			// Auto close sidebar di mobile saat resize
			if (!isMobile) {
				sidebarOpen = false;
			}
		}
	}

	$effect(() => {
		if (!browser) return;
		
		// Inisialisasi check screen size
		checkScreenSize();
		
		// Event listener untuk resize
		window.addEventListener('resize', checkScreenSize);
		
		// Auto close sidebar saat klik di luar di mobile
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

		// 1. If not logged in and not on the public root page, redirect to root
		if (!isLoggedIn && pathname !== '/') {
			goto('/');
			return;
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
		
		// Close sidebar saat pindah halaman di mobile
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
		if (!requiredRoles) return true; // No specific role required, so allow
		return requiredRoles.includes($auth.user.role);
	}
	
	// Toggle sidebar untuk mobile
	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
</svelte:head>

{#if $auth.isLoggedIn}
	<!-- Mobile Header -->
	<div class="md:hidden fixed top-0 left-0 right-0 bg-white border-b border-slate-200 shadow-sm z-50 px-4 py-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<button
					id="mobile-menu-button"
					onclick={toggleSidebar}
					class="p-2 rounded-lg hover:bg-slate-100 transition-colors"
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
					<h2 class="text-lg font-bold text-slate-800">BESQ Portal</h2>
					<p class="text-xs text-slate-500">{$auth.user?.name}</p>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<img
					src="https://i.pravatar.cc/150?u={$auth.user?.nik}"
					alt="User"
					class="w-8 h-8 rounded-full border-2 border-white shadow-sm"
				/>
			</div>
		</div>
	</div>

	<!-- Sidebar Overlay (Mobile only) -->
	{#if sidebarOpen && isMobile}
		<div
			class="fixed inset-0 bg-black/50 z-40 transition-opacity md:hidden"
			onclick={() => sidebarOpen = false}
		></div>
	{/if}

	<!-- Sidebar -->
	<aside
		class="fixed left-0 top-0 h-full w-64 bg-white border-r border-slate-200 shadow-lg z-40 transition-transform duration-300 ease-in-out md:translate-x-0"
		class:translate-x-0={sidebarOpen}
		class:-translate-x-full={!sidebarOpen && isMobile}
	>
		<!-- Sidebar Header (Desktop) -->
		<div class="hidden md:block p-6 border-b border-slate-200">
			<h2 class="text-xl font-bold text-slate-800">BESQ Portal</h2>
			<p class="text-xs text-slate-500 mt-1">Welcome, {$auth.user?.name}</p>
		</div>
		
		<!-- Sidebar Header (Mobile) -->
		<div class="md:hidden p-6 border-b border-slate-200 flex items-center gap-3">
			<img
				src="https://i.pravatar.cc/150?u={$auth.user?.nik}"
				alt="User"
				class="w-10 h-10 rounded-full border-2 border-white shadow-sm"
			/>
			<div>
				<h2 class="text-lg font-bold text-slate-800">BESQ Portal</h2>
				<p class="text-xs text-slate-500">{$auth.user?.name}</p>
				<p class="text-xs text-slate-400">{$auth.user?.role} • {$auth.user?.nik}</p>
			</div>
		</div>

		<!-- Sidebar Navigation -->
		<nav class="flex-1 p-4 space-y-1 overflow-y-auto max-h-[calc(100vh-200px)]">
			{#if canAccess("/")}
				<a
					href="/"
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors group"
					class:bg-blue-50={$page.url.pathname === "/"}
					class:text-blue-600={$page.url.pathname === "/"}
					class:border-l-4={$page.url.pathname === "/"}
					class:border-l-blue-500={$page.url.pathname === "/"}
					class:pl-3={$page.url.pathname === "/"}
				>
					<div class="w-8 h-8 flex items-center justify-center rounded-lg bg-slate-100 group-hover:bg-blue-100 transition-colors"
						class:bg-blue-100={$page.url.pathname === "/"}
						class:text-blue-600={$page.url.pathname === "/"}
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
							/>
						</svg>
					</div>
					<span class="font-medium">Dashboard</span>
					{#if $page.url.pathname === "/"}
						<span class="ml-auto w-2 h-2 rounded-full bg-blue-500"></span>
					{/if}
				</a>
			{/if}
			{#if canAccess("/karyawan")}
				<a
					href="/karyawan"
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors group"
					class:bg-blue-50={$page.url.pathname === "/karyawan"}
					class:text-blue-600={$page.url.pathname === "/karyawan"}
					class:border-l-4={$page.url.pathname === "/karyawan"}
					class:border-l-blue-500={$page.url.pathname === "/karyawan"}
					class:pl-3={$page.url.pathname === "/karyawan"}
				>
					<div class="w-8 h-8 flex items-center justify-center rounded-lg bg-slate-100 group-hover:bg-blue-100 transition-colors"
						class:bg-blue-100={$page.url.pathname === "/karyawan"}
						class:text-blue-600={$page.url.pathname === "/karyawan"}
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
							/>
						</svg>
					</div>
					<span class="font-medium">Karyawan</span>
					{#if $page.url.pathname === "/karyawan"}
						<span class="ml-auto w-2 h-2 rounded-full bg-blue-500"></span>
					{/if}
				</a>
			{/if}
			{#if canAccess("/cutting")}
				<a
					href="/cutting"
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors group"
					class:bg-blue-50={$page.url.pathname === "/cutting"}
					class:text-blue-600={$page.url.pathname === "/cutting"}
					class:border-l-4={$page.url.pathname === "/cutting"}
					class:border-l-blue-500={$page.url.pathname === "/cutting"}
					class:pl-3={$page.url.pathname === "/cutting"}
				>
					<div class="w-8 h-8 flex items-center justify-center rounded-lg bg-slate-100 group-hover:bg-blue-100 transition-colors"
						class:bg-blue-100={$page.url.pathname === "/cutting"}
						class:text-blue-600={$page.url.pathname === "/cutting"}
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M14.121 14.121L19 19m-7-7l7-7m-7 7l-2.879 2.879M12 12L9.121 9.121m0 5.758a3 3 0 10-4.243 4.243 3 3 0 004.243-4.243zm0-5.758a3 3 0 10-4.243-4.243 3 3 0 004.243 4.243z"
							/>
						</svg>
					</div>
					<span class="font-medium">Cutting</span>
					{#if $page.url.pathname === "/cutting"}
						<span class="ml-auto w-2 h-2 rounded-full bg-blue-500"></span>
					{/if}
				</a>
			{/if}
			{#if canAccess("/pressing")}
				<a
					href="/pressing"
					class="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-700 hover:bg-slate-100 transition-colors group"
					class:bg-blue-50={$page.url.pathname === "/pressing"}
					class:text-blue-600={$page.url.pathname === "/pressing"}
					class:border-l-4={$page.url.pathname === "/pressing"}
					class:border-l-blue-500={$page.url.pathname === "/pressing"}
					class:pl-3={$page.url.pathname === "/pressing"}
				>
					<div class="w-8 h-8 flex items-center justify-center rounded-lg bg-slate-100 group-hover:bg-blue-100 transition-colors"
						class:bg-blue-100={$page.url.pathname === "/pressing"}
						class:text-blue-600={$page.url.pathname === "/pressing"}
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 20h16M7 20V4c0-.5.5-1 1-1h8c.5 0 1 .5 1 1v16"
							/>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 8h6m-3-5v5m-5 4h10v3H7v-3z"
							/>
						</svg>
					</div>
					<span class="font-medium">Pressing</span>
					{#if $page.url.pathname === "/pressing"}
						<span class="ml-auto w-2 h-2 rounded-full bg-blue-500"></span>
					{/if}
				</a>
			{/if}
			
			<!-- User Info Card (Mobile) -->
			<div class="md:hidden mt-4 p-4 bg-slate-50 rounded-xl">
				<div class="flex items-center gap-3">
					<img
						src="https://i.pravatar.cc/150?u={$auth.user?.nik}"
						alt="User"
						class="w-12 h-12 rounded-full border-2 border-white shadow-sm"
					/>
					<div class="flex-1 min-w-0">
						<h3 class="font-bold text-slate-800 truncate">{$auth.user?.name}</h3>
						<p class="text-xs text-slate-500 truncate">{$auth.user?.nik}</p>
						<p class="text-xs text-slate-400 capitalize">{$auth.user?.role}</p>
					</div>
				</div>
			</div>
		</nav>

		<!-- Sidebar Footer with Logout Button -->
		<div class="p-4 border-t border-slate-200">
			<button
				onclick={handleLogout}
				class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-rose-600 hover:bg-rose-50 transition-colors font-medium group"
			>
				<div class="w-8 h-8 flex items-center justify-center rounded-lg bg-rose-100 group-hover:bg-rose-200 transition-colors">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
						/>
					</svg>
				</div>
				<span>Keluar / Logout</span>
			</button>
		</div>
	</aside>

	<!-- Main Content with Responsive Offset -->
	<div class="min-h-screen bg-slate-50 transition-all duration-300"
		class:md:ml-64={!isMobile}
		class:pt-16={isMobile}
	>
		<div class="p-4 sm:p-5 md:p-6">
			{@render children()}
		</div>
	</div>
{:else}
	<!-- No Sidebar when not logged in -->
	{@render children()}
{/if}

<style>
	/* Custom Scrollbar untuk Sidebar */
	aside nav::-webkit-scrollbar {
		width: 4px;
	}
	aside nav::-webkit-scrollbar-track {
		background: transparent;
	}
	aside nav::-webkit-scrollbar-thumb {
		background: #cbd5e1;
		border-radius: 10px;
	}
	aside nav::-webkit-scrollbar-thumb:hover {
		background: #94a3b8;
	}
	
	/* Smooth transition untuk semua */
	* {
		transition: background-color 0.2s ease, border-color 0.2s ease;
	}
	
	/* Animasi untuk mobile menu */
	@keyframes slideIn {
		from {
			transform: translateX(-100%);
		}
		to {
			transform: translateX(0);
		}
	}
	
	@keyframes slideOut {
		from {
			transform: translateX(0);
		}
		to {
			transform: translateX(-100%);
		}
	}
	
	/* Touch-friendly untuk mobile */
	@media (max-width: 768px) {
		a, button {
			min-height: 44px;
		}
	}
</style>