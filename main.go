package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: checkWsOrigin,
}

func main() {
	registerHandlers()

	log.Println("Listening on localhost:3000...")
	http.ListenAndServeTLS(":3000", "server.crt", "server.key", nil)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading request")

	if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
		log.Println("Failed top upgrade ws:")
		log.Println(err)
		return
	} else {
		readPump(conn)
	}
}

func readPump(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msgStr := string(bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1)))
		log.Println(fmt.Sprintf("Got Message: %v", msgStr))
	}
}

func checkWsOrigin(r *http.Request) bool {
	return true
}
