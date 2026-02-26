<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <div class="auth-header">
        <h1>Trello Clone</h1>
        <p>Create your account</p>
      </div>
      
      <form @submit.prevent="handleSignup" class="auth-form">
        <div class="form-group">
          <label>Full Name</label>
          <input v-model="name" type="text" placeholder="e.g. Dev Sherkhane" required />
        </div>

        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" placeholder="Enter email" required />
        </div>
        
        <div class="form-group">
          <label>Password</label>
          <input v-model="password" type="password" placeholder="Create password" required />
        </div>

        <button type="submit" :disabled="auth.loading" class="btn-primary">
          {{ auth.loading ? 'Creating account...' : 'Sign Up' }}
        </button>
      </form>
      
      <div class="auth-footer">
        <p>Already have an account? <router-link to="/login">Log In</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import { useToast } from "vue-toastification";

const auth = useAuthStore();
const router = useRouter();
const toast = useToast();

const name = ref('');
const email = ref('');
const password = ref('');

const handleSignup = async () => {
  const result = await auth.register({ 
    username: name.value, 
    email: email.value, 
    password: password.value 
  });
  
  if (result.success) {
    toast.success("Welcome to Trello Clone!");
    router.push('/');
  } else {
    toast.error(result.message);
  }
};
</script>

<style scoped>
/* You can reuse the exact same CSS from your LoginView.vue */
.auth-footer {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
}
.auth-footer a {
  color: var(--trello-blue);
  text-decoration: none;
  font-weight: bold;
}
</style>