package main

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	"net/http"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"strconv"
)

func main() {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	subscribers.PostSubscriber(func(subj, reply string, p *subscribers.Post) {
		fmt.Printf("Received a post on subject %s! with Post ID %s\n", subj, p.Id)
		newPost := mongo.GetPostById(bson.ObjectIdHex(p.Id))
		fmt.Printf("Mongo return for Post %+v\n", newPost)

	})
	subscribers.CommentSubscriber(func(subj, reply string, cmt *subscribers.Comment) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		comment, uniqueCommentator := mongo.GetFullCommentById(cmt.Id)
		if comment.Post.Created.ProfileId == comment.Created.ProfileId {
			//Self comment
			return
		}
		postCreator := mongo.GetProfileById(comment.Post.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Post.Created.ProfileId})
		mongo.CreateNotification(comment.PostId, "comment", "post",comment.Post.Created.ProfileId)

		data := i18n.DataModel{
			Name:  comment.Created.Name,
			Count: uniqueCommentator - 1,
			Post:  comment.Post.Title,
		}
		var msgStr string

		if uniqueCommentator > 1 {
			msgStr = i18n.GetString(postCreator.Language, "comment_multi", data)
		} else {
			msgStr = i18n.GetString(postCreator.Language, "comment_one", data)
		}
		msg := firebase.ManchMessage{
			Title:      "Manch",
			Message:    msgStr,
			Icon:       comment.Created.Avatar,
			DeepLink:   "manch://posts/" + comment.PostId.Hex(),
			BadgeCount: strconv.Itoa(comment.Post.CommentCount),
			Id:         comment.PostId.Hex() + "_comment",
		}
		fmt.Printf("\nMessage %+v", msg)
		//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
	})

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
