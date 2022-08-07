package main

import (
	. "config"
	"myerror"
	"os"
	"time"
)

// 退出函数
func QuitFunc() {
	myerror.Write("程序在" + time.Now().String() + "关闭")
	DB.Close()
	time.Sleep(2 * time.Second)
	myerror.Write("程序在" + time.Now().String() + "完成关闭，即将退出")
	os.Exit(0)
}
