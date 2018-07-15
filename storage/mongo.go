package storage

import (
	"gopkg.in/mgo.v2"
	"log"
	. "github.com/ataboo/gowssrv/config"
)

var Mongo *mgo.Database

func init() {
	Mongo = connectMongoDb()
}

func connectMongoDb() *mgo.Database {
	session, err := mgo.Dial(Config.MongoServer)
	if err != nil {
		log.Println("Failed to connect to Mongo")
		log.Fatal(err)
	}

	return session.DB(Config.MongoDb)
}