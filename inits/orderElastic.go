package inits

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"log"
)

var Elastic *elasticsearch.Client
var ElasticType *elasticsearch.TypedClient

func InitElastic() {
	Elastic, _ = elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(Elastic.Info())

	ElasticType, _ = elasticsearch.NewTypedClient(elasticsearch.Config{})
	log.Println(elasticsearch.Version)
	log.Println(Elastic.Info())
}
