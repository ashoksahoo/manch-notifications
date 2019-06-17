package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notification-service/pkg/constants"
	"notification-service/pkg/mongo"
	"notification-service/pkg/utils"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/globalsign/mgo/bson"
	// "github.com/elastic/go-elasticsearch/v7"
	// "github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	POST_INDEX = constants.IndexNames["POSTS"]
)

type Moderation struct {
	Done        bool `json:"done" bson:"done"`
	ModeratedBy struct {
		ProfileId string `json:"profile_id"`
		Name      string `json:"name"`
		Date      string `json:"date"`
		Avatar    string `json:"avatar"`
		Type      string `json:"type"`
	} `json:"moderated_by"`
}

type CommunityMeta struct {
	CommunityId   string `json:"community_id" bson:"community_id"`
	Name          string `json:"name" bson:"name"`
	Icon          string `json:"icon" bson:"icon"`
	Language      string `json:"language" bson:"language"`
	GeoRegion     string `json:"geo_region"`
	GeoRegionName string `json:"geo_region_name" bson:"geo_region_name"`
	Type          string `json:"type" bson:"type"`
	Scope         string `json:"scope" bson:"scope"`
}

type PostIndex struct {
	ID                string                `json:"id,omitempty"`
	ShortId           string                `json:"short_id,omitempty"`
	Title             string                `json:"title,omitempty"`
	SearchText        string                `json:"search_text,omitempty"`
	MediaUrls         []mongo.Media         `json:"media_urls,omitempty"`
	RepostedPostId    string                `json:"reposted_post_id,omitempty"`
	Anonymous         bool                  `json:"anonymous,omitempty"`
	NoOfReposts       int                   `json:"no_of_reposts,omitempty"`
	Details           string                `json:"details,omitempty"`
	CommunityIds      []bson.ObjectId       `json:"community_ids,omitempty"`
	Communities       []mongo.CommunityMeta `json:"communities,omitempty"`
	PostType          string                `json:"post_type,omitempty"`
	PostFormat        string                `json:"post_format,omitempty"`
	Language          string                `json:"language,omitempty"`
	Source            mongo.Source          `json:"source",omitempty`
	Upvotes           int                   `json:"up_votes,omitempty"`
	Downvotes         int                   `json:"down_votes,omitempty"`
	NoOfViews         int                   `json:"no_of_views,omitempty"`
	NoOfImpressions   int                   `json:"no_of_impressions,omitempty"`
	EngagementToday   int                   `json:"engagement_today,omitempty"`
	TrendingNotified  bool                  `json:"trending_notified,omitempty"`
	IgnoreFromFeed    bool                  `json:"ignore_from_feed,omitempty"`
	IgnoreReason      string                `json:"ignore_reason,omitempty"`
	Tags              []string              `json:"tags,omitempty"`
	TagPositions      mongo.TagPositions    `json:"tag_positions,omitempty"`
	Moderation        mongo.Moderation      `json:"moderation,omitempty"`
	PostLevel         string                `json:"post_level,omitempty"`
	ContentSource     string                `json:"content_source,omitempty"`
	SourcedBy         string                `json:"sourced_by,omitempty"`
	Popularity        int                   `json:"popularity,omitempty"`
	FeedBaseTimeStamp string                `json:"feed_base_ts,omitempty"`
	CommentModerated  bool                  `json:"comment_moderated,omitempty"`
	StoryTime         int                   `json:"story_time,omitempty"`
	NoOfComments      int                   `json:"no_of_comments,omitempty"`
	CreatedAt         string                `json:"createdAt,omitempty"`
	UpdatedAt         string                `json:"updatedAt,omitempty"`
}

func getSearchText(post mongo.PostModel) string {
	title := post.Title
	sourceDescription := post.Source.Description
	var mediaDescriptions strings.Builder

	for _, media := range post.MediaUrls {
		mediaDescriptions.WriteString(media.Description + " ")
	}

	detail := post.Details
	return title + " " + detail + " " + sourceDescription + " " + mediaDescriptions.String() + " "
}

