<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="$emit('close')">
    <div class="share-modal-content">
      <div class="modal-header" style="padding-right: 0;">
        <h3>Share Board</h3>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>
      <div class="modal-body" style="flex-direction: column; gap: 15px;">
        <p class="subtitle" style="margin: 0;">Invite someone to collaborate on this board.</p>
        <div class="form-group">
          <input 
            v-model="localEmail" 
            type="email" 
            placeholder="Email address" 
            @keyup.enter="onInvite"
            class="share-input" 
          />
        </div>
        <button class="btn-primary w-100" style="padding: 12px;" @click="onInvite" :disabled="!localEmail.trim()">
          Send Invite
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';

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
  background: rgba(15, 23, 42, 0.6); backdrop-filter: blur(4px);
  display: flex; justify-content: center; align-items: flex-start;
  z-index: 100; overflow-y: auto; padding: 50px 0;
  animation: fadeIn 0.2s ease-out;
}

.share-modal-content {
  background: white; width: 400px; padding: 24px;
  border-radius: var(--border-radius); box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  animation: modalSlideIn 0.2s ease-out;
}

@keyframes modalSlideIn {
  from { opacity: 0; transform: translateY(20px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.modal-header {
  display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px;
}
.modal-header h3 { margin: 0; font-size: 24px; font-weight: 700; color: var(--text-main); letter-spacing: -0.5px; }
.subtitle { color: var(--text-muted); font-size: 14px; }

.btn-close {
  background: none; border: none; font-size: 20px; cursor: pointer;
  color: var(--text-secondary); padding: 4px 8px; border-radius: var(--border-radius-sm);
}
.btn-close:hover { background: var(--surface-secondary); color: var(--text-main); }

.share-input {
  width: 100%; padding: 12px; border: 1px solid #e2e8f0;
  border-radius: var(--border-radius-sm); font-family: inherit; font-size: 14px;
  box-sizing: border-box; transition: all 0.2s;
}
.share-input:focus { outline: none; border-color: var(--primary-color); box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15); }

.btn-primary {
  background: var(--brand-primary); color: var(--text-on-brand);
  border: none; border-radius: var(--border-radius-sm);
  cursor: pointer; font-weight: 600; transition: all 0.2s;
}
.btn-primary:hover { background: var(--brand-primary-hover); transform: translateY(-1px); }
.btn-primary:disabled { background: #cbd5e1; cursor: not-allowed; transform: none; }
.w-100 { width: 100%; }
</style>
