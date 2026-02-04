import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';

// User interface
export interface User {
	username: string;
	role: string;
	name: string;
}

// Auth store state interface
interface AuthState {
	token: string | null;
	user: User | null;
	isLoggedIn: boolean;
}

// Initialize state from localStorage if in browser
const getInitialState = (): AuthState => {
	if (browser) {
		const storedToken = localStorage.getItem('token');
		const storedUser = localStorage.getItem('user');
		
		if (storedToken && storedUser) {
			try {
				const user = JSON.parse(storedUser) as User;
				return {
					token: storedToken,
					user: user,
					isLoggedIn: true
				};
			} catch (error) {
				console.error('Error parsing stored auth data:', error);
				return {
					token: null,
					user: null,
					isLoggedIn: false
				};
			}
		}
	}
	
	return {
		token: null,
		user: null,
		isLoggedIn: false
	};
};

// Create writable store with initial state
export const auth = writable<AuthState>(getInitialState());

// Login function
export const login = (token: string, userData: User) => {
	auth.set({
		token,
		user: userData,
		isLoggedIn: true
	});
	
	if (browser) {
		localStorage.setItem('token', token);
		localStorage.setItem('user', JSON.stringify(userData));
	}
};

// Logout function
export const logout = () => {
	auth.set({
		token: null,
		user: null,
		isLoggedIn: false
	});
	
	if (browser) {
		localStorage.removeItem('token');
		localStorage.removeItem('user');
		goto('/');
	}
};