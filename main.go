package main

import (
	"gogormlearn/routers"
)

func main() {
	r := routers.InitRoute()
	//inits.InitSentinel()
	//inits.InitRedis()

	//	test.Rand()
	//test.TestShardingProxy()
	r.Run(":8085")
}
