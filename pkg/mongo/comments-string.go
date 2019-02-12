package mongo

import (
	"notification-service/pkg/constants"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

var (
	COMMENT_STRINGS_MODEL = constants.ModelNames["COMMENT_STRINGS"]
)

type CommentString struct {
	Id               bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	ProfileId        bson.ObjectId `json:"profile_id" bson:"profile_id"`
	CommentStringIds []string      `json:"comment_ids" bson:"comment_ids"`
}

func GetCommentStringsByProfileId(profileId bson.ObjectId) (error, CommentString) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENT_STRINGS_MODEL)
	result := CommentString{}
	err := C.Find(bson.M{"profile_id": profileId}).One(&result)
	return err, result
}

func AddCommentStringToProfileId(profileId bson.ObjectId, commentId string) {
	fmt.Println("updating")
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(COMMENT_STRINGS_MODEL)
	count, _ := C.Find(bson.M{"profile_id": profileId}).Count()
	if count > 0 {
		C.Update(bson.M{"profile_id": profileId}, bson.M{
			"$addToSet": bson.M{"comment_ids": commentId},
		})
	} else {
		C.Insert(bson.M{
			"profile_id": profileId,
			"comment_ids": []string{commentId},
		})
	}
}