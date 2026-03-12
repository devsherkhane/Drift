<template>
  <div v-if="card" class="modal-backdrop" @click.self="$emit('close')">
    <div class="card-detail-modal">
      <button class="modal-close" @click="$emit('close')">✕</button>

      <div class="modal-header">
        <div class="header-title">
          <h3>💳 {{ card.title }}</h3>
          <p class="subtitle">in list <u>{{ card.listName }}</u></p>
        </div>
        <button @click="$emit('delete')" class="btn-danger">Delete Card</button>
      </div>

      <div class="modal-body">
        <div class="main-column">
          <!-- Description -->
          <div class="modal-section">
            <h4>📝 Description</h4>
            <textarea 
              :value="card.description" 
              placeholder="Add a more detailed description..."
              @blur="$emit('update-field', 'description', $event.target.value)"
              class="description-input"
            ></textarea>
          </div>

          <!-- Attachments -->
          <div class="modal-section">
            <h4>📎 Attachments</h4>
            <div class="attachments-list">
              <div v-for="file in card.attachments" :key="file.id" class="attachment-item">
                <div class="file-preview">
                  <img v-if="isImage(file.filename)" :src="getAttachmentUrl(file.file_path)" alt="preview"/>
                  <span v-else class="file-icon">📄</span>
                </div>
                <div class="file-info">
                  <a :href="getAttachmentUrl(file.file_path)" target="_blank">{{ file.filename }}</a>
                  <span class="file-date">Added recently</span>
                </div>
              </div>
            </div>
            <input type="file" ref="fileInput" style="display: none" @change="onFileUpload" />
            <button class="btn-secondary" @click="$refs.fileInput.click()">Add Attachment</button>
          </div>

          <!-- Comments -->
          <div class="modal-section comments-section">
            <h4>💬 Activity</h4>
            <div class="comment-input-area">
              <div class="comment-avatar">{{ userInitial }}</div>
              <div class="comment-box">
                <textarea v-model="localCommentText" placeholder="Write a comment..."></textarea>
                <button class="btn-primary" @click="onSubmitComment" :disabled="!localCommentText.trim()">Save</button>
              </div>
            </div>
            
            <div class="comments-list">
              <div v-for="comment in card.comments" :key="comment.id" class="comment-item">
                <div class="comment-avatar">{{ comment.user_name?.charAt(0).toUpperCase() || 'U' }}</div>
                <div class="comment-content">
                  <span class="comment-user">{{ comment.user_name }}</span>
                  <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
                  <p>{{ comment.text }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="sidebar-column">
          <div class="modal-section">
            <h4>📅 Due Date</h4>
            <input 
              type="date" 
              :value="card.due_date" 
              @change="$emit('update-field', 'due_date', $event.target.value)"
              class="date-picker"
            />
          </div>

          <div class="modal-section">
            <h4>🏷️ Labels</h4>
            <div class="labels-picker">
              <div 
                v-for="color in labelColors" 
                :key="color"
                class="label-swatch" 
                :style="{ backgroundColor: color }"
                :class="{ active: card.label_color === color }" 
                @click="$emit('update-field', 'label_color', color)"
              ></div>
              <button class="btn-clear-label" @click="$emit('update-field', 'label_color', null)">Clear</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

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
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  z-index: 100;
  overflow-y: auto;
  padding: 50px 0;
  animation: fadeIn 0.2s ease-out;
}

