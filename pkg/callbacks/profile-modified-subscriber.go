package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func ProfileModifiedCB(subj, reply string, updatedProfile *subscribers.Profile) {
	fmt.Printf("Received a User Profile Update on subject %s! with User %+v\n", subj, updatedProfile)
	isUpdated := false
	update := bson.M{}
	query := bson.M{"created.profile_id": bson.ObjectIdHex(updatedProfile.Id), "anonymous": bson.M{"$ne": true}}
	profile := mongo.GetProfileById(bson.ObjectIdHex(updatedProfile.Id))
	if updatedProfile.Name != "" {
		isUpdated = true
		update["created.name"] = updatedProfile.Name
	}
	if updatedProfile.Avatar != "" {
		isUpdated = true
		update["created.avatar"] = updatedProfile.Avatar
	}

	if isUpdated && updatedProfile.DisplayProfileChangedUpdatedAt == profile.DisplayProfileChangedUpdatedAt {
		// update post
		mongo.UpdateAllPostsByQuery(query, bson.M{"$set": update})
		// update comment
		mongo.UpdateAllCommentsByQuery(query, bson.M{"$set": update})
		// update profile
		mongo.UpdateProfileById(profile.Id,
			bson.M{
				"$set": bson.M{
					"profiles.$.display_profile_changed_updated": true,
				},
			})
	}

	fmt.Printf("Processed a User Profile Update on subject %s! with User Id %s\n", subj, updatedProfile.Id)
}
