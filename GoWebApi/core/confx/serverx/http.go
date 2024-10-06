package serverx

import (
	"errors"
	"fmt"

	"webapi/core/confx"
	"webapi/core/logx"
)

// Checking 检查
func Checking() {
	c := confx.Get().Servers
	if c == nil {
		panic(errors.New("config servers is nil"))
	}
	if c.GetHttp() == nil {
		panic(errors.New("config servers.http is nil"))
	}
	if c.GetHttp().GetPort() == 0 {
		panic(errors.New("config servers.http.port is nil"))
	}

	logx.Debugf("http Config=%v", c)
}

func HttpAddr() string {
	return fmt.Sprintf("%s:%d", confx.Get().Servers.Http.GetIp(), confx.Get().Servers.Http.GetPort())
}
