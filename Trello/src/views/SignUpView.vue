<template>
  <div class="auth-page animate-fade-in">
    <div class="glass-auth-container">
      <div class="auth-card glass-panel shadow-strong">
        <div class="auth-header">
          <div class="logo-area">
              <div class="logo-box">
                  <TrelloIcon :size="32" class="primary-icon" />
              </div>
              <h1>Trello Clone</h1>
          </div>
          <p class="auth-subtitle">Create your account and start organizing your projects beautifully.</p>
        </div>
        
        <form @submit.prevent="handleSignup" class="auth-form-refined">
          <div class="refined-input-group">
            <label><UserIcon :size="12" /> Full Name</label>
            <div class="input-container">
              <input v-model="name" type="text" placeholder="e.g. John Doe" required />
            </div>
          </div>

          <div class="refined-input-group">
            <label><MailIcon :size="12" /> Email Address</label>
            <div class="input-container">
              <input v-model="email" type="email" placeholder="name@company.com" required />
            </div>
          </div>
          
          <div class="refined-input-group">
            <label><LockIcon :size="12" /> Password</label>
            <div class="input-container">
              <input v-model="password" type="password" placeholder="••••••••" required />
            </div>
          </div>

          <button type="submit" :disabled="auth.loading" class="btn-primary-glow w-100 mt-4">
            <UserPlusIcon v-if="!auth.loading" :size="18" />
            {{ auth.loading ? 'Creating account...' : 'Create Account' }}
          </button>
        </form>
        
        <div class="auth-divider">
            <span>or</span>
        </div>

        <div class="auth-footer-refined">
          <p>Already have an account? <router-link to="/login" class="brand-link">Sign In</router-link></p>
        </div>
      </div>
      
      <div class="auth-decoration">
          <div class="blob blob-1"></div>
          <div class="blob blob-2"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import { useToast } from "vue-toastification";
import { 
  Columns as TrelloIcon, 
  User as UserIcon, 
  Mail as MailIcon, 
  Lock as LockIcon, 
  UserPlus as UserPlusIcon 
} from 'lucide-vue-next';

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
.auth-page {
    min-height: 100vh;
    background: var(--bg-gradient);
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px 20px;
    overflow: hidden;
    position: relative;
}

.glass-auth-container {
    position: relative;
    z-index: 10;
    width: 100%;
    max-width: 480px;
}

.auth-card {
    background: var(--surface-primary);
    padding: 56px 48px;
    border-radius: 32px;
    border: 1px solid var(--border-subtle);
    position: relative;
    z-index: 20;
}

.auth-header {
    text-align: center;
    margin-bottom: 40px;
}

.logo-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    margin-bottom: 24px;
}

.logo-box {
    width: 64px;
    height: 64px;
    background: var(--brand-primary-light);
    border-radius: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--brand-primary);
    box-shadow: 0 8px 16px rgba(99, 102, 241, 0.15);
}

.auth-header h1 {
    color: var(--text-main);
    font-size: 32px;
    font-weight: 900;
    letter-spacing: -1.5px;
    margin: 0;
}

.auth-subtitle {
  color: var(--text-muted);
  font-size: 15px;
  line-height: 1.6;
  margin: 0;
  max-width: 320px;
  margin: 0 auto;
}

.auth-form-refined {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.refined-input-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.refined-input-group label {
    font-size: 12px;
    font-weight: 800;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    display: flex;
    align-items: center;
    gap: 8px;
}

.input-container input {
    width: 100%;
    padding: 16px 20px;
    background: var(--surface-secondary);
    border: 2px solid transparent;
    border-radius: 14px;
    font-size: 15px;
    color: var(--text-main);
    transition: all 0.25s ease;
}

.input-container input:focus {
    background: var(--surface-primary);
    border-color: var(--brand-primary);
    outline: none;
    box-shadow: 0 0 0 4px var(--brand-primary-light);
}

.btn-primary-glow {
    background: var(--brand-primary);
    color: white;
    padding: 16px;
    font-size: 16px;
    font-weight: 800;
    border: none;
    border-radius: 16px;
    cursor: pointer;
    box-shadow: 0 6px 20px rgba(99, 102, 241, 0.25);
    transition: all 0.25s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
}

.btn-primary-glow:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 10px 25px rgba(99, 102, 241, 0.35);
}

.auth-divider {
    display: flex;
    align-items: center;
    text-align: center;
    margin: 8px 0;
}
.auth-divider::before, .auth-divider::after {
    content: '';
    flex: 1;
    border-bottom: 1px solid var(--border-subtle);
}
.auth-divider span {
    padding: 0 16px;
    font-size: 12px;
    color: var(--text-extramuted);
    font-weight: 700;
    text-transform: uppercase;
}

.auth-footer-refined {
    text-align: center;
}

.auth-footer-refined p {
    font-size: 14px;
    color: var(--text-muted);
    font-weight: 500;
}

.brand-link {
    color: var(--brand-primary);
    font-weight: 800;
    text-decoration: none;
}
.brand-link:hover {
    text-decoration: underline;
}

/* Decoration */
.auth-decoration {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 600px;
    height: 600px;
    z-index: 5;
    pointer-events: none;
}

.blob {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.15;
    animation: pulse 10s infinite alternate;
}

.blob-1 {
    width: 300px;
    height: 300px;
    background: var(--brand-primary);
    top: -50px;
    left: -50px;
}

.blob-2 {
    width: 250px;
    height: 250px;
    background: #10b981;
    bottom: -50px;
    right: -50px;
    animation-delay: -5s;
}

@keyframes pulse {
    from { transform: scale(1) translate(0, 0); }
    to { transform: scale(1.2) translate(50px, 30px); }
}

.animate-fade-in {
    animation: fadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
}

.w-100 { width: 100%; }
.mt-4 { margin-top: 16px; }
</style>