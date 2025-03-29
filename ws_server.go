package fall

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	connections map[string]*websocket.Conn
	mu          sync.Mutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		connections: make(map[string]*websocket.Conn),
	}
}

func (s *WebSocketServer) AddConnection(clientID string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[clientID] = conn
	fmt.Printf("Client %s connected\n", clientID)
}

func (s *WebSocketServer) RemoveConnection(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.connections[clientID]; ok {
		fmt.Printf("Client %s disconnected\n", clientID)
		delete(s.connections, clientID)
	}
}

func (s *WebSocketServer) SendMessage(clientID string, message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, ok := s.connections[clientID]
	if !ok {
		fmt.Printf("Client %s not found\n", clientID)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		fmt.Printf("Error sending message to client %s: %v\n", clientID, err)
		delete(s.connections, clientID)
	}
}
