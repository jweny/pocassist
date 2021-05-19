package logging

import (
	conf2 "github.com/jweny/pocassist/pkg/conf"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

var GlobalLogger *logrus.Logger

func Setup(debug bool) {
	logName := conf2.GlobalConfig.ServerConfig.LogName
	if logName == "" {
		logName = "debug.log"
	}
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

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("logging.Setup, fail to get current dir")

	}
	file := path.Join(dir, logName)
	fileOutput, err := rotatelogs.New(file)
	if err != nil {
		log.Fatalf("logging.Setup, fail to create '%s': %v", logName, err)
	}
	mv := io.MultiWriter(os.Stdout, fileOutput)
	GlobalLogger.SetOutput(mv)
}

