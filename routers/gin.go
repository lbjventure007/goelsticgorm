package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gogormlearn/Middleware"
	"gogormlearn/controllers"
	// 导入session存储引擎
	"github.com/gin-contrib/sessions/cookie"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("user"))

	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/", controllers.NewTestController().Index())
	r.GET("/test", Middleware.CasbinMiddleware, controllers.NewTestController().Test())

	r.GET("/kafka", controllers.NewTestController().Kafka())
	r.GET("/login", controllers.NewUserController().Login())
	r.GET("/login-out", controllers.NewUserController().LoginOut())
	r.GET("/rocketmq", controllers.NewTestController().Rocketmq())
	r.GET("/rocketmq-tran", controllers.NewTestController().RocketmqTran())

	return r

}
