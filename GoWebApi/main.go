package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	_ "webapi/docs"

	"webapi/core"
	"webapi/routers"
)

var (
	// GitVersion 指定版本号 ( go build -ldflags "-X main.GitVersion=xxx" )
	GitVersion string
	// ConfigFile 指定配置文件 ( webapi --config=./config.yaml )
	ConfigFile string
)

// @title	API接口
// @version	1.0
// @description	API接口文档
func main() {
	flag.StringVar(&ConfigFile, "config", "./config.yaml", "config path: --config config.yaml")
	flag.Parse()

	hostname, _ := os.Hostname()
	fmt.Printf("[go version: %s] [git: %s] [config file: %s] [hostname: %s] \n",
		runtime.Version(), GitVersion, ConfigFile, hostname)

	if err := routers.NewRouter(core.NewServer(ConfigFile)).HttpStart(); err != nil {
		panic(err)
	}
}
