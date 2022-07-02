package myerror

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Raise401(context *gin.Context, err error) {
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("401ERROR:", err.Error())))
	if err != nil {
		fmt.Println("Raise401:无法写入", err)
	}
	context.HTML(401, "myerror/401error", nil)
}

func Raise404(context *gin.Context, err error) {
	fmt.Println("404ERROR:", err.Error())
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("404ERROR:", err.Error())))
	if err != nil {
		fmt.Println("Raise404:无法写入", err)
	}
	context.HTML(404, "myerror/404error", nil)
}

func Raise500(context *gin.Context, err error) {
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("500ERROR:", err.Error())))
	if err != nil {
		fmt.Println("Raise500:无法写入", err)
	}
	context.HTML(500, "myerror/500error", nil)
}

func CRaise404(context *gin.Context) {
	Raise404(context, fmt.Errorf("Underfined path %s", context.Request.RequestURI))
}

func ShowWarning(context *gin.Context, err error, s string) {
	context.HTML(200, "myerror/show_warning", s)
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("ShowWarning:", err.Error())))
	if err != nil {
		fmt.Println("ShowWarning:无法写入", err)
	}
}

func LogError(err error) {
	_, err = gin.DefaultWriter.Write([]byte(fmt.Sprintln("LogError:", err.Error())))
	if err != nil {
		fmt.Println("LogError:无法写入", err)
	}
}
