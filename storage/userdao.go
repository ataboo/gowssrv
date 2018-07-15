package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/ataboo/gowssrv/models"
)

var Users UserDao

func init() {
	Users = UserDao{"users"}
}

type UserDao struct {
	colName string
}

func (m UserDao) Collection() *mgo.Collection {
	return Mongo.C(m.colName)
}

func (m UserDao) Find(id string) (models.User, error) {
	user := models.User{}
	err := m.Collection().FindId(id).One(user)

	return user, err
}

func (m UserDao) ByUsername(username string) (models.User, error) {
	user := models.User{}
	err := m.Collection().Find(bson.M{"username": username}).One(user)

	return user, err
}

func (m UserDao) Store(username string, password string) (models.User, error) {
	user := models.User{}
	user.ID = bson.NewObjectId()
	user.Username = username
	user.SetPassword(password)
	err := m.Save(user)

	return user, err
}

func (m UserDao) Save(user models.User) (error) {
	_, err := m.Collection().UpsertId(user.ID, &user)

	return err
}
