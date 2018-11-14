package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type VoteModelPost struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     PostModel
	Created      Creator `json:"created" bson:"created"`
	Value        int     `json:"vote" bson:"vote"`
}

type VoteModelComment struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     CommentModel
	Created      Creator `json:"created" bson:"created"`
	Value        int     `json:"vote" bson:"vote"`
}

func (c CommentModel) GetVote(Id string) (VoteModelComment) {
	s := session.Clone()
	defer s.Close()
	vote := VoteModelComment{}
	vote.Resource = c
	V := s.DB("manch").C("votes")
	V.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&vote)
	return vote
}

func (p PostModel) GetVote(Id string) (VoteModelPost) {
	s := session.Clone()
	defer s.Close()
	vote := VoteModelPost{}
	vote.Resource = p
	V := s.DB("manch").C("votes")
	V.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&vote)
	fmt.Printf("\nVote from DB %+v", vote)
	return vote
}
