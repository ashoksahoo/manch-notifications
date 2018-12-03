package utils

import (
	"strings"
)

func TruncateTitle(title string, no_of_words int) string{
	var postTitle string
	titleWords := strings.Split(title, " ")
	if len(titleWords) > no_of_words {
		postTitle = strings.Join(titleWords[:no_of_words], " ")
		postTitle = postTitle + "..."
	} else {
		postTitle = title
	}
	return postTitle
}