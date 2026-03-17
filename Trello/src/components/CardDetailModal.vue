<template>
  <Transition name="fade">
    <div v-if="card" class="modal-backdrop" @click.self="$emit('close')">
      <div class="card-detail-modal glass-panel shadow-strong animate-modal-in">
        <button class="modal-close-btn" @click="$emit('close')">
          <XIcon :size="20" />
        </button>

        <div class="modal-header">
          <div class="header-icon-container">
            <CreditCardIcon :size="24" class="primary-icon" />
          </div>
          <div class="header-content">
            <h3 class="modal-title">{{ card.title }}</h3>
            <p class="modal-subtitle">
              <ListIcon :size="14" /> in list <u>{{ card.listName }}</u>
            </p>
          </div>
        </div>

        <div class="modal-layout">
          <div class="main-column">
            <!-- Description -->
            <div class="modal-section">
              <div class="section-header">
                <AlignLeftIcon :size="20" />
                <h4>Description</h4>
              </div>
              <textarea 
                :value="card.description" 
                placeholder="Add a more detailed description..."
                @blur="$emit('update-field', 'description', $event.target.value)"
                class="description-input"
              ></textarea>
            </div>

            <!-- Attachments -->
            <div class="modal-section">
              <div class="section-header">
                <PaperclipIcon :size="20" />
                <h4>Attachments</h4>
              </div>
              <div class="attachments-grid">
                <div v-for="file in card.attachments" :key="file.id" class="attachment-card glass-panel">
                  <div class="file-preview-box">
                    <img v-if="isImage(file.filename)" :src="getAttachmentUrl(file.file_path)" alt="preview"/>
                    <FileIcon v-else :size="24" class="fallback-file-icon" />
                  </div>
                  <div class="file-details">
                    <a :href="getAttachmentUrl(file.file_path)" target="_blank" class="filename-link">{{ file.filename }}</a>
                    <span class="file-meta">Added recently</span>
                  </div>
                </div>
              </div>
              <input type="file" ref="fileInput" style="display: none" @change="onFileUpload" />
              <button class="btn-action-outline mt-4" @click="$refs.fileInput.click()">
                <PlusIcon :size="16" /> Add Attachment
              </button>
            </div>

            <!-- Activity -->
            <div class="modal-section">
              <div class="section-header">
                <MessageSquareIcon :size="20" />
                <h4>Activity</h4>
              </div>
              <div class="comment-composer">
                <div class="user-avatar shadow-sm">{{ userInitial }}</div>
                <div class="composer-box">
                  <textarea v-model="localCommentText" placeholder="Write a comment..."></textarea>
                  <button class="btn-primary-glow" @click="onSubmitComment" :disabled="!localCommentText.trim()">
                    Share Comment
                  </button>
                </div>
              </div>
              
              <div class="comments-timeline">
                <div v-for="comment in card.comments" :key="comment.id" class="comment-thread">
                  <div class="user-avatar small">{{ comment.user_name?.charAt(0).toUpperCase() || 'U' }}</div>
                  <div class="comment-bubble glass-panel">
                    <div class="comment-meta">
                      <span class="comment-author">{{ comment.user_name }}</span>
                      <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                    </div>
                    <p class="comment-body">{{ comment.text }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Sidebar -->
          <div class="sidebar-column">
            <div class="sidebar-section">
              <div class="section-header small">
                <CalendarIcon :size="16" />
                <h5>Due Date</h5>
              </div>
              <input 
                type="date" 
                :value="card.due_date" 
                @change="$emit('update-field', 'due_date', $event.target.value)"
                class="date-input-refined"
              />
            </div>

            <div class="sidebar-section">
              <div class="section-header small">
                <TagIcon :size="16" />
                <h5>Labels</h5>
              </div>
              <div class="labels-grid">
                <div 
                  v-for="color in labelColors" 
                  :key="color"
                  class="label-pill" 
                  :style="{ backgroundColor: color }"
                  :class="{ 'is-active': card.label_color === color }" 
                  @click="$emit('update-field', 'label_color', color)"
                ></div>
              </div>
              <button class="btn-clear-action" @click="$emit('update-field', 'label_color', null)">
                <XIcon :size="12" /> Clear Label
              </button>
            </div>

            <div class="sidebar-section danger-zone-sidebar">
              <h5>Actions</h5>
              <button @click="$emit('delete')" class="btn-danger-soft">
                <Trash2Icon :size="16" /> Delete Card
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref } from 'vue';
import { 
  CreditCard as CreditCardIcon, 
  List as ListIcon, 
  X as XIcon,
  Trash2 as Trash2Icon,
  AlignLeft as AlignLeftIcon,
  Paperclip as PaperclipIcon,
  Plus as PlusIcon,
  MessageSquare as MessageSquareIcon,
  File as FileIcon,
  Calendar as CalendarIcon,
  Tag as TagIcon
} from 'lucide-vue-next';

