<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <div class="auth-header">
        <h1>Reset Password</h1>
        <p>Enter your email to get a reset link.</p>
      </div>

      <form @submit.prevent="handleForgot" class="auth-form">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="Enter your account email" required />
        </div>

        <button type="submit" :disabled="loading" class="btn-primary">
          {{ loading ? 'Sending...' : 'Send Reset Link' }}
        </button>

        <p v-if="message" class="success-text">{{ message }}</p>
        <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
      </form>
      <div class="auth-footer">
        <p>Remembered your password? <router-link to="/login">Log In</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import api from '../api';

const email = ref('');
const message = ref('');
const errorMsg = ref('');
const loading = ref(false);

const handleForgot = async () => {
  loading.value = true;
  message.value = '';
  errorMsg.value = '';
  
  try {
    const response = await api.post('/forgot-password', { email: email.value });
    message.value = response.data.message;
  } catch (error) {
    errorMsg.value = error.response?.data?.error || "Something went wrong";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* Reusing your LoginView styles */
.auth-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #0079bf;
  padding: 20px;
}

.auth-card {
  background: #ffffff;
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

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  font-size: 12px;
  font-weight: 700;
  color: #64748b;
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

.success-text {
  color: #059669; /* Emerald 600 */
  font-size: 14px;
  margin-top: 20px;
  text-align: center;
  font-weight: 600;
  background: #ecfdf5;
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #d1fae5;
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