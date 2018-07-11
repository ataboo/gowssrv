package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ataboo/gowssrv/session"
)

func registerHandlers() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/create_user", handleCreateUser)
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/game", handleGame)
	http.HandleFunc("/ws", handleWs)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if checkSession(w, r, false) {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}

	mustParseTemplates(w, newViewData(w, r), "base.tmpl", "index.tmpl")
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	if checkSession(w, r, false) {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}

	mustParseTemplates(w, newViewData(w, r), "base.tmpl", "signup.tmpl")
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if checkSession(w, r, false) {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err.Error())
	}

	userName := r.Form.Get("username")
	password := r.Form.Get("password")

	var valErrs = []string{}

	if len(password) < 8 {
		valErrs = append(valErrs, "Password must be at least 8 characters")
	}

	if len(userName) < 5 {
		valErrs = append(valErrs, "Username must be at least 5 characters")
	}

	if session.GetByUsername(userName) != nil {
		valErrs = append(valErrs, "A user already exists with that name")
	}

	if len(valErrs) > 0 {
		session.Sessions.SessionStart(w, r).Set("flash_message", strings.Join(valErrs, ", "))
		data := newViewData(w, r)
		mustParseTemplates(w, data, "base.tmpl", "signup.tmpl")
	} else {
		user := session.AddUser(userName, password)

		newSession := session.Sessions.SessionStart(w, r)
		newSession.Set("user_id", user.ID)
		http.Redirect(w, r, "/game", 302)
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err.Error())
	}

	userName := r.Form.Get("username")
	password := r.Form.Get("password")

	user := session.GetByUsername(userName)

	newSession := session.Sessions.SessionStart(w, r)
	if user != nil && user.CheckPassword(password) {
		newSession.Set("user_id", user.ID)
		http.Redirect(w, r, "/game", 302)
		return
	}

	newSession.Set("flash_message", "Invalid username or password.")
	http.Redirect(w, r, "/", 302)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session.Sessions.SessionDestroy(w, r)
	http.Redirect(w, r, "/", 302)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r, true) {
		return
	}

	mustParseTemplates(w, newViewData(w, r), "base.tmpl", "game.tmpl")
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r, true) {
		return
	}

	user := session.GetCurrentUser(w, r)

	log.Println("Upgrading request")

	if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
		log.Println("Failed to upgrade ws:")
		log.Println(err)
	} else {
		readPump(conn, user)
	}
}

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
