import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
	id: number;
	username: string;
	role: string;
	name?: string;
	nik?: string;
	email?: string; // Tambahkan field email
}

interface AuthState {
	isLoggedIn: boolean;
	token: string | null;
	user: User | null;
}

// Inisialisasi state awal dari LocalStorage (Synchronous agar siap pakai)
const storedToken = browser ? localStorage.getItem('token') : null;
const storedUser = browser ? localStorage.getItem('user') : null;

const initialState: AuthState = {
	isLoggedIn: !!storedToken,
	token: storedToken,
	user: storedUser ? JSON.parse(storedUser) : null
};

export const auth = writable<AuthState>(initialState);

export const login = (token: string, user: User) => {
	if (browser) {
		localStorage.setItem('token', token);
		localStorage.setItem('user', JSON.stringify(user));
	}
	auth.set({
		isLoggedIn: true,
		token,
		user
	});
};

export const logout = () => {
	if (browser) {
		localStorage.removeItem('token');
		localStorage.removeItem('user');
        
        // Opsional: Hapus data sesi lain
        localStorage.removeItem('activeMachine');
        localStorage.removeItem('selectedLot');
	}
	auth.set({
		isLoggedIn: false,
		token: null,
		user: null
	});
    
    // Paksa refresh halaman agar memori bersih
    if (browser) {
        window.location.replace('/');
    }
};