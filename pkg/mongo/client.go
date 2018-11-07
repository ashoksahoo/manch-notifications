package mongo

import (
	"github.com/globalsign/mgo"
)

func Session() *mgo.Session {

	//mgo.SetDebug(true)
	//var aLogger *log.Logger
	//aLogger = log.New(os.Stderr, "", log.LstdFlags)
	//mgo.SetLogger(aLogger)

	session, _ := mgo.Dial("mongodb://localhost:27017/")
	return session
}
