package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
	ResurfacedDate     time.Time `json:"resurfaced_date"`
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

func AddTagToIndex(tags []string) {

	hashTagData := HashTag{
		ActualCreationTime: time.Now(),
		LastUpdatedTime:    time.Now(),
		ResurfacedDate:     time.Now(),
	}

	var wg sync.WaitGroup
	for i, tag := range tags {
		wg.Add(1)

		go func(tagName string) {
			defer wg.Done()
			hashTagData.ID = tagName
			hashTagData.Keyword = TypeInput{
				Input: []string{tagName},
			}
			hashTagData.Title = tagName

			// upsert data
			var upsertData StringInterface
			hashTagDataEncoded, _ := json.Marshal(hashTagData)
			json.Unmarshal(hashTagDataEncoded, &upsertData)

			// Build the request body.
			body := esutil.NewJSONReader(StringInterface{
				"script": StringInterface{
					"source": "ctx._source.no_of_posts += params.count;ctx._source.last_updated_time=params.last_updated",
					"lang":   "painless",
					"params": StringInterface{
						"count":        1,
						"last_updated": hashTagData.LastUpdatedTime,
					},
				},
				"upsert": upsertData,
			})
			fmt.Println("requesting", body)
			// create update request
			req := esapi.UpdateRequest{
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
			fmt.Println("response is", res)
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
		Index:          []string{index},
		Body:           body,
		SourceExcludes: []string{"*"},
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
	options := r["suggest"].(map[string]interface{})["hashtags"].([]interface{})[0].(map[string]interface{})["options"]
	var response []map[string]interface{}
	for _, v := range options.([]interface{}) {
		id := v.(map[string]interface{})["_id"].(string)
		response = append(response, map[string]interface{}{"id": id})
	}
	return nil, response
}

func getScore(baseTime time.Time, noOfPost int, additionScore int) int {
	return (noOfPost*10 + additionScore)
}

/**
* update hashtag weight and returns the weight
* it takes tagName and additionScore
 */
func UpdateTagWeight(tag string, additionScore int) (error, int) {
	err, doc := GetDocumentById(tag, "tags")
	if err != nil {
		return err, 0
	}
	source := doc["_source"].(map[string]interface{})
	noOfPost := source["no_of_posts"].(int)
	baseTime := source["resurfaced_date"].(time.Time)
	weight := getScore(baseTime, noOfPost, 0)

	// update score to suggested field
	// Build the request body.
	body := esutil.NewJSONReader(StringInterface{
		"script": StringInterface{
			"source": "ctx._source.keyword.weight = params.weight",
			"params": StringInterface{
				"weight": weight,
			},
		},
	})
	fmt.Println("requesting", body)
	// create update request
	req := esapi.UpdateRequest{
		Index:      "tags",
		DocumentID: tag,
		Body:       body,
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return errors.New("Error getting response"), 0
	}
	defer res.Body.Close()
	fmt.Println("response is", res)
	if res.IsError() {
		return errors.New("Error getting response"), 0
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return errors.New("Error parsing the response body"), 0
	}
	log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	return nil, weight
}
