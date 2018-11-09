package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"log"
	"os"
)

var url = os.Getenv("MONGO_DB")
var session *mgo.Session

func init() {
	var err error
	if url == "" {
		url = "mongodb://localhost:27017/"
	}
	if session, err = mgo.Dial(url); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Initialized Mongo Connected")
	}
}
