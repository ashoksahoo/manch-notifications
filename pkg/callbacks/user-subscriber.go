package callbacks

import (
	"fmt"
	"notification-service/pkg/elasticsearch"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"
)

func UserSubscriberCB(subj, reply string, u *subscribers.User) {
	fmt.Printf("Received a New User on subject %s! with User %+v\n", subj, u)

	// get user from db
	user := mongo.GetUserById(u.Id)

	profile := user.Profiles[0]
	// index user
	currentISOTime := utils.ISOFormat(time.Now())
	elasticsearch.CreateUserIndex(elasticsearch.UserIndex{
		ID:   profile.Id.Hex(),
		Name: profile.Name,
		NameKeyword: elasticsearch.TypeInput{
			Input: []string{profile.Name},
		},
		Avatar:             profile.Avatar,
		AboutMe:            profile.AboutMe,
		Type:               profile.Type,
		NoOfPosts:          profile.NoOfPosts,
		NoOfLikes:          profile.NoOfLikes,
		NoOfComments:       profile.CommentsCount,
		NoOfShares:         profile.NoOfShares,
		NoOfFollowers:      profile.FollowerCount,
		NoOfFollowing:      profile.NoOfFollowing,
		NoOfManchFollowing: profile.NoOfManchFollowing,
		LastActiveHour:     utils.ISOFormat(profile.LastActiveHour),
		TotalCoins:         profile.TotalCoins,
		TotalManchCreated:  profile.TotalManchCreated,
		BranchLink:         profile.BranchLink,
		CreatedAt:          currentISOTime,
		UpdatedAt:          currentISOTime,
	})
	fmt.Printf("Processed a New User on subject %s! with User Id %s\n", subj, u.Id)
}
