package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		// read
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Panicln(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// WS endpoint
func serveWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)

	}

	// listen message
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple server")
	})

	http.HandleFunc("/ws", serveWS)
}

func main() {
	setupRoutes()

	err := http.ListenAndServe(":8071", nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
