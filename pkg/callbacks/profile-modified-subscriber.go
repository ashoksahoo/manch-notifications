package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"

	"github.com/globalsign/mgo/bson"
)

func ProfileModifiedCB(subj, reply string, updatedProfile *subscribers.Profile) {
	fmt.Printf("Received a New User on subject %s! with User %+v\n", subj, updatedProfile)
	isUpdated := false
	update := bson.M{}
	query := bson.M{"created.profile_id": bson.ObjectIdHex(updatedProfile.Id)}
	profile := mongo.GetProfileById(bson.ObjectIdHex(updatedProfile.Id))
	if updatedProfile.Name != "" {
		isUpdated = true
		update["created.name"] = updatedProfile.Name
	}
	if updatedProfile.Avatar != "" {
		isUpdated = true
		update["created.avatar"] = updatedProfile.Avatar
	}
	fmt.Println("requested displayprofilechange", updatedProfile.DisplayProfileChangedUpdatedAt)
	fmt.Println("requested displayprofilechange", profile.DisplayProfileChangedUpdatedAt)
	if isUpdated && updatedProfile.DisplayProfileChangedUpdatedAt == profile.DisplayProfileChangedUpdatedAt {
		// update post
		mongo.UpdatePostByItr(query, bson.M{"$set": update})
		// update comment
		mongo.UpdateCommentByItr(query, bson.M{"$set": update})
		// update profile
		mongo.UpdateProfileById(profile.Id,
			bson.M{
				"$set": bson.M{
					"profiles.$.display_profile_changed_updated":    false,
					"profiles.$.display_profile_changed_updated_at": time.Now(),
				},
			})
	}

	fmt.Printf("Processed a New User on subject %s! with User Id %s\n", subj, updatedProfile.Id)
}
