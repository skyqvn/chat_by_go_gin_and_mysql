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

func CreateLogFile() {
	T = time.Now()
	fmt.Println(fmt.Sprint("create ./log/", T.Year(), ";", T.Month(), ";", T.Day(), " ", T.Hour(), ";", T.Minute(), ";", T.Second(), " chat.log"))
	f, err := os.Create(fmt.Sprint("./log/", T.Year(), ";", T.Month(), ";", T.Day(), " ", T.Hour(), ";", T.Minute(), ";", T.Second(), " chat.log"))
	if err != nil {
		myerror.LogError(err)
		return
	}
	LogFile = io.MultiWriter(f, os.Stdout)
}

func DoSthOnTime() {
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", CreateLogFile)
	if err != nil {
		myerror.LogError(err)
		return
	}
	c.Start()
}
