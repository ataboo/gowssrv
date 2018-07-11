package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"bitbucket.org/ataboo/lasecapgo/atautils"

	"golang.org/x/crypto/bcrypt"
)

const USERS_FILE = ".users.json"

type User struct {
	ID       string
	Username string
	Hashed   []byte
}

func GetByUsername(userName string) *User {
	users, err := loadUserFile()
	if err != nil {
		atautils.Logger().Debug("Failed to load %s.", USERS_FILE)
		return nil
	}

	for _, user := range users {
		if user.Username == userName {
			return &user
		}
	}

	return nil
}

func AddUser(userName string, password string) *User {
	id := atautils.UniqueID()

	user := User{
		ID:       id,
		Username: userName,
	}
	user.SetPassword(password)

	return &user
}

func loadUserFile() (map[string]User, error) {
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

func (user *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(user.Hashed, []byte(password)) == nil
}

func (user *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return SessionError{"Failed to set password."}
	}
	user.Hashed = hashed
	user.Save()

	return nil
}

func (user User) Save() error {
	users, _ := loadUserFile()
	users[user.ID] = user

	raw, jErr := json.Marshal(users)
	if jErr != nil {
		fmt.Println("err")
		return jErr
	}

	return ioutil.WriteFile(USERS_FILE, raw, 0755)
}
