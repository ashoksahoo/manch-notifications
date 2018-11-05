package subscribers

import (
	"github.com/nats-io/go-nats"
)

type Post struct {
	Id    string `json:"_id"`
	GUID  string `json:"guid"`
	Title string `json:"title"`
}

type Comment struct {
	Id    string `json:"_id"`
	GUID  string `json:"guid"`
	Title string `json:"title"`
}

func PostSubscriber(c *nats.EncodedConn, callback func(subj, reply string, m *Post)) {
	c.QueueSubscribe("manch-api:development:post", "manch-notification-service:development", callback)
}

func CommentSubscriber(c *nats.EncodedConn, callback func(subj, reply string, m *Comment)) {
	c.QueueSubscribe("manch-api:development:comment", "manch-notification-service:development", callback)
}
