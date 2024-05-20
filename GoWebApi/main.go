package main

import (
	"flag"
	"time"
	
	"webapi/core/config"
	"webapi/core/gormx"
	"webapi/docs"
	"webapi/router"
)

var (
	// Version 指定版本号 ( go build -ldflags "-X main.Version=10000" )
	Version string
)

// flagConfigFile 指定配置文件
var flagConfigFile string

// @title	API接口
func main() {
	// init config
	flag.StringVar(&flagConfigFile, "config", "./config.yaml", "config path eg: -config config.yaml")
	flag.Parse()
	config.Init(flagConfigFile)
	// api doc
	Version += " " + time.Now().Format(time.DateTime)
	docs.SwaggerInfo.Description = docs.ErrorDescription()
	docs.SwaggerInfo.Version = Version
	config.Get().App.Version = Version
	// init data db
	gormx.Init(config.Get().Mysql)
	//redisx.Init(config.Get().Redis)
	// init http
	if err := router.Router().Run(config.Get().Servers.Http.Addr()); err != nil {
		panic(err)
	}
}
