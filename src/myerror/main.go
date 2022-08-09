package myerror

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

var Months = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

// 报告401错误，显示待登录页面
func Raise401(context *gin.Context, err error) {
	now := time.Now()
	month:=Months[time.Now().Month().String()]
	ts := fmt.Sprint(now.Year(), "/", month, "/", now.Day(), " - ", now.Hour(), ":", now.Minute(), ":", now.Second(), " ")
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("[INFO]", ts, "401ERROR:", err.Error())))
	if err != nil {
		fmt.Println("Raise401:无法写入", err)
	}
	context.HTML(401, "myerror/401error", nil)
}

// 报告404错误
func Raise404(context *gin.Context, err error) {
	now := time.Now()
	month:=Months[time.Now().Month().String()]
	ts := fmt.Sprint(now.Year(), "/", month, "/", now.Day(), " - ", now.Hour(), ":", now.Minute(), ":", now.Second(), " ")
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("[INFO]", ts, "404ERROR:", err.Error())))
	if err != nil {
		fmt.Println("Raise404:无法写入", err)
	}
	context.HTML(404, "myerror/404error", nil)
}

// 报告500错误
func Raise500(context *gin.Context, err error) {
	now := time.Now()
	month:=Months[time.Now().Month().String()]
	ts := fmt.Sprint(now.Year(), "/", month, "/", now.Day(), " - ", now.Hour(), ":", now.Minute(), ":", now.Second(), " ")
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("[INFO]", ts, "500ERROR:", err.Error())))
	if err != nil {
		fmt.Println("Raise500:无法写入", err)
	}
	context.HTML(500, "myerror/500error", nil)
}

// 路由的NoRoute函数
func CRaise404(context *gin.Context) {
	Raise404(context, fmt.Errorf("Underfined path %s", context.Request.RequestURI))
}

// 向用户显示警告并向日志文件和控制台输出
func ShowWarning(context *gin.Context, err error, s string) {
	context.HTML(200, "myerror/show_warning", s)
	now := time.Now()
	month:=Months[time.Now().Month().String()]
	ts := fmt.Sprint(now.Year(), "/", month, "/", now.Day(), " - ", now.Hour(), ":", now.Minute(), ":", now.Second(), " ")
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("[INFO]", ts, "ShowWarning:", err.Error())))
	if err != nil {
		fmt.Println("ShowWarning:无法写入", err)
	}
}

// 用于向日志文件和控制台输出错误
func LogError(err error) {
	now := time.Now()
	month:=Months[time.Now().Month().String()]
	ts := fmt.Sprint(now.Year(), "/", month, "/", now.Day(), " - ", now.Hour(), ":", now.Minute(), ":", now.Second(), " ")
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("[INFO]", ts, "LogError:", err.Error())))
	if err != nil {
		fmt.Println("LogError:无法写入", err)
	}
}

// 用于向日志文件和控制台输出
func Write(s string) {
	now := time.Now()
	month:=Months[time.Now().Month().String()]
	ts := fmt.Sprint(now.Year(), "/", month, "/", now.Day(), " - ", now.Hour(), ":", now.Minute(), ":", now.Second(), " ")
	_, err := gin.DefaultWriter.Write([]byte(fmt.Sprintln("[INFO]" + ts + "Write:" + s + "\n")))
	if err != nil {
		fmt.Println("Write:无法写入", err)
	}
}
