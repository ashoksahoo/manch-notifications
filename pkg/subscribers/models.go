package subscribers

import (
	"notification-service/pkg/mongo"
	"time"
)

type Post struct {
	Id           string               `json:"_id"`
	GUID         string               `json:"guid"`
	Title        string               `json:"title"`
	Type         string               `json:"type"`
	Community    []string             `json:"community_ids"`
	Language     string               `json:"language"`
	New          bool                 `json:"isNew"`
	IsBot        bool                 `json:"is_bot"`
	CreatorType  string               `json:"creator_type"`
	Tags         []string             `json:"tags"`
	TagsPosition []mongo.TagPositions `json:"tag_positions"`
}

type Comment struct {
	Id      string `json:"_id"`
	GUID    string `json:"guid"`
	Content string `json:"content"`
	Type    string `json:"comment_type"`
	New     bool   `json:"isNew"`
}

type User struct {
	Id    string   `json:"_id"`
	GUID  string   `json:"guid"`
	Name  string   `json:"name"`
	Phone string   `json:"phone"`
	Roles []string `json:"roles"`
	New   bool     `json:"isNew"`
}

type Vote struct {
	Id           string `json:"_id"`
	GUID         string `json:"guid"`
	Resource     string `json:"resource"`
	ResourceType string `json:"resource_type"`
	Direction    string `json:"vote"`
	New          bool   `json:"isNew"`
}

type Subscription struct {
	Id           string `json:"_id"`
	GUID         string `json:"guid"`
	Resource     string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	ProfileId    string `json:"profile_id"`
	New          bool   `json:"isNew"`
}

type Share struct {
	Id           string `json:"_id"`
	GUID         string `json:"guid"`
	ResourceId   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	ProfileId    string `json:"profile_id"`
}

type SharePost struct {
	Id         string `json:"_id"`
	GUID       string `json:"guid"`
	ProfileId  string `json:"profile_id"`
	ShareCount int    `json:"no_of_shares"`
}

type LiveTopicComment struct {
	Id         string `json:"_id"`
	PostId     string `json:"post_id"`
	CommentId  string `json:"comment_id"`
	CreatedBy  string `json:"created_by"`
	IsReply    bool   `json:"is_reply"`
	UpVotes    int    `json:"up_votes"`
	ReplyCount int    `json:"reply_count"`
}

type LiveTopicsWinner struct {
	Id           string   `json:"_id"`
	Title        string   `json:"title"`
	BannerImage  string   `json:"banner_image"`
	Winners      []string `json:"winners"`
	Participants []string `json:"participants"`
}

type BlackListProfile struct {
	ProfileId    string    `json:"profile_id"`
	Status       string    `json:"status"`
	BlockedOn    time.Time `json:"blocked_on"`
	BlockedTill  time.Time `json:"blocked_till"`
	Reason       string    `json:"reason"`
	LastWarnedOn time.Time `json:"last_warned_on"`
}

type MileStone struct {
	ProfileId string `json:"profile_id"`
	MileStone string `json:"milestone"`
}

type Community struct {
	Id             string `json:"_id"`
	GUID           string `json:"guid"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	Icon           string `json:"icon"`
	Language       string `json:"language"`
	Status         string `json:"status"`
	Visibility     string `json:"visibility"`
	FollowersCount int    `json:"no_of_followers"`
}

type LiveTopicPoll struct {
	TopicId        string `json:"topic_id"`
	ParticipantId  string `json:"participant_id"`
	UserAnswerId   int    `json:"user_answer_id"`
	UserAnswerText string `json:"user_answer_text"`
	ResultId       int    `json:"result_id"`
	CoinsEarned    int    `json:"coins_earned"`
	CoinsLost      int    `json:"coins_lost"`
}

type Streak struct {
	StartDate    time.Time `json:"start_date" bson:"start_date"`
	EndDate      time.Time `json:"end_date" bson:"end_date"`
	StreakLength int       `json:"streak_length" bson:"streak_length"`
}

type UserStreak struct {
	ProfileId     string `json:"profile_id" bson:"profile_id"`
	CurrentStreak Streak `json:"current_streak" bson:"current_streak"`
	LastStreak    Streak `json:"last_streak" bson:"last_streak"`
	LongestStreak Streak `json:"longest_streak" bson:"longest_streak"`
}

type Profile struct {
	Id                             string    `json:"_id"`
	Name                           string    `json:"name"`
	Avatar                         string    `json:"avatar"`
	DisplayProfileChangedUpdated   bool      `json:"display_profile_changed_updated"`
	DisplayProfileChangedUpdatedAt time.Time `json:"display_profile_changed_at"`
}
