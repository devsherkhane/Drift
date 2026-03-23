<template>
  <main class="content animate-fade-in">
    <nav class="board-header">
      <div class="header-left">
        <h1 class="brand-name">Trello Clone</h1>
      </div>
      <div class="header-right">
        <router-link to="/search" class="btn-icon-nav" title="Search">
          <SearchIcon :size="20" />
        </router-link>
        <router-link to="/activity" class="btn-icon-nav" title="Activity Feed">
          <ActivityIcon :size="20" />
        </router-link>
        <router-link to="/templates" class="btn-icon-nav" title="Templates">
          <CopyIcon :size="20" />
        </router-link>
        <router-link to="/archived" class="btn-icon-nav" title="Archived Boards">
          <ArchiveIcon :size="20" />
        </router-link>
        <router-link to="/profile" class="btn-icon-nav" title="Profile Settings">
          <SettingsIcon :size="20" />
        </router-link>
        <button @click="logout" class="btn-icon-nav btn-logout-icon" title="Logout">
          <LogOutIcon :size="20" />
        </button>
      </div>
    </nav>

    <div v-if="pendingInvitations.length > 0" class="dashboard-section invitations-section">
      <h2 class="section-title">
        <BellIcon :size="22" class="title-icon" />
        Pending Invitations
      </h2>
      <div class="invitation-grid">
        <div v-for="invite in pendingInvitations" :key="invite.id" class="invitation-card glass-panel">
          <div class="invite-info">
            <span class="invite-label">Board Invite</span>
            <h3 class="invite-title">{{ invite.title }}</h3>
          </div>
          <div class="invite-actions">
            <button class="btn-invite btn-accept" @click="respondToInvite(invite.id, 'accept')">Accept</button>
            <button class="btn-invite btn-decline" @click="respondToInvite(invite.id, 'decline')">Decline</button>
          </div>
        </div>
      </div>
    </div>

    <div class="dashboard-section">
      <h2 class="section-title">
        <LayoutDashboardIcon :size="22" class="title-icon" />
        My Boards <span class="badge">{{ myBoards.length }}</span>
      </h2>
      <div class="board-grid">
        <div class="board-tile create-tile" @click="isModalOpen = true">
          <div class="create-tile-content">
            <div class="plus-circle">
              <PlusIcon :size="32" />
            </div>
            <span>Create new board</span>
          </div>
        </div>

        <div v-for="board in myBoards" :key="board.id" class="board-tile glass-panel" @click="openBoard(board.id)">
          <div class="board-tile-content">
            <h3 class="board-title-text">{{ board.title }}</h3>
            <div class="board-meta">
              <span class="meta-item"><UsersIcon :size="14" /> Owner</span>
            </div>
          </div>
          <div class="board-actions">
            <button class="btn-action-small" @click.stop="editBoard(board)" title="Edit Board">
              <Edit2Icon :size="14" />
            </button>
            <button class="btn-action-small btn-delete" @click.stop="deleteBoard(board.id)" title="Delete Board">
              <Trash2Icon :size="14" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="sharedBoards.length > 0" class="dashboard-section">
      <h2 class="section-title">
        <UsersIcon :size="22" class="title-icon" />
        Shared with Me <span class="badge">{{ sharedBoards.length }}</span>
      </h2>
      <div class="board-grid">
        <div v-for="board in sharedBoards" :key="board.id" class="board-tile glass-panel shared-tile" @click="openBoard(board.id)">
          <div class="board-tile-content">
            <div class="shared-badge">Shared</div>
            <h3 class="board-title-text">{{ board.title }}</h3>
          </div>
        </div>
      </div>
    </div>
  </main>

  <Transition name="modal">
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-window glass-panel shadow-2xl">
        <div class="modal-header">
          <h3>{{ editingBoard ? 'Edit Board' : 'Create Board' }}</h3>
          <button class="btn-close-modal" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <label class="input-label">Board Title</label>
          <input 
            v-model="boardFormTitle" 
            placeholder="e.g. Project Roadmap" 
            @keyup.enter="submitBoardForm" 
            ref="modalInput"
            autofocus
          />
        </div>
        <div class="modal-footer">
          <button class="btn-secondary-flat" @click="closeModal">Cancel</button>
          <button class="btn-primary-glow" :disabled="!boardFormTitle.trim()" @click="submitBoardForm">
            {{ editingBoard ? 'Save Changes' : 'Create Board' }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import api from '../api';
import { useToast } from "vue-toastification";

// Lucide Icons
import { 
  LayoutDashboard as LayoutDashboardIcon, 
  Plus as PlusIcon, 
  Settings as SettingsIcon, 
  LogOut as LogOutIcon, 
  Bell as BellIcon,
  Users as UsersIcon,
  Edit2 as Edit2Icon,
  Trash2 as Trash2Icon,
  Search as SearchIcon,
  Activity as ActivityIcon,
  Copy as CopyIcon,
  Archive as ArchiveIcon
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();
const toast = useToast();
const boards = ref([]);

const isModalOpen = ref(false);
const boardFormTitle = ref('');
const editingBoard = ref(null);

const myBoards = computed(() => boards.value.filter(b => b.is_owner));
const sharedBoards = computed(() => boards.value.filter(b => !b.is_owner && b.status === 'accepted'));
const pendingInvitations = computed(() => boards.value.filter(b => !b.is_owner && b.status === 'pending'));

const fetchBoards = async () => {
  try {
    const response = await api.get('/boards');
    boards.value = response.data || []; 
  } catch (error) { 
    toast.error("Failed to fetch boards"); 
  }
};

const logout = () => {
  auth.logout();
};

const openBoard = (id) => router.push(`/board/${id}`);

const respondToInvite = async (boardId, action) => {
  try {
    await api.patch(`/boards/${boardId}/invitation`, { action });
    toast.success(`Invitation ${action}ed`);
    await fetchBoards();
  } catch (err) {
    toast.error(`Failed to ${action} invitation`);
  }
};

const editBoard = (board) => {
  editingBoard.value = board;
  boardFormTitle.value = board.title;
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
  editingBoard.value = null;
  boardFormTitle.value = '';
};

const submitBoardForm = async () => {
  if (!boardFormTitle.value.trim()) return;

  try {
    if (editingBoard.value) {
      await api.put(`/boards/${editingBoard.value.id}`, { title: boardFormTitle.value });
      const index = boards.value.findIndex(b => b.id === editingBoard.value.id);
      boards.value[index].title = boardFormTitle.value;
      toast.success("Board updated");
    } else {
      const response = await api.post('/boards', { title: boardFormTitle.value });
      if (!boards.value) boards.value = []; 
      boards.value.push(response.data);
      toast.success("Board created");
    }
    closeModal();
  } catch (err) { toast.error("Error saving board"); }
};

const deleteBoard = async (id) => {
  if (!confirm("Are you sure you want to delete this board?")) return;
  try {
    await api.delete(`/boards/${id}`);
    boards.value = boards.value.filter(b => b.id !== id);
    toast.success("Board deleted");
  } catch (err) { toast.error("Error deleting board"); }
};

onMounted(fetchBoards);
</script>

<style scoped>
.content {
  padding: 0;
  min-height: 100vh;
  background: var(--bg-gradient);
}

.board-header {
  padding: 16px 48px;
  background: var(--surface-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 50;
  box-shadow: var(--shadow-sm);
}

.brand-name {
  color: var(--brand-primary);
  margin: 0;
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -1px;
}

.btn-icon-nav {
  background: transparent;
  color: var(--text-muted);
  border: 1px solid transparent;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.btn-icon-nav:hover {
  background: var(--brand-primary-light);
  color: var(--brand-primary);
  transform: translateY(-2px);
}

.btn-logout-icon:hover {
  background: #fee2e2;
  color: #ef4444;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.dashboard-section {
  margin-bottom: 48px;
  padding: 0 48px;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-main);
  margin: 40px 0 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  color: var(--brand-primary);
  opacity: 0.8;
}

.badge {
  background: var(--brand-primary-light);
  color: var(--brand-primary);
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 20px;
  font-weight: 600;
}

.board-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.board-tile {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--border-radius);
  padding: 24px;
  height: 160px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--text-main);
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.board-tile:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
  border-color: var(--brand-primary);
}

.board-tile-content {
  flex-grow: 1;
}

.board-title-text {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.board-meta {
  margin-top: 8px;
  display: flex;
  gap: 12px;
}

.meta-item {
  font-size: 12px;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 4px;
}

.board-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s ease, transform 0.2s ease;
  transform: translateY(10px);
}

.board-tile:hover .board-actions {
  opacity: 1;
  transform: translateY(0);
}

.btn-action-small {
  background: var(--surface-secondary);
  border: 1px solid var(--border-subtle);
  color: var(--text-muted);
  border-radius: 8px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-action-small:hover {
  background: var(--brand-primary);
  color: white;
  border-color: var(--brand-primary);
}

.btn-action-small.btn-delete:hover {
  background: #ef4444;
  border-color: #ef4444;
}

.create-tile {
  background: transparent;
  border: 2px dashed var(--brand-primary);
  opacity: 0.7;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.create-tile:hover {
  opacity: 1;
  background: var(--brand-primary-light);
  border-style: solid;
}

.create-tile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--brand-primary);
}

.plus-circle {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--brand-primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.create-tile:hover .plus-circle {
  transform: scale(1.1) rotate(90deg);
  background: var(--brand-primary);
  color: white;
}

/* Modal Enhancements */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-window {
  background: var(--surface-primary);
  width: 95%;
  max-width: 440px;
  border-radius: 20px;
  overflow: hidden;
  padding: 0;
}

.modal-header {
  padding: 24px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-subtle);
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
}

.btn-close-modal {
  background: transparent;
  border: none;
  font-size: 24px;
  color: var(--text-muted);
  cursor: pointer;
}

.modal-body {
  padding: 32px;
}

.input-label {
  display: block;
  font-size: 13px;
  font-weight: 700;
  color: var(--text-muted);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.modal-body input {
  width: 100%;
  padding: 14px 18px;
  background: var(--surface-secondary);
  border: 2px solid transparent;
  border-radius: 12px;
  font-size: 16px;
  color: var(--text-main);
  transition: all 0.2s;
}

.modal-body input:focus {
  outline: none;
  border-color: var(--brand-primary);
  background: var(--surface-primary);
  box-shadow: 0 0 0 4px var(--brand-primary-light);
}

.modal-footer {
  padding: 24px 32px;
  background: var(--surface-secondary);
  display: flex;
  justify-content: flex-end;
  gap: 16px;
}

.btn-primary-glow {
  background: var(--brand-primary);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.btn-primary-glow:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

.btn-secondary-flat {
  background: transparent;
  color: var(--text-muted);
  border: none;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
}

/* Transitions */
.modal-enter-active, .modal-leave-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.modal-enter-from, .modal-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(20px);
}

.animate-fade-in {
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
