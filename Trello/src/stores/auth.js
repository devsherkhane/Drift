import { defineStore } from 'pinia';
import api from '../api';

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        token: localStorage.getItem('token') || null,
        loading: false
    }),

    getters: {
        isAuthenticated: (state) => !!state.token
    },

    actions: {
        async login(email, password) {
            this.loading = true;
            try {
                const response = await api.post('/login', { email, password });

                // CHECK IF 2FA IS REQUIRED
                if (response.data.requires_2fa) {
                    return { success: false, requires_2fa: true, user_id: response.data.user_id };
                }

                // Normal Login
                this.token = response.data.token;
                localStorage.setItem('token', this.token);
                await this.fetchUser();
                return { success: true };
            } catch (error) {
                return { success: false, message: error.response?.data?.error || 'Login failed' };
            } finally {
                this.loading = false;
            }
        },

        logout() {
            this.token = null;
            this.user = null;
            localStorage.removeItem('token');
            window.location.href = '/login'; // Hard redirect to clear state
        },

        async register(userData) {
            this.loading = true;
            try {
                const response = await api.post('/register', userData);
                this.token = response.data.token;
                localStorage.setItem('token', this.token);
                await this.fetchUser();
                return { success: true };
            } catch (error) {
                return {
                    success: false,
                    message: error.response?.data?.error || 'Registration failed'
                };
            } finally {
                this.loading = false;
            }
        },

        async fetchUser() {
            if (!this.token) return;
            try {
                const response = await api.get('/user/me');
                this.user = response.data;
                // Sync theme from profile
                if (this.user.theme) {
                    document.documentElement.setAttribute('data-theme', this.user.theme);
                    localStorage.setItem('theme', this.user.theme);
                }
            } catch (error) {
                console.error("Failed to fetch user:", error);
                if (error.response?.status === 401) {
                    this.logout();
                }
            }
        }
    }
});