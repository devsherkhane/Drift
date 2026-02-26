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

                // Save to state
                this.token = response.data.token;
                this.user = response.data.user;

                // Save to browser storage
                localStorage.setItem('token', this.token);

                return { success: true };
            } catch (error) {
                return {
                    success: false,
                    message: error.response?.data?.error || 'Login failed'
                };
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
    // Add this inside the actions: {} block in your auth.js store
    async register(userData) {
            this.loading = true;
            try {
                const response = await api.post('/register', userData);
                // Usually, you log the user in immediately after registration
                this.token = response.data.token;
                this.user = response.data.user;
                localStorage.setItem('token', this.token);
                return { success: true };
            } catch (error) {
                return {
                    success: false,
                    message: error.response?.data?.error || 'Registration failed'
                };
            } finally {
                this.loading = false;
            }
        }
    }
});