package main

import (
	"flag"

	"webapi/core"
	"webapi/routers"
)

var (
	// ConfigFile 指定配置文件 (webapi --config=./config.yaml)
	ConfigFile string
)

func main() {
	flag.StringVar(&ConfigFile, "config", "./config.yaml", "config path: --config config.yaml")
	flag.Parse()
	if err := routers.NewRouter(core.NewServer(ConfigFile)).HttpStart(); err != nil {
		panic(err)
	}
}
