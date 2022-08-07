package main

import (
	"time"
	"users"
)

// 聊天室结构体
type ChatGroupType struct {
	Id        uint64
	Name      string
	Introduce string
	Password  string
}

// 消息结构体
type ReportType struct {
	ChatGroup uint64
	Owner     uint64
	Value     string
	SendTime  time.Time
}

// 连接结构体，一个用户与每一个他在的聊天群都有一个
type MemberType struct {
	Owner     uint64
	ChatGroup uint64
}

// 一个由Report和User整合的类
type Ru struct {
	R ReportType
	U users.UserType
}
