package routers

import (
	"github.com/gin-gonic/gin"
	"gogormlearn/controllers"
)

func InitRoute() *gin.Engine {
	r := gin.Default()

	r.GET("/", controllers.NewTestController().Index())
	r.GET("/test", controllers.NewTestController().Test())

	r.GET("/kafka", controllers.NewTestController().Kafka())
	r.GET("/rocketmq", controllers.NewTestController().Rocketmq())
	r.GET("/rocketmq-tran", controllers.NewTestController().RocketmqTran())
	return r

}
