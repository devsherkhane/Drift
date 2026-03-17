<template>
  <div class="board-canvas animate-fade-in">
    <BoardHeader 
      :boardTitle="boardTitle"
      :searchResults="searchResults"
      :isSearching="isSearching"
      @search="performSearch"
      @open-share="isShareModalOpen = true"
      @toggle-menu="toggleMenu"
      @open-card="openCardDetails"
      @logout="logout"
    />

    <div v-if="isLoading" class="lists-container">
      <div v-for="i in 3" :key="i" class="list-wrapper skeleton-wrapper">
        <div class="list-content skeleton-list glass-panel">
          <div class="skeleton-header"></div>
          <div v-for="j in 4" :key="j" class="skeleton-card"></div>
        </div>
      </div>
    </div>

    <div v-else class="lists-container">
      <BoardList 
        v-for="list in filteredLists" 
        :key="list.id" 
        :list="list"
        :isEditing="currentlyEditingList === list.id"
        @update-title="updateListTitle"
        @delete-list="deleteList"
        @update-cards="cards => list.cards = cards"
        @card-moved="onCardMove"
        @open-card="openCardDetails"
        @start-add="currentlyEditingList = list.id"
        @cancel-add="currentlyEditingList = null"
        @submit-card="submitCard"
      />

      <div class="list-wrapper">
        <div v-if="!isAddingList" class="add-list-placeholder glass-panel" @click="isAddingList = true">
          <PlusIcon :size="18" /> Add another list
        </div>
        <div v-else class="add-list-form glass-panel">
          <input v-model="newListTitle" placeholder="Enter list title..." @keyup.enter="submitList" autofocus />
          <div class="form-actions">
            <button class="btn-add-primary" @click="submitList">Add list</button>
            <button class="btn-cancel-small" @click="isAddingList = false">
              <XIcon :size="18" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <CardDetailModal 
      v-if="activeCard"
      :card="activeCard"
      :userInitial="authStore.user?.username?.charAt(0).toUpperCase()"
      @close="activeCard = null"
      @delete="deleteCard"
      @update-field="updateCardDetails"
      @upload-file="handleFileUpload"
      @submit-comment="submitComment"
    />

    <BoardMenu 
      :isOpen="isMenuOpen"
      :activityLogs="activityLogs"
      @close="isMenuOpen = false"
      @archive-board="archiveBoard"
    />

    <ShareModal 
      :isOpen="isShareModalOpen"
      @close="isShareModalOpen = false"
      @invite="inviteCollaborator"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../api';
import { useWebsocket } from '../api/websocket';
import { useToast } from "vue-toastification";
import { useBoardStore } from '../stores/board';
import { useAuthStore } from '../stores/auth';

// Icons
import { Plus as PlusIcon, X as XIcon } from 'lucide-vue-next';

// Components
import BoardHeader from '../components/BoardHeader.vue';
import BoardList from '../components/BoardList.vue';
import CardDetailModal from '../components/CardDetailModal.vue';
import BoardMenu from '../components/BoardMenu.vue';
import ShareModal from '../components/ShareModal.vue';

const toast = useToast();
const { connect, disconnect } = useWebsocket();
const route = useRoute();
const router = useRouter();

// Store Initialization
const boardStore = useBoardStore();
const authStore = useAuthStore();

// UI State
const searchResults = ref([]);
const isSearching = ref(false);

const currentlyEditingList = ref(null);
const isAddingList = ref(false);
const newListTitle = ref('');
const activeCard = ref(null);

// Menu State
const isMenuOpen = ref(false);
const activityLogs = ref([]);

// Share State
const isShareModalOpen = ref(false);

// Computed properties mapping to Pinia state
const boardTitle = computed(() => boardStore.boardTitle);
const isLoading = computed(() => boardStore.isLoading);
const lists = computed(() => boardStore.lists);

// The actual lists displayed
const filteredLists = computed(() => lists.value);

// --- Search ---
const performSearch = async (query) => {
  if (!query) {
    searchResults.value = [];
    isSearching.value = false;
    return;
  }
  isSearching.value = true;
  try {
    const response = await api.get(`/search/advanced?q=${encodeURIComponent(query)}&board_id=${route.params.id}`);
    searchResults.value = response.data || [];
  } catch (err) {
    console.error("Search failed:", err);
    searchResults.value = [];
  } finally {
    isSearching.value = false;
  }
};

// Logout Helper
const logout = () => {
  authStore.logout();
  router.push('/login');
};

