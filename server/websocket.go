package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandelConnections(w http.ResponseWriter, r *http.Request) {
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
}
