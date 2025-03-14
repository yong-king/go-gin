package qrcode

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/youngking/gin-blog/pkg/file"
	"github.com/youngking/gin-blog/pkg/setting"
	utills "github.com/youngking/gin-blog/pkg/util"
	"image/jpeg"
)

type QrCode struct {
	URL    string
	Height int
	Width  int
	Ext    string
	Mode   qr.Encoding
	Level  qr.ErrorCorrectionLevel
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

// 获得相对路径
func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

// 获得绝对路径
func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

// 获取URL访问路径
func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + name
}

// 生成文件名
func GetFileName(value string) string {
	return utills.EncodeMD5(value)
}

// 获取文件后缀名
func (q *QrCode) GetFileExt() string {
	return q.Ext
}

// 检测是否存在
func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetFileName(q.URL) + q.GetFileExt()
	if file.CheckNotExist(src) == true {
		return false
	}
	return true
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetFileName(q.URL) + q.GetFileExt()
	src := path + name
	//fmt.Println(src)
	if file.CheckNotExist(src) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}
		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}
		//fmt.Println(code)

		file, err := file.MustOpen(name, path)
		//fmt.Println(file)
		if err != nil {
			return "", "", err
		}
		defer file.Close()

		err = jpeg.Encode(file, code, nil)
		if err != nil {
			return "", "", err
		}

	}
	return name, path, nil
}