// --- Card Updates ---
const updateCardDetails = async (field, value) => {
  if (!activeCard.value) return;

  activeCard.value[field] = value;
  boardStore.updateCardLocally({ id: activeCard.value.id, [field]: value });

  try {
    await api.put(`/cards/${activeCard.value.id}`, {
      title: activeCard.value.title,
      description: activeCard.value.description,
      due_date: activeCard.value.due_date,
      label_color: activeCard.value.label_color
    });
    if (field === 'label_color') toast.success("Label updated");
  } catch (err) {
    toast.error(`Failed to update ${field}`);
    boardStore.fetchBoardDetails(route.params.id); 
  }
};

const onCardMove = async (event, newListId) => {
  const item = event.added || event.moved;
  if (!item) return;

  const card = item.element;
  const newIndex = item.newIndex;

  try {
    await api.patch(`/cards/${card.id}/move`, {
      new_list_id: newListId,
      new_position: newIndex
    });
  } catch (err) {
    toast.error("Failed to sync card position");
    boardStore.fetchBoardDetails(route.params.id);
  }
};

// --- List Actions ---
const deleteList = async (listId) => {
  if (!confirm("Are you sure you want to delete this list and all its cards?")) return;
  try {
    await api.delete(`/lists/${listId}`);
    boardStore.lists = boardStore.lists.filter(l => l.id !== listId);
    toast.success("List deleted");
  } catch (err) { toast.error("Failed to delete list"); }
};

const updateListTitle = async (listId, newTitle) => {
  if (!newTitle.trim()) return;
  try {
    await api.put(`/lists/${listId}`, { title: newTitle });
    toast.success("List title updated");
  } catch (err) { toast.error("Failed to update list"); }
};

const submitList = async () => {
  if (!newListTitle.value.trim()) return;
  try {
    const response = await api.post(`/lists`, { board_id: parseInt(route.params.id), title: newListTitle.value.trim() });
    boardStore.lists.push({ ...response.data, cards: [] });
    newListTitle.value = '';
    isAddingList.value = false;
  } catch (err) { toast.error("List creation failed"); }
};

// --- Card Actions ---
const submitCard = async (listId, title) => {
  console.log("Submitting card:", { listId, title });
  try {
    const response = await api.post('/cards', { list_id: listId, title });
    console.log("Card created successfully:", response.data);
    const targetList = boardStore.lists.find(l => l.id === listId);
    if (targetList) {
      if (!targetList.cards) targetList.cards = [];
      targetList.cards.push(response.data);
    }
    currentlyEditingList.value = null;
  } catch (err) { 
    console.error("Card creation failed error:", err);
    toast.error("Card creation failed"); 
  }
};

const deleteCard = async () => {
  if (!confirm("Are you sure you want to delete this card?")) return;
  try {
    await api.delete(`/cards/${activeCard.value.id}`);
    const list = boardStore.lists.find(l => l.title === activeCard.value.listName);
    if (list) list.cards = list.cards.filter(c => c.id !== activeCard.value.id);
    activeCard.value = null;
    toast.success("Card deleted");
  } catch (err) { toast.error("Failed to delete card"); }
};

const openCardDetails = async (card, listName) => {
  const parsedDate = card.due_date ? new Date(card.due_date).toISOString().split('T')[0] : null;

  activeCard.value = { 
    ...card, 
    listName,
    due_date: parsedDate,
    attachments: [],
    comments: []
  };

  try {
    const [attachmentsRes, commentsRes] = await Promise.all([
      api.get(`/cards/${card.id}/attachments`),
      api.get(`/cards/${card.id}/comments`)
    ]);
    activeCard.value.attachments = attachmentsRes.data || [];
    activeCard.value.comments = commentsRes.data || [];
  } catch (err) {
    console.error("Failed to load card secondary data", err);
  }
};

// --- Attachments ---
const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;
  const formData = new FormData();
  formData.append('file', file); 
  formData.append('card_id', activeCard.value.id);
  
  try {
    const response = await api.post(`/attachments`, formData, { headers: { 'Content-Type': 'multipart/form-data' } });
    if (!activeCard.value.attachments) activeCard.value.attachments = [];
    activeCard.value.attachments.push(response.data);
    toast.success("File attached!");
  } catch (err) { toast.error("Upload failed"); } 
  finally { event.target.value = ''; }
};

// --- Comments ---
const submitComment = async (text) => {
  try {
    const response = await api.post(`/comments`, { card_id: activeCard.value.id, text });
    if (!activeCard.value.comments) activeCard.value.comments = [];
    activeCard.value.comments.unshift(response.data);
  } catch (err) { toast.error("Comment failed"); }
};

