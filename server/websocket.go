package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

type Message struct {
	Type     string `json:"type"` // 'chat' for regular messages, 'notice' for notices
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/*func HandelConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		//read message from browser
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message :: ", err)
			break
		}

		log.Printf("Received : %s \n", message)

		//Echo the message back
		if err := ws.WriteMessage(messageType, message); err != nil {
			log.Println("Error writing message :: ", err)
			break
		}
	}
}*/

func HandelConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	SendNotice("Welcome to the Chat Room!")

	// Register new client
	clients[ws] = true

	for {
		var msg Message
		// Read new message as JSON
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("readJson error : %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("writeJson err : %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func SendNotice(noticeText string) {
	notice := Message{
		Type:     "notice",
		Email:    "",
		Username: "Admin",
		Message:  noticeText,
	}
	broadcast <- notice
}
