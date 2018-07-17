package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/ataboo/gowssrv/storage"
	"bitbucket.org/ataboo/servecapgo/session"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/ataboo/gowssrv/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func registerHandlers(r *mux.Router) {
	r.HandleFunc("/", index).Methods("GET")
	r.Handle("/restricted", tokenMiddleware(restricted)).Methods("GET")
	r.HandleFunc("/user", storeUser).Methods("POST")
	r.HandleFunc("/auth", authorize).Methods("POST")
	r.HandleFunc("/logout", tokenMiddleware(logout)).Methods("POST")
	r.HandleFunc("/ws", tokenMiddleware(upgradeWs)).Methods("GET")
}

func index(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, "Welcome to gowssrv! Keep your hands in plain site.")
}

func restricted(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, "VIP access activated.")
}

func storeUser(w http.ResponseWriter, r *http.Request) {
	userSubmit := models.UserSubmit{}

	err := json.NewDecoder(r.Body).Decode(&userSubmit)
	if err != nil {
		jsonErrorResponse(w, 422, "invalid format")
		return
	}

	if errs := userSubmit.ValidateCreate(); len(errs) > 0 {
		jsonErrorResponse(w, 422, errs)
		return
	}

	if _, err := storage.Users.ByUsername(userSubmit.Username); err == nil {
		jsonErrorResponse(w, 422, []string{"Username already taken"})
		return
	}

	user := userSubmit.ToUser()

	if err := storage.Users.Store(user); err != nil {
		log.Println(fmt.Sprintf("Failed to store user: \n%v", err))
		jsonErrorResponse(w, 500, "failed to create new user, please try again later")
		return
	}

	jsonResponse(w, 200, map[string]string{
		"message": "created user successfully",
	})
}

func authorize(w http.ResponseWriter, r *http.Request) {
	userSubmit := models.UserSubmit{}

	err := json.NewDecoder(r.Body).Decode(&userSubmit)
	if err != nil {
		jsonErrorResponse(w, 422, "invalid format")
		return
	}

	user, err := storage.Users.ByUsername(userSubmit.Username)
	if err != nil || !user.CheckPassword(userSubmit.Password) {
		jsonErrorResponse(w, 401, "Invalid password or user not found.")
		return
	}

	//TODO: rate limit.

	token := makeToken()

	//TODO: boot user from game if existing session.

	user.SessionId = token
	storage.Users.Save(user)

	jsonResponse(w, 200, map[string]interface{}{
		"message": "started new session",
		"token": token,
	})
}

func logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	// This is safe as long as we validate the token in the middleware.
	user.SessionId = ""
	storage.Users.Save(user)

	jsonResponse(w, 200, "Logged out successfully")
}

func upgradeWs(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("Upgrading websocket")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		jsonErrorResponse(w, 400, "Failed to upgrade connection")
		return
	}


}
