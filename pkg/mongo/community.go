package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type CommunityModel struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Name     string        `json:"name"`
	Phone    string        `json:"phone"`
	Profiles []Profile     `json:"profiles"`
}

func GetCommunityById(Id string) CommunityModel {
	s := session.Clone()
	defer s.Close()
	posts := s.DB("manch").C("communities")
	community := CommunityModel{}
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&community)
	return community
}

func GetCommunityByQuery(q bson.M) []CommunityModel {
	s := session.Clone()
	defer s.Close()
	posts := s.DB("manch").C("communities")
	var communities []CommunityModel
	posts.Find(q).One(&communities)
	return communities
}
