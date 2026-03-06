<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <div class="auth-header">
        <h1>Set New Password</h1>
      </div>

      <form @submit.prevent="handleReset" class="auth-form">
        <div class="form-group">
          <label>New Password</label>
          <input v-model="newPassword" type="password" placeholder="At least 6 characters" required />
        </div>

        <button type="submit" :disabled="loading" class="btn-primary">
          {{ loading ? 'Updating...' : 'Update Password' }}
        </button>

        <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api';
import { useToast } from 'vue-toastification';

const route = useRoute();
const router = useRouter();
const toast = useToast();

const newPassword = ref('');
const errorMsg = ref('');
const loading = ref(false);

const handleReset = async () => {
  const token = route.query.token; // Gets "?token=..." from URL
  if (!token) {
    errorMsg.value = "Invalid or missing token in URL.";
    return;
  }

  loading.value = true;
  errorMsg.value = '';

  try {
    await api.post('/reset-password', { 
      token: token, 
      new_password: newPassword.value 
    });
    toast.success("Password updated! Please log in.");
    router.push('/login');
  } catch (error) {
    errorMsg.value = error.response?.data?.error || "Failed to reset password";
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* Reuse the exact same styles as ForgotPasswordView */
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
</style>