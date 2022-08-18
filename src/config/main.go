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
var DB *sql.DB
var Engine = gin.New()
var Srv = &http.Server{
	Handler: Engine,
}
var R = rand.New(rand.NewSource(time.Now().Unix()))
var T = time.Now()
var F, e2 = os.Create(fmt.Sprint("./log/", T.Year(), ";", T.Month(), ";", T.Day(), " ", T.Hour(), ";", T.Minute(), ";", T.Second(), " chat.log"))
var LogFile = io.MultiWriter(F, os.Stdout)
var SourceName string
var Quit = make(chan os.Signal)

func init() {
	var err error
	ReadConfig()
	DB, err = sql.Open("mysql", SourceName)
	if err != nil {
		fmt.Println("数据库错误：", err.Error())
		return
	}
	if e2 != nil {
		fmt.Println("文件打开错误：", e2.Error())
		return
	}
}
