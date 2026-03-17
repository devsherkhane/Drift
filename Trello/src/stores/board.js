import { defineStore } from 'pinia';
import api from '../api';

export const useBoardStore = defineStore('board', {
  state: () => ({
    boardId: null,
    boardTitle: 'Loading...',
    lists: [],
    isLoading: false,
  }),

  actions: {
    async fetchBoardDetails(boardId) {
      this.boardId = boardId;
      this.isLoading = true;
      try {
        const response = await api.get(`/boards/${boardId}`);
        this.boardTitle = response.data.title;
        // Ensure every list has a cards array to prevent draggable errors
        this.lists = response.data.lists?.map(list => ({
          ...list,
          cards: list.cards || []
        })) || [];
      } catch (error) {
        console.error("Failed to load board", error);
        throw error;
      } finally {
        this.isLoading = false;
      }
    },

    // --- SURGICAL WEBSOCKET UPDATES ---
    
    // Updates a specific card's data without re-fetching the board
    updateCardLocally(updatedCard) {
      for (const list of this.lists) {
        const cardIndex = list.cards.findIndex(c => c.id === updatedCard.id);
        if (cardIndex !== -1) {
          // Merge new data into the existing card
          list.cards[cardIndex] = { ...list.cards[cardIndex], ...updatedCard };
          return;
        }
      }
    },

    // Handles drag-and-drop from another user
    moveCardLocally(cardId, fromListId, toListId, newPosition, updatedCardData = null) {
      let movedCard = null;

      // 1. Remove from source list
      const sourceList = this.lists.find(l => l.id === fromListId);
      if (sourceList) {
        const cardIndex = sourceList.cards.findIndex(c => c.id === cardId);
        if (cardIndex !== -1) {
          movedCard = sourceList.cards.splice(cardIndex, 1)[0];
        }
      }

      // 2. Insert into destination list
      const destList = this.lists.find(l => l.id === toListId);
      if (destList && movedCard) {
        // If the backend sent fresh card data with the WS event, use it
        const finalCard = updatedCardData ? { ...movedCard, ...updatedCardData } : movedCard;
        destList.cards.splice(newPosition, 0, finalCard);
      }
    }
  }
});