const props = defineProps({
  card: { type: Object, default: null },
  userInitial: { type: String, default: 'U' }
});

const emit = defineEmits(['close', 'delete', 'update-field', 'upload-file', 'submit-comment']);

const localCommentText = ref('');
const fileInput = ref(null);

const labelColors = ['#61bd4f', '#f2d600', '#ff9f1a', '#eb5a46', '#c377e0', '#0079bf'];

const getAttachmentUrl = (filePath) => {
  if (!filePath) return '';
  return `http://localhost:8080/${filePath.replace(/\\/g, '/')}`;
};

const isImage = (filename) => {
  if (!filename) return false;
  return /\.(jpg|jpeg|png|gif|webp)$/i.test(filename);
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
};

const onFileUpload = (event) => {
  emit('upload-file', event);
};

const onSubmitComment = () => {
  if (!localCommentText.value.trim()) return;
  emit('submit-comment', localCommentText.value);
  localCommentText.value = '';
};
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  z-index: 2000;
  overflow-y: auto;
  padding: 60px 20px;
}

.card-detail-modal {
  background: var(--surface-primary);
  width: 100%;
  max-width: 800px;
  padding: 40px;
  border-radius: 24px;
  position: relative;
  box-shadow: var(--shadow-strong);
  border: 1px solid var(--border-subtle);
}

.modal-close-btn {
  position: absolute; top: 24px; right: 24px;
  background: var(--surface-secondary); border: none;
  width: 40px; height: 40px; border-radius: 50%;
  cursor: pointer; color: var(--text-muted);
  display: flex; align-items: center; justify-content: center;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-close-btn:hover { 
  background: #fee2e2; 
  color: #ef4444; 
  transform: rotate(90deg) scale(1.1);
}

.modal-header {
  display: flex; gap: 20px; align-items: flex-start; margin-bottom: 40px;
}

.header-icon-container {
  background: var(--brand-primary-light);
  width: 48px; height: 48px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}

.primary-icon { color: var(--brand-primary); }

.header-content { flex-grow: 1; }

.modal-title { 
  margin: 0; font-size: 24px; font-weight: 800; 
  color: var(--text-main); letter-spacing: -0.5px;
}

.modal-subtitle { 
  margin: 6px 0 0 0; color: var(--text-muted); 
  font-size: 14px; display: flex; align-items: center; gap: 6px;
}

.modal-layout { display: flex; gap: 48px; }
.main-column { flex: 1; min-width: 0; }
.sidebar-column { width: 220px; flex-shrink: 0; }

.section-header { margin-bottom: 16px; display: flex; align-items: center; gap: 12px; color: var(--text-main); }
.section-header h4 { margin: 0; font-size: 17px; font-weight: 700; }
.section-header.small { color: var(--text-muted); gap: 8px; margin-bottom: 8px;}
.section-header.small h5 { margin: 0; font-size: 13px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px; }

.modal-section { margin-bottom: 48px; }

.description-input {
  width: 100%; min-height: 140px; padding: 18px;
  background: var(--surface-secondary); border: 2px solid transparent;
  border-radius: 16px; resize: vertical;
  font-family: inherit; font-size: 15px; box-sizing: border-box; transition: all 0.25s;
  color: var(--text-main);
}
.description-input:focus {
  background: var(--surface-primary); border-color: var(--brand-primary);
  box-shadow: 0 0 0 4px var(--brand-primary-light); outline: none;
}

.attachments-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 16px; }

.attachment-card {
  display: flex; gap: 14px; padding: 12px;
  border: 1px solid var(--border-subtle); border-radius: 12px;
  transition: all 0.2s;
}
.attachment-card:hover { transform: translateY(-2px); border-color: var(--brand-primary); }

.file-preview-box {
  width: 70px; height: 50px; background: var(--surface-secondary);
  border-radius: 8px; display: flex; align-items: center; justify-content: center; overflow: hidden;
}
.file-preview-box img { width: 100%; height: 100%; object-fit: cover; }
.fallback-file-icon { color: var(--text-extramuted); }

