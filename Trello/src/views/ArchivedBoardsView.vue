<template>
  <div class="archived-page animate-fade-in">
    <nav class="page-header">
      <div class="header-left">
        <ArchiveIcon :size="20" class="primary-icon" />
        <h2>Archived Boards</h2>
      </div>
      <div class="header-right">
        <button class="btn-action-premium" @click="$router.push('/')">
          <ArrowLeftIcon :size="16" /> Dashboard
        </button>
      </div>
    </nav>

    <div class="archived-content">
      <!-- Loading -->
      <div v-if="isLoading" class="state-center">
        <div class="loading-spinner"></div>
        <p>Loading archived boards...</p>
      </div>

      <!-- Boards Grid -->
      <div v-else-if="archivedBoards.length > 0">
        <p class="results-count">
          <ArchiveIcon :size="16" />
          <strong>{{ archivedBoards.length }}</strong> archived board{{ archivedBoards.length !== 1 ? 's' : '' }}
        </p>
        <div class="board-grid">
          <div v-for="board in archivedBoards" :key="board.id" class="board-tile glass-panel">
            <div class="board-tile-content">
              <div class="archived-badge">
                <ArchiveIcon :size="12" /> Archived
              </div>
              <h3 class="board-title-text">{{ board.title }}</h3>
            </div>
            <div class="board-actions">
              <button class="btn-restore" @click="restoreBoard(board.id)" :disabled="restoringId === board.id">
                <RotateCcwIcon :size="14" />
                {{ restoringId === board.id ? 'Restoring...' : 'Restore' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="state-center empty-state">
        <div class="empty-icon-box">
          <ArchiveIcon :size="48" />
        </div>
        <h3>No archived boards</h3>
        <p>When you archive a board, it will appear here. You can restore it at any time.</p>
        <button class="btn-back-dash" @click="$router.push('/')">
          <ArrowLeftIcon :size="16" /> Go to Dashboard
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api';
import { useToast } from 'vue-toastification';
import {
  Archive as ArchiveIcon,
  ArrowLeft as ArrowLeftIcon,
  RotateCcw as RotateCcwIcon
} from 'lucide-vue-next';

const toast = useToast();
const archivedBoards = ref([]);
const isLoading = ref(true);
const restoringId = ref(null);

const fetchArchived = async () => {
  isLoading.value = true;
  try {
    const res = await api.get('/boards?archived=true');
    archivedBoards.value = res.data || [];
  } catch (err) {
    toast.error('Failed to load archived boards');
  } finally {
    isLoading.value = false;
  }
};

const restoreBoard = async (boardId) => {
  restoringId.value = boardId;
  try {
    await api.patch(`/boards/${boardId}/archive`);
    archivedBoards.value = archivedBoards.value.filter(b => b.id !== boardId);
    toast.success('Board restored successfully');
  } catch (err) {
    toast.error('Failed to restore board');
  } finally {
    restoringId.value = null;
  }
};

onMounted(fetchArchived);
</script>

<style scoped>
.archived-page {
  min-height: 100vh;
  background: var(--bg-gradient);
  display: flex;
  flex-direction: column;
}

.page-header {
  padding: 12px 32px;
  background: var(--surface-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left { display: flex; align-items: center; gap: 12px; }
.header-left h2 { color: var(--text-main); margin: 0; font-size: 18px; font-weight: 800; letter-spacing: -0.5px; }
.primary-icon { color: var(--brand-primary); }

.btn-action-premium {
  background: var(--brand-primary-light); color: var(--brand-primary);
  border: none; padding: 8px 16px; border-radius: 10px;
  font-size: 13px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; gap: 8px; transition: all 0.2s;
}
.btn-action-premium:hover { background: var(--brand-primary); color: white; transform: translateY(-1px); }

.archived-content {
  flex: 1; padding: 40px 48px;
}

.results-count {
  font-size: 14px; color: var(--text-muted); font-weight: 600;
  margin-bottom: 24px; display: flex; align-items: center; gap: 8px;
}

.board-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.board-tile {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-height: 140px;
}

.board-tile:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--brand-primary-light);
}

.archived-badge {
  display: inline-flex; align-items: center; gap: 6px;
  font-size: 11px; font-weight: 700; text-transform: uppercase;
  letter-spacing: 0.5px; color: var(--text-extramuted);
  background: var(--surface-secondary); padding: 4px 10px;
  border-radius: 8px; margin-bottom: 8px;
}

.board-title-text {
  font-size: 18px; font-weight: 700; color: var(--text-main);
  letter-spacing: -0.5px; margin: 0;
}

.board-actions { display: flex; justify-content: flex-end; }

.btn-restore {
  background: var(--brand-accent); color: white;
  border: none; padding: 8px 18px; border-radius: 10px;
  font-size: 13px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; gap: 6px; transition: all 0.2s;
  box-shadow: 0 4px 10px rgba(16, 185, 129, 0.2);
}

.btn-restore:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 14px rgba(16, 185, 129, 0.3);
}

.btn-restore:disabled { opacity: 0.7; cursor: not-allowed; }

/* States */
.state-center {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 100px 20px; text-align: center;
}

.loading-spinner {
  width: 40px; height: 40px; border: 4px solid var(--border-subtle);
  border-top-color: var(--brand-primary); border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.state-center p { color: var(--text-muted); font-weight: 600; margin-top: 16px; }

.empty-icon-box {
  width: 100px; height: 100px; background: var(--brand-primary-light);
  border-radius: 28px; display: flex; align-items: center;
  justify-content: center; color: var(--brand-primary); margin-bottom: 24px;
}

.empty-state h3 { font-size: 20px; font-weight: 800; color: var(--text-main); margin-bottom: 8px; }
.empty-state p { font-size: 15px; color: var(--text-muted); max-width: 380px; line-height: 1.6; margin-bottom: 24px; }

.btn-back-dash {
  background: var(--brand-primary-light); color: var(--brand-primary);
  border: none; padding: 10px 24px; border-radius: 12px;
  font-weight: 700; cursor: pointer; display: flex; align-items: center;
  gap: 8px; transition: all 0.2s;
}
.btn-back-dash:hover { background: var(--brand-primary); color: white; }

.animate-fade-in { animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
