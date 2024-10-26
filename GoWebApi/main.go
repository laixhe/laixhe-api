package main

import (
	"flag"
	
	"webapi/core"
	"webapi/core/confx/serverx"
	"webapi/router"
)

var (
	// GitVersion 指定版本号 ( go build -ldflags "-X main.GitVersion=10000" )
	GitVersion string
)

// flagConfigFile 指定配置文件
var flagConfigFile string

// @title	API接口
func main() {
	// init config
	flag.Parse()
	flag.StringVar(&flagConfigFile, "config", "./config.yaml", "config path eg: -config config.yaml")
	core.Init(flagConfigFile, GitVersion)
	// init http
	if err := router.Router().Run(serverx.HttpAddr()); err != nil {
		panic(err)
	}
}
