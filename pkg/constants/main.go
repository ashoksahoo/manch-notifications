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
	"COMMENT":              "comment",             // comment on a post
	"REPLY":                "reply",               // reply to comments
	"MULTI_REPLY":          "multi_reply",         // Someone else also replied to a comment
	"MULTI_COMMENT":        "multi_comment",       // Someone else also comment on a post
	"POST_REMOVE":          "remove",              // post removed from feed
	"POST_SHARE":           "share",               // post shared
	"USER_FOLLOW":          "follow",              // user followed another another user
	"VOTE":                 "vote",                // vote on comment or post
	"LIVE_TOPIC_WINNER":    "live_topic_winner",   // live topic winners
	"USER_BLOCKED":         "user.blocked",        // user blocked
	"USER_UNBLOCKED":       "user.unblocked",      // user unblocked
	"USER_WARNED":          "user.warned",         // user warned
	"100_COIN_MILESTONE":   "100_coin_milestone",  // 100 coin milestone achieved
	"JOIN_MANCH_REQUEST":   "join_manch_request",  // join manch reqeust
	"JOIN_MANCH_APPROVED":  "join_manch_approved", // join manch approved
	"MANCH_100_MEMBERS":    "manch_100_members",   // manch hits 100 users
	"MANCH_ACTIVATION":     "manch_activation",    // manch activated
	"STREAK_MILESTONE":     "streak_milestone",    // streak milestone
	"100_COIN_REFERRAL":    "100_coin_referral",
	"500_COIN_MILESTONE":   "500_coin_milestone",
	"2500_COIN_MILESTONE":  "2500_coin_milestone",
	"5000_COIN_MILESTONE":  "5000_coin_milestone",
	"10000_COIN_MILESTONE": "10000_coin_milestone",
	"25000_COIN_MILESTONE": "25000_coin_milestone",
	"1_STREAK_MILESTONE":   "1_streak_milestone",
	"7_STREAK_MILESTONE":   "7_streak_milestone",
	"30_STREAK_MILESTONE":  "30_streak_milestone",
	"100_STREAK_MILESTONE": "100_streak_milestone",
	"REPOSTED_POST":        "reposted_post",
	"USER_REVIEW":          "user.review",
	"STREAK_MISSING":       "streak_missing",
	"POST_ON_MANCH":        "post_on_manch",
	"KARMA_POINTS":         "karma_points",
}

var NotificationPurposeResource = map[string]string{
	"comment":              "ic_nc_comment",  // comment on a post
	"reply":                "ic_nc_comment",  // reply to comments
	"multi_reply":          "ic_nc_comment",  // Someone else also replied to a comment
	"multi_comment":        "ic_nc_comment",  // Someone else also comment on a post
	"remove":               "ic_nc_remove",   // post removed from feed
	"share":                "ic_nc_share",    // post shared
	"follow":               "ic_nc_follow",   // user followed another another user
	"vote":                 "ic_nc_vote",     // vote on comment or post
	"live_topic_winner":    "ic_nc_trophy",   // live topic winners notification
	"join_manch_request":   "ic_nc_my_manch", // join manch request
	"join_manch_approved":  "ic_nc_my_manch", // join manch approved
	"manch_100_members":    "ic_nc_my_manch", // manch achieved 100 members
	"manch_activation":     "ic_nc_my_manch", // manch activated
	"100_coin_milestone":   "ic_earn_coin",
	"100_coin_referral":    "ic_earn_coin",
	"500_coin_milestone":   "ic_badge_contributor",
	"2500_coin_milestone":  "ic_badge_rising_star",
	"5000_coin_milestone":  "ic_milestone_super_user",
	"10000_coin_milestone": "ic_badge_super_star",
	"25000_coin_milestone": "ic_badge_elite",
	"1_streak_milestone":   "ic_milestone_1_day_steak",
	"7_streak_milestone":   "ic_milestone_7_day_steak",
	"30_streak_milestone":  "ic_milestone_30_day_steak",
	"100_streak_milestone": "ic_milestone_100_day_steak",
	"post_on_manch":        "ic_nc_my_manch",
}

