<template>
    <div class="auth-wrapper">
        <div class="auth-card">
            <div class="auth-header">
                <h1>Trello Clone</h1>
                <p>Log in to continue</p>
            </div>

            <form @submit.prevent="handleLogin" class="auth-form">
                <div class="form-group">
                    <label>Email</label>
                    <input v-model="email" type="email" placeholder="Enter email" required />
                </div>

                <div class="form-group">
                    <label>Password</label>
                    <input v-model="password" type="password" placeholder="Enter password" required />
                </div>

                <button type="submit" :disabled="auth.loading" class="btn-primary">
                    {{ auth.loading ? 'Logging in...' : 'Log In' }}
                </button>

                <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
            </form>
            <div class="auth-footer">
                <p>Don't have an account? <router-link to="/signup">Sign Up</router-link></p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();

const email = ref('');
const password = ref('');
const errorMsg = ref('');

const handleLogin = async () => {
    const result = await auth.login(email.value, password.value);
    if (result.success) {
        router.push('/');
    } else {
        errorMsg.value = result.message;
    }
};
</script>

<style scoped>
.auth-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: var(--trello-blue);
}

.auth-card {
    background: white;
    padding: 40px;
    border-radius: var(--border-radius);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    width: 100%;
    max-width: 400px;
}

.auth-header {
    text-align: center;
    margin-bottom: 25px;
}

.auth-header h1 {
    color: var(--trello-blue);
    font-size: 28px;
    margin-bottom: 5px;
}

.form-group {
    margin-bottom: 15px;
}

.form-group label {
    display: block;
    font-size: 12px;
    font-weight: 700;
    color: #5e6c84;
    margin-bottom: 4px;
}

input {
    width: 100%;
    padding: 10px;
    border: 2px solid #dfe1e6;
    border-radius: 3px;
    background-color: #fafbfc;
    font-size: 14px;
}

input:focus {
    border-color: #4c9aff;
    outline: none;
    background: white;
}

.btn-primary {
    width: 100%;
    padding: 10px;
    background-color: #5aac44;
    color: white;
    border: none;
    border-radius: 3px;
    font-weight: bold;
    cursor: pointer;
    margin-top: 10px;
}

.btn-primary:hover {
    background-color: #61bd4f;
}

.error-text {
    color: #eb5a46;
    font-size: 13px;
    margin-top: 10px;
    text-align: center;
}
</style>