package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notification-service/pkg/mongo"

	// "github.com/elastic/go-elasticsearch/v7"
	// "github.com/elastic/go-elasticsearch/v7/esapi"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/globalsign/mgo/bson"
)

type StringInterface map[string]interface{}

type TypeInput struct {
	Input  []string `json:"input"`
	Weight int      `json:"weight"`
}

type HashTag struct {
	ID                 string    `json:"id"`
	Keyword            TypeInput `json:"keyword"`
	Title              string    `json:"title"`
	Image              string    `json:"image"`
	NoOfPosts          int       `json:"no_of_posts"`
	NoOfLikes          int       `json:"no_of_likes"`
	NoOfComments       int       `json:"no_of_comments"`
	ActualCreationTime time.Time `json:"actual_creation_time"`
	LastUpdatedTime    time.Time `json:"last_updated_time"`
	Resurfaced         bool      `json:"resurfaced"`
}

func GetDocumentById(id, index string) (error, map[string]interface{}) {
	var r map[string]interface{}
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err, r
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error on getting data")
		return errors.New("Error on getting data"), r
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return errors.New("Error parsing the response body"), r
	}
	log.Println("response", r)
	return nil, r
}

func AddTagIndexFromPost(post mongo.PostModel) {

	hashTagData := HashTag{
		Title:              post.Title,
		ActualCreationTime: time.Now(),
		LastUpdatedTime:    time.Now(),
	}

	var wg sync.WaitGroup
	for i, tag := range post.Tags {
		wg.Add(1)

		go func(tagName string) {
			defer wg.Done()
			hashTagData.ID = tagName
			hashTagData.Keyword = TypeInput{
				Input: []string{tagName},
			}
			// Build the request body.
			body := esutil.NewJSONReader(hashTagData)
			fmt.Println("requesting", body)
			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "tags",
				DocumentID: hashTagData.ID,
				Body:       body,
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
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
		}(tag)
	}
	wg.Wait()
}

func SearchHashTags(query bson.M) (error, interface{}) {
	var x interface{}
	fmt.Println("Query is ", query)
	limit, limitOk := query["limit"]
	delete(query, "limit")
	// skip, skipOk := query["skip"]
	// delete(query, "skip")
	if !limitOk {
		limit = 20
	}
	// if !skipOk {
	// 	skip = 0
	// }
	keyword := query["keyword"]
	body := esutil.NewJSONReader(StringInterface{
		"suggest": StringInterface{
			"hashtags": StringInterface{
				"prefix": keyword,
				"completion": StringInterface{
					"field": "keyword",
					"size":  limit,
					"fuzzy": StringInterface{
						"fuzziness": 1,
					},
				},
			},
		},
	})
	index := "tags"
	var r StringInterface
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  body,
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err, x
	}

	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Error on getting data"), x
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return errors.New("Error parsing the response body"), x
	}
	// log.Printf("response\n%+v", r)
	response := r["suggest"].(map[string]interface{})["hashtags"].([]interface{})[0].(map[string]interface{})
	return nil, response["options"]
}
