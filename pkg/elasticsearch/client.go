package elasticsearch

import (
	"log"
	"os"
	"github.com/elastic/go-elasticsearch/v7"
)

var url = os.Getenv("ELASTICSEARCH_URL")

var es *elasticsearch.Client

func init() {
	var err error
	if url == "" {
		url = "http://localhost:9200/"
	}
	es, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Println("Error While connecting elastic search", err)
	} else {
		log.Println("Initialized elasticsearch at", url);
	}
}

