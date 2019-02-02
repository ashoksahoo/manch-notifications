package constants

const (
	MANCH_OFFICIAL_PROFILE_TE = "5c3c3bfd89ac4a794d45b14d"
	MANCH_OFFICIAL_PROFILE_HE = "5c1c92c8eda9bd1771bcf0a7"
)

var NotificationStatus = map[string]string{
	"NA":        "na",
	"PENDING":   "pending",
	"SENT":      "sent",
	"DELIVERED": "delivered",
	"READ":      "read",
	"FAILED":    "failed",
}

var NotificationPurpose = map[string]string{
	"COMMENT":       "comment",       // comment on a post
	"REPLY":         "reply",         // reply to comments
	"MULTI_REPLY":   "multi_reply",   // Someone else also replied to a comment
	"MULTI_COMMENT": "multi_comment", // Someone else also comment on a post
	"POST_REMOVE":   "remove",        // post removed from feed
	"POST_SHARE":    "share",         // post shared
	"USER_FOLLOW":   "follow",        // user followed another another user
	"VOTE":          "vote",          // vote on comment or post
}

var NotificationTemplate = map[string]string{
	"TRANSACTIONAL": "transactional",
	"PROMOTIONAL":   "promotional",
}

var ModelNames = map[string]string{
	"USERS":                   "users",
	"POSTS":                   "posts",
	"COMMUNITIES":             "communities",
	"COMMENTS":                "comments",
	"VOTES":                    "votes",
	"RESOURCE_SCORE":          "RESOURCE_SCORE",
	"USER_FOLLOWS":            "user_follows",
	"SHARES":                  "SHARES",
	"FCM_TOKENS":              "fcm_tokens",
	"EVENTS":                  "events",
	"EXPLORABLE_ENTITIES":     "explorable_entities",
	"MEDIA_SOURCES":           "media_sources",
	"USER_SCORE":              "user_score",
	"COMMENT_STRINGS":         "comment_strings",
	"COMMENT_SCHEDULEDS":      "comment_scheduleds",
	"NOTIFICATION_V2":         "notificationsv2",
	"USER_FOLLOWS_SCHEDULEDS": "user_follows_scheduleds",
	"VOTE_SCHEDULEDS":         "vote_scheduleds",
	"SHARES_SCHEDULEDS":       "shares_scheduleds",
	"USER_SCORES":             "user_scores",
}
