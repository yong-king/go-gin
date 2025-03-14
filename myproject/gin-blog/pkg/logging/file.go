package logging

import (
	"fmt"
	"github.com/youngking/gin-blog/pkg/file"
	"github.com/youngking/gin-blog/pkg/setting"
	"log"
	"os"
	"time"
)

//var (
//	LogSavePath = "runtime/logs/"
//	LogSaveName = "log"
//	LogFileExt  = "log"
//	TimeFormat  = "2006-01-02"
//)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

//	func getLogFileFullPath() string {
//		return LogSavePath + LogSaveName + time.Now().Format(TimeFormat) + "." + LogFileExt
//	}
func getLogFilename() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt)
}

func openLogFile(filePath string, fileName string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExitMkDir err: %v", err)
	}
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("file.Open err: %v", err)
	}
	return f, nil
}

func mkDir(path string) {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+path, os.ModePerm)
	if err != nil {
		log.Fatalf("Fail to create dir :%v!", err)
	}
}
