package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	POSTS_MODEL = constants.ModelNames["POSTS"]
)
type Type struct{
	Height float64 `json:"height" bson:"height"`
	Width float64	`json:"width" bson:"width"`
	AspectRatio float64 `json:"aspect_ratio" bson:"aspect_ratio"`
}
type Source struct{
	Title string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	MediaType string `json:"media_type" bson:"media_type"`
	Url string `json:"url" bson:"url"`
}
type Dimensions struct{
	Type Type `json:"type" bson:"type"`
	Required bool `json:"required" bson:"required"`
}	
type Media struct {
	Url       string `json:"url" bson:"url"`
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`
	Type 	  string  `json:"type" bson:"type"`
	UploadedOn string `json:"uploadedOn" bson:"uploadedOn"`
	Description string `json:"description" bson:"description"`
	FileType 	string `json:"file_type" bson:"file_type"`
	Dimensions Dimensions `json:"dimensions" bson:"dimensions"`
}

type BlockReason struct {
	DeleteReason     string `json:"delete_reason" bson:"delete_reason"`
	IgnoreFeedReason string `json:"ignore_from_feed_reason" bson:"ignore_from_feed_reason"`
}

type TagPositions struct {
	Tag      string `json:"tag" bson:"tag"`
	TagStart int    `json:"tag_start" bson:"tag_end"`
	TagEnd   int    `json:"tag_end" bson:"tag_end"`
}
type Moderation struct {
	Done bool `json:"done" bson:"done"`
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
	PostLevel      string          `json:"post_level" bson:"post_level"`
	Reason         BlockReason     `json:"reason" bson:"reason"`
	Language       string          `json:"language" bson:"language"`
	SourcedBy      string          `json:"sourced_by" bson:"sourced_by"`
	RepostedPostId bson.ObjectId   `json:"reposted_post_id" bson:"reposted_post_id"`
	RepostCount    int             `json:"no_of_reposts" bson:"no_of_reposts"`
	Tags           []string        `json:"tags" bson:"tags"`
	TagsPosition   TagPositions    `json:"tag_positions" bson:"tag_positions"`
	Anonymous	   bool			   `json:"anonymous" bson:"anonymous"`
	ContentTextForm string		   `json:"content_text_form" bson:"content_text_form"`
	Details        string 		   `json:"details" bson:"details"`
	ShortID 	   string 		   `json:"short_id" bson:"short_id"`
	Source         Source		   `json:"source" bson:"source"`
	Moderation 	   Moderation	   `json:"moderation" bson:"moderation"`
}

func GetPost(Id bson.ObjectId) PostModel {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	P := s.DB("manch").C(POSTS_MODEL)
	P.Find(bson.M{"_id": Id}).One(&post)
	//FIXME: Do we need to get all the comments from there or its available of the post itself.
	//post.CommentCount = GetCommentCount(Id)
	return post
}

func GetPostById(Id string) (error, PostModel) {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	P := s.DB("manch").C(POSTS_MODEL)
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

func GetPostByQuery(query bson.M) (error, PostModel) {
	s := session.Clone()
	defer s.Close()
	post := PostModel{}
	P := s.DB("manch").C("posts")
	err := P.Find(query).One(&post)
	fmt.Println("error", err)
	return err, post
}

func GetPostCountByQuery(query bson.M) int {
	s := session.Clone()
	defer s.Close()
	P := s.DB("manch").C(POSTS_MODEL)
	n, err := P.Find(query).Count()
	if err != nil {
		return 0
	}
	return n
}

func UpdateAllPostsByQuery(query, update bson.M) {
	s := session.Clone()
	defer s.Close()
	P := s.DB("manch").C(POSTS_MODEL)
	info, err := P.UpdateAll(query, update)
	if err != nil {
		fmt.Println("error on updating post", err)
	} else {
		fmt.Println("post update info", info)
	}
}

func UpdateOnePostsByQuery(query, update bson.M) {
	s := session.Clone()
	defer s.Close()
	P := s.DB("manch").C(POSTS_MODEL)
	err := P.Update(query, update)
	if err != nil {
		fmt.Println("error on updating post", err)
	} else {
		fmt.Println("updated post")
	}
}

func GetUniquePostCreatorOnManch(communityId bson.ObjectId, adminIds []bson.ObjectId, startAt time.Time) int {
	s := session.Clone()
	defer s.Close()
	P := s.DB("manch").C(POSTS_MODEL)
	var result []bson.ObjectId
	fmt.Println("communityId", communityId)
	fmt.Printf("admin ids \n%+v\n", adminIds)
	fmt.Println("start AT", startAt)
	P.Find(bson.M{"community_ids": communityId, "createdAt": bson.M{"$gte": startAt}, "created.profile_id": bson.M{"$nin": adminIds}}).Distinct("created.profile_id", &result)
	return len(result)
}
