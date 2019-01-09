package mongo

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

type Media struct {
	Url       string `json:"url" bson:"url"`
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`
}
type PostModel struct {
	Id             bson.ObjectId   `json:"_id" bson:"_id"`
	Title          string          `json:"title" bson:"title"`
	Created        Creator         `json:"created" bson:"created"`
	CommunityIds   []bson.ObjectId `json:"community_ids" bson:"community_ids"`
	CommentCount   int             `json:"no_of_comments" bson:"no_of_comments"`
	UpVotes        int             `json:"up_votes" bson:"up_votes"`
	DownVotes      int             `json:"down_votes" bson:"down_votes"`
	Views          int             `json:"no_of_views" bson:"no_of_views"`
	Impressions    int             `json:"no_of_impressions" bson:"no_of_impressions"`
	MediaUrls      []Media         `json:"media_urls" bson:"media_urls"`
	IgnoreFromFeed bool            `json:"ignore_from_feed" bson:"ignore_from_feed"`
	IgnoreReason   string          `json:"ignore_reason" bson:"ignore_reason"`
}

func GetPost(Id bson.ObjectId) PostModel {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	P := s.DB("manch").C("posts")
	P.Find(bson.M{"_id": Id}).One(&post)
	//FIXME: Do we need to get all the comments from there or its available of the post itself.
	//post.CommentCount = GetCommentCount(Id)
	return post
}

func GetPostById(Id string) (error, PostModel) {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	P := s.DB("manch").C("posts")
	err := P.Find(bson.M{"_id": bson.ObjectIdHex(Id), "deleted": false}).One(&post)
	fmt.Println("error", err)
	return err, post
}

func ExtractThumbNailFromPost(post PostModel) (thumb string) {
	if len(post.MediaUrls) > 0 {
		thumb = post.MediaUrls[0].Thumbnail
	}
	return
}
