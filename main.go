package main

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	"net/http"
	"notification-service/pkg/firebase"
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
		comment, uniqueCommentator := mongo.GetFullCommentById(cmt.Id)
		if comment.Post.Created.ProfileId == comment.Created.ProfileId {
			//Self comment
			return
		}
		profiles := []bson.ObjectId{comment.Post.Created.ProfileId}
		tokens := mongo.GetTokensByProfiles(profiles)
		var msgStr string
		if uniqueCommentator > 2 {
			msgStr = comment.Created.Name + " & " + strconv.Itoa(uniqueCommentator-1) + " others commented on Your Post"
		} else if uniqueCommentator == 2 {
			msgStr = comment.Created.Name + " & " + "one other commented on Your Post"
		} else {
			msgStr = comment.Created.Name + " commented on Your Post"
		}
		msg := firebase.ManchMessage{
			Title:      "New Comment",
			Message:    msgStr,
			Icon:       comment.Created.Avatar,
			DeepLink:   "manch://posts/" + comment.PostId.Hex(),
			BadgeCount: strconv.Itoa(comment.Post.CommentCount),
			Id:         comment.PostId.Hex() + "_comment",
		}
		fmt.Printf("\n Avatar: %+v", comment.Created)
		fmt.Printf("\n Message: %+v", msg)
		firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
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
