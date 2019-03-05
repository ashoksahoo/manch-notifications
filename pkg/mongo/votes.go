package mongo

import (
	"notification-service/pkg/constants"
	"fmt"

	"github.com/globalsign/mgo/bson"
)

var (
	VOTES_MODEL = constants.ModelNames["VOTES"]
)

type VoteModelPost struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     PostModel
	ResourceId   bson.ObjectId `json:"resource" bson:"resource"`
	Created      Creator       `json:"created" bson:"created"`
	Value        int           `json:"vote" bson:"vote"`
}

type VoteModelComment struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     CommentModel
	Created      Creator `json:"created" bson:"created"`
	Value        int     `json:"vote" bson:"vote"`
}

func (c CommentModel) GetVote(Id string) VoteModelComment {
	s := session.Clone()
	defer s.Close()
	vote := VoteModelComment{}
	vote.Resource = c
	V := s.DB("manch").C(VOTES_MODEL)
	V.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&vote)
	return vote
}

func (p PostModel) GetVote(Id string) VoteModelPost {
	s := session.Clone()
	defer s.Close()
	vote := VoteModelPost{}
	vote.Resource = p
	V := s.DB("manch").C(VOTES_MODEL)
	V.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&vote)
	fmt.Printf("\nVote from DB %+v", vote)
	return vote
}

func GetAllVoteByQuery(query bson.M) []VoteModelPost {
	fmt.Println("query for get all vote", query)
	s := session.Clone()
	defer s.Close()
	V := s.DB("manch").C(VOTES_MODEL)
	votes := []VoteModelPost{}
	V.Find(query).All(&votes)
	return votes
}
