package callbacks

import (
	"time"
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

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
		update["create.name"] = updatedProfile.Name
	}
	if updatedProfile.Avatar != "" {
		isUpdated = true
		update["created.avatar"] = updatedProfile.Avatar
	}
	if isUpdated && updatedProfile.DisplayProfileUpdatedAt == profile.DisplayProfileUpdatedAt{
		// update post
		mongo.UpdatePostByItr(query, bson.M{"$set": update})
		// update comment
		mongo.UpdateCommentByItr(query, bson.M{"$set": update})
		// update profile
		mongo.UpdateProfileById(profile.Id,
			bson.M{
				"$set": bson.M{
					"profiles.$.display_profile_changes": false,
					"profiles.$.display_profile_updated_at": time.Now(),
				},
			})
	}

	fmt.Printf("Processed a New User on subject %s! with User Id %s\n", subj, updatedProfile.Id)
}
