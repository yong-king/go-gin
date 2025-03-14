package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type Leave int

var (
	F                  *os.File // 文件指针变量，用来存储
	DefaultPrefix      = ""     // 默认日志前缀
	DefaultCallerDepth = 2      // 默认调用深度
	logger             *log.Logger
	logPrefix          = "" // 当前日志前缀
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Leave = iota
	INFO
	WARN
	ERROR
	FATAL
)

func Setup() {
	// 获取文件路径
	filePath := getLogFilePath()
	fileName := getLogFilename()
	F, err := openLogFile(filePath, fileName)
	if err != nil {
		log.Fatalln(err)
	}
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	SetPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	SetPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	SetPrefix(WARN)
	logger.Println(v)
}

func Error(v ...interface{}) {
	SetPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	SetPrefix(FATAL)
	logger.Fatalln(v)
}

func SetPrefix(level Leave) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d] ", levelFlags[level], file, line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
