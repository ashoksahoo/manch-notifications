package mongo

import (
	"notification-service/pkg/constants"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	COMMENT_SCHEDULEDS_MODEL = constants.ModelNames["COMMENT_SCHEDULEDS"]
	COMMENTS_MODEL = constants.ModelNames["COMMENTS"]
)

type CommentModel struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	Content      string        `json:"content"`
	PostId       bson.ObjectId `json:"post_id" bson:"post_id"`
	Post         PostModel
	Created      Creator         `json:"created" bson:"created"`
	UpVotes      int             `json:"up_votes" bson:"up_votes"`
	DownVotes    int             `json:"down_votes" bson:"down_vote"`
	CommentId    bson.ObjectId   `json:"comment_id" bson:"comment_id"`
	Parents      []bson.ObjectId `json:"parents" bson:"parents"`
	CommentCount int             `json:"no_of_comments" bson:"no_of_comments"`
}

type CommentScheduleModel struct {
	Id          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Content     string        `json:"content"`
	PostId      bson.ObjectId `json:"post_id" bson:"post_id"`
	Schedule    Schedule      `json:"schedule" bson:"schedule"`
	Created     Creator       `json:"created" bson:"created"`
	CommentType string        `json:"comment_type" bson:"comment_type"`
}

func CreateCommentSchedule(content string, postId bson.ObjectId, commentCreator Creator, scheduleTime time.Time) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENT_SCHEDULEDS_MODEL)
	schedule := Schedule{
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      commentCreator,
	}
	commentScheduleData := CommentScheduleModel{
		Content:     content,
		PostId:      postId,
		Schedule:    schedule,
		Created:     commentCreator,
		CommentType: "TEXT",
	}
	err := C.Insert(commentScheduleData)
	if err == nil {
		fmt.Printf("inserted comment schedule: %+v\n", commentScheduleData)
	} else {
		fmt.Println(err)
	}
}

func GetFullCommentById(Id string) (CommentModel, int) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	c := CommentModel{}
	C.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&c)
	var uniqCommentator int
	c.Post = GetPost(c.PostId)
	uniqCommentator = GetCommentatorCount(c.PostId, c.Post.Created.ProfileId)
	return c, uniqCommentator
}

func GetCommentCount(postId bson.ObjectId) int {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	count, _ := C.Find(bson.M{"post_id": postId}).Count()
	return count
}

func GetCommentatorCount(postId bson.ObjectId, opId bson.ObjectId) int {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	var result []bson.ObjectId
	C.Find(bson.M{"post_id": postId, "created.profile_id": bson.M{"$ne": opId}}).Distinct("created.profile_id", &result)
	//fmt.Printf("R %+v", result)
	//fmt.Printf("P %+v", postId)
	//fmt.Printf("OP %+v", opId)
	return len(result)
}

func GetReplierCount(commentId, commentCreator bson.ObjectId) int {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	var result []bson.ObjectId
	C.Find(bson.M{"comment_id": commentId, "created.profile_id": bson.M{"$ne": commentCreator}}).Distinct("created.profile_id", &result)
	return len(result)
}

func GetCommentById(Id string) CommentModel {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	c := CommentModel{}
	C.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&c)
	return c
}

func GetCommentsByPostId(postId, commentCreator bson.ObjectId) []bson.ObjectId {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	var result []bson.ObjectId
	C.Find(bson.M{
		"post_id":            postId,
		"created.profile_id": bson.M{"$ne": commentCreator},
		"parents":            bson.M{"$exists": true},
		"$where":             "this.parents.length<2",
	}).Distinct("created.profile_id", &result)
	return result
}

func GetRepliesByCommentId(postId, commentId, replyCreator bson.ObjectId) []bson.ObjectId {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENTS_MODEL)
	var result []bson.ObjectId
	C.Find(bson.M{
		"post_id":            postId,
		"comment_id":         commentId,
		"created.profile_id": bson.M{"$ne": replyCreator},
	}).Distinct("created.profile_id", &result)
	return result
}

func RemoveCommentScheduleByPostId(pId bson.ObjectId) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENT_SCHEDULEDS_MODEL)
	info, err := C.RemoveAll(bson.M{"post_id": pId})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("deleted comment schedule info", info)
	}
}
