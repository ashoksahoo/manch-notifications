package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type PostModel struct {
	Id           bson.ObjectId   `json:"_id" bson:"_id"`
	Title        string          `json:"title" bson:"title"`
	Created      Creator         `json:"created" bson:"created"`
	CommunityIds []bson.ObjectId `json:"community_ids" bson:"community_ids"`
	CommentCount int
}

func GetPostById(Id bson.ObjectId) (PostModel) {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	P := s.DB("manch").C("posts")
	P.Find(bson.M{"_id": Id}).One(&post)
	post.CommentCount = GetCommentCount(Id)
	return post
}
