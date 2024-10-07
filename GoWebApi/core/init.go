package core

import (
	"time"

	"webapi/core/confx"
	"webapi/core/confx/appx"
	"webapi/core/confx/serverx"
	"webapi/core/gormx"
	"webapi/core/jwtx"
	"webapi/core/redisx"
	"webapi/docs"
)

func Init(configFile, gitVersion string) {
	// init config
	confx.Init(configFile)
	// init api doc
	docs.SwaggerInfo.Description = docs.ErrorDescription()
	docs.SwaggerInfo.Version = gitVersion + "-" + time.Now().Format(time.DateTime)

	// config check
	appx.Checking()
	jwtx.Checking()
	serverx.Checking()

	// init db
	gormx.Init(confx.Get().GetDb())
	//mongox.Init(confx.Get().GetMongodb())
	redisx.Init(confx.Get().GetRedis())
}
