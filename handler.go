package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type UsersResponse struct {
	User  interface{}
	User1 interface{}
}

func connect(w http.ResponseWriter, r *http.Request) {
	user, user1 := server.createUser()
	response := UsersResponse{
		User:  user,
		User1: user1,
	}

	json.NewEncoder(w).Encode(response)
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	server.mu.RLock()
	defer server.mu.RUnlock()

	rooms := make([]*ChatRoom, 0, len(server.ChatRooms))

	for _, room := range server.ChatRooms {
		rooms = append(rooms, room)
	}
	json.NewEncoder(w).Encode(rooms)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	room, err := server.createChatRoom(req.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(room)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("room")
	userID := r.URL.Query().Get("user")
	if err := server.joinChatRoom(roomName, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	server.mu.RLock()
	sender, ok := server.Users[msg.SenderID]
	server.mu.RUnlock()

	if !ok {
		http.Error(w, "Sender not found", http.StatusBadRequest)
		return
	}

	if msg.MessageType == "public" {
		if err := server.broadcastMessage(r.URL.Query().Get("room"), msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if msg.MessageType == "private" {
		if err := server.privateMessage(sender, msg.ReceiverID, msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

// SSE Emplementations
func stream(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	server.mu.RLock()
	user, ok := server.Users[userID]

	server.mu.RUnlock()

	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case message := <-user.MessageChan:
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			log.Printf("User %s disconnected", userID)
			return
		case <-time.After(30 * time.Second):
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		}
	}
}

func getRoomMembers(w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("room")
	server.mu.RLock()
	defer server.mu.RUnlock()
	room, ok := server.ChatRooms[roomName]
	if !ok {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(room.Members)
}
