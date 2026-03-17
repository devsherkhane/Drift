<template>
  <div class="profile-canvas animate-fade-in">
    <nav class="profile-header glass-panel">
      <div class="header-left">
        <SettingsIcon :size="20" class="primary-icon" />
        <h2>User Settings</h2>
      </div>
      <div class="header-right">
        <button class="btn-action-premium" @click="goBack">
          <ArrowLeftIcon :size="16" /> Back to Dashboard
        </button>
      </div>
    </nav>

    <div class="profile-content">
      <div class="settings-card glass-panel shadow-strong">
        
        <section class="premium-section">
          <div class="section-badge">Avatar</div>
          <h3>Profile Picture</h3>
          <div class="avatar-manager">
            <div class="avatar-orbit shadow-md">
              <span v-if="!avatarUrl">{{ auth.user?.username?.charAt(0).toUpperCase() || 'U' }}</span>
              <img v-else :src="avatarUrl" alt="Avatar" />
            </div>
            <div class="upload-interaction">
              <p class="section-hint">Choose a photo to personalize your profile.</p>
              <div class="upload-actions">
                <input type="file" ref="fileInput" accept="image/*" style="display: none" @change="onFileSelected" />
                <button class="btn-action-outline-small" @click="$refs.fileInput.click()">
                  <UploadIcon :size="16" /> Choose Image
                </button>
                <button class="btn-primary-glow-small" :disabled="!selectedFile" @click="uploadAvatar">
                  Save Changes
                </button>
              </div>
              <p class="selected-file-name" v-if="selectedFile">
                <FileIcon :size="12" /> {{ selectedFile.name }}
              </p>
            </div>
          </div>
        </section>

        <section class="premium-section">
          <div class="section-badge">Security</div>
          <h3>Account Information</h3>
          <form @submit.prevent="updateProfile" class="refined-form">
            <div class="input-group">
              <label><MailIcon :size="12" /> Email Address</label>
              <input type="email" :value="auth.user?.email" disabled class="input-readonly" />
              <small class="lock-hint">Email cannot be changed.</small>
            </div>
            <div class="input-group">
              <label><UserIcon :size="12" /> Display Name</label>
              <input type="text" v-model="profileForm.username" placeholder="New username" required />
            </div>
            
            <div class="form-divider"></div>
            
            <div class="input-group">
              <label><LockIcon :size="12" /> Current Password</label>
              <input type="password" v-model="profileForm.current_password" placeholder="Verify to save changes" required />
            </div>
            <div class="input-group">
              <label><KeyIcon :size="12" /> New Password (Optional)</label>
              <input type="password" v-model="profileForm.new_password" placeholder="Leave blank to keep current" />
            </div>
            <button type="submit" class="btn-primary-glow w-100">Update Profile</button>
          </form>
        </section>

        <section class="premium-section">
          <div class="section-badge">Appearance</div>
          <h3>Display Theme</h3>
          <p class="section-hint">Switch between light and dark modes.</p>
          
          <div class="theme-grid">
            <button 
              class="theme-toggle-btn" 
              :class="{ active: currentTheme === 'light' }" 
              @click="setTheme('light')"
            >
              <SunIcon :size="18" /> Light
            </button>
            <button 
              class="theme-toggle-btn" 
              :class="{ active: currentTheme === 'dark' }" 
              @click="setTheme('dark')"
            >
              <MoonIcon :size="18" /> Dark
            </button>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useToast } from 'vue-toastification';
import api from '../api';

// Icons
import { 
  Settings as SettingsIcon, 
  ArrowLeft as ArrowLeftIcon, 
  Upload as UploadIcon, 
  User as UserIcon, 
  Mail as MailIcon, 
  Lock as LockIcon, 
  Key as KeyIcon,
  Sun as SunIcon,
  Moon as MoonIcon,
  File as FileIcon
} from 'lucide-vue-next';

const router = useRouter();
const auth = useAuthStore();
const toast = useToast();

// State
const selectedFile = ref(null);
const fileInput = ref(null);

