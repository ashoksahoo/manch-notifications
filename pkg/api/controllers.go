package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notification-service/pkg/mongo"
	"reflect"
	"strconv"

	"github.com/globalsign/mgo/bson"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
)

var (
	DEFAULT_LIMIT = 20
	DEFAULT_SKIP = 0
)


func GetAllNotification(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
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
	for k, v := range queries {
		fmt.Println(k, v[0])
		if bson.IsObjectIdHex(v[0]) {
			bsonQuery[k] = bson.ObjectIdHex(v[0])
		} else {
			bsonQuery[k] = v[0]
		}
	}

	// get notification
	_, notifications := mongo.GetNotificationByQuery(bsonQuery)

	// create meta info
	metaInfo := bson.M{"limit": limit, "skip": skip}
	metaInfo["offset"] = skip + len(notifications)

	if len(notifications) == 0 {
		metaInfo["has_next"] = false
	} else {
		metaInfo["has_next"] = true
	}
	fmt.Println(reflect.TypeOf(notifications))
	fmt.Println("notification is ", notifications)

	render.JSON(w, r, bson.M{"data": notifications, "meta": metaInfo})
}

func GetNotificationById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var objectId bson.ObjectId
	if bson.IsObjectIdHex(id) {
		objectId = bson.ObjectIdHex(id)
	}
	_, notification := mongo.GetNotificationById(objectId)
	render.JSON(w, r, bson.M{"data": notification})
}

func UpdateNotificationById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	notification := mongo.NotificationModel{}
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		panic(err)
	}
	fmt.Println(notification)
	render.JSON(w, r, map[string]string{"message": "Update Notification By ID", "id": id})
}

func DeleteNotificationById(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"message": "Delete Notification By Id"})
}
