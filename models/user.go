package models

import (
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type User struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	Username string `bson:"user_name" json:"user_name"`
	AuthHashed []byte `bson:"auth_hashed" json:"auth_hashed"`
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