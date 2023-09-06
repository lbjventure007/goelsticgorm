package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	searchs "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gogormlearn/inits"
	"gogormlearn/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"strings"
)

func main() {
	inits.InitElastic()
	dsn := "root:1234qwer@tcp(127.0.0.1:3306)/seata_client?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db.Use(sharding.Register(sharding.Config{
	//	ShardingKey:         "user_id",
	//	NumberOfShards:      64,
	//	PrimaryKeyGenerator: sharding.PKSnowflake,
	//}, "orders").Register(sharding.Config{
	//	ShardingKey:         "user_id",
	//	NumberOfShards:      256,
	//	PrimaryKeyGenerator: sharding.PKSnowflake,
	//	// This case for show up give notifications, audit_logs table use same sharding rule.
	//}, Notification{}, AuditLog{}))
	db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "user_id",
		NumberOfShards:      2,
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, "orders"))

	//db.Create(model.Orders{UserID: 1, ProductID: 1})
	//db.Create(model.Orders{UserID: 2, ProductID: 1})

	orders := model.Orders{}
	db.Where("id", int64(1698268938120466432)).Find(&orders)
	fmt.Println(orders)

	//create, err := inits.Elastic.Indices.Create("myindex")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(create)
	//document := struct {
	//	Name string `json:"name"`
	//}{"test"}
	//data, _ := json.Marshal(document)
	//index, err := inits.Elastic.Index("myindex", bytes.NewReader(data))
	//if err != nil {
	//	return
	//}
	//fmt.Println(index)

	//document1 := struct {
	//	Name string `json:"name"`
	//}{"test1"}

	//do, err := inits.Elastic.Index("myindex").Id("2").Request(document1).Do(context.TODO())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(do)

	//response, err := inits.Elastic.Mget().Do(context.TODO())

	//json, _ := response.Source_.MarshalJSON()
	//fmt.Println(string(json), response.Fields, response.Id_, response.Index_)
	//fmt.Println(response)

	//inits.Elastic.Search().
	//	Index("myindex").
	//	Request(&search.Request{
	//		Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
	//	}).
	//	Do(context.TODO())

	search, err := inits.Elastic.Search(func(request *esapi.SearchRequest) {
		request.Index = append(request.Index, "myindex")

	})

	fmt.Println(search, err)

	source, err := inits.Elastic.GetSource("myindex", "2")
	fmt.Println(source, err)

	//inits.Elastic.GetSource
	//response, err := inits.Elastic.Search(func(request *esapi.SearchRequest) {
	//	request.Index = []string{"pinyintest"}
	//
	//	//query := `{"query": {"match": {"info": "ldh"}}}`
	//
	//	request.Query = "2"
	//
	//})

	//querys := `{"query": {"match": {"info": "刘德华"}}}`

	//query := `{"query": {"match": {"info": "ldh"}}}`
	//r, err := inits.Elastic.Search(
	//	inits.Elastic.Search.WithIndex("pinyintest"),
	//	inits.Elastic.Search.WithBody(strings.NewReader(query)),
	//	//inits.Elastic.Search.WithSource("info"),
	//	inits.Elastic.Search.WithHuman(),
	//	inits.Elastic.Search.WithSize(1),
	//	inits.Elastic.Search.WithFrom(0),
	//	inits.Elastic.Search.WithSourceIncludes("info,username"),
	//)

	var b = true
	do, err := inits.ElasticType.Search().Index("pinyintest").
		//Raw(strings.NewReader(query)).
		Request(&searchs.Request{
			Highlight: &types.Highlight{
				PreTags: []string{"<font color='red'>"},

				Fields: map[string]types.HighlightField{

					"info": {
						//PreTags:       []string{"<font color='red'>"},
						//PostTags:      []string{"</font>"},
						//MatchedFields: []string{"info"},
					},
				},

				PostTags:          []string{"</font>"},
				RequireFieldMatch: &b,
			},

			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"info": {
						Query: "ldh",
					},
				},
			},
		}).Do(context.TODO())

	if len(do.Hits.Hits) > 0 {
		h := do.Hits.Hits[0]
		fmt.Println(111111, string(h.Source_), h.Highlight["info"][0])
	}

	querys := `{
				"query": 
					{
						"match": 
						{
							"info": "ldh"	
						}
					},
				"highlight":{
                    "pre_tags": "<b class='key' style='color:red'>",
					"post_tags": "</b>",
					"fields": {
					  "info": {}
					}
				}
				}`
	response, _ := inits.Elastic.Search(func(request *esapi.SearchRequest) {
		request.Index = []string{"pinyintest"}
		request.Body = strings.NewReader(querys)
	})
	fmt.Println(err, "222222", response)

	//自动补全查询
	q := `{
  "suggest": {
    "info_suggest": {
      "text": "l w",
      "completion":{
        "field":"info",
        "skip_duplicates":false,
        "size": 10
  
      }
      
    }
  }
}`
	r, err := inits.Elastic.Search(func(request *esapi.SearchRequest) {
		request.Body = strings.NewReader(q)
		request.Pretty = true

	})

	fmt.Println(33333, r, err)
	s := "lw"
	q1 := `
    {
  "suggest": {
    "info_suggest": {
      "text": "l d",
      "completion":{
        "field":"info",
        "skip_duplicates":false,
        "size": 10
  
      }
      
    }
  }
}
`
	bb := true
	sizel := 10
	completeion := &types.CompletionSuggester{}
	completeion.Field = "info"
	completeion.SkipDuplicates = &bb
	completeion.Size = &sizel
	r2, err := inits.ElasticType.Search().Suggest(&types.Suggester{
		Text:       &s,
		Suggesters: map[string]types.FieldSuggester{
			//	Completion:completeion,
			//	Text:
		},
	}).Index("test5").Do(context.TODO())
	fmt.Println(4444, r2, err)

	r3, err := inits.ElasticType.Search().Index("test5").
		Raw(strings.NewReader(q1)).Do(context.TODO())
	suggest1 := r3.Suggest["info_suggest"][0]

	fmt.Println(55555, err)

	suggest2, ok := suggest1.(*types.CompletionSuggest) //类型断言

	//tests, ok1 := suggest1.(*Tests)
	//fmt.Println("----", tests, ok1, "----")
	if ok {
		completionSuggestOptions := suggest2.Options
		if len(completionSuggestOptions) > 0 {
			for _, options := range completionSuggestOptions {
				source_ := options.Source_
				json, _ := source_.MarshalJSON()
				fmt.Printf("%d, %s ", 1111, json)
			}
		}
	}

}
