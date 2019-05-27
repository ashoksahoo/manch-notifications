package api

import (
	"fmt"
	"net/http"
	"notification-service/pkg/elasticsearch"
	"strconv"

	"github.com/globalsign/mgo/bson"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

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
	metaInfo := bson.M{"limit": limit, "skip": skip}
	metaInfo["offset"] = skip + len(response.([]interface{}))

	if len(response.([]interface{})) == 0 {
		metaInfo["has_next"] = false
	} else {
		metaInfo["has_next"] = true
	}

	render.JSON(w, r, bson.M{"data": response, "meta": metaInfo})
}

func GetTagByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err, response := elasticsearch.GetDocumentById(id, "tags")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, bson.M{"data": response})
}
