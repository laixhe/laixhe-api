package core

import (
	"fmt"
	"time"

	"github.com/laixhe/gonet"

	"webapi/core/config"
	"webapi/docs"
)

type Core struct {
	Config *config.Config
}

var coreData *Core

func Init(configFile, gitVersion string) {
	// init config
	c := config.Init(configFile).AppChecking().HttpChecking().JwtChecking()
	// init api doc
	docs.SwaggerInfo.Description = docs.ErrorDescription()
	docs.SwaggerInfo.Version = gitVersion + "-" + time.Now().Format(time.DateTime)

	// init db
	if err := gonet.GormInit(c.Gorm); err != nil {
		panic(err)
	}
	// if err := gonet.MongoInit(c.Mongodb); err != nil {
	// 	panic(err)
	// }
	if err := gonet.RedisInit(c.Redis); err != nil {
		panic(err)
	}

	coreData = &Core{Config: c}
}

func Config() *config.Config {
	return coreData.Config
}

func HttpAddr() string {
	return fmt.Sprintf(":%d", coreData.Config.Http.Port)
}
