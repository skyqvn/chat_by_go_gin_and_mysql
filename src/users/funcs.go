package users

import (
	. "config"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"myerror"
	"strconv"
	"time"
)

// 如果有正确的cookies,则返回用户的ID和nil
// 否则进行报错处理,并返回错误
func GetUser(context *gin.Context) (uint64, error) {
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
		myerror.Raise401(context, err)
		return 0, err
	}
	rows, err = DB.Query("select id from user where login_code=?", loginCode)
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
		err = rows.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return 0, err
		}
		rows, err = DB.Query("select last_login_time from user where id=?", id)
		if err != nil {
			myerror.Raise500(context, err)
			return 0, err
		}
		var lastLoginTime time.Time
		rows.Next()
		err = rows.Scan(&lastLoginTime)
		if err != nil {
			myerror.Raise500(context, err)
			return 0, err
		}
		err = rows.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return 0, err
		}
		if time.Now().Sub(lastLoginTime).Hours() > 168 {
			_, err = DB.Exec("update user set login_code=0 where id=?", id)
			if err != nil {
				myerror.Raise500(context, err)
				return 0, err
			}
			myerror.Raise500(context, fmt.Errorf("cookies已过期"))
			return 0, fmt.Errorf("cookies已过期")
		}
		return id, nil
	}
	myerror.Raise401(context, fmt.Errorf("未登录"))
	return 0, fmt.Errorf("未登录")
}

// 判断用户是否可以按cookies登录的中间件
// 如果可以,就进行下一步操作
// 否则进入等待登录界面
func IsLogin(context *gin.Context) {
	_, err := GetUser(context)
	if err != nil {
		context.Abort()
	} else {
		context.Next()
	}
}

// 用于登录的函数
func LoginFunc(context *gin.Context, form UserType, next string) error {
	var rows *sql.Rows
	var password string
	var err error
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	rows, err = DB.Query("select password from user where name=?", form.Name)
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
	err = rows.Close()
	if err != nil {
		myerror.Raise500(context, err)
		return err
	}
	if password == form.Password {
		var loginCode uint64
		for {
			loginCode = R.Uint64()
			rows, err = DB.Query("select * from user where login_code=?", loginCode)
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
			err = rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
				return err
			}
		}
		_, err = DB.Exec("update user set login_code=?,last_login_time=? where name=? and password=?", loginCode, time.Now(), form.Name, form.Password)
		if err != nil {
			return err
		}
		context.SetCookie("login_code", strconv.FormatUint(loginCode, 10), 604800, "/", URL, false, true)
		context.Redirect(302, next)
		return nil
	}
	context.HTML(200, "users/login", gin.H{
		"form":    form,
		"warning": "密码错误",
	})
	return fmt.Errorf("LoginFunc:密码错误")
}

// 验证密码是否合法的函数
func IsAValidChatGroupPassword(password string) error {
	charList := [...]byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'k', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'K', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_', '-'}
	if len(password) > 20 {
		return fmt.Errorf("密码长度不能超过二十个字符")
	}
	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于八个字符")
	}
	for i := 0; i < len(password); i++ {
		ok := false
		for j := 0; j < len(charList); j++ {
			if password[i] == charList[j] {
				ok = true
				break
			}
		}
		if ok == false {
			return fmt.Errorf("密码必须是由英文字母大小写、数字、\"_\"、\"-\"组成的不多于二十个字符、不少于八个字符的字符串")
		}
	}
	return nil
}
