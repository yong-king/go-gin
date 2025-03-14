package article_service

import (
	"github.com/golang/freetype"
	"github.com/youngking/gin-blog/pkg/file"
	"github.com/youngking/gin-blog/pkg/qrcode"
	"github.com/youngking/gin-blog/pkg/setting"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

// 文章海报结构体
type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrCode
}

// 创建二维码海报实例函数
func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

// 获取海报的前缀
func GetPosterFlag() string {
	return "poster"
}

// 检查合成的海报是否存在
func (a *ArticlePoster) CheckMergeImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}
	return true
}

// 打开合成的海报文件
func (a *ArticlePoster) OpenMergeImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 背景图片
type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

// 二维码偏移
type Pt struct {
	X int
	Y int
}

// 海报扩展，包含背景
type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

func NewArticlePosterBg(name string, articlePoster *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: articlePoster,
		Rect:          rect,
		Pt:            pt,
	}
}

// 生成合成图片
func (a *ArticlePosterBg) Generate() (string, string, error) {
	// 获取文件路径
	fullPath := qrcode.GetQrCodeFullPath()
	// 获取二维码文件名和路径
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	// 检查是否已经有合成的图片
	// 没有
	if !a.CheckMergeImage(path) {
		// 打开合成图片文件
		mergedF, err := a.OpenMergeImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		// 打开背景图片文件
		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		// 打开二维码图片文件
		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		// 解码背景图片和二维码图片
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}

		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		// 创建jpg图片
		jpg := image.NewRGBA(image.Rect(a.X0, a.Y0, a.X1, a.Y1))
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.X, a.Y)), draw.Over)

		err = a.DrawPoster(&DrawText{
			JPG:     jpg,
			mergedF: mergedF,
			Title:   "Golang",
			X0:      80,
			Y0:      160,
			Size0:   42,

			SubTitle: "---YoungKing",
			X1:       320,
			Y1:       220,
			Size1:    36,
		}, "msyhbd.ttc")
		if err != nil {
			return "", "", err
		}
	}
	return fileName, path, nil
}

// 加文字结构体
type DrawText struct {
	JPG     draw.Image
	mergedF *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	// 获取字体路径
	fontScorce := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + fontName
	// 读取字体文件
	fontBytes, err := os.ReadFile(fontScorce)
	if err != nil {
		return err
	}
	// 解析字体，获取字体轮廓等
	trueTypeFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// 创建一个freetype.Context对象，所有的字体绘制操作都将基于这个上下文进行。
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetFontSize(d.Size0)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.Black)
	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.mergedF, d.JPG, nil)
	if err != nil {
		return err
	}
	return nil
}
