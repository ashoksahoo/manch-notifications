package constants

const (
	NA                        = "NA"
	PENDING                   = "pending"
	SENT                      = "sent"
	DELIVERED                 = "delivered"
	READ                      = "read"
	FAILED                    = "failed"
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
	"PROMOTIONAL": "promotional",
}
