package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notification-service/pkg/mongo"
	"notification-service/pkg/utils"
	"strings"

	// "github.com/elastic/go-elasticsearch/v7"
	// "github.com/elastic/go-elasticsearch/v7/esapi"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/globalsign/mgo/bson"
)

var (
	TrendingBgs = []string{
		"https://manch-dev.s3.ap-south-1.amazonaws.com/app-banners/bg_trending_blue.jpg'",
		"https://manch-dev.s3.ap-south-1.amazonaws.com/app-banners/bg_trending_green.jpg'",
		"https://manch-dev.s3.ap-south-1.amazonaws.com/app-banners/bg_trending_red.jpg'",
		"https://manch-dev.s3.ap-south-1.amazonaws.com/app-banners/bg_trending_voilet.jpg",
	}
)

type StringInterface map[string]interface{}

type TypeInput struct {
	Input  []string `json:"input"`
	Weight int      `json:"weight"`
}

type Resurfaced struct {
	ResurfacedStart string `json:"resurfaced_start"`
	ResurfacedEnd   string `json:"resurfaced_end"`
	NoOfPosts       int    `json:"no_of_posts"`
	AdditionalScore int    `json:"additional_score"`
}

type HashTag struct {
	ID                 string       `json:"id"`
	Keyword            TypeInput    `json:"keyword"`
	TagName            string       `json:"tagname"`
	Title              string       `json:"title"`
	Image              string       `json:"image"`
	NoOfPosts          int          `json:"no_of_posts"`
	NoOfLikes          int          `json:"no_of_likes"`
	NoOfComments       int          `json:"no_of_comments"`
	ActualCreationTime string       `json:"actual_creation_time"`
	LastUpdatedTime    string       `json:"last_updated_time"`
	Resurfaced         bool         `json:"resurfaced"`
	ResurfacedDate     string       `json:"resurfaced_date"`
	AdditionalScore    int          `json:"additional_score"`
	ResurfacedArchive  []Resurfaced `json:"resurfaced_archive"`
}

func getTrendingBg() string {
	randomIndex := utils.Random(0, len(TrendingBgs))
	return TrendingBgs[randomIndex]
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
	fmt.Println("response getdocumentbyid", res)
	if res.IsError() {
		log.Printf("Error on getting data")
		return errors.New("Error on getting data"), r
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return errors.New("Error parsing the response body"), r
	}
	return nil, r
}

