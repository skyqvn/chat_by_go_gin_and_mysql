package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"math/rand"
	"myerror"
	"net/http"
	"os"
	"time"
	"users"
)

const HostURL = "192.168.31.177"

var DB, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chat?parseTime=true")
var Engine = gin.New()
var R = rand.New(rand.NewSource(time.Now().Unix()))
var t = time.Now()
var f, err = os.Create(fmt.Sprint("./log/", t.Year(), ";", t.Month(), ";", t.Day(), " ", t.Hour(), ";", t.Minute(), ";", t.Second(), " chat.log"))
var LogFile = io.MultiWriter(f, os.Stdout)

func main() {
	CreateLogOnTime()
	gin.SetMode(gin.ReleaseMode)
	if err != nil {
		fmt.Println("文件打开错误：", err)
		return
	}
	gin.DefaultWriter = LogFile
	gin.DefaultErrorWriter = LogFile
	defer DB.Close()
	defer f.Close()
	Engine.Use(func(context *gin.Context) {
		gin.LoggerWithWriter(LogFile)(context)
	})
	Engine.LoadHTMLGlob("templates/**/*")
	users.DB = DB
	users.R = R
	users.LocalHost = HostURL
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
	err = Engine.Run("0.0.0.0:80")
	if err != nil {
		_, err = gin.DefaultWriter.Write([]byte(err.Error()))
		if err != nil {
			fmt.Println("无法写入")
		}
	}
}
