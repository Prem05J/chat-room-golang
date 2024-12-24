package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (s *Server) createUser() (*User, *User) {
	user := &User{
		ID:          uuid.New().String(),
		Name:        "Prem Kumar",
		MessageChan: make(chan []byte, 100),
	}

	user1 := &User{
		ID:          uuid.New().String(),
		Name:        "Dhoni",
		MessageChan: make(chan []byte, 100),
	}
	s.mu.Lock()
	s.Users[user.ID] = user
	s.Users[user1.ID] = user1
	s.mu.Unlock()
	return user, user1
}

func (s *Server) createChatRoom(name string) (*ChatRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.ChatRooms[name]; exists {
		return nil, fmt.Errorf("chat room %s already exists", name)
	}
	room := &ChatRoom{
		Name:    name,
		Members: make(map[string]*User),
	}
	s.ChatRooms[name] = room
	return room, nil
}

func (s *Server) joinChatRoom(roomName, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	room, ok := s.ChatRooms[roomName]
	if !ok {
		return fmt.Errorf("room %s not found", roomName)
	}
	user, ok := s.Users[userID]
	if !ok {
		return fmt.Errorf("user %s not found", userID)
	}
	room.Members[userID] = user
	return nil
}

func (s *Server) broadcastMessage(roomName string, msg Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	room, ok := s.ChatRooms[roomName]
	if !ok {
		return fmt.Errorf("room %s not found", roomName)
	}
	messageBytes, _ := json.Marshal(msg)

	for _, member := range room.Members {
		member.MessageChan <- messageBytes
	}
	return nil
}

func (s *Server) privateMessage(sender *User, receiverID string, msg Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	receiver, ok := s.Users[receiverID]

	if !ok {
		return fmt.Errorf("receiver %s not found", receiverID)
	}

	messageBytes, _ := json.Marshal(msg)
	receiver.MessageChan <- messageBytes
	return nil
}
