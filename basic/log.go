package basic

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
)

var GlobalLogger *logrus.Logger


func InitLog(debug bool, logFile string) error {
	GlobalLogger = &logrus.Logger{
		Formatter: &prefixed.TextFormatter{
			ForceColors:     true,
			ForceFormatting: true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	GlobalLogger.SetLevel(logrus.InfoLevel)
	if debug == true {
		GlobalLogger.SetLevel(logrus.DebugLevel)
	}
	fileOutput, err := rotatelogs.New(logFile)
	if err != nil {
		return err
	}
	mv := io.MultiWriter(os.Stdout, fileOutput)
	GlobalLogger.SetOutput(mv)
	return nil
}

// todo
//// GoodF print good message
//func Green(format string, args ...interface{}) {
//	good := color.HiGreenString("[+]")
//	fmt.Printf("%s %s\n", good, fmt.Sprintf(format, args...))
//}
//
//func Yellow(format string, args ...interface{}) {
//	good := color.YellowString("[!]")
//	fmt.Printf("%s %s\n", good, fmt.Sprintf(format, args...))
//}
//
//// InforF print info message
//func InforF(format string, args ...interface{}) {
//	GlobalLogger.Info(fmt.Sprintf(format, args...))
//}
//
//func Info(args ...interface{}) {
//	GlobalLogger.Infoln(args)
//}
//
//// ErrorF print good message
//func ErrorF(format string, args ...interface{}) {
//	GlobalLogger.Error(fmt.Sprintf(format, args...))
//}
//
//func Error(args ...interface{}) {
//	GlobalLogger.Errorln(args)
//}
//
//func WarningF(format string, args ...interface{}) {
//	GlobalLogger.Warningf(fmt.Sprintf(format, args...))
//}
//
//func Warning(args ...interface{}) {
//	GlobalLogger.Warningln(args)
//}
//
//// DebugF print debug message
//func DebugF(format string, args ...interface{}) {
//	GlobalLogger.Debug(fmt.Sprintf(format, args...))
//}
//
//func Debug(args ...interface{}) {
//	GlobalLogger.Debugln(args)
//}
