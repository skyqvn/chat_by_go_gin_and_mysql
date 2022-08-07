package main

import (
	. "config"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"myerror"
	"os"
	"time"
)

// 创建文件的函数，由DoSthOnTime()在午夜十二点定时调用
func CreateLogFile() {
	T = time.Now()
	fn := fmt.Sprint("./log/", T.Year(), "-", T.Month(), "-", T.Day(), " ", T.Hour(), "-", T.Minute(), "-", T.Second(), " chat.log")
	fmt.Printf("create log file:%s", fn)
	f, err := os.Create(fn)
	if err != nil {
		myerror.LogError(err)
		return
	}
	LogFile = io.MultiWriter(f, os.Stdout)
}

// 用于定时任务
func DoSthOnTime() {
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", CreateLogFile)
	if err != nil {
		myerror.LogError(err)
		return
	}
	c.Start()
}
