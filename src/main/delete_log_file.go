package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"myerror"
	"os"
	"time"
)

func CreateLogOnTime() {
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", func() {
		t = time.Now()
		fmt.Println(fmt.Sprint("create ./log/", t.Year(), ";", t.Month(), ";", t.Day(), " ", t.Hour(), ";", t.Minute(), ";", t.Second(), " chat.log"))
		f, err = os.Create(fmt.Sprint("./log/", t.Year(), ";", t.Month(), ";", t.Day(), " ", t.Hour(), ";", t.Minute(), ";", t.Second(), " chat.log"))
		if err != nil {
			myerror.LogError(err)
			return
		}
		LogFile = io.MultiWriter(f, os.Stdout)
	})
	if err != nil {
		myerror.LogError(err)
		return
	}
	c.Start()
}
