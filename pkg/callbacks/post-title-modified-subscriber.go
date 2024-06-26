package callbacks

import (
	"fmt"
	"notification-service/pkg/elasticsearch"
	"notification-service/pkg/subscribers"
)

func PostTitleModifiedCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post title modified on subject %s! with Post %+v\n", subj, p)
	var additionalScore int
	if p.CreatorType == "bot" {
		additionalScore = 50 * 60
	} else {
		additionalScore = 5 * 60
	}
	elasticsearch.AddTagToIndex(p.Tags, additionalScore, p.TagsPosition)

	fmt.Printf("Processed a post title modified on subject %s! with Post Id%s\n", subj, p.Id)
}
