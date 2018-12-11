package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
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
			fmt.Println("Self Comment")
			return
		}
		postCreator := mongo.GetProfileById(comment.Post.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Post.Created.ProfileId})
		notification := mongo.CreateNotification(comment.PostId, "comment", "post", comment.Post.Created.ProfileId)

		postTitle := utils.TruncateTitle(comment.Post.Title, 4)
		data := i18n.DataModel{
			Name:  comment.Created.Name,
			Count: uniqueCommentator - 1,
			Post:  postTitle,
		}
		var msgStr string

		if uniqueCommentator > 1 {
			msgStr = i18n.GetString(postCreator.Language, "comment_multi", data)
		} else {
			msgStr = i18n.GetString(postCreator.Language, "comment_one", data)
		}
		title := i18n.GetAppTitle(postCreator.Language)
		msg := firebase.ManchMessage{
			Title:      title,
			Message:    msgStr,
			Icon:       mongo.ExtractThumbNailFromPost(comment.Post),
			DeepLink:   "manch://posts/" + comment.PostId.Hex(),
			BadgeCount: strconv.Itoa(comment.Post.CommentCount),
			Id:         notification.Identifier,
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
			fmt.Println("Invalid Vote")
			//Do not process downvotes and unvote
			return
		}
		post := mongo.GetPostById(v.Resource)
		vote := post.GetVote(v.Id)
		if vote.Created.ProfileId == post.Created.ProfileId {
			//Self Vote
			fmt.Println("Self Vote")
			return
		}
		postCreator := mongo.GetProfileById(post.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		notification := mongo.CreateNotification(post.Id, "like", "post", vote.Created.ProfileId)

		postTitle := utils.TruncateTitle(post.Title, 4)
		data := i18n.DataModel{
			Name:  vote.Created.Name,
			Post:  postTitle,
			Count: post.UpVotes,
		}
		var msgStr string
		if post.UpVotes > 1 {
			msgStr = i18n.GetString(postCreator.Language, "post_like_multi", data)
		} else {
			msgStr = i18n.GetString(postCreator.Language, "post_like_one", data)
		}
		title := i18n.GetAppTitle(postCreator.Language)
		msg := firebase.ManchMessage{
			Title:    title,
			Message:  msgStr,
			Icon:     mongo.ExtractThumbNailFromPost(post),
			DeepLink: "manch://posts/" + post.Id.Hex(),
			Id:       notification.Identifier,
		}

		fmt.Printf("\nGCM Message %+v", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}

	})
	subscribers.VoteCommentSubscriber(func(subj, reply string, v *subscribers.Vote) {
		fmt.Printf("\nNats MSG %+v", v)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.VoteCommentSubscriber", r)
			}
		}()
		dir, err := strconv.Atoi(v.Direction)
		if err != nil || dir < 1 {
			fmt.Println("Invalid Vote")
			//Do not process downvotes and unvote
			return
		}
		comment := mongo.GetCommentById(v.Resource)
		vote := comment.GetVote(v.Id)
		if vote.Created.ProfileId == comment.Created.ProfileId {
			//Self Vote
			fmt.Println("Self Vote")
			return
		}
		post := mongo.GetPostById(comment.PostId.Hex())
		commentCreator := mongo.GetProfileById(comment.Created.ProfileId)
		notification := mongo.CreateNotification(comment.Id, "like", "comment", vote.Created.ProfileId)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Created.ProfileId})

		commentTitle := utils.TruncateTitle(comment.Content, 4)

		data := i18n.DataModel{
			Name:    vote.Created.Name,
			Comment: commentTitle,
			Count:   comment.UpVotes,
		}
		var msgStr string
		if comment.UpVotes > 1 {
			msgStr = i18n.GetString(commentCreator.Language, "comment_like_multi", data)
		} else {
			msgStr = i18n.GetString(commentCreator.Language, "comment_like_one", data)
		}
		title := i18n.GetAppTitle(commentCreator.Language)
		msg := firebase.ManchMessage{
			Title:    title,
			Message:  msgStr,
			Icon:     mongo.ExtractThumbNailFromPost(post),
			DeepLink: "manch://posts/" + comment.PostId.Hex(),
			Id:       notification.Identifier,
		}

		fmt.Printf("\nGCM Message %+v", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}

	})

	subscribers.UserSubscriber(func(subj, reply string, u *subscribers.User) {
		fmt.Printf("Received user %v", u)
		fmt.Println("received new user")
		// create follow schedule for this user
		// get all bot profiles
		botUsers := mongo.GetBotUsers()
		var resourceId bson.ObjectId

		// shuffle this array
		var botProfilesIds [100]string
		i := 0
		for _, botUser := range botUsers {
			profiles := botUser.Profiles
			for _, profile := range profiles {
				if i == 100 {
					break
				}
				// append profile.Id
				fmt.Println("profile:", profile.Id.Hex())
				botProfilesIds[i] = profile.Id.Hex()
				i++
			}
		}

		// shuffle the bot profiles ids
		fmt.Println("bot profile ids:", botProfilesIds)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(i, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })
		fmt.Println("after shuffle:", botProfilesIds)

		user := mongo.GetUserById(u.Id)
		userProfileId := user.Profiles[0].Id
		resourceId = userProfileId


		c := mongo.Creator{
			Id:        bson.NewObjectId(),
			ProfileId: userProfileId,
			Name:      user.Name,
			UserType:  user.UserType,
		}

		current := time.Now()
		fmt.Println("current time:", current)

		// 0-5th minute - +5 followes
		rMinute := utils.Random(0, 5)
		fmt.Println("random minute:", rMinute)
		t := current.Add(time.Duration(rMinute) * time.Minute)
		fmt.Println("time:", t)
		j := 0
		followers := 5
		for ; j < followers; j++ {
			doc := mongo.CreateFollowSchedule(t, bson.ObjectIdHex(botProfilesIds[j]), resourceId, c)
			fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 5 minuts to 1 hours - +5
		followers += 5
		rMinute = utils.Random(5, 60)
		fmt.Println("random minute:", rMinute)
		t = current.Add(time.Duration(rMinute) * time.Minute)
		fmt.Println("time:", t)
		for ; j < followers; j++ {
			doc := mongo.CreateFollowSchedule(t, bson.ObjectIdHex(botProfilesIds[j]), resourceId, c)
			fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 1 Hr to 6Hr +5 followers
		followers += 5
		rHour := utils.Random(1, 6)
		fmt.Println("random minute:", rMinute)
		t = current.Add(time.Duration(rHour) * time.Hour)
		fmt.Println("time:", t)
		for ; j < followers; j++ {
			doc := mongo.CreateFollowSchedule(t, bson.ObjectIdHex(botProfilesIds[j]), resourceId, c)
			fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 6 Hr to 24 Hr +5 followers
		followers += 5
		rHour = utils.Random(6, 24)
		t = current.Add(time.Duration(rHour) * time.Hour)
		fmt.Println("time:", t)
		for ; j < followers; j++ {
			doc := mongo.CreateFollowSchedule(t, bson.ObjectIdHex(botProfilesIds[j]), resourceId, c)
			fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 1st to 3rd day +10 followers
		followers += 10
		rDay := utils.Random(1, 3)
		t = current.AddDate(0, 0, rDay)
		for ; j < followers; j++ {
			doc := mongo.CreateFollowSchedule(t, bson.ObjectIdHex(botProfilesIds[j]), resourceId, c)
			mongo.AddFollowSchedule(doc)
		}

		// 3rd to 7th day +10 followers
		followers += 5
		rDay = utils.Random(3, 7)
		t = current.AddDate(0, 0, rDay)
		for ; j < followers; j++ {
			doc := mongo.CreateFollowSchedule(t, bson.ObjectIdHex(botProfilesIds[j]), resourceId, c)
			mongo.AddFollowSchedule(doc)
		}
	})

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
