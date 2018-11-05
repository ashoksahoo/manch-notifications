package mongo

import "github.com/globalsign/mgo"

func Session() *mgo.Session{
	session, _ := mgo.Dial("mongodb://localhost:27017/")
	return session
}