var NotificationPurposeIcon = map[string]string{
	"comment":              "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_comment-min.png",     // comment on a post
	"reply":                "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_comment-min.png",     // reply to comments
	"multi_reply":          "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_comment-min.png",     // Someone else also replied to a comment
	"multi_comment":        "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_comment-min.png",     // Someone else also comment on a post
	"remove":               "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_default_mic-min.png", // post removed from feed
	"share":                "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_default_mic-min.png", // post shared
	"follow":               "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_follow-min.png",      // user followed another another user
	"vote":                 "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_like-min.png",        // vote on comment or post
	"live_topic_winner":    "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_trophy-min.png",      // live topics winners notification
	"join_manch_request":   "https://s3.ap-south-1.amazonaws.com/manch-dev/app-banners/MyManch_Notification_PNG.png",    // join manch request
	"join_manch_approved":  "https://s3.ap-south-1.amazonaws.com/manch-dev/app-banners/MyManch_Notification_PNG.png",    // join manch approved
	"manch_100_members":    "https://s3.ap-south-1.amazonaws.com/manch-dev/app-banners/MyManch_Notification_PNG.png",    // manch achieved 100 members
	"manch_activation":     "https://s3.ap-south-1.amazonaws.com/manch-dev/app-banners/MyManch_Notification_PNG.png",    // manch activated
	"500_coin_milestone":   "https://s3.ap-south-1.amazonaws.com/manch-dev/profile_levels/ic_badge_contributor.png",
	"2500_coin_milestone":  "https://s3.ap-south-1.amazonaws.com/manch-dev/profile_levels/ic_badge_rising_star.png",
	"5000_coin_milestone":  "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_super_user.png",
	"10000_coin_milestone": "https://s3.ap-south-1.amazonaws.com/manch-dev/profile_levels/ic_badge_super_star.png",
	"25000_coin_milestone": "https://s3.ap-south-1.amazonaws.com/manch-dev/profile_levels/ic_badge_elite.png",
	"1_streak_milestone":   "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_1_day_steak.png",
	"7_streak_milestone":   "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_7_day_steak.png",
	"30_streak_milestone":  "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_30_day_steak.png",
	"100_streak_milestone": "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_100_day_steak.png",
	"100_coin_milestone":   "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_default_mic-min.png",
	"100_coin_referral":    "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_default_mic-min.png",
	"reposted_post":        "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/repost_icon.png",
	"streak_missing":       "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/streak_miss_notification.png",
	"post_on_manch":        "https://s3.ap-south-1.amazonaws.com/manch-dev/app-banners/MyManch_Notification_PNG.png",
	"karma_points":         "https://s3.ap-south-1.amazonaws.com/manch-dev/app-resource-icons/ic_n_default_mic-min.png",
}

var NotificationTemplate = map[string]string{
	"TRANSACTIONAL": "transactional",
	"PROMOTIONAL":   "promotional",
}

var ModelNames = map[string]string{
	"USERS":                     "users",
	"POSTS":                     "posts",
	"COMMUNITIES":               "communities",
	"COMMENTS":                  "comments",
	"VOTES":                     "votes",
	"RESOURCE_SCORE":            "RESOURCE_SCORE",
	"USER_FOLLOWS":              "user_follows",
	"SHARES":                    "shares",
	"FCM_TOKENS":                "fcm_tokens",
	"EVENTS":                    "events",
	"EXPLORABLE_ENTITIES":       "explorable_entities",
	"MEDIA_SOURCES":             "media_sources",
	"USER_SCORE":                "user_score",
	"COMMENT_STRINGS":           "comment_strings",
	"COMMENT_SCHEDULEDS":        "comment_scheduleds",
	"NOTIFICATION_V2":           "notificationsv3",
	"USER_FOLLOWS_SCHEDULEDS":   "user_follows_scheduleds",
	"VOTE_SCHEDULEDS":           "vote_scheduleds",
	"SHARES_SCHEDULEDS":         "shares_scheduleds",
	"USER_SCORES":               "user_scores",
	"USER_COINS":                "user_coins",
	"USER_LEADERBOARDS":         "user_leaderboards",
	"WHATSAPP_SCHEDULEDS":       "whatsapp_scheduleds",
	"USER_COINS_SCHEDULE_MODEL": "user_coins_scheduleds",
	"COMMUNITIES_STATS_MODEL":   "community_stats",
	"REFERRALS":                 "referrals",
	"USER_TAGS":                 "user_tags",
}

var BlackListStatus = map[string]string{
	"WARNING":    "warning",
	"BLOCKED":    "blocked",
	"UN_BLOCKED": "unblocked",
}

var MileStones = map[string]string{
	"100_COIN_MILESTONE":   "100",
	"500_COIN_MILESTONE":   "500",
	"1000_COIN_MILESTONE":  "1000",
	"5000_COIN_MILESTONE":  "5000",
	"10000_COIN_MILESTONE": "10000",
	"25000_COIN_MILESTONE": "25000",
}

var CommunityStatus = map[string]string{
	"CREATED":          "created",
	"PENDING_APPROVAL": "pending_approval",
	"APPROVED":         "approved",
	"ACTIVATED":        "activated",
}
