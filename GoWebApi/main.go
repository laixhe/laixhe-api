package main

import (
	"flag"

	"webapi/core"
	_ "webapi/docs"
	"webapi/routers"
)

var (
	// ConfigFile 指定配置文件 (webapi --config=./config.yaml)
	ConfigFile string
)

// @title	API接口
// @version	1.0
// @description	API接口文档
func main() {
	flag.StringVar(&ConfigFile, "config", "./config.yaml", "config path: --config config.yaml")
	flag.Parse()
	if err := routers.NewRouter(core.NewServer(ConfigFile)).HttpStart(); err != nil {
		panic(err)
	}
}
