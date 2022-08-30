package config

import (
	"myerror"
	"os"
)

type WebSiteConfig struct {
	Port       int
	ServerAddr string
	URL        string
	SourceName string
}

func ReadConfig() {
	URL = os.Getenv("URL")
	ServerAddr = os.Getenv("ServerAddress") + ":" + os.Getenv("ServerPort")
	Srv.Addr = ServerAddr
	SourceName = os.Getenv("SourceName")
	myerror.Write("服务地址:" + ServerAddr)
}
