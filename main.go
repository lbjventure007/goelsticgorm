package main

import (
	"gogormlearn/inits"
	"gogormlearn/routers"
)

func main() {
	r := routers.InitRoute()
	inits.InitSentinel()

	r.Run()
}
