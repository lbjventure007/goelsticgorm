package main

import (
	"gogormlearn/routers"
)

func main() {
	r := routers.InitRoute()
	//inits.InitSentinel()

	r.Run()
}
