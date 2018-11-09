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
	return "manch-api:" + env + s

}

func queue() string {
	env := os.Getenv("env")
	if env == "" {
		env = "development"
	}
	return "manch-notification-service:" + env
}
func PostSubscriber(callback func(subj, reply string, m *Post)) {
	c.QueueSubscribe(subject("post"), queue(), callback)
}

func CommentSubscriber(callback func(subj, reply string, m *Comment)) {
	c.QueueSubscribe(subject("comment"), queue(), callback)
}

func UserSubscriber(callback func(subj, reply string, m *User)) {
	c.QueueSubscribe(subject("user"), queue(), callback)
}

func VoteSubscriber(callback func(subj, reply string, m *User)) {
	c.QueueSubscribe(subject("vote"), queue(), callback)
}

func SubsSubscriber(callback func(subj, reply string, m *Subscription)) {
	c.QueueSubscribe(subject("sub"), queue(), callback)
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
