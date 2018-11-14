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
		newPost := mongo.GetPostById(p.Id)
		fmt.Printf("Mongo return for Post %+v\n", newPost)

	})
	subscribers.CommentSubscriber(func(subj, reply string, c *subscribers.Comment) {
		fmt.Printf("\nNats MSG %+v", c)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.CommentSubscriber", r)
			}
		}()
		comment, uniqueCommentator := mongo.GetFullCommentById(c.Id)
		if comment.Post.Created.ProfileId == comment.Created.ProfileId {
			//Self comment
			return
		}
		postCreator := mongo.GetProfileById(comment.Post.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Post.Created.ProfileId})
		notification := mongo.CreateNotification(comment.PostId, "comment", "post", comment.Post.Created.ProfileId)

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
			Id:         notification.Id.Hex(),
		}
		//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
	})
	subscribers.VoteSubscriberPost(func(subj, reply string, v *subscribers.Vote) {
		//fmt.Printf("\nNats MSG %+v", v)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.VoteSubscriberPost", r)
			}
		}()
		dir, err := strconv.Atoi(v.Direction)
		if err == nil || dir < 1 {
			print(dir)
			//Do not process downvotes and unvote
			return
		}
		post := mongo.GetPostById(v.Resource)
		vote := post.GetVote(v.Id)
		if vote.Created.ProfileId == post.Created.ProfileId {
			//Self Vote
			return
		}
		postCreator := mongo.GetProfileById(post.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		print(tokens)
		notification := mongo.CreateNotification(post.Id, "vote", "post", vote.Created.ProfileId)

		data := i18n.DataModel{
			Name: vote.Created.Name,
			Post: post.Title,
		}
		var msgStr string
		msgStr = i18n.GetString(postCreator.Language, "post_like_one", data)
		msg := firebase.ManchMessage{
			Title:    "Manch",
			Message:  msgStr,
			Icon:     vote.Created.Avatar,
			DeepLink: "manch://posts/" + post.Id.Hex(),
			Id:       notification.Id.Hex(),
		}

		fmt.Printf("\nMessage %+v", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}

	})
	subscribers.VoteSubscriberComment(func(subj, reply string, v *subscribers.Vote) {
		//fmt.Printf("\nNats MSG %+v", v)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.VoteSubscriberComment", r)
			}
		}()
		dir, err := strconv.Atoi(v.Direction)
		if err == nil || dir < 1 {
			//Do not process downvotes and unvote
			return
		}
		comment := mongo.GetCommentById(v.Resource)
		vote := comment.GetVote(v.Id)
		commentCreator := mongo.GetProfileById(comment.Created.ProfileId)
		notification := mongo.CreateNotification(comment.Id, "vote", "comment", vote.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Created.ProfileId})
		data := i18n.DataModel{
			Name: vote.Created.Name,
			Post: comment.Content,
		}
		var msgStr string

		msgStr = i18n.GetString(commentCreator.Language, "comment_like_one", data)
		msg := firebase.ManchMessage{
			Title:    "Manch",
			Message:  msgStr,
			Icon:     vote.Created.Avatar,
			DeepLink: "manch://posts/" + comment.PostId.Hex(),
			Id:       notification.Id.Hex(),
		}

		fmt.Printf("\nMessage %+v", msg)
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
