package subscribers

type Post struct {
	Id          string   `json:"_id"`
	GUID        string   `json:"guid"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Community   []string `json:"community_ids"`
	Language    string   `json:"language"`
	New         bool     `json:"isNew"`
	IsBot       bool     `json:"is_bot"`
	CreatorType string   `json:"creator_type"`
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
	Id string `json:"_id"`
	GUID string `json:"guid"`
	Resource string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	ProfileId string `json:"profile_id"`
}