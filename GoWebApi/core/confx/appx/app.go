package appx

import (
	"webapi/api/gen/config/capp"
	"webapi/core/confx"
	"webapi/core/logx"
)

func Checking() {
	c := confx.Get().GetApp()
	if c.Version == "" {
		c.Version = "v0.1"
	}
	if c.Mode == "" {
		c.Mode = capp.ModeType_debug.String()
	} else {
		c.Mode = capp.ModeType_name[capp.ModeType_value[c.Mode]]
	}
	logx.Debugf("app Config=%v", c)
}
