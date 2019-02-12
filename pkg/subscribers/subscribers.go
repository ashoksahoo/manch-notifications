package subscribers

import (
	"fmt"
	"log"
	"os"

	nats "github.com/nats-io/go-nats"
)

var c *nats.EncodedConn
var nc *nats.Conn
var err error
var url = os.Getenv("NATS")

func subject(s string) string {
	env := os.Getenv("env")
	if env == "" {
		env = "development"
	}
	s = "manch-api." + env + "." + s
	fmt.Println("listening for ", s)
	return s
}

func queue() string {
	env := os.Getenv("env")
	if env == "" {
		env = "development"
	}
	q := "manch-notification-service." + env
	return q
}

func PostSubscriber(callback func(subj, reply string, m *Post)) {
	go c.QueueSubscribe(subject("post.new"), queue(), callback)
}

func PostUpdateSubscriber(callback func(subj, reply string, m *Post)) {
	go c.QueueSubscribe(subject("post.update"), queue(), callback)
}

func PostRemovedSubscriber(callback func(subj, reply string, m *Post)) {
	go c.QueueSubscribe(subject("post.ignore-from-feed"), queue(), callback)
}

func RepostSubscriber(callback func(subj, reply string, m *Post)) {
	go c.QueueSubscribe(subject("post.re-post"), queue(), callback)
}

func PostModeratedSubscriber(callback func(subj, reply string, m *Post)) {
	go c.QueueSubscribe(subject("post.moderated"), queue(), callback)
}

func PostDeletedSubscriber(callback func(subj, reply string, m *Post)) {
	go c.QueueSubscribe(subject("post.deleted"), queue(), callback)	
}

func CommentSubscriber(callback func(subj, reply string, m *Comment)) {
	go c.QueueSubscribe(subject("comment.new"), queue(), callback)
}

func CommentUpdateSubscriber(callback func(subj, reply string, m *Comment)) {
	go c.QueueSubscribe(subject("comment.update"), queue(), callback)
}

func UserSubscriber(callback func(subj, reply string, m *User)) {
	go c.QueueSubscribe(subject("user.new"), queue(), callback)
}

func VotePostSubscriber(callback func(subj, reply string, m *Vote)) {
	go c.QueueSubscribe(subject("vote.post"), queue(), callback)
}

func VoteCommentSubscriber(callback func(subj, reply string, m *Vote)) {
	go c.QueueSubscribe(subject("vote.comment"), queue(), callback)
}

func UserFollowSubscriber(callback func(subj, reply string, m *Subscription)) {
	go c.QueueSubscribe(subject("user.follow"), queue(), callback)
}

func UserFollowRemovedSubscriber(callback func(subj, reply string, m *Subscription)) {
	go c.QueueSubscribe(subject("user.follow.removed"), queue(), callback)
}

func SharePostSubscriber(callback func(subj, reply string, m *SharePost)) {
	go c.QueueSubscribe(subject("share.post"), queue(), callback)
}

func ShareSubscriber(callback func(subj, reply string, m *Share))  {
	go c.QueueSubscribe(subject("share"), queue(), callback)
}

func init() {
	if url == "" {
		url = nats.DefaultURL
	}
	if nc, err = nats.Connect(url); err != nil {
		log.Fatal(err)
	}
	if c, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Initialized NATS Connection to " + url)
	}

}