.card-detail-modal {
  background: var(--surface-primary);
  width: 768px;
  padding: 32px;
  border-radius: 20px;
  position: relative;
  box-shadow: var(--shadow-strong);
  border: 1px solid var(--border-subtle);
  animation: modalSlideIn 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalSlideIn {
  from { opacity: 0; transform: translateY(20px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.modal-close {
  position: absolute; top: 20px; right: 20px;
  background: #f1f5f9; border: none; font-size: 18px;
  width: 36px; height: 36px; border-radius: 50%;
  cursor: pointer; color: var(--text-secondary);
  display: flex; align-items: center; justify-content: center;
  transition: all 0.2s;
}
.modal-close:hover { background: #e2e8f0; color: var(--text-primary); }

.modal-header {
  display: flex; justify-content: space-between;
  align-items: flex-start; margin-bottom: 24px; padding-right: 40px;
}
.modal-header h3 { margin: 0; font-size: 24px; font-weight: 700; color: var(--text-main); letter-spacing: -0.5px; }
.modal-header .subtitle { margin: 4px 0 0 32px; color: var(--text-muted); font-size: 14px; }

.btn-danger {
  background: #fef2f2; color: #ef4444; border: 1px solid #fee2e2;
  padding: 8px 16px; border-radius: var(--border-radius-sm);
  cursor: pointer; font-weight: 600; transition: all 0.2s;
}
.btn-danger:hover { background: #fee2e2; border-color: #fca5a5; }

.modal-body { display: flex; gap: 32px; }
.main-column { flex: 3; }
.sidebar-column { flex: 1; min-width: 180px; }
.modal-section { margin-bottom: 32px; }
.modal-section h4 {
  margin: 0 0 16px 0; color: var(--text-primary);
  font-size: 16px; font-weight: 600; display: flex; align-items: center; gap: 8px;
}

.description-input {
  width: 100%; min-height: 120px; padding: 16px;
  background: #f8fafc; border: 1px solid #e2e8f0;
  border-radius: var(--border-radius-sm); resize: vertical;
  font-family: inherit; font-size: 14px; box-sizing: border-box; transition: all 0.2s;
}
.description-input:focus {
  background: #fff; border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15); outline: none;
}

.attachment-item {
  display: flex; gap: 16px; margin-bottom: 16px; padding: 16px;
  border: 1px solid #e2e8f0; border-radius: var(--border-radius-sm);
  background: #f8fafc; transition: all 0.2s;
}
.attachment-item:hover { border-color: #cbd5e1; background: #fff; }

.file-preview {
  width: 80px; height: 60px; background: #e2e8f0;
  border-radius: var(--border-radius-sm); display: flex;
  align-items: center; justify-content: center; overflow: hidden;
}
.file-preview img { width: 100%; height: 100%; object-fit: cover; }
.file-icon { font-size: 24px; }
.file-info { display: flex; flex-direction: column; justify-content: center; gap: 4px; }
.file-info a { color: var(--text-primary); font-weight: 600; font-size: 14px; text-decoration: none; word-break: break-all; }
.file-info a:hover { color: var(--primary-color); text-decoration: underline; }
.file-date { font-size: 12px; color: var(--text-secondary); }

.btn-secondary {
  background: rgba(255, 255, 255, 0.2); border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 8px 16px; border-radius: var(--border-radius-sm);
  cursor: pointer; color: white; font-weight: 600; transition: all 0.2s;
}
.btn-secondary:hover { background: rgba(255, 255, 255, 0.3); transform: translateY(-1px); }

.btn-primary {
  background: var(--brand-primary); color: var(--text-on-brand);
  border: none; padding: 8px 16px; border-radius: var(--border-radius-sm);
  cursor: pointer; font-weight: 600; transition: all 0.2s;
}
.btn-primary:hover { background: var(--brand-primary-hover); transform: translateY(-1px); }
.btn-primary:disabled { background: #cbd5e1; cursor: not-allowed; transform: none; }

.comment-input-area { display: flex; gap: 16px; margin-bottom: 24px; }
.comment-avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: var(--bg-gradient); color: white;
  display: flex; align-items: center; justify-content: center;
  font-weight: 600; font-size: 14px; flex-shrink: 0; box-shadow: var(--shadow-sm);
}
.comment-box { flex-grow: 1; display: flex; flex-direction: column; gap: 12px; }
.comment-box textarea {
  width: 100%; height: 80px; padding: 12px;
  border: 1px solid #e2e8f0; border-radius: var(--border-radius-sm);
  background: #f8fafc; font-family: inherit; font-size: 14px;
  resize: vertical; box-sizing: border-box; transition: all 0.2s;
}
.comment-box textarea:focus { background: #fff; border-color: var(--primary-color); box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15); outline: none; }
.comment-box .btn-primary { align-self: flex-start; }

.comment-item { display: flex; gap: 16px; margin-bottom: 20px; }
.comment-content {
  background: #fff; border: 1px solid #e2e8f0;
  padding: 16px; border-radius: var(--border-radius-sm); width: 100%;
}
.comment-user { font-weight: 600; color: var(--text-primary); }
.comment-date { font-size: 12px; color: var(--text-secondary); }
.comment-content p { margin: 0; color: var(--text-primary); font-size: 14px; line-height: 1.5; white-space: pre-wrap; }

.date-picker {
  width: 100%; padding: 10px; border: 1px solid #e2e8f0;
  border-radius: var(--border-radius-sm); background: #f8fafc;
  font-family: inherit; color: var(--text-primary); box-sizing: border-box; transition: all 0.2s;
}
.date-picker:focus { outline: none; border-color: var(--primary-color); box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15); background: #fff; }

.labels-picker { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 12px; }
.label-swatch { width: 44px; height: 32px; border-radius: var(--border-radius-sm); cursor: pointer; transition: transform 0.2s; }
.label-swatch:hover { transform: scale(1.05); }
.label-swatch.active { border: 2px solid var(--text-primary); box-shadow: 0 0 0 2px white inset; }
.btn-clear-label {
  background: #f1f5f9; border: none; padding: 8px 12px;
  border-radius: var(--border-radius-sm); cursor: pointer;
  margin-top: 8px; width: 100%; text-align: center;
  font-weight: 600; color: var(--text-secondary); transition: all 0.2s;
}
.btn-clear-label:hover { background: #e2e8f0; color: var(--text-primary); }
</style>
