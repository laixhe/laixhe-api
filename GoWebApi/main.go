package main

import (
	"flag"
	"fmt"

	"webapi/core"
	"webapi/router"
)

var (
	// GitVersion 指定版本号 ( go build -ldflags "-X main.GitVersion=xxx" )
	GitVersion string
)

// flagConfigFile 指定配置文件 (webapi --config=./config.yaml)
var flagConfigFile string

// @title	API接口
func main() {
	// init config
	flag.StringVar(&flagConfigFile, "config", "./config.yaml", "config path: --config config.yaml")
	flag.Parse()
	fmt.Println("main show", flagConfigFile, GitVersion)
	core.Init(flagConfigFile, GitVersion)
	// init http
	if err := router.Router().Run(core.HttpAddr()); err != nil {
		panic(err)
	}
}
