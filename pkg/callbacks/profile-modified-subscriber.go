package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func ProfileModifiedCB(subj, reply string, profile *subscribers.Profile) {
	fmt.Printf("Received a New User on subject %s! with User %+v\n", subj, profile)
	isUpdated := false
	update := bson.M{}
	query := bson.M{"created.profile_id": bson.ObjectIdHex(profile.Id)}
	if profile.Name != "" {
		isUpdated = true
		update["name"] = profile.Name
	}
	if profile.Avatar != "" {
		isUpdated = true
		update["avatar"] = profile.Avatar
	}
	if isUpdated {
		// update post
		mongo.UpdatePostByItr(query, update)
		// update comment
		mongo.UpdateCommentByItr(query, update)
	}

	fmt.Printf("Processed a New User on subject %s! with User Id %s\n", subj, profile.Id)
}
