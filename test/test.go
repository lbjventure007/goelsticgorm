package test

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	searchs "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/go-zookeeper/zk"
	"gogormlearn/common"
	"gogormlearn/inits"
	"gogormlearn/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"strings"
	"time"
)

type X struct {
	i int
}

func Test() {

	//zookeeper 分布式锁
	inits.InitZk()
	l := zk.NewLock(inits.ZkConn, "/test", zk.WorldACL(zk.PermAll))

	x := &X{10}

	go func(h *X) {

		err2 := l.Lock()
		fmt.Println("111")
		if err2 != nil {
			panic(err2)
		}

		i := h.i - 1

		fmt.Println(i)
		time.Sleep(2 * time.Second)
		l.Unlock()
	}(x)

	fmt.Println("222")
	time.Sleep(5 * time.Second)
	return
	//

	inits.InitRedis()
	//
	add := inits.RedisClient.BFAdd(context.TODO(), "bfloom", "a")
	if add.Err() != nil {
		fmt.Println(add.Err().Error())
	} else {

		fmt.Println(add.Val())
	}
	val := inits.RedisClient.BFExists(context.TODO(), "bfloom", "a").Val()
	fmt.Println(val)
	return
	//结合redis 和 lua 模拟解决超卖问题

	//redis单线程  把script当作一个整体 在redis单线程里面操作

	//scrip := `local value = redis.call("Get", KEYS[1])
	//	print("当前值为 " .. value);
	//	if( value - KEYS[2] >= 0 ) then
	//		local leftStock = redis.call("DecrBy" , KEYS[1],KEYS[2])
	//		print("剩余值为" .. leftStock );
	//		return leftStock
	//	else
	//		print("数量不够，无法扣减");
	//		return value - KEYS[2]
	//	end
	//	return -1`
	script := inits.Script(`
		local value = redis.call("Get", KEYS[1])
		
		if( value - KEYS[2] >= 0 ) then
			local leftStock = redis.call("DecrBy" , KEYS[1],KEYS[2])
			
			return leftStock
		else
			
			return value - KEYS[2]
		end
		return -1`)

	schan := make(chan int, 100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			keys := make([]string, 0)
			keys = append(keys, "stock") //要查询的key
			keys = append(keys, "1")     //要减少的库存

			run := script.Run(context.TODO(), inits.RedisClient, keys)
			if run.Err() != nil {
				fmt.Println(run.Err().Error())
			} else {
				result, _ := run.Result()
				i2 := result.(int64)
				if i2 >= 0 {
					inits.RedisClient.LPush(context.TODO(), "succ", i)
					fmt.Println(i, "操作成功")
				} else {
					fmt.Println(i, "操作失败")
				}

			}
			schan <- i
		}(i)
	}

	for {
		select {
		case <-schan:

		case <-time.After(10 * time.Second):
			goto loop
		}
	}

loop:
	fmt.Println(inits.RedisClient.LRange(context.TODO(), "succ", 0, -1))
	fmt.Println("over")
	return

	//一致性hash

	cHashRing := common.NewConsistent()

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		cHashRing.Add(common.NewNode(i, "172.18.1."+si, 8080, "host_"+si, 1))
	}

	for k, v := range cHashRing.Nodes {
		fmt.Println("Hash:", k, " IP:", v.Ip)
	}

	ipMap := make(map[string]int, 0)
	for i := 0; i < 1000; i++ {
		si := fmt.Sprintf("key%d", i)
		k := cHashRing.Get(si)
		if _, ok := ipMap[k.Ip]; ok {
			ipMap[k.Ip] += 1
		} else {
			ipMap[k.Ip] = 1
		}
	}

	for k, v := range ipMap {
		fmt.Println("Node IP:", k, " count:", v)
	}

	return

	//end

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
