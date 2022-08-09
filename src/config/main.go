package config

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var URL string
var ServerAddr string
var DB, e1 = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chat?parseTime=true")
var Engine = gin.New()
var Srv = &http.Server{
	Handler: Engine,
}
var R = rand.New(rand.NewSource(time.Now().Unix()))
var T = time.Now()
var F, e2 = os.Create(fmt.Sprint("./log/", T.Year(), ";", T.Month(), ";", T.Day(), " ", T.Hour(), ";", T.Minute(), ";", T.Second(), " chat.log"))
var LogFile = io.MultiWriter(F, os.Stdout)
var Quit = make(chan os.Signal)

func init() {
	if e1 != nil {
		fmt.Println("数据库错误：", e1.Error())
		return
	}
	if e2 != nil {
		fmt.Println("文件打开错误：", e2.Error())
		return
	}
	ReadConfig()
}
