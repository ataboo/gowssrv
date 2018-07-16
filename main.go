package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/ataboo/gowssrv/api"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkWsOrigin,
}

func main() {
	api.Start()
}

//func readPump(conn *websocket.Conn, user *session.User) {
//	defer func() {
//		conn.Close()
//		fmt.Printf("\nFinal User State: %v", user.GameObj)
//		user.Save()
//	}()
//
//	for {
//		_, message, err := conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("error: %v", err)
//			}
//
//			break
//		}
//		msgStr := string(bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1)))
//		log.Println(fmt.Sprintf("Got Message: %v", msgStr))
//
//		splitMsg := strings.Split(msgStr, "|")
//
//		if len(splitMsg) == 2 {
//			switch splitMsg[0] {
//			case "player_update":
//				gObj := session.GameObject{}
//				err := json.Unmarshal([]byte(splitMsg[1]), &gObj)
//
//				if err == nil {
//					user.GameObj = gObj
//				} else {
//					log.Println(fmt.Sprintf("Invalid game object: %s", splitMsg[1]))
//					fmt.Printf("ERR: %s", err)
//				}
//			}
//		}
//
//		// err = conn.WriteMessage(websocket.TextMessage, []byte("Message received"))
//
//		if err != nil {
//			log.Println(fmt.Sprintf("Error writing message: %v", err))
//		}
//	}
//}

func checkWsOrigin(r *http.Request) bool {
	return true
}
