package mongo

import (
	"fmt"
	"log"
	"os"

	"github.com/globalsign/mgo"
)

var url = os.Getenv("MONGO_DB")
var session *mgo.Session

func init() {
	var err error
	if url == "" {
		url = "mongodb://mongo:27017/"
	}
	if session, err = mgo.Dial(url); err != nil {
		log.Fatal("Mongo Error:", err)
	} else {
		fmt.Println("Initialized Mongo Connected")
	}
	mgo.SetDebug(true)
}
