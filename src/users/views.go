package users

import (
	. "config"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"myerror"
)

func Login(context *gin.Context) {
	var form UserType
	var ok bool
	if context.Request.Method == "POST" {
		form.Name, ok = context.GetPostForm("name")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Login:无name字段"))
			return
		}
		form.Password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Login:无password字段"))
			return
		}
		next, ok := context.GetPostForm("next")
		if !ok {
			next = "/"
		}
		err := LoginFunc(context, form, next)
		if err != nil {
			myerror.LogError(err)
		}
		return
	}
	context.HTML(200, "users/login", gin.H{
		"form":    form,
		"warning": "",
	})
}

func LoggedOut(context *gin.Context) {
	var rows *sql.Rows
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	loginCode, err := context.Cookie("login_code")
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	rows, err = DB.Query("select * from user where login_code=?", loginCode)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	if rows.Next() { //如果有此login_code
		context.SetCookie("login_code", "", 0, "/", Address, false, true)
		context.HTML(200, "users/logged_out", nil)
		return
	}
	err = rows.Close()
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	myerror.ShowWarning(context, fmt.Errorf("LoggedOut:登出失败"), "登出失败")
}

func Register(context *gin.Context) {
	var rows *sql.Rows
	var form UserType
	var ok bool
	var err error
	var password2 string
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	if context.Request.Method == "POST" {
		form.Name, ok = context.GetPostForm("name")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Register:无name字段"))
			return
		}
		form.Password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Register:无password字段"))
			return
		}
		rows, err = DB.Query("select * from user where name=?", form.Name)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		if rows.Next() { //如果已经有此名用户了
			context.HTML(200, "users/register", gin.H{
				"form":    form,
				"warning": "已有此名用户",
			})
			return
		}
		err = rows.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		password2, ok = context.GetPostForm("password2")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Register:无password2字段"))
			return
		}

		if password2 != form.Password {
			context.HTML(200, "users/register", gin.H{
				"form":    form,
				"warning": "密码不匹配",
			})
			return
		}
		err = IsAValidChatGroupPassword(form.Password)
		if err != nil {
			context.HTML(200, "users/register", gin.H{
				"form":    form,
				"warning": err.Error(),
			})
			return
		}
		form.Introduce, ok = context.GetPostForm("introduce")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Register:无introduce字段"))
			return
		}
		next, ok := context.GetPostForm("next")
		if !ok {
			next = "/"
		}

		_, err = DB.Exec("insert into user(name,password,introduce) values(?,?,?)", form.Name, form.Password, form.Introduce)
		if err == nil {
			//注册之后自动登录
			err = LoginFunc(context, form, next)
			if err == nil {
				context.Redirect(302, "/")
				return
			}
		}
	}
	context.HTML(200, "users/register", gin.H{
		"form":    form,
		"warning": "",
	})
}
