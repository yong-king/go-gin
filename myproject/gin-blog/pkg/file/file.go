package file

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

// 获取文件长度
func GetSize(f multipart.File) (int, error) {
	coutent, err := ioutil.ReadAll(f)
	return len(coutent), err
}

// 获取文件后缀
func GetExt(fileName string) string {
	ext := path.Ext(fileName)
	return ext
}

// 检测是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// 检测权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err) // 判断是否是权限错误
}

// 不存在就创建
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// 创建
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func MustOpen(fileName string, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	fmt.Println("dir:", dir)
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	fmt.Println("prem:", perm)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExitMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("file.Open src: %s, err: %v", src, err)
	}
	return f, nil
}
