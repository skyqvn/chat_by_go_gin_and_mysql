package config

import (
	"encoding/json"
	"myerror"
	"os"
	"strconv"
)

type WebSiteConfig struct {
	Port int
	ServerAddr string
	URL string
}

func ReadConfig()  {
	configFile,err:=os.Open("../../conf.json")
	if err!=nil {
		myerror.LogError(err)
		return
	}
	defer configFile.Close()
	decoder:=json.NewDecoder(configFile)
	conf:=WebSiteConfig{}
	err=decoder.Decode(&conf)
	if err!=nil {
		myerror.LogError(err)
		return
	}
	URL =conf.URL
	ServerAddr=conf.ServerAddr+":"+strconv.FormatInt(int64(conf.Port),10)
	Srv.Addr=ServerAddr
}
