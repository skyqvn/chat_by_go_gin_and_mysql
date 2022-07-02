package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"myerror"
	"strconv"
	"users"
)

func Index(context *gin.Context) {
	user, err := users.GetUser(context)
	if err != nil {
		return
	}
	var sli []ChatGroupType
	rows, err := DB.Query("select chatgroup from member where owner=?", user)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	var i int
	var r *sql.Rows
	var cg ChatGroupType
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		r, err = DB.Query("select name,password,introduce,id from chatgroup where id=?", i)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		r.Next()
		err = r.Scan(&cg.Name, &cg.Password, &cg.Introduce, &cg.Id)
		if err != nil {
			myerror.Raise404(context, err)
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
	groupId, err := strconv.ParseUint(context.Param("group_id"), 10, 64)
	if err != nil {
		myerror.Raise500(context, err)
		return
	}
	var group ChatGroupType
	var reports []Ru
	rows, err := DB.Query("select name,password,introduce,id from chatgroup where id=?", groupId)
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
	rows, err = DB.Query("select chatgroup,userid,value,send_time from report where chatgroup=? order by send_time desc", groupId)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	var r Ru
	for rows.Next() {
		r = Ru{}
		err := rows.Scan(&r.R.ChatGroup, &r.R.UserID, &r.R.Value, &r.R.SendTime)
		rs, err := DB.Query("select id,name,introduce from user where id=?", r.R.UserID)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		rs.Next()
		err = rs.Scan(&r.U.Id, &r.U.Name, &r.U.Introduce)
		if err != nil {
			myerror.Raise500(context, err)
		}
		reports = append(reports, r)
	}
	context.HTML(200, "main/chatgroup", gin.H{
		"group":   group,
		"reports": reports,
	})
}

func SendMessage(context *gin.Context) {
	if context.Request.Method == "POST" {
		var ok bool
		var err error
		form := ReportType{}
		form.ChatGroup, err = strconv.ParseUint(context.Param("group_id"), 10, 64)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		var group uint64
		rows, err := DB.Query("select id from chatgroup where id=?", form.ChatGroup)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		rows.Next()
		err = rows.Scan(&group)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		form.Value, ok = context.GetPostForm("value")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("SendMessage:无value字段"))
			return
		}
		form.UserID, err = users.GetUser(context)
		if err != nil {
			return
		}
		_, err = DB.Exec("insert into report(chatgroup,userid,value) values (?,?,?)", form.ChatGroup, form.UserID, form.Value)
		if err != nil {
			myerror.Raise404(context, err)
			return
		}
		context.Redirect(302, "/chatgroup/"+context.Param("group_id"))
		return
	}
	myerror.Raise500(context, fmt.Errorf("SendMessage:请求方法错误"))
}

func JoinGroup(context *gin.Context) {
	var password string
	var group ChatGroupType
	user, err := users.GetUser(context)
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
		rows, err := DB.Query("select name,introduce,password from chatgroup where id=?", group.Id)
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
		rows, err = DB.Query("select * from member where chatgroup=? and owner=?", group.Id, user)
		if err != nil {
			myerror.Raise500(context, err)
			return
		}
		if rows.Next() {
			context.Redirect(302, "/")
			return
		}
		password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("JoinGroup:无password字段"))
			return
		}
		if password == group.Password {
			var member MemberType
			member.Owner = user
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
			myerror.Raise404(context, fmt.Errorf("CreateGroup:无name字段"))
			return
		}
		form.Password, ok = context.GetPostForm("password")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("CreateGroup:无password字段"))
			return
		}
		password, ok := context.GetPostForm("password2")
		if !ok {
			myerror.Raise404(context, fmt.Errorf("CreateGroup:无password2字段"))
			return
		}
		if form.Password != password {
			context.HTML(200, "main/create_group", gin.H{
				"form":    form,
				"warning": "密码不匹配",
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
	rows, err := DB.Query("select owner,chatgroup from member where owner=? and chatgroup=?", user, groupId)
	if err != nil {
		myerror.Raise404(context, err)
		return
	}
	//查看rows长度是否为零
	if rows.Next() == false { //如果长度为零
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
		if !rows.Next() { //如果聊天室没有人了
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
