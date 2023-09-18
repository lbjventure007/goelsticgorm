package test

import (
	"fmt"
	"gogormlearn/inits"
	"gogormlearn/model"
)

func TestShardingProxy() {
	inits.InitDb()

	//-----全部字段指定
	//orders := model.Orders{}
	//orders.ID = 1
	//orders.UserID = 2
	//orders.OrderId = 121324432423425
	//orders.Name = "aa"
	//orders.CreateTime = time.Now()
	//reuslt := inits.Db.Create(&orders)
	//fmt.Println(reuslt.Error, reuslt.RowsAffected, orders.OrderId)

	//----指定部分字段
	//orders := model.Orders{}
	//orders.ID = 13
	//orders.UserID = 3
	////orders.OrderId = 121324432423425 //如果不指定 则 使用shardidngsphere proxy配置的order_id的分配生成算法
	//orders.Name = "aacccccccc"
	//orders.CreateTime = time.Now()
	//reuslt := inits.Db.Select("ID", "UserID", "Name", "CreateTime").Create(&orders)
	//fmt.Println("inser------", reuslt.Error, reuslt.RowsAffected, orders.OrderId)

	var o []model.Orders
	re := inits.Db.Where("user_id in ?", []interface{}{2, 3}).Where("name like ?", "aac%").Limit(2).Find(&o)

	var o1 model.Orders
	inits.Db.First(&o1, "order_id = ?", 910469871765028864)

	res := inits.Db.Model(&model.Orders{}).Where("order_id = ?", 910469871765028864).Update("name", "testcccddd11")
	var o2 model.Orders
	inits.Db.First(&o2, "order_id = ?", 910469871765028864)

	inits.Db.Where("order_id = ?", 910469871765028864).Delete(&model.Orders{})

	var o3 model.Orders
	inits.Db.First(&o3, "order_id = ?", 910469871765028864)

	fmt.Println(re.Error, o, o1, res.RowsAffected, o3)
	fmt.Println("批量查询", o)
	fmt.Println("查询910469871765028864", o1)
	fmt.Println("查询修后910469871765028864", o2)
	fmt.Println("查询删除后910469871765028864", o3)

}
