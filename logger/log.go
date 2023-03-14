package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level    string `json:"level" yaml:"level"`
	Prefix   string `json:"prefix" yaml:"prefix"`
	Director string `json:"director" yaml:"director"`
	Savetime int64  `json:"savetime" yaml:"savetime"`
}

var logger *zap.SugaredLogger
var prefix string
var log_level zapcore.Level

func NewLogger(c Config) {
	var l *zap.Logger
	prefix = c.Prefix
	switch c.Level { // 初始化配置文件的Level
	case "debug":
		log_level = zap.DebugLevel
	case "info":
		log_level = zap.InfoLevel
	case "warn":
		log_level = zap.WarnLevel
	case "error":
		log_level = zap.ErrorLevel
	default:
		log_level = zap.InfoLevel
	}

	if log_level == zap.DebugLevel || log_level == zap.ErrorLevel {
		//New 从提供的 zapcore.Core 和 Options 构造一个新的 Logger。如果传递的 zapcore.Core 为零，则它会回退到使用无操作实现。这是构建 Logger 最灵活的方式，但也是最冗长的。
		//AddStacktrace 将 Logger 配置为记录处于或高于给定级别的所有消息的堆栈跟踪。
		l = zap.New(getEncoderCore(c.Director, c.Savetime), zap.AddStacktrace(log_level))
	} else {
		l = zap.New(getEncoderCore(c.Director, c.Savetime))
	}

	// 记录行号
	l = l.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	logger = l.Sugar()
	logger.Info("logger init done...")
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  config.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return config
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(path string, t int64) (core zapcore.Core) {
	writer, err := GetWriteSyncer(path, t) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	//NewCore 创建一个将日志写入 WriteSyncer 的 Core。
	return zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig()), writer, log_level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(prefix + "2006/01/02 - 15:04:05"))
}

func GetWriteSyncer(p string, t int64) (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(p, "%Y-%m-%d.log"), //日志的路径和文件名
		// zaprotatelogs.WithLinkName(CONFIG.Zap.LinkName), // 生成软链，指向最新日志文件
		zaprotatelogs.WithMaxAge(time.Duration(t*24)*time.Hour), //保存日期的时间
		zaprotatelogs.WithRotationTime(24*time.Hour),            //每天分割一次日志
	)
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}
