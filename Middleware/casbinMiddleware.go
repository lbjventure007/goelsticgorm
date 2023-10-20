package Middleware

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func CasbinMiddleware(ctx *gin.Context) {

	//e, err := casbin.NewEnforcer("./conf/model.conf", "./conf/policy.csv")

	a, _ := gormadapter.NewAdapter("mysql", "root:1234qwer@tcp(localhost:3306)/gozero", true) // Your driver and data source.
	e, _ := casbin.NewEnforcer("./conf/model.conf", a)

	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	e.LoadPolicy()

	//sub := "alice"
	session := sessions.Default(ctx)
	user := session.Get("user")
	fmt.Println("user:--", user)
	username, ok := user.(string)

	if !ok || username == "" {
		ctx.AbortWithError(400, errors.New("没有登陆"))
		// deny the request, show an error
		ctx.Abort()
		return
	}

	sub := ""

	sub = username

	//obj := "data1"
	obj := ctx.Request.RequestURI
	act := ctx.Request.Method
	fmt.Println("sub: ", sub, " obj: ", obj, " act: ", act)
	if res, _ := e.Enforce(sub, obj, act); res {
		// permit alice to read data1
		fmt.Println(res)
		ctx.Next()
	} else {
		ctx.AbortWithError(400, errors.New("没有权限"))
		// deny the request, show an error
		ctx.Abort()
		return
	}
}
