package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notification-service/pkg/constants"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/globalsign/mgo/bson"
	// "github.com/elastic/go-elasticsearch/v7"
	// "github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	USERS_INDEX = constants.IndexNames["USERS"]
)

type UserIndex struct {
	ID                 string    `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	NameKeyword        TypeInput `json:"name_keyword,omitempty"`
	Avatar             string    `json:"avatar,omitempty"`
	AboutMe            string    `json:"about_me,omitempty"`
	Type               string    `json:"type,omitempty"`
	NoOfPosts          int       `json:"no_of_posts,omitempty"`
	NoOfLikes          int       `json:"no_of_likes,omitempty"`
	NoOfComments       int       `json:"no_of_comments,omitempty"`
	NoOfShares         int       `json:"no_of_shares,omitempty"`
	NoOfFollowers      int       `json:"no_of_followers,omitempty"`
	NoOfFollowing      int       `json:"no_of_following,omitempty"`
	NoOfManchFollowing int       `json:"no_of_manch_following,omitempty"`
	LastActiveHour     time.Time `json:"last_active_hour,omitempty"`
	TotalCoins         int       `json:"total_coins,omitempty"`
	TotalManchCreated  int       `json:"total_manch_created,omitempty"`
	BranchLink         string    `json:"branch_link,omitempty"`
	CreatedAt          string    `json:"createdAt,omitempty"`
	UpdatedAt          string    `json:"updatedAt,omitempty"`
}

func CreateUserIndex(userIndex UserIndex) {
	// upsert data
	var insertData StringInterface
	encodedUserIndex, _ := json.Marshal(userIndex)
	json.Unmarshal(encodedUserIndex, &insertData)

	fmt.Println("insert data", insertData)
	// Build the request body.
	body := esutil.NewJSONReader(insertData)
	// create update request
	req := esapi.IndexRequest{
		Index:      USERS_INDEX,
		DocumentID: userIndex.ID,
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
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), userIndex.ID)
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

func UpdateUserIndex(userIndex UserIndex) {
	var updateData StringInterface
	encodedUserIndex, _ := json.Marshal(userIndex)
	json.Unmarshal(encodedUserIndex, &updateData)

	fmt.Println("update data", updateData)
	// Build the request body.
	body := esutil.NewJSONReader(StringInterface{
		"doc": updateData,
	})
	// create update request
	req := esapi.UpdateRequest{
		Index:      USERS_INDEX,
		DocumentID: userIndex.ID,
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
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), userIndex.ID)
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

func SuggestUsers(query bson.M) (error, map[string]interface{}) {
	var r StringInterface
	fmt.Println("Query is ", query)
	limit, limitOk := query["limit"]
	delete(query, "limit")
	if !limitOk {
		limit = 20
	}
	keyword := query["name_keyword"]
	body := esutil.NewJSONReader(StringInterface{
		"suggest": StringInterface{
			"users": StringInterface{
				"prefix": keyword,
				"completion": StringInterface{
					"field": "name_keyword",
					"size":  limit,
					"fuzzy": StringInterface{
						"fuzziness": 1,
					},
				},
			},
		},
	})

	req := esapi.SearchRequest{
		Index:          []string{USERS_INDEX},
		Body:           body,
		SourceIncludes: []string{"tagname"},
	}
	res, err := req.Do(context.Background(), es)
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
	// log.Printf("response\n%+v", r)
	options := r["suggest"].(map[string]interface{})["users"].([]interface{})[0];
	return nil, options.(map[string]interface{})

}
