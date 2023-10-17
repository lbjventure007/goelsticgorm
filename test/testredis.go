package test

import (
	"context"
	"fmt"
	"gogormlearn/inits"
	"time"
)

func Rand() {

	//for i := 0; i < 1000000; i++ {
	//	_, err := inits.RedisClient.SAdd(context.TODO(), "userlistrandprize", i).Result()
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//}

	start := time.Now().Second()
	for i := 0; i < 100; i++ {
		resu, err := inits.RedisClient.SPopN(context.TODO(), "userlistrandprize", 100).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(resu)
	}
	end := time.Now().Second()
	cost := end - start
	fmt.Println(cost)
}
