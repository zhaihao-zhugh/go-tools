package log

import (
	"fmt"
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LOGCONFIG struct {
	Level    string `json:"level" yaml:"level"`
	Prefix   string `json:"prefix" yaml:"prefix"`
	Director string `json:"director"  yaml:"director"`
	Savetime int64  `json:"savetime" yaml:"savetime"`
}

type LOGGER struct {
	*zap.Logger
}

var level zapcore.Level
var l *zap.Logger

func NewLogger(lever, prefix, director string, savetime int64) *LOGGER {
	switch lever { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		//New 从提供的 zapcore.Core 和 Options 构造一个新的 Logger。如果传递的 zapcore.Core 为零，则它会回退到使用无操作实现。这是构建 Logger 最灵活的方式，但也是最冗长的。
		//AddStacktrace 将 Logger 配置为记录处于或高于给定级别的所有消息的堆栈跟踪。
		l = zap.New(getEncoderCore(director, savetime), zap.AddStacktrace(level))
	} else {
		l = zap.New(getEncoderCore(director, savetime))
	}

	// 记录行号
	l = l.WithOptions(zap.AddCaller())

	l.Info("logger init done...")
	return &LOGGER{l}
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
		EncodeCaller:   zapcore.FullCallerEncoder,
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
	return zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig()), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.Prefix + "2006/01/02 - 15:04:05"))
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
