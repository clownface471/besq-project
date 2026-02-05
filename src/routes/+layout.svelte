<script lang="ts">
    import './layout.css';
    import favicon from '$lib/assets/favicon.svg';
    import { auth, logout } from '$lib/stores/auth';
    import { browser } from '$app/environment';
    import { page } from '$app/stores';
    import Swal from 'sweetalert2';

    let { children } = $props();
    
    let sidebarOpen = $state(false);
    let isMobile = $state(false);
    let isChecking = $state(true); // State untuk loading awal (mencegah flash content)
    let isAuthorized = $state(false);

    // KONFIGURASI HAK AKSES (Centralized Config)
    // Key: Prefix URL, Value: Array Role yang diizinkan
    const routePermissions: Record<string, string[]> = {
        "/admin": ["ADMIN"],
        "/karyawan": ["ADMIN"],
        "/cutting": ["ADMIN", "OPERATOR_CUTTING"],
        "/pressing": ["ADMIN", "OPERATOR_PRESSING"],
        "/scan-mesin": ["ADMIN", "OPERATOR_PRESSING", "OPERATOR_CUTTING"],
        "/scan-barcode-prs": ["ADMIN", "OPERATOR_PRESSING"],
        "/kpcp-detail": ["ADMIN", "OPERATOR_PRESSING"],
        "/barcode-ctng": ["ADMIN", "OPERATOR_CUTTING"],
    };

    // Halaman yang tidak butuh sidebar
    const hideSidebarRoutes = ['/oprator', '/scan-mesin', '/scan-barcode-prs', '/kpcp-detail']; 
    let showSidebar = $derived($auth.isLoggedIn && !hideSidebarRoutes.includes($page.url.pathname));

    // === LOGIKA UTAMA SECURITY CHECK ===
    $effect(() => {
        if (!browser) return;

        const pathname = $page.url.pathname;
        const user = $auth.user;
        const isLoggedIn = $auth.isLoggedIn;

        // 1. Jika di Halaman Login (root), biarkan
        if (pathname === '/') {
            // Jika sudah login, lempar ke dashboard masing-masing (UX friendly)
            if (isLoggedIn && user) {
                redirectBasedOnRole(user.role);
            } else {
                isAuthorized = true; // Izinkan render halaman login
            }
            isChecking = false;
            return;
        }

        // 2. Jika Halaman Lain tapi BELUM Login -> TENDANG KELUAR
        if (!isLoggedIn) {
            window.location.href = '/';
            return;
        }

        // 3. Cek Permission Role
        if (user) {
            // Cari aturan permission yang cocok dengan URL saat ini
            const matchedRoute = Object.keys(routePermissions).find(route => pathname.startsWith(route));
            
            if (matchedRoute) {
                const requiredRoles = routePermissions[matchedRoute];
                // Jika role user TIDAK ada di daftar yang diizinkan
                if (!requiredRoles.includes(user.role)) {
                    // BLOCK AKSES!
                    Swal.fire({
                        icon: 'error',
                        title: 'Akses Ditolak',
                        text: `Anda tidak memiliki izin untuk mengakses halaman ini.`,
                        timer: 2000,
                        showConfirmButton: false
                    }).then(() => {
                        // Kembalikan ke jalan yang benar
                        redirectBasedOnRole(user.role);
                    });
                    isAuthorized = false; // Jangan render konten
                } else {
                    isAuthorized = true; // Boleh lewat
                }
            } else {
                // Jika route tidak ada di daftar 'routePermissions', anggap AMAN (atau mau diblokir default?)
                // Saya set AMAN untuk halaman umum lain (jika ada)
                isAuthorized = true; 
            }
        }

        isChecking = false; // Selesai pengecekan
        
        // Handle Mobile Sidebar
        isMobile = window.innerWidth < 768;
        if (!isMobile) sidebarOpen = false;
    });

function redirectBasedOnRole(role: string) {
        if (role === 'ADMIN') {
            window.location.href = '/admin'; // Admin tetap ke dashboard admin
        } else if (role.startsWith('OPERATOR_')) {
            window.location.href = '/oprator'; // SEMUA Operator masuk ke portal dulu!
        } else {
            window.location.href = '/';
        }
    }

    function toggleSidebar() {
        sidebarOpen = !sidebarOpen;
    }

    function getInitials(name: string): string {
        return name ? name.split(' ').slice(0, 2).map(n => n.charAt(0).toUpperCase()).join('') : 'U';
    }
    
    // Helper untuk Link Sidebar (Hanya tampilkan jika boleh akses)
    function canAccess(path: string): boolean {
        if (!$auth.user) return false;
        const requiredRoles = routePermissions[path];
        if (!requiredRoles) return true;
        return requiredRoles.includes($auth.user.role);
    }
    
    async function handleLogoutConfirm() {
        const result = await Swal.fire({
            title: "Logout?",
            text: "Anda akan keluar dari sesi ini.",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: '#e11d48', // Rose-600
            confirmButtonText: "Ya, Keluar",
            cancelButtonText: "Batal"
        });

        if (result.isConfirmed) {
            logout();
        }
    }
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
</svelte:head>

