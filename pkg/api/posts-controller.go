package api

import (
	"fmt"
	"net/http"
	"notification-service/pkg/elasticsearch"

	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/render"
)

func HandlePostSearch(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	fmt.Println("queries are", queries)
	bsonQuery := bson.M{}
	limit, skip, _ := ExtractStandardOptions(r)
	bsonQuery["skip"] = skip
	bsonQuery["limit"] = limit

	if len(queries["q"]) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("q is required"))
		return
	}
	bsonQuery["q"] = queries["q"][0]
	// get notification
	err, response := elasticsearch.SearchPost(bsonQuery)

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
