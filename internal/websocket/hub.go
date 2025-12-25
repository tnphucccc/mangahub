// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package websocket

import (
	"context" // Ensure context is explicitly imported
	"encoding/json"
	"fmt"
	"log"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients per room (room -> clients)
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan models.WebSocketMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan models.WebSocketMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Context was cancelled, shut down the hub
			for _, clients := range h.rooms {
				for client := range clients {
					close(client.send)
				}
			}
			return
		case client := <-h.register:
			// Client joins their designated room
			if client.room == "" {
				client.room = "general"
			}
			if h.rooms[client.room] == nil {
				h.rooms[client.room] = make(map[*Client]bool)
			}
			h.rooms[client.room][client] = true

			// Send join notification to room
			joinMsg := models.NewSystemMessage(client.room, fmt.Sprintf("%s joined the chat", client.username))
			h.broadcastToRoom(client.room, joinMsg)

			log.Printf("[Hub] User %s joined room %s", client.username, client.room)

		case client := <-h.unregister:
			if clients, ok := h.rooms[client.room]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)

					// Clean up empty rooms
					if len(clients) == 0 {
						delete(h.rooms, client.room)
					}

					// Send leave notification to room
					leaveMsg := models.NewSystemMessage(client.room, fmt.Sprintf("%s left the chat", client.username))
					h.broadcastToRoom(client.room, leaveMsg)

					log.Printf("[Hub] User %s left room %s", client.username, client.room)
				}
			}

		case message := <-h.Broadcast:
			// Broadcast message to all clients in the specified room
			h.broadcastToRoom(message.Room, message)
		}
	}
}

// broadcastToRoom sends a message to all clients in a specific room
func (h *Hub) broadcastToRoom(room string, message models.WebSocketMessage) {
	clients, ok := h.rooms[room]
	if !ok {
		return
	}

	// Marshal message to JSON
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("[Hub] Error marshaling message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(clients, client)
		}
	}
}
