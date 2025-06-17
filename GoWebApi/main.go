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
	server := core.NewServer(ConfigFile)
	if err := routers.NewRouter(server).Listen(); err != nil {
		panic(err)
	}
}