func AddTagToIndex(tags []string, additionalScore int, tagsPositions []mongo.TagPositions) {

	currentISOTime := utils.ISOFormat(time.Now())
	image := getTrendingBg()
	hashTagData := HashTag{
		ActualCreationTime: currentISOTime,
		LastUpdatedTime:    currentISOTime,
		ResurfacedDate:     currentISOTime,
		Image:              image,
		AdditionalScore:    additionalScore,
		NoOfPosts:          1,
	}

	for i, tag := range tags {

		func(tagName string) {
			fmt.Println("indexing..", tagName)

			// git title
			var title string
			for _, tagPosition := range tagsPositions {
				if strings.ToLower(tagPosition.Tag) == tagName {
					title = tagPosition.Tag
				}
			}
			hashTagData.ID = strings.ToLower(tagName)
			tokenizeText := utils.TokenizeText(tagName, 4)

			hashTagData.Keyword = TypeInput{
				Input: tokenizeText,
			}
			hashTagData.Title = title
			hashTagData.TagName = title
			// upsert data
			var upsertData StringInterface
			hashTagDataEncoded, _ := json.Marshal(hashTagData)
			json.Unmarshal(hashTagDataEncoded, &upsertData)

			fmt.Println("uspert data", upsertData)
			// Build the request body.
			body := esutil.NewJSONReader(StringInterface{
				"script": StringInterface{
					"source": "ctx._source.no_of_posts += params.count;ctx._source.last_updated_time=params.last_updated;if(ctx._source.additional_score != null) {ctx._source.additional_score += params.additional_score} else {ctx._source.additional_score = params.additional_score}",
					"lang":   "painless",
					"params": StringInterface{
						"count":            1,
						"last_updated":     hashTagData.LastUpdatedTime,
						"additional_score": hashTagData.AdditionalScore,
					},
				},
				"upsert": upsertData,
			})
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
		SourceIncludes: []string{"tagname"},
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
		source := v.(map[string]interface{})["_source"].(map[string]interface{})
		tagname := source["tagname"]
		response = append(response, map[string]interface{}{"_id": id, "tagname": tagname})
	}
	return nil, response
}

func getScore(baseTime time.Time, noOfPost int, additionScore int) int {
	return (int(baseTime.Unix()) + noOfPost*10*60 + additionScore)
}

/*
* update hashtag weight and returns the weight
* it takes tagName and additionScore
 */
func UpdateTagWeight(tag string, additionScore int, isTrending bool) (error, map[string]interface{}) {
	response := map[string]interface{}{}
	err, doc := GetDocumentById(tag, "tags")

	if err != nil {
		return err, response
	}
	source := doc["_source"].(map[string]interface{})
	noOfPost := source["no_of_posts"].(float64)
	tagname := source["tagname"].(string)
	baseTime := source["resurfaced_date"].(string)
	resurfacedDate := utils.ParseISOToTime(baseTime)
	currentTime := time.Now()
	diff := currentTime.Sub(resurfacedDate)
	hours := diff.Hours()
	resurfaced := false

	if hours > 24 && !isTrending {
		resurfaced = true
	}

	fmt.Println("source is", source["additional_score"])
	if source["additional_score"] != nil {
		additionScore = int(source["additional_score"].(float64))
	}

	weight := getScore(utils.ParseISOToTime(baseTime), int(noOfPost), additionScore)
	body := esutil.NewJSONReader(StringInterface{
		"script": StringInterface{
			"source": "ctx._source.keyword.weight = params.weight",
			"params": StringInterface{
				"weight": weight,
			},
		},
	})

	if resurfaced {
		weight = getScore(currentTime, 0, 0)
		noOfPost = 0
		body = esutil.NewJSONReader(StringInterface{
			"script": StringInterface{
				"source": "ctx._source.keyword.weight = params.weight;ctx._source.no_of_posts=params.count;ctx._source.additional_score=params.count;ctx._source.resurfaced_date=params.current_date;ctx._source.resurfaced=true;if(ctx._source.resurfaced_archive != null){ctx._source.resurfaced_archive.add(params.resurfaced_archive)}else{ctx._source.resurfaced_archive=[params.resurfaced_archive]}",
				"params": StringInterface{
					"weight":       weight,
					"count":        0,
					"current_date": utils.ISOFormat(currentTime),
					"resurfaced_archive": StringInterface{
						"resurfaced_start": baseTime,
						"resurfaced_end":   utils.ISOFormat(currentTime),
						"no_of_posts":      noOfPost,
						"additional_score": additionScore,
					},
				},
			},
		})
	}

	fmt.Println("weight is", weight)

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
		return errors.New("Error getting response"), response
	}
	defer res.Body.Close()
	fmt.Println("response is", res)
	if res.IsError() {
		return errors.New("Error getting response"), response
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return errors.New("Error parsing the response body"), response
	}
	log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	response["weight"] = weight
	response["tagname"] = tagname
	response["no_of_posts"] = int(noOfPost)
	return nil, response
}

func GetImageById(id string) (error, string) {
	err, doc := GetDocumentById(id, "tags")
	if err != nil {
		return err, ""
	}
	return nil, doc["_source"].(map[string]interface{})["image"].(string)
}

func UpdateImageById(id, imageUrl string) (error, string) {

	body := esutil.NewJSONReader(StringInterface{
		"script": StringInterface{
			"source": "ctx._source.image = params.image_url;",
			"lang":   "painless",
			"params": StringInterface{
				"image_url": imageUrl,
			},
		},
	})
	// create update request
	req := esapi.UpdateRequest{
		Index:      "tags",
		DocumentID: id,
		Body:       body,
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err, ""
	}
	defer res.Body.Close()
	fmt.Println("response is", res)
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), id)
		return err, ""
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return err, ""
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			return nil, imageUrl
		}
	}
}
