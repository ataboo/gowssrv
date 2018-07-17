package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "github.com/ataboo/gowssrv/models"
	"sync"
	"fmt"
)

const userCollection = "users"

var Users userDao

func init() {
	Users = userDao{userCollection, sync.Mutex{}}
}

// userDao Data Access Object for users
type userDao struct {
	colName string
	lock sync.Mutex
}

func (m userDao) Collection() *mgo.Collection {
	return Mongo.C(m.colName)
}

func (m userDao) Find(id string) (User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	user := User{}

	if !bson.IsObjectIdHex(id) {
		return user, fmt.Errorf("invalid object id")
	}

	err := m.Collection().FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (m userDao) ByUsername(username string) (User, error) {
	user := User{}
	err := m.Collection().Find(bson.M{"user_name": username}).One(&user)

	return user, err
}

func (m userDao) BySession(sessionId string) (User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	user := User{}
	err := m.Collection().Find(bson.M{"session_id": sessionId}).One(&user)

	return user, err
}

func (m userDao) Store(user User) (error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, err := m.ByUsername(user.Username); err == nil {
		return fmt.Errorf("username already taken")
	}
	err := m.Collection().Insert(&user)

	return err
}

func (m userDao) Save(user User) (error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.Collection().UpdateId(user.ID, &user)

	return err
}
