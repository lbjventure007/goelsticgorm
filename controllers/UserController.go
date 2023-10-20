package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gogormlearn/model"
)

type UserController struct {
	BaseController
}

func NewUserController() UserController {
	return UserController{}
}

func (u UserController) Login() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		user := &model.User{}
		if err := ctx.Bind(user); err != nil {
			ctx.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		fmt.Println("bind--", user)
		user, errs := user.Login()
		if errs != "" {
			ctx.JSON(400, gin.H{
				"message1": errs,
			})
			return
		}

		session := sessions.Default(ctx)
		session.Set("user", user.Username)
		err := session.Save()
		fmt.Println("save session:", err)
		ctx.JSON(200, gin.H{
			"message": "ok",
			"data":    u,
		})
		return

	}
}

func (u UserController) LoginOut() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		session := sessions.Default(ctx)
		session.Delete("user")
		err := session.Save()
		fmt.Println("save session:", err)
		ctx.JSON(200, gin.H{
			"message": "退出成功",
		})
		return

	}
}