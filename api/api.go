package api

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()
	registerHandlers(r)

	log.Println("Attempting to Listen on localhost:3000...")
	err := http.ListenAndServeTLS(":3000", "server.crt", "server.key", r)
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