{#if isChecking}
    <div class="min-h-screen flex items-center justify-center bg-slate-50">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

{:else if isAuthorized}
    {#if showSidebar} 
        <div class="md:hidden fixed top-0 left-0 right-0 bg-white/95 backdrop-blur-xl border-b border-slate-100 shadow-md z-50 px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                    <button onclick={toggleSidebar} class="p-2.5 rounded-lg hover:bg-slate-100">
                        <svg class="w-6 h-6 text-slate-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                        </svg>
                    </button>
                    <span class="font-bold text-slate-800">BESQ</span>
                </div>
                <div class="w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-white text-xs font-bold">
                    {getInitials($auth.user?.name || '')}
                </div>
            </div>
        </div>

        {#if sidebarOpen && isMobile}
            <div class="fixed inset-0 bg-black/40 z-40 md:hidden" onclick={() => sidebarOpen = false}></div>
        {/if}

        <aside class="fixed left-0 top-0 h-screen w-64 bg-white border-r border-slate-200 z-50 transition-transform duration-300 md:translate-x-0 pt-16 md:pt-0 flex flex-col"
            class:translate-x-0={sidebarOpen}
            class:-translate-x-full={!sidebarOpen && isMobile}
        >
            <div class="hidden md:flex p-6 border-b border-slate-100 items-center gap-3">
                <div class="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center text-white font-bold">B</div>
                <span class="font-bold text-lg text-slate-800">BESQ Portal</span>
            </div>

            <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
                <a href="/" onclick={(e) => { e.preventDefault(); redirectBasedOnRole($auth.user?.role || ''); }}
                   class="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-600 hover:bg-slate-50 font-medium">
                    <i class="fa-solid fa-house w-5 text-center"></i>
                    Dashboard
                </a>

                {#if canAccess("/admin")}
                    <div class="mt-4 mb-2 px-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Admin</div>
                    <a href="/admin" class="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-600 hover:bg-slate-50 font-medium"
                       class:text-indigo-600={$page.url.pathname.startsWith('/admin')} class:bg-indigo-50={$page.url.pathname.startsWith('/admin')}>
                        <i class="fa-solid fa-users-gear w-5 text-center"></i>
                        Management User
                    </a>
                {/if}

                <div class="mt-4 mb-2 px-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Produksi</div>
                
                {#if canAccess("/pressing")}
                    <a href="/pressing" class="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-600 hover:bg-slate-50 font-medium"
                       class:text-indigo-600={$page.url.pathname.startsWith('/pressing')} class:bg-indigo-50={$page.url.pathname.startsWith('/pressing')}>
                        <i class="fa-solid fa-compress w-5 text-center"></i>
                        Station Pressing
                    </a>
                {/if}

                {#if canAccess("/cutting")}
                    <a href="/cutting" class="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-600 hover:bg-slate-50 font-medium"
                       class:text-indigo-600={$page.url.pathname.startsWith('/cutting')} class:bg-indigo-50={$page.url.pathname.startsWith('/cutting')}>
                        <i class="fa-solid fa-scissors w-5 text-center"></i>
                        Station Cutting
                    </a>
                {/if}
            </nav>

            <div class="p-4 border-t border-slate-100">
                <button onclick={handleLogoutConfirm} class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-rose-600 hover:bg-rose-50 font-medium transition-colors">
                    <i class="fa-solid fa-right-from-bracket w-5 text-center"></i>
                    Keluar
                </button>
            </div>
        </aside>
    {/if}

    <div class="min-h-screen bg-slate-50 transition-all duration-300"
         class:md:ml-64={!isMobile && showSidebar} 
         class:pt-16={isMobile && showSidebar}
    >
        <div class={showSidebar ? "p-4 sm:p-6" : "p-0"}>
            {@render children()}
        </div>
    </div>

{:else}
    <div class="min-h-screen bg-slate-50"></div>
{/if}

<style>
    /* Styling scrollbar halus untuk sidebar */
    nav::-webkit-scrollbar { width: 4px; }
    nav::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
</style>