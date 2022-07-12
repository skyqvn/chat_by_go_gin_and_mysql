package main

import (
	. "config"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"myerror"
	"net/http"
	"users"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	DoSthOnTime()
	gin.DefaultWriter = LogFile
	gin.DefaultErrorWriter = LogFile
	defer DB.Close()
	defer F.Close()
	Engine.Use(func(context *gin.Context) {
		gin.LoggerWithWriter(LogFile)(context)
		context.Next()
	})
	Engine.LoadHTMLGlob("templates/**/*")
	Engine.GET("/favicon.ico", func(context *gin.Context) {
		context.File("./static_file/group_icon.ico")
	})
	homePage := Engine.Group("/", users.IsLogin)
	{
		homePage.GET("", Index)
		homePage.GET("chatgroup/:group_id", ChatGroup)
		homePage.POST("send_message/:group_id", SendMessage)
		homePage.Any("join_group/:group_id", JoinGroup)
		homePage.Any("create_group", CreateGroup)
		homePage.GET("delete_member/:group_id", DeleteMember)
		homePage.GET("search_group", Search)
	}
	user := Engine.Group("user/")
	{
		user.Any("login", users.Login)
		user.Any("logged_out", users.LoggedOut)
		user.Any("register", users.Register)
	}
	Engine.StaticFS("/static_file", http.Dir("static_file"))
	Engine.NoRoute(myerror.CRaise404)
	err := Engine.Run("0.0.0.0:80")
	if err != nil {
		_, err = gin.DefaultWriter.Write([]byte(err.Error()))
		if err != nil {
			fmt.Println("无法写入")
		}
	}
}
