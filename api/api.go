package api

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"log"
	"github.com/ataboo/gowssrv/config"
	"github.com/ataboo/gowssrv/models"
	"github.com/gorilla/websocket"
)

var Logger *logging.Logger

type WsConnectible interface {
	AddUser(user models.User, conn websocket.Conn)
	RemoveUser(user models.User)
}

func Start() (stopChan chan int) {
	Logger = logging.MustGetLogger("gowssrv_api")
	logging.SetLevel(logging.INFO, "gowssrv_api")

	stopChan = make(chan int)
	func () {
		go startTLS()
		Logger.Debug("Started TLS server")

		for {
			select {
			case <-stopChan:
				Logger.Debug("Got stop chan.  Stopping servers.")
				return
			default:
				//
			}
		}
	}()

	return stopChan
}

func startTLS() {
	r := mux.NewRouter()
	registerHandlers(r)

	hostAddress := config.Config.Api.HostAddress

	Logger.Info("TLS Listening on "+hostAddress+"...")
	err := http.ListenAndServeTLS(hostAddress, "server.crt", "server.key", r)
	if err != nil {
		log.Fatal(err)
	}
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{})  {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func jsonErrorResponse(w http.ResponseWriter, code int, payload interface{}) {
	jsonResponse(w, code, map[string]interface{}{"error": payload})
}
