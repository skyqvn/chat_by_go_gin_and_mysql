package main

import "time"

type ChatGroupType struct {
	Id        uint64
	Name      string
	Introduce string
	Password  string
}

type ReportType struct {
	ChatGroup uint64
	UserId    uint64
	Value     string
	SendTime  time.Time
}

type MemberType struct {
	Owner     uint64
	ChatGroup uint64
}
