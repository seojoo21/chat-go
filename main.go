package main

import (
	"log"
	"net/http"

	"github.com/seojoo21/go-chat/server"
)

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", server.HandelConnections)

	go server.HandleMessages()

	log.Println("Listening on http://localhost:8080...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
