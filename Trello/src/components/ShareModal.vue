<template>
  <Transition name="fade">
    <div v-if="isOpen" class="modal-backdrop" @click.self="$emit('close')">
      <div class="share-modal-content glass-panel shadow-strong animate-modal-in">
        <div class="modal-header">
          <div class="header-title-group">
            <div class="icon-box">
              <UsersIcon :size="20" class="primary-icon" />
            </div>
            <h3>Share Board</h3>
          </div>
          <button class="modal-close-btn" @click="$emit('close')">
            <XIcon :size="18" />
          </button>
        </div>
        <div class="modal-body">
          <p class="share-subtitle">Invite someone to collaborate on this board.</p>
          <div class="input-with-icon">
            <MailIcon :size="16" class="input-inner-icon" />
            <input 
              v-model="localEmail" 
              type="email" 
              placeholder="Enter email address" 
              @keyup.enter="onInvite"
              class="premium-input-field" 
            />
          </div>
          <button class="btn-primary-glow w-100 mt-4" @click="onInvite" :disabled="!localEmail.trim()">
            Send Invitation
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, watch } from 'vue';
import { X as XIcon, Users as UsersIcon, Mail as MailIcon } from 'lucide-vue-next';

const props = defineProps({
  isOpen: { type: Boolean, default: false }
});

const emit = defineEmits(['close', 'invite']);

const localEmail = ref('');

watch(() => props.isOpen, (val) => {
  if (val) localEmail.value = '';
});

const onInvite = () => {
  if (!localEmail.value.trim()) return;
  emit('invite', localEmail.value.trim());
};
</script>

<style scoped>
.modal-backdrop {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(15, 23, 42, 0.4); backdrop-filter: blur(8px);
  display: flex; justify-content: center; align-items: flex-start;
  z-index: 2000; overflow-y: auto; padding: 100px 20px;
}

.share-modal-content {
  background: var(--surface-primary); 
  width: 100%; max-width: 420px; 
  padding: 32px;
  border-radius: 24px; 
  box-shadow: var(--shadow-strong);
  border: 1px solid var(--border-subtle);
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;
}

.header-title-group { display: flex; align-items: center; gap: 12px; }
.icon-box { 
  background: var(--brand-primary-light); 
  width: 40px; height: 40px; border-radius: 12px; 
  display: flex; align-items: center; justify-content: center;
}
.primary-icon { color: var(--brand-primary); }

.modal-header h3 { margin: 0; font-size: 20px; font-weight: 800; color: var(--text-main); }

.modal-close-btn {
  background: var(--surface-secondary); border: none; padding: 8px;
  border-radius: 50%; cursor: pointer; color: var(--text-muted);
  display: flex; align-items: center; justify-content: center;
  transition: all 0.2s;
}
.modal-close-btn:hover { background: #fee2e2; color: #ef4444; }

.share-subtitle { color: var(--text-muted); font-size: 14px; margin: 0 0 20px 0; line-height: 1.5; }

.input-with-icon { position: relative; display: flex; align-items: center; }
.input-inner-icon { position: absolute; left: 16px; color: var(--text-muted); pointer-events: none; }

.premium-input-field {
  width: 100%; padding: 14px 16px 14px 44px; 
  background: var(--surface-secondary);
  border: 2px solid transparent; border-radius: 14px; 
  font-family: inherit; font-size: 14px;
  box-sizing: border-box; transition: all 0.2s;
  color: var(--text-main);
}
.premium-input-field:focus { 
  outline: none; border-color: var(--brand-primary); 
  background: var(--surface-primary);
  box-shadow: 0 0 0 4px var(--brand-primary-light); 
}

.btn-primary-glow {
  background: var(--brand-primary); color: white;
  border: none; padding: 14px; border-radius: 14px;
  font-weight: 700; cursor: pointer; transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}
.btn-primary-glow:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 16px rgba(99, 102, 241, 0.4); }
.btn-primary-glow:disabled { background: var(--text-extramuted); cursor: not-allowed; box-shadow: none; }

.w-100 { width: 100%; }
.mt-4 { margin-top: 16px; }

.animate-modal-in {
  animation: modalScaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes modalScaleIn {
  from { opacity: 0; transform: scale(0.95) translateY(10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
