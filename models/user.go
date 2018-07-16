package models

import (
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type User struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Username   string        `bson:"user_name" json:"user_name"`
	AuthHashed []byte        `bson:"auth_hashed" json:"auth_hashed"`
	SessionId  string        `bson:"session_id" json:"session_id"`
}

func (user *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(user.AuthHashed, []byte(password)) == nil
}

func (user *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("failed to set password")
	}
	user.AuthHashed = hashed

	return nil
}

type UserSubmit struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func (u *UserSubmit) ValidateCreate() (valErrs []string) {
	if len(u.Password) < 8 {
		valErrs = append(valErrs, "Password must be at least 8 characters")
	}

	if len(u.Username) < 5 {
		valErrs = append(valErrs, "Username must be at least 5 characters")
	}

	return valErrs
}

func (u *UserSubmit) ToUser() User {
	user := User{}

	user.ID = bson.NewObjectId()
	user.Username = u.Username
	user.SetPassword(u.Password)

	return user
}