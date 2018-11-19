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
	/**
		This processes Comments from Posts
		1) Get Comment Details and Unique commentator count
		2) Validate self comment
		3) Get Who created the post -> He gets the notification and we need his current lang
		4) Get tokens from the above profile (Supports multiple device tokens.)
		5) Create/Update Notification Table which has the meta info for the notificaiotn
		6) Construct Data for i18n template
		7) Generate template using template data and String Formatter
		8) Create push notification
		9) Fire the notifications in routines.

	 */
	subscribers.CommentSubscriber(func(subj, reply string, c *subscribers.Comment) {
		//fmt.Printf("\nNats MSG %+v", c)
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
	/**
		This processes Upvotes from Posts
		1) Get Voting Details
		2) Validate only upvote & self vote
		3) Get Who created the post -> He gets the notification and we need his current lang
		4) Get tokens from the above profile (Supports multiple device tokens.)
		5) Create/Update Notification Table which has the meta info for the notificaiotn
		6) Construct Data for i18n template
		7) Generate template using template data and String Formatter
		8) Create push notification
		9) Fire the notifications in routines.

	 */
	subscribers.VotePostSubscriber(func(subj, reply string, v *subscribers.Vote) {
		//fmt.Printf("\nNats MSG %+v", v)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.VotePostSubscriber", r)
			}
		}()
		dir, err := strconv.Atoi(v.Direction)
		if err != nil || dir < 1 {
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
		notification := mongo.CreateNotification(post.Id, "vote", "post", vote.Created.ProfileId)
		//fmt.Printf("\nPost %+v", post)
		data := i18n.DataModel{
			Name:  vote.Created.Name,
			Post:  post.Title,
			Count: post.UpVotes,
		}
		var msgStr string
		if post.UpVotes > 1 {
			msgStr = i18n.GetString(postCreator.Language, "post_like_multi", data)
		} else {
			msgStr = i18n.GetString(postCreator.Language, "post_like_one", data)
		}
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
	subscribers.VoteCommentSubscriber(func(subj, reply string, v *subscribers.Vote) {
		//fmt.Printf("\nNats MSG %+v", v)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.VoteCommentSubscriber", r)
			}
		}()
		dir, err := strconv.Atoi(v.Direction)
		if err != nil || dir < 1 {
			//Do not process downvotes and unvote
			return
		}
		comment := mongo.GetCommentById(v.Resource)
		vote := comment.GetVote(v.Id)
		if vote.Created.ProfileId == comment.Created.ProfileId {
			//Self Vote
			return
		}
		commentCreator := mongo.GetProfileById(comment.Created.ProfileId)
		notification := mongo.CreateNotification(comment.Id, "vote", "comment", vote.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Created.ProfileId})
		data := i18n.DataModel{
			Name:    vote.Created.Name,
			Comment: comment.Content,
			Count:   comment.UpVotes,
		}
		var msgStr string
		if comment.UpVotes > 1 {
			msgStr = i18n.GetString(commentCreator.Language, "comment_like_multi", data)
		} else {
			msgStr = i18n.GetString(commentCreator.Language, "comment_like_one", data)
		}
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
