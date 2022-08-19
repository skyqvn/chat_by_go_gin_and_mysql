package main

import (
	. "config"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"myerror"
	"net/http"
	"os/signal"
	"syscall"
	"users"
)

func main() {
	var err error
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
	// 测试服务器连通
	Engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	homePage := Engine.Group("/", users.IsLogin)
	{
		homePage.GET("", Index)
		homePage.GET("chatgroup/:group_id", ChatGroup)
		homePage.POST("send_message/:group_id", SendMessage)
		homePage.Any("join_group/:group_id", JoinGroup)
		homePage.Any("create_group", CreateGroup)
		homePage.GET("delete_member/:group_id", DeleteMember)
		homePage.Any("search", Search)
	}
	user := Engine.Group("user/")
	{
		user.Any("login", users.Login)
		user.Any("logged_out", users.LoggedOut)
		user.Any("register", users.Register)
	}
	Engine.StaticFS("/static_file", http.Dir("static_file"))
	Engine.NoRoute(myerror.CRaise404)
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		for {
			<-Quit
			myerror.Write("检测到退出信号")
		invalidInput:
			var s string
			fmt.Println("确定退出吗？(y,n):")
			fmt.Scanln(&s)
			if s == "y" || s == "Y" {
				QuitFunc()
			} else if s == "n" || s == "N" {
				fmt.Println("取消了退出")
				continue
			} else {
				goto invalidInput
			}
		}
	}()
	err = Srv.ListenAndServe()
	if err != nil {
		myerror.Write(err.Error())
	}
	select {}
}