const avatarUrl = computed(() => {
  if (!auth.user?.avatar_url) return '';
  if (auth.user.avatar_url.startsWith('http')) return auth.user.avatar_url;
  return `http://localhost:8080/${auth.user.avatar_url.replace(/\\/g, '/')}`;
});

const profileForm = ref({
  username: '',
  current_password: '',
  new_password: ''
});

const currentTheme = ref('light');

// Initialization
onMounted(async () => {
  if (!auth.user) {
    await auth.fetchUser();
  }
  
  if (auth.user) {
    profileForm.value.username = auth.user.username || '';
    
    // Load existing theme
    const savedTheme = localStorage.getItem('theme') || auth.user.theme || 'light';
    currentTheme.value = savedTheme;
    document.documentElement.setAttribute('data-theme', savedTheme);
  }
});

const goBack = () => {
  router.push('/');
};

// Avatar Handling
const onFileSelected = (event) => {
  const file = event.target.files[0];
  if (file) {
    selectedFile.value = file;
  }
};

const uploadAvatar = async () => {
  if (!selectedFile.value) return;
  
  const formData = new FormData();
  formData.append('avatar', selectedFile.value); // Backend expects 'avatar'

  try {
    const response = await api.post('/user/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    toast.success("Profile picture updated!");
    selectedFile.value = null;
    
    // Update local state by re-fetching user data
    await auth.fetchUser();
  } catch (error) {
    toast.error(error.response?.data?.error || "Failed to upload avatar");
  }
};

// Profile Update Handling
const updateProfile = async () => {
  try {
    await api.put('/profile', {
      username: profileForm.value.username,
      current_password: profileForm.value.current_password,
      new_password: profileForm.value.new_password
    });
    
    toast.success("Profile updated successfully");
    
    // Clear passwords after successful update
    profileForm.value.current_password = '';
    profileForm.value.new_password = '';
    
    // Refresh user data in store
    await auth.fetchUser();
  } catch (error) {
    toast.error(error.response?.data?.error || "Failed to update profile");
  }
};

// Theme Toggle Logic
const setTheme = async (theme) => {
  currentTheme.value = theme;
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);

  try {
    await api.put('/user/theme', { theme });
    toast.success(`Theme updated to ${theme} mode`);
  } catch (err) {
    toast.error("Failed to save theme preference");
  }
};
</script>

<style scoped>
.profile-canvas {
  min-height: 100vh;
  background: var(--bg-gradient);
  display: flex;
  flex-direction: column;
}

.profile-header {
  padding: 12px 32px;
  background: var(--surface-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-subtle);
  display: flex; justify-content: space-between; align-items: center;
  position: sticky; top: 0; z-index: 100;
}

.header-left { display: flex; align-items: center; gap: 12px; }
.header-left h2 {
  color: var(--text-main); margin: 0; font-size: 18px; font-weight: 800; letter-spacing: -0.5px;
}
.primary-icon { color: var(--brand-primary); }

.btn-action-premium {
  background: var(--brand-primary-light); color: var(--brand-primary);
  border: none; padding: 8px 16px; border-radius: 10px;
  font-size: 13px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; gap: 8px; transition: all 0.2s;
}
.btn-action-premium:hover { background: var(--brand-primary); color: white; transform: translateY(-1px); }

.profile-content {
  flex-grow: 1; display: flex; justify-content: center; align-items: flex-start;
  padding: 60px 20px; overflow-y: auto;
}

.settings-card {
  background: var(--surface-primary); width: 100%; max-width: 680px;
  padding: 48px; border-radius: 24px; border: 1px solid var(--border-subtle);
}

