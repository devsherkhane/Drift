<template>
  <div class="board-canvas">
    <nav class="board-header">
      <div class="header-left">
        <h2>{{ boardTitle }}</h2>
      </div>
      <div class="search-wrapper">
        <input 
          v-model="searchQuery" 
          placeholder="Filter cards..." 
          class="search-input" 
        />
      </div>
    </nav>

    <div class="lists-container">
      <div v-for="list in filteredLists" :key="list.id" class="list-wrapper">
        <div class="list-content">
          <div class="list-header">
            <h4>{{ list.title }}</h4>
          </div>

          <draggable 
            v-model="list.cards" 
            group="cards" 
            item-key="id" 
            class="cards-list"
            ghost-class="ghost-card" 
            @change="onCardMove($event, list.id)"
          >
            <template #item="{ element }">
              <div class="card" @click="openCardDetails(element, list.title)">
                <div 
                  v-if="element.label_color" 
                  class="card-label-bar" 
                  :style="{ backgroundColor: element.label_color }"
                ></div>
                {{ element.title }}
              </div>
            </template>
          </draggable>

          <div v-if="currentlyEditingList === list.id" class="add-card-form">
            <textarea 
              v-model="newCardTitle" 
              placeholder="Enter a title for this card..."
              @keyup.enter="submitCard(list.id)" 
              ref="cardInput"
            ></textarea>
            <div class="form-actions">
              <button class="btn-add" @click="submitCard(list.id)">Add card</button>
              <button class="btn-close" @click="currentlyEditingList = null">✕</button>
            </div>
          </div>

          <button v-else class="add-card-btn" @click="openAddCard(list.id)">
            + Add a card
          </button>
        </div>
      </div>

      <div class="list-wrapper">
        <div v-if="!isAddingList" class="add-list-placeholder" @click="isAddingList = true">
          + Add another list
        </div>
        <div v-else class="add-list-form">
          <input 
            v-model="newListTitle" 
            placeholder="Enter list title..." 
            @keyup.enter="submitList" 
            autofocus 
          />
          <div class="form-actions">
            <button class="btn-add" @click="submitList">Add list</button>
            <button class="btn-close" @click="isAddingList = false">✕</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="activeCard" class="modal-backdrop" @click.self="activeCard = null">
      <div class="card-detail-modal">
        <button class="modal-close" @click="activeCard = null">✕</button>

        <div class="modal-section">
          <h3>💳 {{ activeCard.title }}</h3>
          <p class="subtitle">in list <u>{{ activeCard.listName }}</u></p>
        </div>

        <div class="modal-section">
          <h4>Labels</h4>
          <div class="labels-picker">
            <div 
              v-for="color in ['#61bd4f', '#f2d600', '#ff9f1a', '#eb5a46', '#c377e0', '#0079bf']" 
              :key="color"
              class="label-swatch"
              :style="{ backgroundColor: color }"
              :class="{ active: activeCard.label_color === color }"
              @click="updateCardLabel(color)"
            ></div>
            <button class="btn-clear-label" @click="updateCardLabel(null)">None</button>
          </div>
        </div>

        <div class="modal-section">
          <h4>Description</h4>
          <textarea 
            v-model="activeCard.description" 
            placeholder="Add a more detailed description..."
            @blur="updateCardDescription"
          ></textarea>
        </div>

        <div class="modal-section comments-section">
          <h4>Activity</h4>
          <div class="comment-input-area">
            <textarea v-model="newCommentText" placeholder="Write a comment..."></textarea>
            <button 
              class="btn-add" 
              @click="submitComment" 
              :disabled="!newCommentText.trim()"
            >Save</button>
          </div>

          <div class="comments-list">
            <div v-for="comment in activeCard.comments" :key="comment.id" class="comment-item">
              <div class="comment-avatar">{{ comment.user_name?.charAt(0) || 'U' }}</div>
              <div class="comment-content">
                <span class="comment-user">{{ comment.user_name }}</span>
                <p>{{ comment.text }}</p>
                <span class="comment-date">{{ new Date(comment.created_at).toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import draggable from 'vuedraggable';
import api from '../api';
import { useWebsocket } from '../api/websocket';

const { connect, disconnect } = useWebsocket();
const route = useRoute();

const boardTitle = ref('Loading...');
const lists = ref([]);
const searchQuery = ref('');

// UI State
const currentlyEditingList = ref(null);
const newCardTitle = ref('');
const isAddingList = ref(false);
const newListTitle = ref('');
const activeCard = ref(null);
const newCommentText = ref('');

// --- COMPUTED: Real-time Search ---
const filteredLists = computed(() => {
  if (!searchQuery.value.trim()) return lists.value;
  const query = searchQuery.value.toLowerCase();
  return lists.value.map(list => ({
    ...list,
    cards: list.cards?.filter(card => 
      card.title.toLowerCase().includes(query) || 
      card.description?.toLowerCase().includes(query)
    )
  }));
});

// --- API ACTIONS ---

const fetchBoardDetails = async () => {
  const boardId = route.params.id;
  try {
    const response = await api.get(`/boards/${boardId}`);
    boardTitle.value = response.data.title;
    lists.value = response.data.lists || [];
  } catch (err) {
    console.error("Failed to fetch board details:", err);
  }
};

const updateCardLabel = async (color) => {
  activeCard.value.label_color = color;
  try {
    await api.patch(`/cards/${activeCard.value.id}`, { label_color: color });
    // Local update to keep it snappy
    const boardId = route.params.id;
    fetchBoardDetails(); 
  } catch (err) { console.error("Label update failed"); }
};

const submitComment = async () => {
  try {
    const response = await api.post(`/cards/${activeCard.value.id}/comments`, {
      text: newCommentText.value
    });
    if (!activeCard.value.comments) activeCard.value.comments = [];
    activeCard.value.comments.unshift(response.data);
    newCommentText.value = '';
  } catch (err) {
    alert("Failed to post comment");
  }
};

const submitList = async () => {
  if (!newListTitle.value.trim()) return;
  try {
    const response = await api.post(`/boards/${route.params.id}/lists`, {
      title: newListTitle.value.trim()
    });
    lists.value.push({ ...response.data, cards: [] });
    newListTitle.value = '';
    isAddingList.value = false;
  } catch (err) { alert("Could not create list."); }
};

const openAddCard = (listId) => {
  currentlyEditingList.value = listId;
  newCardTitle.value = '';
};

const submitCard = async (listId) => {
  if (!newCardTitle.value.trim()) return;
  try {
    const response = await api.post('/cards', {
      list_id: listId,
      title: newCardTitle.value.trim()
    });
    const targetList = lists.value.find(l => l.id === listId);
    if (!targetList.cards) targetList.cards = [];
    targetList.cards.push(response.data);
    newCardTitle.value = '';
    currentlyEditingList.value = null;
  } catch (err) { alert("Could not save card."); }
};

const onCardMove = async (event, newListId) => {
  if (event.added || event.moved) {
    const card = event.added ? event.added.element : event.moved.element;
    const newIndex = event.added ? event.added.newIndex : event.moved.newIndex;
    try {
      await api.patch(`/cards/${card.id}/move`, {
        list_id: newListId,
        position: newIndex
      });
    } catch (err) {
      console.error("Move sync failed:", err);
      fetchBoardDetails();
    }
  }
};

const openCardDetails = (card, listName) => {
  activeCard.value = { ...card, listName };
};

const updateCardDescription = async () => {
  try {
    await api.patch(`/cards/${activeCard.value.id}`, {
      description: activeCard.value.description
    });
  } catch (err) { console.error("Update failed"); }
};

// Lifecycle
onMounted(() => {
  fetchBoardDetails();
  connect(route.params.id);
});

onUnmounted(() => {
  disconnect();
});
</script>

<style scoped>
/* Main Canvas */
.board-canvas {
  height: 100vh;
  background-color: var(--trello-blue);
  display: flex;
  flex-direction: column;
}

.board-header {
  padding: 10px 20px;
  background: rgba(0, 0, 0, 0.15);
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-input {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 3px;
  padding: 6px 12px;
  color: white;
  width: 250px;
  transition: width 0.3s;
}

.search-input:focus { background: white; color: #172b4d; outline: none; width: 300px; }

.lists-container {
  display: flex;
  align-items: flex-start;
  padding: 10px;
  gap: 12px;
  overflow-x: auto;
  flex-grow: 1;
}

.list-wrapper { width: 272px; flex-shrink: 0; }

.list-content {
  background: var(--trello-gray);
  border-radius: 3px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.card {
  background: white;
  padding: 10px;
  margin-bottom: 8px;
  border-radius: 3px;
  box-shadow: 0 1px 0 rgba(9, 30, 66, 0.25);
  cursor: grab;
  font-size: 14px;
  position: relative;
}

.card-label-bar {
  height: 4px; width: 40px; border-radius: 2px; margin-bottom: 4px;
}

.labels-picker { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 10px; }
.label-swatch {
  width: 40px; height: 32px; border-radius: 3px; cursor: pointer;
}
.label-swatch.active { border: 2px solid #172b4d; }

.btn-clear-label { background: none; border: 1px solid #ddd; padding: 4px 8px; border-radius: 3px; cursor: pointer; }

.modal-backdrop {
  position: fixed;
  top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.6);
  display: flex; justify-content: center; align-items: center;
  z-index: 100;
}

.card-detail-modal {
  background: #f4f5f7;
  width: 600px;
  padding: 30px;
  border-radius: 3px;
  position: relative;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-close { position: absolute; top: 15px; right: 15px; border: none; font-size: 20px; cursor: pointer; }

.comment-item { display: flex; gap: 12px; margin-bottom: 16px; }
.comment-avatar {
  width: 32px; height: 32px; background: #dfe1e6;
  border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: bold;
}

.comment-content p {
  background: white; padding: 8px 12px; border-radius: 3px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}
</style>