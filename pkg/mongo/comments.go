package mongo

import (
	"github.com/globalsign/mgo/bson"
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

func GetFullCommentById(Id string) (CommentModel, int) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C("comments")
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
	C := s.DB("manch").C("comments")
	count, _ := C.Find(bson.M{"post_id": postId}).Count()
	return count
}

func GetCommentatorCount(postId bson.ObjectId, opId bson.ObjectId) int {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C("comments")
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
	C := s.DB("manch").C("comments")
	var result []bson.ObjectId
	C.Find(bson.M{"comment_id": commentId, "created.profile_id": bson.M{"$ne": commentCreator}}).Distinct("created.profile_id", &result)
	return len(result)
}

func GetCommentById(Id string) CommentModel {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C("comments")
	c := CommentModel{}
	C.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&c)
	return c
}
