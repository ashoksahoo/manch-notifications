package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	COMMUNITIES_MODEL       = constants.ModelNames["COMMUNITIES"]
	COMMUNITIES_STATS_MODEL = constants.ModelNames["COMMUNITIES_STATS_MODEL"]
)

type CommunityParent struct {
	CommunityId bson.ObjectId `json:"community_id" bson:"community_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Type        string        `json:"type" bson:"type"`
	Scope       string        `json:"scope" bson:"scope"`
}

type CommunityModel struct {
	Id           bson.ObjectId     `json:"_id" bson:"_id"`
	Name         string            `json:"name" bson:"name"`
	Description  string            `json:"description" bson:"description"`
	Icon         string            `json:"icon" bson:"icon"`
	Type         string            `json:"type" bson:"type"`
	Language     string            `json:"language" bson:"language"`
	Category     string            `json:"category" bson:"category"`
	Scope        string            `json:"scope" bson:"scope"`
	Order        string            `json:"order" bson:"order"`
	Parents      []CommunityParent `json:"parents" bson:"parents"`
	DirectParent CommunityParent   `json:"direct_parent" bson:"direct_parent"`
	status       string            `json:"status" bson:"status"`
	visibility   string            `json:"visibility" bson:"visibility"`
}

type CommunityStatsModel struct {
	Id           bson.ObjectId     `json:"_id,omitempty" bson:"_id,omitempty"`
	CommunityId  bson.ObjectId     `json:"community_id" bson:"community_id"`
	Type         string            `json:"type" bson:"type"`
	Language     string            `json:"language" bson:"language"`
	InterestId   bson.ObjectId     `json:"interest_id" bson:"interest_id"`
	DirectParent CommunityParent   `json:"direct_parent" bson:"direct_parent"`
	Parents      []CommunityParent `json:"parents" bson:"parents"`
	Action       string            `json:"action" bson:"action"`
	EntityId     bson.ObjectId     `json:"entity_id" bson:"entity_id"`
	EntityType   string            `json:"entity_type" bson:"entity_type"`
	CreatedAt    time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt" bson:"updatedAt"`
	ProfileId    bson.ObjectId     `json:"profile_id" bson:"profile_id"`
}

func GetCommunityById(Id string) CommunityModel {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMUNITIES_MODEL)
	community := CommunityModel{}
	C.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&community)
	return community
}

func GetCommunityByQuery(q bson.M) []CommunityModel {
	s := session.Clone()
	defer s.Close()
	posts := s.DB("manch").C(COMMUNITIES_MODEL)
	var communities []CommunityModel
	posts.Find(q).One(&communities)
	return communities
}

func CreateCommunityStats(communityStats CommunityStatsModel) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMUNITIES_STATS_MODEL)
	community := GetCommunityById(communityStats.CommunityId.Hex())
	communityStats.InterestId = community.DirectParent.CommunityId
	now := time.Now()
	communityStats.Type = community.Type
	communityStats.DirectParent = community.DirectParent
	communityStats.Parents = community.Parents
	communityStats.CreatedAt = now
	communityStats.UpdatedAt = now
	communityStats.Language = community.Language

	fmt.Printf("community stats %+v\n\n", communityStats)
	err := C.Insert(communityStats)
	if err != nil {
		fmt.Println("error while creating community stats", err)
	} else {
		fmt.Println("created community stats successfully", communityStats.Id)
	}
}
