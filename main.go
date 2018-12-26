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
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(utils.SplitTimeInRange(0, 5*60, 6, time.Second))
		w.Write([]byte("pong"))
	})

	subscribers.PostSubscriber(func(subj, reply string, p *subscribers.Post) {
		fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
		// get all bot users
		botUsers := mongo.GetBotUsers()

		// array of bot profiles ids
		var botProfilesIds [100]string
		// no. of profiles counter
		i := 0
		for _, botUser := range botUsers {
			profiles := botUser.Profiles
			for _, profile := range profiles {
				if i == 100 {
					break
				}
				botProfilesIds[i] = profile.Id.Hex()
				i++
			}
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(i, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })
		var no_of_votes int
		if p.CreatorType == "bot" {
			no_of_votes = utils.Random(30, 45)
		} else if p.CreatorType == "super_level_1" {
			no_of_votes = utils.Random(20, 25)
		} else {
			no_of_votes = utils.Random(5, 10)
		}

		j := 0
		fmt.Println("no_of_votes: ", no_of_votes)
		t := utils.SplitTimeInRange(1*60, 30*60, no_of_votes, time.Second)
		for k := 0; j < no_of_votes; j, k = j+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
			mongo.AddVoteSchedule(vote)
		}

		randomVotes := utils.Random(5, 20)
		no_of_votes += randomVotes
		fmt.Println("no_of_votes: ", no_of_votes)
		t = utils.SplitTimeInRange(30, 2*24*60, randomVotes, time.Minute)
		for k := 0; j < no_of_votes; j, k = j+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
			mongo.AddVoteSchedule(vote)
		}
		fmt.Printf("Processed a post on subject %s! with Post Id%s\n", subj, p.Id)
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
		fmt.Printf("Received a comment on subject %s! with Comment %+v\n", subj, c)
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

		// comment on comment
		if len(comment.Parents) >= 2 {
			fmt.Println("reply on comments")
			// get comment id
			commentId := comment.CommentId
			replyOnComment := mongo.GetCommentById(commentId.Hex())
			fmt.Printf("reply on comments %+v\n", replyOnComment)
			// get replied on comment creator
			replyOnCommentCreator := mongo.GetProfileById(replyOnComment.Created.ProfileId)
			notification1 := mongo.CreateNotification(replyOnComment.Id, "comment", "comment", replyOnComment.Created.ProfileId)
			tokens1 := mongo.GetTokensByProfiles([]bson.ObjectId{replyOnComment.Created.ProfileId})
			// comment title
			count := len(notification.UniqueUsers) - 1
			fmt.Println("Comment count: ", count)
			commentTitle := replyOnComment.Content
			utils.TruncateTitle(commentTitle, 4)
			data1 := i18n.DataModel{
				Name:    comment.Created.Name,
				Comment: commentTitle,
				Count: count,
			}
			var msgStr1 string
			if count > 0 {
				msgStr1 = i18n.GetString(replyOnCommentCreator.Language, "comment_reply_multi", data1)				
			} else {
				msgStr1 = i18n.GetString(replyOnCommentCreator.Language, "comment_reply_one", data1)
			}
			fmt.Println("message string: ", msgStr1)
			title := i18n.GetAppTitle(replyOnCommentCreator.Language)
			msgStr1 = strings.Replace(msgStr1, "\"\" ", "", 1)
			msg := firebase.ManchMessage{
				Title:      title,
				Message:    msgStr1,
				Icon:       mongo.ExtractThumbNailFromPost(comment.Post),
				DeepLink:   "manch://posts/" + comment.PostId.Hex(),
				BadgeCount: strconv.Itoa(replyOnComment.CommentCount),
				Id:         notification1.Identifier,
			}
			if tokens1 != nil {
				for _, token := range tokens {
					go firebase.SendMessage(msg, token.Token)
				}
			} else {
				fmt.Printf("No token")
			}
			fmt.Println("end reply on comments")
		}

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
		msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
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
		fmt.Printf("Processed a comment on subject %s! with Comment %s\n", subj, c.Id)

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
		fmt.Printf("Received a vote on subject %s! with vote %+v\n", subj, v)
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
		msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
		msg := firebase.ManchMessage{
			Title:    title,
			Message:  msgStr,
			Icon:     mongo.ExtractThumbNailFromPost(post),
			DeepLink: "manch://posts/" + post.Id.Hex(),
			Id:       notification.Identifier,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
		fmt.Printf("Processed a vote on subject %s! with vote Id %s\n", subj, v.Id)

	})

	subscribers.VoteCommentSubscriber(func(subj, reply string, v *subscribers.Vote) {
		fmt.Printf("Received a Vote on subject %s! with Vote %+v\n", subj, v)

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
		msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
		title := i18n.GetAppTitle(commentCreator.Language)
		msg := firebase.ManchMessage{
			Title:    title,
			Message:  msgStr,
			Icon:     mongo.ExtractThumbNailFromPost(post),
			DeepLink: "manch://posts/" + comment.PostId.Hex(),
			Id:       notification.Identifier,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token\n")
		}
		fmt.Printf("Processed a Vote on subject %s! with Vote Id %s\n", subj, v.Id)

	})

	subscribers.UserSubscriber(func(subj, reply string, u *subscribers.User) {
		fmt.Printf("Received a New User on subject %s! with User %+v\n", subj, u)
		// create follow schedule for this user

		// get all bot users
		botUsers := mongo.GetBotUsers()
		var resourceId bson.ObjectId

		// array of bot profiles ids
		var botProfilesIds [100]string

		// no. of profiles counter
		i := 0
		for _, botUser := range botUsers {
			profiles := botUser.Profiles
			for _, profile := range profiles {
				if i == 100 {
					break
				}
				botProfilesIds[i] = profile.Id.Hex()
				i++
			}
		}

		// shuffle the bot profiles ids
		// fmt.Println("bot profile ids before shuffle:", botProfilesIds)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(i, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })
		// fmt.Println("after shuffle:", botProfilesIds)

		// get user from db
		user := mongo.GetUserById(u.Id)
		userProfileId := user.Profiles[0].Id

		// set user to resource
		resourceId = userProfileId

		// 0-5th minute - +5 followes
		j := 0
		randomFollowers := utils.Random(3, 10)
		followers := randomFollowers
		t := utils.SplitTimeInRange(1, 5, randomFollowers, time.Minute)
		for k := 0; j < followers; j, k = j+1, k+1 {
			doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
			// fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 5 minuts to 1 hours - +5
		randomFollowers = utils.Random(5, 10)
		t = utils.SplitTimeInRange(5, 59, randomFollowers, time.Minute)
		followers += randomFollowers
		for k := 0; j < followers; j, k = j+1, k+1 {
			doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
			// fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 1 Hr to 6Hr +5-10 followers
		randomFollowers = utils.Random(5, 10)
		t = utils.SplitTimeInRange(1, 6, randomFollowers, time.Hour)
		followers += randomFollowers
		for k := 0; j < followers; j, k = j+1, k+1 {
			doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
			// fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 6 Hr to 24 Hr +5-10 followers
		randomFollowers = utils.Random(5, 10)
		t = utils.SplitTimeInRange(6, 24, randomFollowers, time.Hour)
		followers += randomFollowers
		for k := 0; j < followers; j, k = j+1, k+1 {
			doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
			// fmt.Printf("saving doc:%+v\n", doc)
			mongo.AddFollowSchedule(doc)
		}

		// 1st to 3rd day +10-15 followers
		randomFollowers = utils.Random(20, 30)
		t = utils.SplitTimeInRange(24, 72, randomFollowers, time.Hour)
		followers += randomFollowers
		for k := 0; j < followers; j, k = j+1, k+1 {
			doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
			mongo.AddFollowSchedule(doc)
		}

		// 3rd to 7th day +10-20 followers
		randomFollowers = utils.Random(20, 30)
		t = utils.SplitTimeInRange(72, 168, randomFollowers, time.Hour)
		followers += randomFollowers
		for k := 0; j < followers; j, k = j+1, k+1 {
			doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
			mongo.AddFollowSchedule(doc)
		}
		fmt.Println("total followers added:", followers)
		fmt.Printf("Processed a New User on subject %s! with User Id %s\n", subj, u.Id)
	})

	subscribers.UserFollowSubscriber(func(subj, reply string, uf *subscribers.Subscription) {
		fmt.Printf("Received a User follow on subject %s! with user follow %+v\n", subj, uf)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in subscribers.UserFollowSubscriber", uf)
			}
		}()

		if uf.ResourceType != "user" {
			fmt.Println("Not a user resource follows")
			return
		}

		userFollow := mongo.GetUserFollowById(uf.Id)
		// fmt.Printf("\nuser follow %+v\n", userFollow)
		follower := mongo.GetProfileById(userFollow.ProfileId)
		// fmt.Printf("\nfollower %+v\n", follower)
		followsTo := mongo.GetProfileById(userFollow.ResourceId)
		// fmt.Printf("\nfollowsTo %+v\n", followsTo)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{userFollow.ResourceId})
		notification := mongo.CreateNotification(followsTo.Id, "follows", "user", follower.Id)

		count := len(notification.UniqueUsers) - 1
		data := i18n.DataModel{
			Name:  follower.Name,
			Count: count,
		}
		var msgStr string

		if len(notification.UniqueUsers) > 1 {
			msgStr = i18n.GetString(followsTo.Language, "follow_user_multi", data)
		} else {
			msgStr = i18n.GetString(followsTo.Language, "follow_user_one", data)
		}

		title := i18n.GetAppTitle(followsTo.Language)
		msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
		msg := firebase.ManchMessage{
			Title:    title,
			Message:  msgStr,
			DeepLink: "manch://profile/" + followsTo.Id.Hex(),
			Id:       notification.Identifier,
		}
		//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token\n")
		}

		fmt.Printf("Processed a User follow on subject %s! with user follow ID %s\n", subj, uf.Id)

	})

	subscribers.PostRemovedSubscriber(func(subj, reply string, p *subscribers.Post) {
		fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
		post := mongo.GetPostById(p.Id)

		postCreator := mongo.GetProfileById(post.Created.ProfileId)

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		notification := mongo.CreateNotification(post.Id, "delete", "post", postCreator.Id)

		reason := post.IgnoreReason
		language := postCreator.Language
		deleteReason := i18n.DeleteReason[language][reason]

		postTitle := utils.TruncateTitle(post.Title, 4)
		data := i18n.DataModel{
			Name:         postCreator.Name,
			Post:         postTitle,
			DeleteReason: deleteReason,
		}
		var msgStr string
		msgStr = i18n.GetString(language, "post_removed", data)
		fmt.Println(msgStr)
		title := i18n.GetAppTitle(language)
		msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
		msg := firebase.ManchMessage{
			Title:    title,
			Message:  msgStr,
			Icon:     mongo.ExtractThumbNailFromPost(post),
			DeepLink: "manch://profile/" + postCreator.Id.Hex(),
			Id:       notification.Identifier,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
		fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)

	})

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
