package main

import (
	"gogormlearn/routers"
	"gogormlearn/test"
)

func main() {
	r := routers.InitRoute()
	//inits.InitSentinel()

	test.TestShardingProxy()
	r.Run()
}
