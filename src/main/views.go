package main

import (
	. "config"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"myerror"
	"strconv"
	"users"
)

func Index(context *gin.Context) {
	var rows *sql.Rows
	var i int
	var rs *sql.Rows
	var cg ChatGroupType
	var err error
	defer func() {
		if rows != nil {
			err = rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
		if rs != nil {
			err = rs.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	userId, err := users.GetUser(context)
	if err != nil {
		return
	}
	var sli []ChatGroupType
	rows, err = DB.Query("select chatgroup from member where owner=?", userId)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	for rows.Next() { //得到所有与当前用户ID匹配的聊天群
		err = rows.Scan(&i)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		rs, err = DB.Query("select name,password,introduce,id from chatgroup where id=?", i)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		rs.Next()
		err = rs.Scan(&cg.Name, &cg.Password, &cg.Introduce, &cg.Id)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		err = rs.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		sli = append(sli, cg)
	}
	context.HTML(200, "main/index", gin.H{
		"chatgroups": sli,
	})
}

func ChatGroup(context *gin.Context) {
	type Ru struct {
		R ReportType
		U users.UserType
	}
	var rows *sql.Rows
	var r Ru
	var rs *sql.Rows
	var err error
	defer func() {
		if rows != nil {
			err = rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
		if rs != nil {
			err = rs.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	groupId, err := strconv.ParseUint(context.Param("group_id"), 10, 64)
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	userId, err := users.GetUser(context)
	if err != nil {
		return
	}
	rows, err = DB.Query("select * from member where chatgroup=? and owner=?", groupId, userId)
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	if !rows.Next() {
		myerror.Raise404(context, fmt.Errorf("此聊天群中没有此用户"))
		return
	}
	err = rows.Close()
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	var group ChatGroupType
	var reports []Ru
	rows, err = DB.Query("select name,password,introduce,id from chatgroup where id=?", groupId)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	rows.Next()
	err = rows.Scan(&group.Name, &group.Password, &group.Introduce, &group.Id)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	err = rows.Close()
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	rows, err = DB.Query("select chatgroup,owner,value,send_time from report where chatgroup=? order by send_time desc", groupId)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	for rows.Next() { //获得所有与当前群关联的消息以及发送消息的用户
		r = Ru{}
		err := rows.Scan(&r.R.ChatGroup, &r.R.Owner, &r.R.Value, &r.R.SendTime)
		rs, err = DB.Query("select id,name,introduce from user where id=?", r.R.Owner)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		rs.Next()
		err = rs.Scan(&r.U.Id, &r.U.Name, &r.U.Introduce)
		if err != nil {
			myerror.Raise500(context, err)
		}
		err = rs.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		reports = append(reports, r)
	}
	context.HTML(200, "main/chatgroup", gin.H{
		"group":   group,
		"reports": reports,
	})
}

func SendMessage(context *gin.Context) {
	var rows *sql.Rows
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	if context.Request.Method == "POST" {
		var ok bool
		var err error
		form := ReportType{}
		form.ChatGroup, err = strconv.ParseUint(context.Param("group_id"), 10, 64)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		form.Owner, err = users.GetUser(context)
		if err != nil {
			return
		}
		rows, err = DB.Query("select * from member where chatgroup=? and owner=?", form.ChatGroup, form.Owner)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		if !rows.Next() {
			myerror.Raise404(context, fmt.Errorf("此聊天群中没有此用户"))
			return
		}
		form.Value, ok = context.GetPostForm("value")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("SendMessage:无value字段"))
			return
		}
		_, err = DB.Exec("insert into report(chatgroup,owner,value) values (?,?,?)", form.ChatGroup, form.Owner, form.Value)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		context.Redirect(302, "/chatgroup/"+context.Param("group_id"))
		return
	}
	//如果请求不是post
	myerror.Raise500(context, fmt.Errorf("SendMessage:请求方法错误"))
}

func JoinGroup(context *gin.Context) {
	var rows *sql.Rows
	var password string
	var group ChatGroupType
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	userId, err := users.GetUser(context)
	if err != nil {
		return
	}
	group.Id, err = strconv.ParseUint(context.Param("group_id"), 10, 64)
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	if context.Request.Method == "POST" {
		var ok bool
		rows, err = DB.Query("select name,introduce,password from chatgroup where id=?", group.Id)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		if !rows.Next() {
			myerror.Raise404(context, fmt.Errorf("JoinGroup:无此聊天群"))
			return
		}
		err = rows.Scan(&group.Name, &group.Introduce, &group.Password)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		err = rows.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		rows, err = DB.Query("select * from member where chatgroup=? and owner=?", group.Id, userId)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		if rows.Next() {
			context.Redirect(302, "/")
			return
		}
		err = rows.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("JoinGroup:无password字段"))
			return
		}
		if password == group.Password {
			var member MemberType
			member.Owner = userId
			member.ChatGroup = group.Id
			_, err := DB.Exec("insert into member(owner,chatgroup) values(?,?)", member.Owner, member.ChatGroup)
			if err != nil {
				myerror.Raise404(context, err)
				return
			}
			context.Redirect(302, "/")
			return
		}
	}

	context.HTML(200, "main/join_group", gin.H{
		"password": password,
		"group":    group,
	})
}

func CreateGroup(context *gin.Context) {
	var form ChatGroupType
	if context.Request.Method == "POST" {
		var ok bool
		form.Name, ok = context.GetPostForm("name")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("CreateGroup:无name字段"))
			return
		}
		form.Password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("CreateGroup:无password字段"))
			return
		}
		password, ok := context.GetPostForm("password2")
		if !ok {
			myerror.Raise500(context, fmt.Errorf("CreateGroup:无password2字段"))
			return
		}
		if form.Password != password {
			context.HTML(200, "main/create_group", gin.H{
				"form":    form,
				"warning": "密码不匹配",
			})
			return
		}
		err := IsAValidChatGroupPassword(form.Password)
		if err != nil {
			context.HTML(200, "main/create_group", gin.H{
				"form":    form,
				"warning": err.Error(),
			})
			return
		}
		form.Introduce, ok = context.GetPostForm("introduce")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("CreateGroup:无introduce字段"))
			return
		}
		result, err := DB.Exec("insert into chatgroup(name,password,introduce) values(?,?,?)", form.Name, form.Password, form.Introduce)
		if err == nil {
			id, err := result.LastInsertId()
			if err != nil {
				myerror.Raise404(context, err)
			}
			m := MemberType{}
			m.ChatGroup = uint64(id)
			m.Owner, err = users.GetUser(context)
			if err != nil {
				return
			}
			_, err = DB.Exec("insert into member(chatgroup,owner) values (?,?)", m.ChatGroup, m.Owner)
			if err != nil {
				myerror.Raise404(context, err)
				return
			}
			context.Redirect(302, "/")
			return
		}
	}
	context.HTML(200, "main/create_group", gin.H{
		"form":    form,
		"warning": "",
	})
}

func DeleteMember(context *gin.Context) {
	var rows *sql.Rows
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				myerror.Raise500(context, err)
			}
		}
	}()
	groupId, err := strconv.ParseUint(context.Param("group_id"), 10, 64)
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	user, err := users.GetUser(context)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	rows, err = DB.Query("select owner,chatgroup from member where owner=? and chatgroup=?", user, groupId)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	//查看rows长度是否为零
	ok := rows.Next()
	err = rows.Close()
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	if !ok { //如果长度为零
		myerror.Raise404(context, fmt.Errorf("DeleteMember:长度为零"))
		return
	} else { //如果长度不为零
		_, err = DB.Exec("delete from member where chatgroup=?", groupId)
		if err != nil {

			myerror.Raise500(context, err)
			return
		}
		rows, err = DB.Query("select * from member where chatgroup=?", groupId)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		ok = rows.Next()
		err = rows.Close()
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		if !ok { //如果聊天室没有人了
			_, err = DB.Exec("delete from chatgroup where id=?", groupId) //删除聊天室
			if err != nil {
				myerror.Raise500(context, err)
				return
			}
		}
		context.Redirect(302, "/")
	}
}

func Search(context *gin.Context) {
	context.HTML(200, "main/search", nil)
}