// --- Board Menu Actions ---
const toggleMenu = async () => {
  isMenuOpen.value = !isMenuOpen.value;
  if (isMenuOpen.value) {
    try {
      const response = await api.get(`/boards/${route.params.id}/activity`);
      activityLogs.value = response.data || [];
    } catch (err) { toast.error("Failed to fetch activity logs"); }
  }
};

const archiveBoard = async () => {
  if (!confirm("Are you sure you want to archive this board? It will no longer appear on your dashboard.")) return;
  try {
    await api.patch(`/boards/${route.params.id}/archive`);
    toast.success("Board archived");
    router.push('/');
  } catch (err) { toast.error("Failed to archive board"); }
};

// --- Collaboration Actions ---
const inviteCollaborator = async (email) => {
  try {
    await api.post(`/boards/${route.params.id}/collaborators`, { email, role: 'member' });
    toast.success(`Invite sent to ${email}`);
    isShareModalOpen.value = false;
  } catch (err) { toast.error(err.response?.data?.error || "Failed to invite collaborator"); }
};

onMounted(async () => {
  await boardStore.fetchBoardDetails(route.params.id);
  connect(route.params.id);
});

onUnmounted(() => disconnect());
</script>

<style scoped>
.board-canvas {
  height: 100vh;
  background: var(--bg-gradient);
  display: flex;
  flex-direction: column;
}

.lists-container {
  display: flex; align-items: flex-start; padding: 24px 32px; gap: 20px;
  overflow-x: auto; flex-grow: 1; scrollbar-width: thin; scrollbar-color: rgba(99, 102, 241, 0.2) transparent;
}
.lists-container::-webkit-scrollbar { height: 10px; }
.lists-container::-webkit-scrollbar-track { background: rgba(0, 0, 0, 0.05); border-radius: 10px; }
.lists-container::-webkit-scrollbar-thumb {
  background-color: rgba(99, 102, 241, 0.2); border-radius: 10px;
  border: 2px solid transparent; background-clip: padding-box;
}
.lists-container::-webkit-scrollbar-thumb:hover { background-color: rgba(99, 102, 241, 0.4); }

.list-wrapper { width: 300px; flex-shrink: 0; }

.btn-add-primary {
  background: var(--brand-primary); color: white;
  border: none; padding: 10px 18px; border-radius: 10px;
  font-weight: 700; cursor: pointer; transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}
.btn-add-primary:hover { background: var(--brand-primary-hover); transform: translateY(-1px); }

.btn-cancel-small {
  background: transparent; border: none; padding: 8px; border-radius: 8px;
  color: var(--text-muted); cursor: pointer; display: flex; align-items: center; justify-content: center;
}
.btn-cancel-small:hover { background: var(--surface-secondary); color: var(--text-main); }

.add-list-placeholder {
  padding: 14px 18px; border-radius: 14px;
  color: var(--text-muted); font-weight: 600; cursor: pointer; 
  display: flex; align-items: center; gap: 8px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(8px);
  border: 1px dashed var(--brand-primary);
}
.add-list-placeholder:hover { 
  background: var(--surface-primary); 
  color: var(--brand-primary); 
  border-style: solid;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.add-list-form { 
  padding: 16px; border-radius: 16px;
  background: var(--surface-primary);
}

.add-list-form input {
  width: 100%; border: 2px solid var(--surface-secondary); border-radius: 10px;
  padding: 12px; margin-bottom: 12px; box-sizing: border-box; font-family: inherit; font-size: 14px;
  transition: all 0.2s;
}
.add-list-form input:focus { outline: none; border-color: var(--brand-primary); box-shadow: 0 0 0 3px var(--brand-primary-light); }
.form-actions { display: flex; align-items: center; gap: 10px; }

/* Skeletons */
.skeleton-list {
  height: 480px; background: rgba(255, 255, 255, 0.4);
  backdrop-filter: blur(10px); padding: 16px; border-radius: 16px;
}
.skeleton-header { height: 24px; width: 50%; background: rgba(0,0,0,0.05); margin-bottom: 24px; border-radius: 6px; }
.skeleton-card { height: 60px; background: white; margin-bottom: 12px; border-radius: 10px; box-shadow: 0 2px 4px rgba(0,0,0,0.02); }
.skeleton-wrapper { animation: pulse 2s infinite ease-in-out; }
@keyframes pulse { 0% { opacity: 0.8; } 50% { opacity: 0.5; } 100% { opacity: 0.8; } }

.animate-fade-in {
  animation: fadeIn 0.5s ease-out forwards;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>