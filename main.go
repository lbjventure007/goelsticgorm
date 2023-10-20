package main

import (
	"gogormlearn/routers"
)

func main() {
	r := routers.InitRoute()
	//inits.InitSentinel()
	//inits.InitRedis()

	// Modify the policy.
	//e.AddPolicy("alice", "data1", "read")//这里可以理解为角色
	//e.AddPolicy("alice", "data2", "read")//这里可以理解为角色  表示alice 有data1 data2 的read权限
	// e.RemovePolicy(...)

	//	e.AddGroupingPolicy("1", "alice") //这里可以理解为 把用户添加到对应的角色分组里
	// Save the policy back to DB.
	//e.SavePolicy()
	//	test.Rand()
	//test.TestShardingProxy()
	r.Run(":8085")
}
