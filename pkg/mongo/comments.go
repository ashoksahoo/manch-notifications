package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type CommentModel struct {
	Id      bson.ObjectId `json:"_id" bson:"_id"`
	Content string        `json:"content"`
	PostId  bson.ObjectId `json:"post_id" bson:"post_id"`
	Post    PostModel
	Created Creator `json:"created" bson:"created"`
}

func GetCommentById(Id string) (CommentModel, int) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C("comments")
	c := CommentModel{}
	C.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&c)
	var uniqCommentator int
	c.Post, uniqCommentator = GetPostById(c.PostId)
	return c, uniqCommentator
}

func GetCommentCount(postId bson.ObjectId) (int, int) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C("comments")
	count, _ := C.Find(bson.M{"post_id": postId}).Count()
	uniqueCommentator := GetCommentatorCount(postId)
	return count, uniqueCommentator
}

func GetCommentatorCount(postId bson.ObjectId) int {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C("comments")
	var result []bson.ObjectId
	C.Find(bson.M{"post_id": postId}).Distinct("post_id", &result)
	return len(result)
}
