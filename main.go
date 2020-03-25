package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urriel/url-shortener-go/controllers"
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2"
)

func main() {
	mgoSession := getMongoSession()
	mongoDB := getMongoDatabase(mgoSession)
	mux := web.New()

	controllers.InitControllers(mux, mongoDB)

	fmt.Println("Listening port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getMongoDatabase(session *mgo.Session) *mgo.Database {
	var DB *mgo.Database

	mongoDatabase, gotVariable := os.LookupEnv("MONGO_DATABASE")
	if gotVariable {
		fmt.Printf("Getting Mongo Database : '%s'\n", mongoDatabase)
		DB = session.DB(mongoDatabase)
	} else {
		fmt.Println("Getting Mongo Database : 'urldb'", mongoDatabase)
		DB = session.DB("urldb")
	}

	return DB
}

func getMongoSession() *mgo.Session {
	var session *mgo.Session
	var err error

	mongoURL, gotVariable := os.LookupEnv("MONGO_URL")
	if gotVariable {
		fmt.Println("Connecting to MONGO_URL")
		session, err = mgo.Dial(mongoURL)
	} else {
		fmt.Println("Connecting to the default mongo url")
		session, err = mgo.Dial("mongodb://localhost:27017/urldb")
	}

	fmt.Println("Connected")

	if err != nil {
		log.Fatal("Cannot connect to mongodb server : " + err.Error())
		panic(err)
	}

	return session
}
