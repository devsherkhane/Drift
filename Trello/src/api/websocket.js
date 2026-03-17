import { useBoardStore } from '../stores/board'; // We'll move board logic here soon

export const useWebsocket = () => {
  let socket = null;

  const connect = (boardId) => {
    const token = localStorage.getItem('token');
    socket = new WebSocket(`ws://localhost:8080/api/ws?boardID=${boardId}&token=${token}`);

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      const boardStore = useBoardStore();

      // IMPORTANT: Your Go backend needs to attach a "sender_id" to the websocket payload.
      // If the event was caused by the current user, we IGNORE the websocket event 
      // because vuedraggable already updated their screen instantly!
      const currentUserId = JSON.parse(localStorage.getItem('user'))?.id; // Assuming you store user info
      if (data.sender_id === currentUserId) return;

      switch (data.type) {
        case 'card_moved':
          // Assuming your Go backend sends the old list, new list, and new position
          if (data.card_id && data.from_list_id && data.to_list_id) {
            boardStore.moveCardLocally(data.card_id, data.from_list_id, data.to_list_id, data.position, data.card);
          } else {
            // Fallback if payload is incomplete
            boardStore.fetchBoardDetails(boardId);
          }
          break;

        case 'card_added':
        case 'comment_added':
          // For now, fallback to refetch, but you can add surgical methods for these too!
          boardStore.fetchBoardDetails(boardId);
          break;

        case 'user_online':
          console.log(`User ${data.user_id} joined the board`);
          break;
      }
    };

    socket.onclose = () => console.log("WebSocket Disconnected");
  };

  const disconnect = () => {
    if (socket) socket.close();
  };

  return { connect, disconnect };
};