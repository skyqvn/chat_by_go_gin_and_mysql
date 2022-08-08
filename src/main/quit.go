package main

import (
	. "config"
	"context"
	"myerror"
	"os"
	"time"
)

// 退出函数
// 用于在退出时进行资源的关闭并等待剩余服务完成
func QuitFunc() {
	myerror.Write("程序在" + time.Now().String() + "关闭")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := Srv.Shutdown(ctx); err != nil {
		myerror.Write(err.Error())
	}
	DB.Close()
	cancel()
	myerror.Write("程序在" + time.Now().String() + "完成关闭，即将退出")
	os.Exit(0)
}
