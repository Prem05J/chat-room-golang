package main

import "sync"

type User struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	MessageChan chan []byte `json:"-"`
}

type ChatRoom struct {
	Name    string           `json:"name"`
	Members map[string]*User `json:"members"`
}

type Message struct {
	SenderID    string `json:"senderID"`
	SenderName  string `json:"senderName"`
	Content     string `json:"content"`
	MessageType string `json:"messageType"`
	ReceiverID  string `json:"receiverId,omitempty"`
}

type Server struct {
	Users     map[string]*User
	ChatRooms map[string]*ChatRoom
	mu        sync.RWMutex
}
