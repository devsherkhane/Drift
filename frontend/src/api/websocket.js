import { useBoardStore } from '../stores/board'; // We'll move board logic here soon

export const useWebsocket = () => {
  let socket = null;

  const connect = (boardId) => {
    const token = localStorage.getItem('token');
    // Using the WebSocket endpoint we built in Go
    socket = new WebSocket(`ws://localhost:8080/ws?boardID=${boardId}&token=${token}`);

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      const boardStore = useBoardStore();

      // Handle different real-time events
      switch (data.type) {
        case 'card_moved':
        case 'card_added':
        case 'comment_added':
          boardStore.fetchBoardDetails(boardId); // Refresh data
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