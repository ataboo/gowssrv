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
	"github.com/dgrijalva/jwt-go"
	"strings"
)

func registerHandlers(r *mux.Router) {


	r.HandleFunc("/", index).Methods("GET")
	r.Handle("/restricted", jwtRequired(restricted)).Methods("GET")
	r.HandleFunc("/user", storeUser).Methods("POST")
	r.HandleFunc("/auth", authorize).Methods("POST")
	r.HandleFunc("/logout", jwtRequired(logout)).Methods("POST")
	//r.HandleFunc("/ws", upgradeWs).Methods("GET")
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

	raw, err := createToken(user)
	if err != nil {
		log.Println("Error creating token: " + err.Error())
		jsonResponse(w, 500, "Failed to authorize, please try again later.")
		return
	}

	signature := strings.Split(raw, ".")[2]

	//TODO: boot user from game if existing session.

	user.SessionId = signature
	storage.Users.Save(user)

	jsonResponse(w, 200, map[string]string{
		"message": "started new session",
		"token": raw,
	})
}

func logout(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("user_token").(*jwt.Token)

	fmt.Printf("\nRaw: %v", token.Raw)
	fmt.Printf("\nClaims: %v", token.Claims)
	fmt.Printf("\nSignature: %v", token.Signature)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r, true) {
		return
	}

	mustParseTemplates(w, newViewData(w, r), "base.tmpl", "game.tmpl")
}

//func handleWs(w http.ResponseWriter, r *http.Request) {
//	if !checkSession(w, r, true) {
//		return
//	}
//
//	user, err := session.GetCurrentUser(w, r)
//
//	if (err != nil) {
//
//	}
//
//	log.Println("Upgrading request")
//
//	if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
//		log.Println("Failed to upgrade ws:")
//		log.Println(err)
//	} else {
//		gObjString, _ := json.Marshal(user.GameObj)
//		update := append([]byte("player_update|"), gObjString...)
//
//		conn.WriteMessage(websocket.TextMessage, update)
//		readPump(conn, user)
//	}
//}

func checkSession(w http.ResponseWriter, r *http.Request, redirect bool) bool {
	sessionStore := session.Sessions.SessionStart(w, r)
	authenticated := sessionStore.Get("user_id") != nil

	if !authenticated && redirect {
		sessionStore.Set("flash_message", "Session expired. Please log in.")
		http.Redirect(w, r, "/", 302)
	}

	return authenticated
}

func mustParseTemplates(w http.ResponseWriter, data interface{}, files ...string) {
	temp, err := parseTemplates(files...)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := temp.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func parseTemplates(files ...string) (*template.Template, error) {
	for i, name := range files {
		files[i] = tmplPath(name)
	}

	return template.ParseFiles(files...)
}

func tmplPath(filename string) string {
	return fmt.Sprintf("resources/templates/%s", filename)
}

type ViewData struct {
	Flash interface{}
	Auth  bool
}

func newViewData(w http.ResponseWriter, r *http.Request) ViewData {
	currentSession := session.Sessions.SessionStart(w, r)
	auth := currentSession.Get("user_id") != nil
	flash := currentSession.Get("flash_message")

	if flash != nil {
		currentSession.Delete("flash_message")
	}

	fmt.Println(fmt.Sprintf("Found flash: %v", flash))

	return ViewData{
		Flash: flash,
		Auth:  auth,
	}
}
