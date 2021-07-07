package logging

import (
	"fmt"
	"github.com/jweny/pocassist/pkg/conf"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config interface {
	Level() zapcore.Level
	LogPath() string
	MaxLogSize() int
	ServiceName() string
	InfoOutput() string
	ErrorOutput() string
	DebugOutput() string
}

func Setup(){
	loggerCfg := conf.GlobalConfig.LogConfig
	NewLogger(conf.GlobalConfig.Level(), loggerCfg.MaxSize, loggerCfg.MaxBackups, loggerCfg.MaxAge,
		loggerCfg.Compress, conf.GlobalConfig)
}


var (
	sugarInfoLogger  *zap.SugaredLogger
	sugarInfoPath    string
	sugarDebugLogger *zap.SugaredLogger
	sugarDebugPath   string
	sugarErrorLogger *zap.SugaredLogger
	sugarErrPath     string
)

func GetInfoLogPath() string {
	return sugarInfoPath
}

func GetDebugLogPath() string {
	return sugarDebugPath
}

func GetErrLogPath() string {
	return sugarErrPath
}

func formatArgs(v ...interface{}) string {
	var formatStrings []string
	for i := 0; i < len(v); i++ {
		t := v[i]
		switch reflect.TypeOf(t).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				formatStrings = append(formatStrings, `%v`)
			}
		}

	}
	// fmt.Println(v, len(v), formatStrings)
	return strings.Join(formatStrings, " ")
}

func Info(v ...interface{}) {
	format := formatArgs(v)
	sugarInfoLogger.Info("", fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	format := formatArgs(v)
	sugarErrorLogger.Error("", fmt.Sprintf(format, v...))
}

func Errorln(v ...interface{}) {
	sugarErrorLogger.Error("", fmt.Sprintln(v...))
}

func Warn(v ...interface{}) {
	format := formatArgs(v)
	sugarErrorLogger.Warn("", fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	format := formatArgs(v)
	sugarDebugLogger.Debug("", fmt.Sprintf(format, v...))
}

func Panic(v ...interface{}) {
	format := formatArgs(v)
	sugarErrorLogger.Panic("", fmt.Sprintf(format, v...))
}

func createLogger(path string, level zapcore.Level, maxSize int, maxBackups int,
	maxAge int, compress bool) *zap.SugaredLogger {
	core := newCore(path, level, maxSize, maxBackups, maxAge, compress)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger.Sugar()
}

func NewLogger(level zapcore.Level, maxSize int, maxBackups int,
	maxAge int, compress bool, config Config) {
	var infoPath, debugPath, errPath, logPath string

	if _, err := os.Stat(config.LogPath()); os.IsNotExist(err) {
		os.Mkdir(config.LogPath(), 0755)
	}

	if len(config.LogPath()) == 0 {
		logPath = "logs"
	}
	infoPath = fmt.Sprintf("%s/%v.log", logPath, config.ServiceName())
	if config.InfoOutput() != "" {
		infoPath = config.InfoOutput()
	}
	sugarInfoLogger = createLogger(infoPath, level, maxSize, maxBackups, maxAge, compress)
	sugarInfoPath = infoPath

	sugarDebugLogger = sugarInfoLogger
	sugarDebugPath = infoPath

	sugarErrorLogger = sugarInfoLogger
	sugarErrPath = infoPath

	if config.DebugOutput() != "" {
		debugPath = config.DebugOutput()
		sugarDebugLogger = createLogger(debugPath, level, maxSize, maxBackups, maxAge, compress)
		sugarDebugPath = debugPath
	}

	if config.ErrorOutput() != "" {
		errPath = config.ErrorOutput()
		sugarErrorLogger = createLogger(errPath, level, maxSize, maxBackups, maxAge, compress)
		sugarErrPath = errPath
	}

	// logger = zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
}

/**
 * zapcore构造
 */
func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	//日志文件路径配置2
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
