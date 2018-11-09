package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type PostModel struct {
	Id           bson.ObjectId   `json:"_id" bson:"_id"`
	Title        string          `json:"title"`
	Creator         `json:"created"`
	CommunityIds []bson.ObjectId `json:"community_ids" bson:"community_ids"`
	//Communities  []CommunityModel
}

func GetPostById(Id string) PostModel {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	posts := s.DB("manch").C("posts")
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&post)
	//post.Communities = GetCommunityByQuery(bson.M{"profile_id": bson.M{"$in": post.CommunityIds}, "deleted": false})
	return post
}
