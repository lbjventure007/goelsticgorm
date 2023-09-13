package controllers

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
)

type TestController struct {
}

func NewTestController() TestController {
	return TestController{}
}
func (t TestController) Index() gin.HandlerFunc {
	return func(context *gin.Context) {
		sentil, blockError := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))

		if blockError != nil {
			context.JSON(400, gin.H{
				"message": "访问太频繁,已被限流",
			})
			context.Abort()
			return
		}
		//正常请求后需要关闭
		defer func() {
			sentil.Exit()
		}()

		context.JSON(200, gin.H{
			"message": "index",
			"data":    "haha",
		})

	}
}

func (t TestController) Test() gin.HandlerFunc {
	return func(context *gin.Context) {

		context.JSON(200, gin.H{
			"message": "test",
			"data":    "haha",
		})
	}
}
