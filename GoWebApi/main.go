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
	// FlagConfigFile 指定配置文件 (webapi --config=./config.yaml)
	FlagConfigFile string
)

// @title	API接口
func main() {
	// init flag
	flag.StringVar(&FlagConfigFile, "config", "./config.yaml", "config path: --config config.yaml")
	flag.Parse()
	fmt.Println("main show", FlagConfigFile, GitVersion)
	// init config
	core.Init(FlagConfigFile, GitVersion)
	// init http
	if err := router.Router().Run(core.HttpAddr()); err != nil {
		panic(err)
	}
}
