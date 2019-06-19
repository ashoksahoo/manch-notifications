package mongo

import (
	"fmt"
	"notification-service/pkg/constants"

	"github.com/globalsign/mgo/bson"
)

var (
	VOTES_MODEL           = constants.ModelNames["VOTES"]
	VOTES_SCHEDULED_MODEL = constants.ModelNames["VOTE_SCHEDULEDS"]
)

type VoteModelPost struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     PostModel
	ResourceId   bson.ObjectId `json:"resource_id" bson:"resource_id"`
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
	fmt.Printf("\nVote from DB %+v\n\n", vote)
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

func CountVoteByQuery(query bson.M) int {
	s := session.Clone()
	defer s.Close()
	V := s.DB("manch").C(VOTES_MODEL)
	count, _ := V.Find(query).Count()
	return count
}

func CountScheduledVotesByQuery(query bson.M) int {
	s := session.Clone()
	defer s.Close()
	V := s.DB("manch").C(VOTES_SCHEDULED_MODEL)
	count, _ := V.Find(query).Count()
	return count
}

func GetAllVotedUserIncludingScheduled(query bson.M) []string {
	s := session.Clone()
	defer s.Close()
	VoteCollection := s.DB("manch").C(VOTES_MODEL)
	VoteScheduledCollection := s.DB("manch").C(VOTES_SCHEDULED_MODEL)

	votes := []VoteModelPost{}
	scheduledVotes := []VoteScheduleModelPost{}

	VoteCollection.Find(query).Select(bson.M{"created": 1}).All(&votes)
	VoteScheduledCollection.Find(query).Select(bson.M{"created": 1}).All(&scheduledVotes)

	votedUsersProfileIds := []string{}

	for _, vote := range votes {
		votedUsersProfileIds = append(votedUsersProfileIds, vote.Created.ProfileId.Hex())
	}

	for _, scheduledVote := range scheduledVotes {
		votedUsersProfileIds = append(votedUsersProfileIds, scheduledVote.Created.ProfileId.Hex())
	}

	return votedUsersProfileIds
}
