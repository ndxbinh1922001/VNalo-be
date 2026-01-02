package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains active WebSocket connections
type Hub struct {
	// Registered clients (userID -> connection)
	clients map[int64]*Client

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to clients
	broadcast chan *BroadcastMessage

	mu sync.RWMutex
}

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	RecipientIDs []int64
	Data         []byte
	Type         string
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan *BroadcastMessage, 1024),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("WebSocket: User %d connected. Total connections: %d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
				log.Printf("WebSocket: User %d disconnected. Total connections: %d", client.UserID, len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, userID := range message.RecipientIDs {
				if client, ok := h.clients[userID]; ok {
					select {
					case client.send <- message.Data:
						// Message sent successfully
					default:
						// Send buffer full, close and remove client
						close(client.send)
						delete(h.clients, userID)
						log.Printf("WebSocket: Client buffer full, disconnecting user %d", userID)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int64, data interface{}) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.clients[userID]
	if !ok {
		return nil // User not connected, skip
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	select {
	case client.send <- jsonData:
		return nil
	default:
		return nil // Channel full, skip
	}
}

// BroadcastToUsers broadcasts a message to multiple users
func (h *Hub) BroadcastToUsers(userIDs []int64, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	h.broadcast <- &BroadcastMessage{
		RecipientIDs: userIDs,
		Data:         jsonData,
	}

	return nil
}

// IsUserOnline checks if a user is online
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clients[userID]
	return ok
}

// GetOnlineUsers returns all online user IDs
func (h *Hub) GetOnlineUsers() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]int64, 0, len(h.clients))
	for userID := range h.clients {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

// GetConnectionCount returns the number of active connections
func (h *Hub) GetConnectionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.clients)
}