.premium-section {
  position: relative; padding: 32px; border-radius: 20px;
  background: var(--surface-secondary); margin-bottom: 32px;
  border: 1px solid var(--border-subtle); transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.premium-section:hover { transform: translateY(-4px); box-shadow: var(--shadow-md); border-color: var(--brand-primary-light); }

.section-badge {
  position: absolute; top: 20px; right: 24px; font-size: 10px;
  font-weight: 800; text-transform: uppercase; letter-spacing: 1px;
  color: var(--brand-primary); background: var(--brand-primary-light);
  padding: 4px 10px; border-radius: 20px;
}

.premium-section h3 { margin: 0 0 24px 0; font-size: 20px; font-weight: 800; color: var(--text-main); }
.section-hint { color: var(--text-muted); font-size: 14px; margin: 0 0 20px 0; line-height: 1.6; }

/* Avatar Manager */
.avatar-manager { display: flex; align-items: center; gap: 32px; }
.avatar-orbit {
  width: 100px; height: 100px; border-radius: 30px;
  background: var(--brand-primary); display: flex; justify-content: center; align-items: center;
  font-size: 40px; font-weight: 800; color: white; overflow: hidden; flex-shrink: 0;
  border: 4px solid var(--surface-primary);
}
.avatar-orbit img { width: 100%; height: 100%; object-fit: cover; }
.upload-interaction { display: flex; flex-direction: column; gap: 12px; }
.upload-actions { display: flex; gap: 12px; }

.btn-action-outline-small {
  background: transparent; border: 1.5px dashed var(--brand-primary);
  padding: 8px 16px; border-radius: 10px; color: var(--brand-primary);
  font-weight: 700; font-size: 13px; cursor: pointer; transition: all 0.2s;
  display: flex; align-items: center; gap: 6px;
}
.btn-action-outline-small:hover { background: var(--brand-primary-light); border-style: solid; }

.btn-primary-glow-small {
  background: var(--brand-primary); color: white; border: none; padding: 8px 20px;
  border-radius: 10px; font-weight: 700; font-size: 13px; cursor: pointer; transition: all 0.2s;
  box-shadow: 0 4px 10px rgba(99, 102, 241, 0.2);
}
.btn-primary-glow-small:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 14px rgba(99, 102, 241, 0.3); }

.selected-file-name { font-size: 12px; color: var(--brand-primary); font-weight: 600; display: flex; align-items: center; gap: 4px; }

/* Refined Form */
.refined-form { display: flex; flex-direction: column; gap: 24px; }
.input-group { display: flex; flex-direction: column; gap: 10px; }
.input-group label { font-size: 12px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.5px; display: flex; align-items: center; gap: 6px; }
.input-group input { 
  padding: 14px 18px; background: var(--surface-primary); 
  border: 2px solid transparent; border-radius: 14px; 
  font-size: 15px; color: var(--text-main); transition: all 0.25s;
}
.input-group input:focus { border-color: var(--brand-primary); outline: none; box-shadow: 0 0 0 4px var(--brand-primary-light); }
.input-readonly { background: var(--surface-secondary) !important; color: var(--text-extramuted) !important; cursor: not-allowed; border: 1px solid var(--border-subtle) !important; }

.form-divider { height: 1px; background: var(--border-subtle); margin: 8px 0; }
.lock-hint { font-size: 11px; color: var(--text-extramuted); margin-top: -4px; font-weight: 500; }

.btn-primary-glow {
  background: var(--brand-primary); color: white; border: none; padding: 16px;
  border-radius: 14px; font-weight: 700; font-size: 15px; cursor: pointer; transition: all 0.25s;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}
.btn-primary-glow:hover { transform: translateY(-1px); box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3); }
.w-100 { width: 100%; }

/* Theme Grid */
.theme-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.theme-toggle-btn {
  padding: 16px; background: var(--surface-primary); border: 2px solid transparent;
  border-radius: 16px; font-size: 15px; font-weight: 700; color: var(--text-muted);
  cursor: pointer; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex; align-items: center; justify-content: center; gap: 10px;
}
.theme-toggle-btn:hover { background: var(--brand-primary-light); color: var(--brand-primary); }
.theme-toggle-btn.active { border-color: var(--brand-primary); color: var(--brand-primary); background: var(--brand-primary-light); }

.animate-fade-in { animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>