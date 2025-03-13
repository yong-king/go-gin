package upload

import (
	"fmt"
	"github.com/youngking/gin-blog/pkg/file"
	"github.com/youngking/gin-blog/pkg/logging"
	"github.com/youngking/gin-blog/pkg/setting"
	"github.com/youngking/gin-blog/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utills.EncodeMD5(fileName)
	return fileName + ext
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(name string) bool {
	ext := file.GetExt(name)
	for _, e := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(ext) == strings.ToUpper(e) {
			return true
		}
	}
	logging.Info("ext:", ext)
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	//logging.Info("size:", size)
	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExitMkDir err: %v", err)
	}

	permission := file.CheckPermission(src)
	if permission == true {
		return fmt.Errorf("file.CheckPermission err: %v", err)
	}

	return nil

}
