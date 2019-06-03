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

type Media struct {
	Url         string `json:"url" bson:"url"`
	Thumbnail   string `json:"thumbnail" bson:"thumbnail"`
	Description string `json:"description" bson:"description"`
	Filetype    string `json:"file_type" bson:"file_type"`
	Duration    int    `json:"duration" bson:"duration"`
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

type Source struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	MediaType   string `json:"media_type" bson:"media_type"`
	Url         string `json:"url" bson:"url"`
}

type Moderation struct {
	Done        bool  `json:"done" bson:"done"`
	ModeratedBy Actor `json:"moderated_by" bson:"moderated_by"`
}

type CommunityMeta struct {
	CommunityId   bson.ObjectId `json:"community_id" bson:"community_id"`
	Name          string        `json:"name" bson:"name"`
	Icon          string        `json:"icon" bson:"icon"`
	Language      string        `json:"language" bson:"language"`
	GeoRegion     string        `json:"geo_region"`
	GeoRegionName string        `json:"geo_region_name" bson:"geo_region_name"`
	Type          string        `json:"type" bson:"type"`
	Scope         string        `json:"scope" bson:"scope"`
}

type PostModel struct {
	Id                bson.ObjectId   `json:"_id" bson:"_id"`
	ShortId           string          `json:"short_id" bson:"short_id"`
	Title             string          `json:"title" bson:"title"`
	Details           string          `json:"details" bson:"details"`
	Created           Creator         `json:"created" bson:"created"`
	CommunityIds      []bson.ObjectId `json:"community_ids" bson:"community_ids"`
	CommentCount      int             `json:"no_of_comments" bson:"no_of_comments"`
	UpVotes           int             `json:"up_votes" bson:"up_votes"`
	DownVotes         int             `json:"down_votes" bson:"down_votes"`
	Views             int             `json:"no_of_views" bson:"no_of_views"`
	Impressions       int             `json:"no_of_impressions" bson:"no_of_impressions"`
	MediaUrls         []Media         `json:"media_urls" bson:"media_urls"`
	IgnoreFromFeed    bool            `json:"ignore_from_feed" bson:"ignore_from_feed"`
	IgnoreReason      string          `json:"ignore_reason" bson:"ignore_reason"`
	PostLevel         string          `json:"post_level" bson:"post_level"`
	Reason            BlockReason     `json:"reason" bson:"reason"`
	Language          string          `json:"language" bson:"language"`
	SourcedBy         string          `json:"sourced_by" bson:"sourced_by"`
	RepostedPostId    bson.ObjectId   `json:"reposted_post_id" bson:"reposted_post_id"`
	RepostCount       int             `json:"no_of_reposts" bson:"no_of_reposts"`
	Tags              []string        `json:"tags" bson:"tags"`
	TagsPosition      TagPositions    `json:"tag_positions" bson:"tag_positions"`
	Source            Source          `json:"source" bson:"source"`
	Moderation        Moderation      `json:"moderation" bson:"moderation"`
	Communities       []CommunityMeta `json:"communities" bson:"communities"`
	Anonymous         bool            `json:"anonymous" bson:"anonymous"`
	PostType          string          `json:"post_type" bson:"post_type"`
	PostFormat        string          `json:"post_format" bson:"post_format"`
	EngagementToday   int             `json:"engagement_today,omitempty" bson:"engagement_today"`
	TrendingNotified  bool            `json:"trending_notified,omitempty"`
	ContentSource     string          `json:"content_source,omitempty" bson:"content_source"`
	Popularity        int             `json:"popularity,omitempty" bson:"popularity"`
	FeedBaseTimeStamp time.Time       `json:"feed_base_ts,omitempty" bson:"feed_base_ts"`
	CommentModerated  bool            `json:"comment_moderated,omitempty" bson:"comment_moderated"`
	StoryTime         int             `json:"story_time,omitempty" bson:"story_time"`
	CreatedAt         time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt,omitempty" bson:"updatedAt"`
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
