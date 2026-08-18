package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	_ "webapi/docs"

	"webapi/app/controllers"
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
// @description	用户认证与用户管理 API 服务
func main() {
	flag.StringVar(&ConfigFile, "config", "./config.yaml", "config path: --config config.yaml")
	flag.Parse()

	hostname, _ := os.Hostname()
	fmt.Printf("[go version: %s] [git: %s] [config file: %s] [hostname: %s] \n",
		runtime.Version(), GitVersion, ConfigFile, hostname)

	// 将构建注入的 GitVersion 同步给健康检查接口 (统一版本来源, 见 app/controllers/health.go)
	if GitVersion != "" {
		controllers.Version = GitVersion
	}

	if err := routers.NewRouter(core.NewServer(ConfigFile)).HttpStart(); err != nil {
		panic(err)
	}
}
