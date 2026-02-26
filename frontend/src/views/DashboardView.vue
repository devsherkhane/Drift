<template>
  <div class="dashboard-container">
    <nav class="navbar">
      <div class="logo">Trello Clone</div>
      <button @click="auth.logout" class="btn-logout">Logout</button>
    </nav>

    <main class="content">
      <h2>Your Boards</h2>
      <div class="board-grid">
        <div class="board-tile create-tile" @click="isModalOpen = true">
          <span>Create new board</span>
        </div>

        <div 
          v-for="board in boards" 
          :key="board.id" 
          class="board-tile"
          @click="openBoard(board.id)"
        >
          <h3>{{ board.title }}</h3>
        </div>
      </div>
    </main>

    <div v-if="isModalOpen" class="modal-backdrop" @click.self="isModalOpen = false">
      <div class="modal-content">
        <h3>Create Board</h3>
        <input 
          v-model="newBoardTitle" 
          placeholder="Enter board title" 
          @keyup.enter="handleCreateBoard" 
          ref="modalInput"
        />
        <div class="modal-actions">
          <button class="btn-cancel" @click="isModalOpen = false">Cancel</button>
          <button 
            class="btn-create" 
            :disabled="!newBoardTitle.trim()" 
            @click="handleCreateBoard"
          >
            Create
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import api from '../api';

// State Management
const auth = useAuthStore();
const router = useRouter();
const boards = ref([]);

// Modal State
const isModalOpen = ref(false);
const newBoardTitle = ref('');

// Fetch all boards from Go backend
const fetchBoards = async () => {
  try {
    const response = await api.get('/boards');
    boards.value = response.data;
  } catch (error) {
    console.error("Failed to fetch boards", error);
  }
};

// Navigate to a specific board
const openBoard = (id) => {
  router.push(`/board/${id}`);
};

// Create a new board via API
const handleCreateBoard = async () => {
  if (!newBoardTitle.value.trim()) return;
  
  try {
    const response = await api.post('/boards', { title: newBoardTitle.value });
    // Add the new board to the list immediately
    boards.value.push(response.data); 
    
    // Reset and close modal
    isModalOpen.value = false;
    newBoardTitle.value = '';
  } catch (err) {
    alert("Error creating board. Check your backend connection.");
  }
};

// Load data when component mounts
onMounted(fetchBoards);
</script>

<style scoped>
/* Dashboard & Navbar Styles */
.navbar {
  height: 44px;
  background-color: var(--trello-blue);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  color: white;
}

.logo {
  font-weight: bold;
  font-size: 18px;
}

.btn-logout {
  background: rgba(255,255,255,0.2);
  border: none;
  color: white;
  padding: 6px 12px;
  border-radius: 3px;
  cursor: pointer;
}

.btn-logout:hover {
  background: rgba(255,255,255,0.3);
}

.content {
  padding: 40px;
  max-width: 1200px;
  margin: 0 auto;
}

h2 {
  margin-bottom: 20px;
  color: #172b4d;
}

/* Grid Layout */
.board-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.board-tile {
  background-color: var(--trello-blue);
  height: 100px;
  border-radius: 3px;
  padding: 10px;
  color: white;
  cursor: pointer;
  transition: transform 0.1s, opacity 0.2s;
}

.board-tile:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}

.create-tile {
  background-color: #f0f2f5;
  color: #172b4d;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
  border: 2px dashed #dfe1e6;
}

.create-tile:hover {
  background-color: #e2e4e9;
}

/* Modal Styles */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 24px;
  border-radius: 3px;
  width: 320px;
  box-shadow: 0 12px 24px rgba(0,0,0,0.5);
}

.modal-content h3 {
  margin-bottom: 15px;
  color: #172b4d;
}

.modal-content input {
  width: 100%;
  margin-bottom: 20px;
  padding: 10px;
  border: 2px solid #dfe1e6;
  border-radius: 3px;
}

.modal-content input:focus {
  border-color: #0079bf;
  outline: none;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-create {
  background: #5aac44;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 3px;
  cursor: pointer;
  font-weight: bold;
}

.btn-create:disabled {
  background: #ebecf0;
  cursor: not-allowed;
  color: #a5adba;
}

.btn-cancel {
  background: transparent;
  border: none;
  color: #5e6c84;
  padding: 8px 12px;
  cursor: pointer;
}

.btn-cancel:hover {
  text-decoration: underline;
}
</style>