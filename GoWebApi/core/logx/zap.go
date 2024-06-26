package logx

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

/**

log:
  # 日志文件路径
  path: logs.log
  # 日志模式 console file
  run_type: console
  # 日志级别 debug  info  error
  level: debug
  # 每个日志文件保存大小 20M
  max_size: 20
  # 保留 N 个备份
  max_backups: 20
  # 保留 N 天
  max_age: 7

*/

var (
	once        sync.Once
	zapLogger   *zap.Logger
	sugarLogger *zap.SugaredLogger
	config      *Config
)

// Config 配置
type Config struct {
	RunType    LogRun   `mapstructure:"run_type"`    // 日志模式
	Level      LogLevel `mapstructure:"level"`       // 日志级别
	Path       string   `mapstructure:"path"`        // 日志文件路径
	MaxSize    uint     `mapstructure:"max_size"`    // 每个日志文件保存大小
	MaxBackups uint     `mapstructure:"max_backups"` // 保留N个备份，默认不限
	MaxAge     uint     `mapstructure:"max_age"`     // 保留N天，默认不限
	Compress   bool     `mapstructure:"compress"`    // 是否压缩，默认不压缩
	CallerSkip uint     `mapstructure:"caller_skip"` // 提升的堆栈帧数,默认 0，0=当前函数，1=上一层函数，....
}

// Init 初始日志
func Init(c *Config) {
	once.Do(func() {
		if c.RunType == LogRunFile {
			initSizeFile(c)
		} else {
			initConsole(c)
		}
		config = c
	})
}

// initSizeFile 初始 zap 日志，按大小切割和备份个数、文件有效期
func initSizeFile(c *Config) {
	// 日志分割
	hook := &lumberjack.Logger{
		Filename:   c.Path, // 日志文件路径，默认 os.TempDir()
		MaxSize:    int(c.MaxSize),
		MaxBackups: int(c.MaxBackups),
		MaxAge:     int(c.MaxAge),
		Compress:   c.Compress,
	}
	// 打印到文件
	write := zapcore.AddSync(hook)
	// 初始 zap 日志
	zapInit(write, c.Level, int(c.CallerSkip)+callerSkip)
}

// initConsole 初始 zap 日志，输出到终端
func initConsole(c *Config) {
	// 打印到控制台
	write := zapcore.AddSync(os.Stdout)
	// 初始 zap 日志
	zapInit(write, c.Level, int(c.CallerSkip)+callerSkip)
}

// zapInit 初始化 zap 基本信息
// write       文件描述符
// serviceName 服务名
// logLevel    日志级别
// callerSkip 提升的堆栈帧数，0=当前函数，1=上一层函数，....
func zapInit(write zapcore.WriteSyncer, logLevel LogLevel, callerSkip int) {
	// 设置日志级别
	var level zapcore.Level
	switch logLevel {
	case LogLevelDebug:
		level = zap.DebugLevel
	case LogLevelInfo:
		level = zap.InfoLevel
	case LogLevelWarn:
		level = zap.WarnLevel
	case LogLevelError:
		level = zap.ErrorLevel
	default:
		level = zap.DebugLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "call",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapTimeEncoder,                 // 日志时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行消耗的时间转化成浮点型的秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器,格式化调用堆栈
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		write,
		level,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 提升打印的堆栈帧数
	addCallerSkip := zap.AddCallerSkip(callerSkip)
	// 开启文件及行号
	development := zap.Development()
	// 添加字段-服务器名称
	//filed := zap.Fields(zap.String("service", serviceName))
	// 构造日志
	//zapLogger = zap.New(core, caller, addCallerSkip, development, filed)
	zapLogger = zap.New(core, caller, addCallerSkip, development)
	sugarLogger = zapLogger.Sugar()
}

func GetLevel() LogLevel {
	return config.Level
}

// zapTimeEncoder 日志时间格式
func zapTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.DateTime))
}

// Debug 调试
func Debug(msg string, args ...zap.Field) {
	zapLogger.Debug(msg, args...)
}

// Debugf 调试
func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

// Info 信息
func Info(msg string, args ...zap.Field) {
	zapLogger.Info(msg, args...)
}

// Infof 信息
func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

// Warn 警告
func Warn(msg string, args ...zap.Field) {
	zapLogger.Warn(msg, args...)
}

// Warnf 警告
func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

// Error 错误
func Error(msg string, args ...zap.Field) {
	zapLogger.Error(msg, args...)
}

// Errorf 错误
func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}
