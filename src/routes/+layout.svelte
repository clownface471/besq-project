<script lang="ts">
    import './layout.css';
    import favicon from '$lib/assets/favicon.svg';
    import { auth, logout } from '$lib/stores/auth';
    import { browser } from '$app/environment';
    import { page } from '$app/stores';
    import Swal from 'sweetalert2';
    import { ROLES } from '$lib/roles'; // Import centralized roles

    let { children } = $props();
    
    let sidebarOpen = $state(false);
    let isMobile = $state(false);
    let isChecking = $state(true);
    let isAuthorized = $state(false);

    // --- ACCESS CONTROL CONFIGURATION ---
    const routePermissions: Record<string, string[]> = {
        // Layer 1: Global Overview (Manager Only)
        "/report-chart": [ROLES.ADMIN, ROLES.MANAGER], 

        // Layer 2: Line Overview (Manager & Team Leader)
        // Note: Logic below allows Manager to inherit access, but we list explicit permissions here
        "/report-chart/line": [ROLES.ADMIN, ROLES.MANAGER, ROLES.TEAM_LEADER],

        // Layer 3: Detail (Manager, Team Leader, Operator)
        "/report-chart/detail": [ROLES.ADMIN, ROLES.MANAGER, ROLES.TEAM_LEADER, ROLES.OPERATOR],

        // Other routes
        "/admin": [ROLES.ADMIN],
        "/cutting": [ROLES.ADMIN, "OPERATOR_CUTTING"], 
        "/pressing": [ROLES.ADMIN, "OPERATOR_PRESSING"],
    };

    const hideSidebarRoutes = ['/oprator', '/scan-mesin', '/scan-barcode-prs', '/kpcp-detail'];
    let showSidebar = $derived($auth.isLoggedIn && !hideSidebarRoutes.includes($page.url.pathname));

    $effect(() => {
        if (!browser) return;

        const pathname = $page.url.pathname;
        const user = $auth.user;
        const isLoggedIn = $auth.isLoggedIn;

        if (pathname === '/') {
            if (isLoggedIn && user) {
                redirectBasedOnRole(user.role);
            } else {
                isAuthorized = true;
            }
            isChecking = false;
            return;
        }

        if (!isLoggedIn) {
            window.location.href = '/';
            return;
        }

        if (user) {
            // Strict Hierarchy Check for Report Chart
            if (pathname.startsWith('/report-chart')) {
                 // Check specific sub-paths first (Deepest to Shallowest)
                 if (pathname.includes('/detail')) {
                    checkAccess([ROLES.ADMIN, ROLES.MANAGER, ROLES.TEAM_LEADER, ROLES.OPERATOR], user.role);
                 } else if (pathname.includes('/line')) {
                    checkAccess([ROLES.ADMIN, ROLES.MANAGER, ROLES.TEAM_LEADER], user.role);
                 } else {
                    // Root /report-chart
                    checkAccess([ROLES.ADMIN, ROLES.MANAGER], user.role);
                 }
            } else {
                // Standard check for other routes
                const matchedRoute = Object.keys(routePermissions).find(route => pathname.startsWith(route));
                if (matchedRoute) {
                    checkAccess(routePermissions[matchedRoute], user.role);
                } else {
                    isAuthorized = true;
                }
            }
        }
        isChecking = false;
        isMobile = window.innerWidth < 768;
        if (!isMobile) sidebarOpen = false;
    });

    function checkAccess(allowedRoles: string[], userRole: string) {
        // Helper to handle partial matches (e.g. "OPERATOR_PRESSING" contains "OPERATOR")
        const hasAccess = allowedRoles.includes(userRole) || allowedRoles.some(r => userRole.includes(r));
        
        if (!hasAccess) {
            Swal.fire({
                icon: 'error',
                title: 'Restricted Access',
                text: `Your role (${userRole}) cannot access this layer.`,
                timer: 2000,
                showConfirmButton: false
            }).then(() => {
                redirectBasedOnRole(userRole);
            });
            isAuthorized = false;
        } else {
            isAuthorized = true;
        }
    }

    function redirectBasedOnRole(role: string) {
        if (role === ROLES.ADMIN) window.location.href = '/admin';
        else if (role === ROLES.MANAGER) window.location.href = '/report-chart'; // Layer 1
        else if (role === ROLES.TEAM_LEADER) window.location.href = '/report-chart/line'; // Layer 2
        else if (role.includes(ROLES.OPERATOR)) window.location.href = '/oprator'; // Portal
        else window.location.href = '/';
    }

    function toggleSidebar() { sidebarOpen = !sidebarOpen; }
    function getInitials(name: string): string { return name ? name.split(' ').slice(0, 2).map(n => n.charAt(0).toUpperCase()).join('') : 'U'; }
    async function handleLogoutConfirm() {
        const result = await Swal.fire({
            title: "Logout?", text: "Anda akan keluar dari sesi ini.", icon: "warning",
            showCancelButton: true, confirmButtonColor: '#e11d48', confirmButtonText: "Ya, Keluar", cancelButtonText: "Batal"
        });
        if (result.isConfirmed) logout();
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
                        <i class="fa-solid fa-bars text-slate-700 text-xl"></i>
                    </button>
                    <span class="font-bold text-slate-800">BESQ</span>
                </div>
            </div>
        </div>

        <aside class="fixed left-0 top-0 h-screen w-64 bg-white border-r border-slate-200 z-50 transition-transform duration-300 md:translate-x-0 pt-16 md:pt-0 flex flex-col"
            class:translate-x-0={sidebarOpen} class:-translate-x-full={!sidebarOpen && isMobile}>
            
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

                {#if $auth.user?.role === ROLES.ADMIN || $auth.user?.role === ROLES.MANAGER}
                    <div class="mt-4 mb-2 px-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Reports</div>
                    <a href="/report-chart" class="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-600 hover:bg-slate-50 font-medium"
                       class:text-indigo-600={$page.url.pathname === '/report-chart'} class:bg-indigo-50={$page.url.pathname === '/report-chart'}>
                        <i class="fa-solid fa-chart-pie w-5 text-center"></i>
                        Factory Overview
                    </a>
                {/if}

                {#if [ROLES.ADMIN, ROLES.MANAGER, ROLES.TEAM_LEADER].includes($auth.user?.role || '')}
                     <a href="/report-chart/line" class="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-600 hover:bg-slate-50 font-medium"
                       class:text-indigo-600={$page.url.pathname.includes('/report-chart/line')} class:bg-indigo-50={$page.url.pathname.includes('/report-chart/line')}>
                        <i class="fa-solid fa-industry w-5 text-center"></i>
                        Line Status
                    </a>
                {/if}

                <div class="p-4 border-t border-slate-100 mt-auto">
                    <button onclick={handleLogoutConfirm} class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-rose-600 hover:bg-rose-50 font-medium transition-colors">
                        <i class="fa-solid fa-right-from-bracket w-5 text-center"></i>
                        Keluar
                    </button>
                </div>
            </nav>
        </aside>
    {/if}

    <div class="min-h-screen bg-slate-50 transition-all duration-300" class:md:ml-64={!isMobile && showSidebar} class:pt-16={isMobile && showSidebar}>
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