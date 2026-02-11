<script lang="ts">
    import { onMount } from 'svelte';
    import { auth } from '$lib/stores/auth';
    import Swal from 'sweetalert2';

    const API_URL = 'http://localhost:8080';
    let users: any[] = [];
    let isLoading = false;

    // Form State
    let isEditing = false;
    let showModal = false;
    let formData = { id: 0, username: '', password: '', role: 'OPERATOR_PRESSING' };

    async function fetchUsers() {
        isLoading = true;
        try {
            const res = await fetch(`${API_URL}/admin/users`, {
                headers: { Authorization: `Bearer ${$auth.token}` }
            });
            const result = await res.json();
            if (res.ok) users = result.data;
        } catch (error) {
            console.error(error);
        } finally {
            isLoading = false;
        }
    }

    function openModal(user: any = null) {
        if (user) {
            isEditing = true;
            formData = { ...user, password: '' }; // Password kosong saat edit
        } else {
            isEditing = false;
            formData = { id: 0, username: '', password: '', role: 'OPERATOR_PRESSING' };
        }
        showModal = true;
    }

    async function handleSubmit() {
        const endpoint = isEditing ? `${API_URL}/admin/users/${formData.id}` : `${API_URL}/admin/add-operator`;
        const method = isEditing ? 'PUT' : 'POST';

        try {
            const res = await fetch(endpoint, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${$auth.token}`
                },
                body: JSON.stringify(formData)
            });

            if (!res.ok) throw new Error("Gagal menyimpan data");

            Swal.fire('Sukses', 'Data berhasil disimpan', 'success');
            showModal = false;
            fetchUsers();
        } catch (err) {
            Swal.fire('Error', 'Terjadi kesalahan', 'error');
        }
    }

    async function handleDelete(id: number) {
        const result = await Swal.fire({
            title: 'Hapus User?',
            text: "Data tidak bisa dikembalikan!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Ya, Hapus'
        });

        if (result.isConfirmed) {
            await fetch(`${API_URL}/admin/users/${id}`, {
                method: 'DELETE',
                headers: { Authorization: `Bearer ${$auth.token}` }
            });
            fetchUsers();
            Swal.fire('Terhapus!', '', 'success');
        }
    }

    onMount(() => {
        fetchUsers();
    });
</script>

<div class="p-6 max-w-6xl mx-auto">
    <div class="flex justify-between items-center mb-6">
        <div>
            <h1 class="text-2xl font-bold text-slate-800">Manajemen Operator</h1>
            <p class="text-slate-500 text-sm">Kelola akun dan hak akses karyawan</p>
        </div>
        <button on:click={() => openModal()} class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2">
            <i class="fa-solid fa-plus"></i> Tambah Operator
        </button>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
                <thead class="bg-slate-50 border-b border-slate-200">
                    <tr>
                        <th class="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Username / NIK</th>
                        <th class="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Role</th>
                        <th class="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">Aksi</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                    {#each users as user}
                        <tr class="hover:bg-slate-50 transition-colors">
                            <td class="px-6 py-4 font-medium text-slate-800">{user.username}</td>
                            <td class="px-6 py-4">
                                <span class={`px-2 py-1 rounded-md text-xs font-bold uppercase 
                                    ${user.role === 'ADMIN' ? 'bg-purple-100 text-purple-700' : 
                                      user.role.includes('CUTTING') ? 'bg-blue-100 text-blue-700' : 
                                      'bg-emerald-100 text-emerald-700'}`}>
                                    {user.role}
                                </span>
                            </td>
                            <td class="px-6 py-4 text-right space-x-2">
                                <button on:click={() => openModal(user)} class="text-indigo-600 hover:text-indigo-800 font-medium text-sm">Edit</button>
                                <button on:click={() => handleDelete(user.id)} class="text-rose-600 hover:text-rose-800 font-medium text-sm">Hapus</button>
                            </td>
                        </tr>
                    {/each}
                    {#if users.length === 0 && !isLoading}
                        <tr><td colspan="3" class="text-center py-8 text-slate-400">Belum ada data user.</td></tr>
                    {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>

{#if showModal}
    <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
            <h2 class="text-xl font-bold mb-4 text-slate-800">{isEditing ? 'Edit User' : 'Tambah User Baru'}</h2>
            
            <div class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Username / NIK</label>
                    <input type="text" bind:value={formData.username} class="w-full border border-slate-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-indigo-500 outline-none" placeholder="Contoh: OP-001">
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Password</label>
                    <input type="password" bind:value={formData.password} class="w-full border border-slate-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-indigo-500 outline-none" placeholder={isEditing ? "(Biarkan kosong jika tidak ubah)" : "Password..."}>
                </div>

                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Role / Divisi</label>
                    <select bind:value={formData.role} class="w-full border border-slate-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-indigo-500 outline-none bg-white">
                        <option value="OPERATOR_PRESSING">OPERATOR PRESSING</option>
                        <option value="OPERATOR_CUTTING">OPERATOR CUTTING</option>
                        <option value="ADMIN">ADMIN</option>
                    </select>
                </div>
            </div>

            <div class="flex justify-end gap-3 mt-6">
                <button on:click={() => showModal = false} class="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg font-medium">Batal</button>
                <button on:click={handleSubmit} class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium">Simpan</button>
            </div>
        </div>
    </div>
{/if}