package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ws.Close()
	for {
		msgType, message, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err.Error())
			break
		}
		log.Printf("recv: %s", message)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(5)
		time.Sleep(time.Duration(n) * time.Second)
		err = ws.WriteMessage(msgType, message)
		if err != nil {
			log.Fatal(err.Error())
			break
		}
	}
}
