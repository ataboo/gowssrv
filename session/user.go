package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const USERS_FILE = ".users.json"

type GameObject struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	W float32 `json:"w"`
	H float32 `json:"h"`
}

type User struct {
	ID       string
	Username string
	Hashed   []byte
	GameObj  GameObject
}

//TODO: this will use another storage method.
func GetByUsername(userName string) *User {
	users, err := loadUsers()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to load %s.", USERS_FILE))
		return nil
	}

	for _, user := range users {
		if user.Username == userName {
			return &user
		}
	}

	return nil
}

//TODO: this will use another storage method.
func FindUser(userId string) *User {
	users, err := loadUsers()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to load users"))
		return nil
	}

	user, ok := users[userId]

	if ok {
		return &user
	}

	return nil
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) *User {
	currentSession := Sessions.SessionStart(w, r)

	return FindUser(currentSession.Get("user_id").(string))
}

func AddUser(userName string, password string) *User {
	id := UniqueID()

	user := User{
		ID:       id,
		Username: userName,
	}
	user.SetPassword(password)
	user.Save()

	return &user
}

func (user *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(user.Hashed, []byte(password)) == nil
}

func (user *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return SessionError{"Failed to set password."}
	}
	user.Hashed = hashed

	return nil
}

// TODO: this will use a different storage method
func (user User) Save() error {
	users, _ := loadUsers()
	users[user.ID] = user

	raw, jErr := json.Marshal(users)
	if jErr != nil {
		fmt.Println("err")
		return jErr
	}

	return ioutil.WriteFile(USERS_FILE, raw, 0755)
}

// TODO: this will use a different storage method
func loadUsers() (map[string]User, error) {
	var users map[string]User
	raw, err := ioutil.ReadFile(USERS_FILE)
	if err != nil {
		return make(map[string]User, 0), err
	}

	err = json.Unmarshal(raw, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}
