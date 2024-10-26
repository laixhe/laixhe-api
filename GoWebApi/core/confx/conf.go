package confx

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"webapi/api/gen/config"
	"webapi/api/gen/config/capp"
	"webapi/api/gen/config/cauth"
	"webapi/api/gen/config/cdb"
	"webapi/api/gen/config/clog"
	"webapi/api/gen/config/cmongodb"
	"webapi/api/gen/config/credis"
	"webapi/api/gen/config/cserver"
	"webapi/core/logx"
)

var cc *config.Config

// Get Config
func Get() *config.Config {
	return cc
}

// splitConfigFile 通过文件路径获取目录、文件名、扩展名
func splitConfigFile(configFile string) (dir string, fileName string, extName string, err error) {
	if len(configFile) == 0 {
		err = errors.New(configFile + " is empty")
		return
	}
	configFiles := strings.Split(configFile, "/")
	lens := len(configFiles) - 1
	if lens == 0 {
		dir = "."
	} else {
		dir = strings.Join(configFiles[:lens], "/")
	}
	files := strings.Split(configFiles[lens], ".")
	if len(files) <= 1 {
		err = errors.New(configFile + " file name is empty")
		return
	}
	fileName = files[0]
	extName = files[1]
	return
}

// InitViper 初始化配置文件
// configFile 配置文件
// isEnv      是否获取环境变量环境
// loadData   装载的数据结构指针类型
func InitViper(configFile string, isEnv bool, loadData interface{}) error {
	dir, fileName, extName, err := splitConfigFile(configFile)
	if err != nil {
		return err
	}

	v := viper.New()
	// 设置配置文件的名字
	v.SetConfigName(fileName)
	// 添加配置文件所在的路径
	v.AddConfigPath(dir)
	// 设置配置文件类型
	v.SetConfigType(extName)
	if err = v.ReadInConfig(); err != nil {
		return err
	}

	// 优先替换环境变量
	if isEnv {
		envConfigs := ListEnvConfig()
		for k := range envConfigs {
			envConfig := envConfigs[k]
			env := os.Getenv(envConfig.env)
			if env != "" {
				v.Set(envConfig.configKey, env) // 替换
			}
		}
	}

	if err = v.Unmarshal(loadData); err != nil {
		return err
	}
	return nil
}

// Init 初始化配置
func Init(configFile string) {
	cc = &config.Config{
		App:     &capp.App{},
		Log:     &clog.Log{},
		Servers: &cserver.Servers{},
		Jwt:     &cauth.Jwt{},
		Db:      &cdb.DB{},
		Mongodb: &cmongodb.MongoDB{},
		Redis:   &credis.Redis{},
	}
	if err := InitViper(configFile, true, cc); err != nil {
		panic(err)
	}

	// 初始化日志
	logx.Init(cc.Log)

	if cc.App.GetPid() != "" {
		pid := os.Getpid()
		if err := os.WriteFile(cc.App.GetPid(), []byte(fmt.Sprintf("%d", pid)), 0666); err != nil {
			panic(err)
		}
	}
}