.file-details { display: flex; flex-direction: column; justify-content: center; min-width: 0; }
.filename-link { 
  color: var(--text-main); font-weight: 700; font-size: 14px; 
  text-decoration: none; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.filename-link:hover { color: var(--brand-primary); }
.file-meta { font-size: 11px; color: var(--text-muted); margin-top: 2px; }

.btn-action-outline {
  background: transparent; border: 1.5px dashed var(--brand-primary);
  padding: 10px 18px; border-radius: 12px; color: var(--brand-primary);
  font-weight: 700; font-size: 14px; cursor: pointer; transition: all 0.2s;
  display: flex; align-items: center; gap: 8px;
}
.btn-action-outline:hover { background: var(--brand-primary-light); border-style: solid; }

.comment-composer { display: flex; gap: 16px; margin-bottom: 32px; }
.user-avatar {
  width: 40px; height: 40px; border-radius: 12px;
  background: var(--brand-primary); color: white;
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; flex-shrink: 0;
}
.user-avatar.small { width: 32px; height: 32px; font-size: 12px; border-radius: 8px; }

.composer-box { flex-grow: 1; display: flex; flex-direction: column; gap: 12px; }
.composer-box textarea {
  width: 100%; height: 100px; padding: 14px;
  border: 2px solid var(--surface-secondary); border-radius: 14px;
  background: var(--surface-secondary); font-family: inherit; font-size: 14px;
  resize: vertical; box-sizing: border-box; transition: all 0.2s;
}
.composer-box textarea:focus { 
  background: var(--surface-primary); border-color: var(--brand-primary); 
  box-shadow: 0 0 0 4px var(--brand-primary-light); outline: none; 
}

.comments-timeline { display: flex; flex-direction: column; gap: 20px; }
.comment-thread { display: flex; gap: 16px; }
.comment-bubble { 
  padding: 16px; border-radius: 16px; border: 1px solid var(--border-subtle); 
  background: var(--surface-primary); flex-grow: 1;
}
.comment-meta { margin-bottom: 8px; display: flex; align-items: center; gap: 10px; }
.comment-author { font-weight: 700; color: var(--text-main); font-size: 14px; }
.comment-time { font-size: 11px; color: var(--text-muted); }
.comment-body { margin: 0; color: var(--text-muted); font-size: 14px; line-height: 1.6; white-space: pre-wrap; }

.sidebar-column { display: flex; flex-direction: column; gap: 32px; }
.sidebar-section { padding-top: 8px; }

.date-input-refined {
  width: 100%; padding: 12px; border: 2px solid var(--surface-secondary);
  border-radius: 12px; background: var(--surface-secondary);
  font-family: inherit; color: var(--text-main); font-weight: 600; transition: all 0.2s;
}
.date-input-refined:focus { outline: none; border-color: var(--brand-primary); background: var(--surface-primary); }

.labels-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; margin-bottom: 12px; }
.label-pill { 
  height: 28px; border-radius: 6px; cursor: pointer; transition: all 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border: 2px solid transparent;
}
.label-pill:hover { transform: scale(1.1); }
.label-pill.is-active { border-color: var(--text-main); transform: scale(1.05); }

.btn-clear-action {
  background: transparent; border: none; padding: 6px;
  font-size: 12px; font-weight: 700; color: var(--text-muted); 
  cursor: pointer; display: flex; align-items: center; gap: 4px;
}
.btn-clear-action:hover { color: #ef4444; }

.btn-danger-soft {
  width: 100%; display: flex; align-items: center; gap: 10px;
  background: #fef2f2; color: #ef4444; border: 1px solid transparent;
  padding: 12px; border-radius: 12px; font-weight: 700; font-size: 14px;
  cursor: pointer; transition: all 0.2s;
}
.btn-danger-soft:hover { background: #fee2e2; transform: translateY(-1px); }

.danger-zone-sidebar { margin-top: 24px; padding-top: 24px; border-top: 1px solid var(--border-subtle); }
.danger-zone-sidebar h5 { margin: 0 0 12px 0; font-size: 13px; color: var(--text-muted); }

.animate-modal-in {
  animation: modalScaleIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes modalScaleIn {
  from { opacity: 0; transform: scale(0.9) translateY(30px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.mt-4 { margin-top: 16px; }
.shadow-strong { box-shadow: var(--shadow-strong); }
.glass-panel { background: var(--surface-glass); backdrop-filter: blur(16px); }

.btn-primary-glow {
  background: var(--brand-primary); color: white;
  border: none; padding: 10px 20px; border-radius: 12px;
  font-weight: 700; cursor: pointer; transition: all 0.2s;
  align-self: flex-start;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}
.btn-primary-glow:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 16px rgba(99, 102, 241, 0.4); }
.btn-primary-glow:disabled { background: var(--text-extramuted); cursor: not-allowed; box-shadow: none; }
</style>
