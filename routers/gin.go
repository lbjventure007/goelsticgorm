package routers

import (
	"github.com/gin-gonic/gin"
	"gogormlearn/controllers"
)

func InitRoute() *gin.Engine {
	r := gin.Default()

	r.GET("/", controllers.NewTestController().Index())
	r.GET("/test", controllers.NewTestController().Test())
	return r

}
