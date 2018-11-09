package subscribers

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"log"
	"os"
)

var c *nats.EncodedConn
var nc *nats.Conn
var err error

type Post struct {
	Id        string   `json:"_id"`
	GUID      string   `json:"guid"`
	Title     string   `json:"title"`
	Type      string   `json:"type"`
	Community []string `json:"community_ids"`
	Language  string   `json:"language"`
	New       bool     `json:"isNew"`
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
	New          bool   `json:"isNew"`
}

type Subscription struct {
	Id           string `json:"_id"`
	GUID         string `json:"guid"`
	Resource     string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	New          bool   `json:"isNew"`
}

func subject(s string) string {
	env := os.Getenv("env")
	if env == "" {
		env = "development"
	}
	s = "manch-api." + env + "." + s
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

func CommentSubscriber(callback func(subj, reply string, m *Comment)) {
	go c.QueueSubscribe(subject("comment.new"), queue(), callback)
}

func CommentUpdateSubscriber(callback func(subj, reply string, m *Comment)) {
	go c.QueueSubscribe(subject("comment.update"), queue(), callback)
}

func UserSubscriber(callback func(subj, reply string, m *User)) {
	go c.QueueSubscribe(subject("user"), queue(), callback)
}

func VoteSubscriber(callback func(subj, reply string, m *User)) {
	go c.QueueSubscribe(subject("vote"), queue(), callback)
}

func SubsSubscriber(callback func(subj, reply string, m *Subscription)) {
	go c.QueueSubscribe(subject("sub"), queue(), callback)
}

func init() {

	if nc, err = nats.Connect(nats.DefaultURL); err != nil {
		log.Fatal(err)
	}
	if c, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Initialized NATS Connection.")
	}

}
