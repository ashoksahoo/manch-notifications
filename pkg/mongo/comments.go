package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type CommentModel struct {
	Id      bson.ObjectId `json:"_id" bson:"_id"`
	Content string        `json:"content"`
	PostId    bson.ObjectId `json:"post_id" bson:"post_id"`
	Post    PostModel
	Created Creator `json:"created" bson:"created"`
}

func GetCommentById(Id string) CommentModel {
	s := session.Clone()
	defer s.Close()
	comments := s.DB("manch").C("comments")
	c := CommentModel{}
	comments.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&c)
	c.Post = GetPostById(c.PostId)
	return c
}
