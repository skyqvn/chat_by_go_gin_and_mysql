package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myerror"
	"strconv"
)

func GetUser(context *gin.Context) (uint64, error) {
	loginCode, err := context.Cookie("login_code")
	if err != nil {
		myerror.Raise401(context, err)
		return 0, err
	}
	rows, err := DB.Query("select id from user where login_code=?", loginCode)
	if err != nil {
		myerror.Raise500(context, err)
		return 0, err
	}
	var id uint64
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			myerror.LogError(err)
			myerror.Raise500(context, err)
			return 0, err
		}
		return id, nil
	}
	myerror.Raise401(context, fmt.Errorf("未登录"))
	return 0, fmt.Errorf("未登录")
}

func IsLogin(context *gin.Context) {
	_, err := GetUser(context)
	if err != nil {
		context.Abort()
	} else {
		context.Next()
	}
}

func LoginFunc(context *gin.Context, form UserType, next string) error {
	var password string
	rows, err := DB.Query("select password from user where name=?", form.Name)
	if err != nil {
		myerror.Raise500(context, err)
		return err
	}
	if !rows.Next() {
		fmt.Println("3")
		context.HTML(200, "users/login", gin.H{
			"form":    form,
			"warning": "无此用户信息",
		})
		return fmt.Errorf("LoginFunc:无此用户信息")
	}
	err = rows.Scan(&password)
	if err != nil {
		myerror.Raise404(context, err)
		return err
	}
	if password == form.Password {
		var loginCode uint64
		for {
			loginCode = R.Uint64()
			rows, err = DB.Query("select login_code from user where login_code=?", loginCode)
			if err != nil {
				myerror.Raise500(context, err)
				return err
			}
			if loginCode == 0 {
				continue
			}
			if !rows.Next() {
				break
			}
		}
		_, err = DB.Exec("update user set login_code=? where name=? and password=?", loginCode, form.Name, form.Password)
		if err != nil {
			return err
		}
		context.SetCookie("login_code", strconv.FormatUint(loginCode, 10), 604800, "/", LocalHost, false, true)
		context.Redirect(302, next)
		return nil
	}
	context.HTML(200, "users/login", gin.H{
		"form":    form,
		"warning": "密码错误",
	})
	return fmt.Errorf("LoginFunc:密码错误")
}

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
	loginCode, err := context.Cookie("login_code")
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	rows, err := DB.Query("select * from user where login_code=?", loginCode)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	if rows.Next() {
		context.SetCookie("login_code", "", 0, "/", LocalHost, false, true)
		context.HTML(200, "users/logged_out", nil)
		return
	}
	myerror.ShowWarning(context, fmt.Errorf("LoggedOut:登出失败"), "登出失败")
}

func Register(context *gin.Context) {
	var form UserType
	var ok bool
	var password2 string
	if context.Request.Method == "POST" {
		form.Name, ok = context.GetPostForm("name")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Register:无name字段"))
			return
		}
		rows, err := DB.Query("select * from user where name=?", form.Name)
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

		form.Password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("Register:无password字段"))
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
