package main

import (
	"fmt"
	"log"
	"net/http"
)

var server = Server{
	Users:     make(map[string]*User),
	ChatRooms: make(map[string]*ChatRoom),
}

const (
	port = ":8080"
)

func main() {
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/rooms", getRooms)
	http.HandleFunc("/create-room", createRoom)
	http.HandleFunc("/join-room", joinRoom)
	http.HandleFunc("/send-message", sendMessage)
	http.HandleFunc("/stream", stream)
	http.HandleFunc("/room-members", getRoomMembers)

	fmt.Printf("Server listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
