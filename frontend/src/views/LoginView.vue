<template>
    <div class="auth-wrapper">
        <div class="auth-card">
            <div class="auth-header">
                <h1>Trello Clone</h1>
                <p>Log in to continue</p>
            </div>

            <form @submit.prevent="handleLogin" class="auth-form">
                <div class="form-group">
                    <label>EMAIL</label>
                    <input v-model="email" type="email" placeholder="Enter email" required />
                </div>

                <div class="form-group">
                    <label>PASSWORD</label>
                    <input v-model="password" type="password" placeholder="Enter password" required />
                    <div class="forgot-link-container">
                        <router-link to="/forgot-password" class="forgot-btn">Forgot password?</router-link>
                    </div>
                </div>

                <button type="submit" :disabled="auth.loading" class="btn-primary">
                    {{ auth.loading ? 'Logging in...' : 'Log In' }}
                </button>

                <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
                <div class="auth-footer">
                    <p>Don't have an account? <router-link to="/signup">Sign Up</router-link></p>
                </div>
            </form>
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
    errorMsg.value = '';
    const result = await auth.login(email.value, password.value);
    
    if (result.success) {
        router.push('/');
    } else {
        errorMsg.value = result.message || 'Login failed';
    }
};
</script>

<style scoped>
.auth-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: #14b8a6;
    padding: 20px;
}

.auth-card {
    background:white;
    padding: 60px 40px;
    border-radius: 24px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.04);
    width: 100%;
    max-width: 440px;
    border: 1px solid #f1f5f9;
}

.auth-header {
    text-align: center;
    margin-bottom: 32px;
}

.auth-header h1 {
    color: #115e59; /* Darker Teal */
    font-size: 38px;
    font-weight: 900;
    letter-spacing: -1.5px;
    margin-bottom: 8px;
}

.auth-header p {
  color: #94a3b8; /* Slate 400 */
  font-size: 16px;
  font-weight: 500;
}

.form-group {
    margin-bottom: 24px;
}

.form-group label {
    display: block;
    font-size: 12px;
    font-weight: 700;
    color: #64748b; /* Slate 500 */
    margin-bottom: 12px;
    text-transform: uppercase;
    letter-spacing: 1.5px;
}

input {
    width: 100%;
    padding: 14px 20px;
    border: 1px solid #f1f5f9;
    border-radius: 12px;
    background-color: #fafafa;
    font-size: 15px;
    color: #0f172a;
    transition: all 0.2s ease;
}

input:focus {
    border-color: #14b8a6;
    outline: none;
    box-shadow: 0 0 0 4px rgba(20, 184, 166, 0.1);
    background: white;
}

.forgot-link-container {
    text-align: right;
    margin-top: 10px;
}

.forgot-btn {
    font-size: 12px;
    color: #0ea5e9;
    text-decoration: none;
    font-weight: 600;
    transition: color 0.2s;
}

.forgot-btn:hover {
    color: #0ea5e9;
}

.btn-primary {
    width: 100%;
    margin-top: 30px;
    padding: 16px;
    font-size: 16px;
    background: #14b8a6;
    color: white;
    border: none;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(20, 184, 166, 0.3);
    transition: all 0.2s;
}

.btn-primary:hover {
    background: #0d9488;
    transform: translateY(-1px);
}

.error-text {
    color: #ef4444;
    font-size: 14px;
    margin-top: 20px;
    text-align: center;
    background: #fef2f2;
    padding: 12px;
    border-radius: 10px;
    border: 1px solid #fee2e2;
}

.auth-footer {
  margin-top: 32px;
  text-align: center;
  font-size: 14px;
  color: #64748b;
}

.auth-footer a {
  color: #14b8a6;
  font-weight: 700;
  text-decoration: none;
  transition: all 0.2s;
}

.auth-footer a:hover {
  color: #0d9488;
  text-decoration: underline;
}
</style>