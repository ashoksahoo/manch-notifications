package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notification-service/pkg/elasticsearch"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type UpdateMeta struct {
	AdditionalScore int    `json:"additional_score"`
	ImageUrl        string `json:"image_url"`
	IsTrending      bool   `json:"is_trending"`
	Title           string `json:"title"`
}

func SearchHashTags(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	fmt.Println("queries are", queries)
	bsonQuery := bson.M{}
	_, skipOK := queries["skip"]
	_, offsetOK := queries["offset"]
	_, limitOk := queries["limit"]
	var limit, skip int
	if skipOK {
		skip, _ = strconv.Atoi(queries["skip"][0])
		delete(queries, "skip")
	} else if offsetOK {
		skip, _ = strconv.Atoi(queries["offset"][0])
		delete(queries, "offset")
	} else {
		skip = DEFAULT_SKIP
	}

	if limitOk {
		limit, _ = strconv.Atoi(queries["limit"][0])
		delete(queries, "limit")
	} else {
		limit = DEFAULT_LIMIT
	}
	bsonQuery["skip"] = skip
	bsonQuery["limit"] = limit

	if len(queries["keyword"]) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("keyword is required"))
		return
	}
	bsonQuery["keyword"] = queries["keyword"][0]
	// get notification
	err, response := elasticsearch.SearchHashTags(bsonQuery)

	fmt.Println("error ", err)
	fmt.Printf("response is \n%+v", response)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	// create meta info
	responseLength := len(response.([]map[string]interface{}))
	metaInfo := bson.M{"limit": limit, "skip": skip}
	metaInfo["offset"] = skip + responseLength

	if responseLength == 0 {
		metaInfo["has_next"] = false
	} else {
		metaInfo["has_next"] = true
	}

	render.JSON(w, r, bson.M{"data": response, "meta": metaInfo})
}

func GetTagByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id = strings.ToLower(id)
	err, response := elasticsearch.GetDocumentById(id, "tags")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, bson.M{"data": response})
}

func UpdateHashtagWeight(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id = strings.ToLower(id)
	updateMeta := UpdateMeta{}
	err := json.NewDecoder(r.Body).Decode(&updateMeta)
	err, response := elasticsearch.UpdateTagWeight(id, updateMeta.AdditionalScore, updateMeta.IsTrending)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, bson.M{"data": response})
}

func GetHashTagImageById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id = strings.ToLower(id)
	err, response := elasticsearch.GetImageById(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, bson.M{"data": map[string]interface{}{"image_url": response}})
}

func UpdateHashtagImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id = strings.ToLower(id)
	updateMeta := UpdateMeta{}
	err := json.NewDecoder(r.Body).Decode(&updateMeta)
	err, response := elasticsearch.UpdateImageById(id, updateMeta.ImageUrl)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, bson.M{"data": map[string]interface{}{"image_url": response}})
}

func UpdateHashtagTitle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	id = strings.ToLower(id)
	updateMeta := UpdateMeta{}
	err := json.NewDecoder(r.Body).Decode(&updateMeta)
	err, title, currentTitle := elasticsearch.UpdateTitleById(id, updateMeta.Title)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, bson.M{"data": map[string]interface{}{"title": title, "current_title": currentTitle}})
}
