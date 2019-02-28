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
	"COMMENT":           "comment",           // comment on a post
	"REPLY":             "reply",             // reply to comments
	"MULTI_REPLY":       "multi_reply",       // Someone else also replied to a comment
	"MULTI_COMMENT":     "multi_comment",     // Someone else also comment on a post
	"POST_REMOVE":       "remove",            // post removed from feed
	"POST_SHARE":        "share",             // post shared
	"USER_FOLLOW":       "follow",            // user followed another another user
	"VOTE":              "vote",              // vote on comment or post
	"LIVE_TOPIC_WINNER": "live_topic_winner", // live topic winners
	"USER_BLOCKED":      "user.blocked",      // user blocked
	"USER_UNBLOCKED":    "user.unblocked",    // user unblocked
	"USER_WARNED":       "user.warned",       // user warned
}

var NotificationPurposeResource = map[string]string{
	"comment":       "ic_nc_comment",       // comment on a post
	"reply":         "ic_nc_reply",         // reply to comments
	"multi_reply":   "ic_nc_multi_reply",   // Someone else also replied to a comment
	"multi_comment": "ic_nc_multi_comment", // Someone else also comment on a post
	"remove":        "ic_nc_remove",        // post removed from feed
	"share":         "ic_nc_share",         // post shared
	"follow":        "ic_nc_follow",        // user followed another another user
	"vote":          "ic_nc_vote",          // vote on comment or post
}

var NotificationPurposeIcon = map[string]string{
	"comment":       "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // comment on a post
	"reply":         "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // reply to comments
	"multi_reply":   "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // Someone else also replied to a comment
	"multi_comment": "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // Someone else also comment on a post
	"remove":        "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // post removed from feed
	"share":         "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // post shared
	"follow":        "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // user followed another another user
	"vote":          "https://cdn2.iconfinder.com/data/icons/instagram-ui/48/jee-68-512.png", // vote on comment or post
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
	"VOTES":                   "votes",
	"RESOURCE_SCORE":          "RESOURCE_SCORE",
	"USER_FOLLOWS":            "user_follows",
	"SHARES":                  "shares",
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
	"USER_COINS":              "user_coins",
	"USER_LEADERBOARDS":       "user_leaderboards",
}

var BlackListStatus = map[string]string{
	"WARNING":    "warning",
	"BLOCKED":    "blocked",
	"UN_BLOCKED": "unblocked",
}
