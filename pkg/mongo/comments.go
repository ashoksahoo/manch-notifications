package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type CommentModel struct {
	Id      bson.ObjectId `json:"_id" bson:"_id"`
	Content string        `json:"content"`
	Post    bson.ObjectId `json:"post_id" bson:"post_id"`
	Creator
}

func GetCommentById(Id string) CommentModel {
	s := session.Clone()
	defer s.Close()
	comments := s.DB("manch").C("comments")
	newComment := CommentModel{}
	comments.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&newComment)
	return newComment
}