func CreatePostIndex(post mongo.PostModel) {
	// upsert data
	searchText := getSearchText(post)
	var postIndex = PostIndex{
		ID:                post.Id.Hex(),
		ShortId:           post.ShortId,
		Title:             post.Title,
		SearchText:        searchText,
		MediaUrls:         post.MediaUrls,
		RepostedPostId:    post.RepostedPostId.Hex(),
		Anonymous:         post.Anonymous,
		NoOfReposts:       post.RepostCount,
		Details:           post.Details,
		CommunityIds:      post.CommunityIds,
		Communities:       post.Communities,
		PostType:          post.PostType,
		PostFormat:        post.PostFormat,
		Language:          post.Language,
		Source:            post.Source,
		Upvotes:           post.UpVotes,
		Downvotes:         post.DownVotes,
		NoOfViews:         post.Views,
		NoOfImpressions:   post.Impressions,
		EngagementToday:   post.EngagementToday,
		TrendingNotified:  post.TrendingNotified,
		IgnoreFromFeed:    post.IgnoreFromFeed,
		IgnoreReason:      post.IgnoreReason,
		Tags:              post.Tags,
		TagPositions:      post.TagsPosition,
		Moderation:        post.Moderation,
		PostLevel:         post.PostLevel,
		ContentSource:     post.ContentSource,
		SourcedBy:         post.SourcedBy,
		Popularity:        post.Popularity,
		FeedBaseTimeStamp: utils.ISOFormat(post.FeedBaseTimeStamp),
		CommentModerated:  post.CommentModerated,
		StoryTime:         post.StoryTime,
		NoOfComments:      post.CommentCount,
		CreatedAt:         utils.ISOFormat(post.CreatedAt),
		UpdatedAt:         utils.ISOFormat(post.UpdatedAt),
	}
	var insertData StringInterface
	encodedPostIndex, _ := json.Marshal(postIndex)
	json.Unmarshal(encodedPostIndex, &insertData)

	fmt.Println("insert data", insertData)
	// Build the request body.
	body := esutil.NewJSONReader(insertData)
	// create update request
	req := esapi.IndexRequest{
		Index:      POST_INDEX,
		DocumentID: post.Id.Hex(),
		Body:       body,
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("response is", res)
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), post.Id.Hex())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

func UpdatePostIndex(postIndex PostIndex) {
	currentISOTime := utils.ISOFormat(time.Now())
	postIndex.UpdatedAt = currentISOTime
	var updateData StringInterface
	encodedUserIndex, _ := json.Marshal(postIndex)
	json.Unmarshal(encodedUserIndex, &updateData)

	fmt.Println("update data", updateData)
	// Build the request body.
	body := esutil.NewJSONReader(StringInterface{
		"doc": updateData,
	})
	// create update request
	req := esapi.UpdateRequest{
		Index:      POST_INDEX,
		DocumentID: postIndex.ID,
		Body:       body,
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("response is", res)
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), postIndex.ID)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

func SearchPost(query bson.M) (error, interface{}) {
	var r StringInterface
	q := query["q"]
	fmt.Println("Query is ", query)
	limit, limitOk := query["limit"]
	delete(query, "limit")
	if !limitOk {
		limit = 20
	}
	skip, skipOk := query["skip"]
	delete(query, "skip")
	if !skipOk {
		skip = 0
	}

	var size, from int
	size = limit.(int)
	from = skip.(int)

	body := esutil.NewJSONReader(StringInterface{
		"query": StringInterface{
			"match": StringInterface{
				"search_text": StringInterface{
					"query":     q,
					"fuzziness": "2",
				},
			},
		},
	})

	req := esapi.SearchRequest{
		Index: []string{POST_INDEX},
		Body:  body,
		Size:  &size,
		From:  &from,
	}

	res, err := req.Do(context.Background(), es)

	log.Printf("Response got\n%+v", res)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err, r
	}

	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error on getting data"), r
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return errors.New("Error parsing the response body"), r
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var response []map[string]interface{}
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		response = append(response, source)
	}
	return nil, response
}
