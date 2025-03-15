package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/youngking/gin-blog/pkg/app"
	"github.com/youngking/gin-blog/pkg/e"
	"github.com/youngking/gin-blog/pkg/export"
	"github.com/youngking/gin-blog/pkg/qrcode"
	"github.com/youngking/gin-blog/pkg/setting"
	utills "github.com/youngking/gin-blog/pkg/util"
	"github.com/youngking/gin-blog/service/article_service"
	"github.com/youngking/gin-blog/service/tag_service"
	"net/http"
)

// @Summary Get a signal article
// @Produce json
// @Param id ptah int true "ID"
// @Success 200 {object} app.Response
// @Failuer 500 {object} app.Response
// @Router /api/v1/articles/{id}[get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	// 解析请求参数
	id := com.StrTo(c.Param("id")).MustInt()

	// 验证参数
	val := validation.Validation{}
	val.Min(id, 1, "id").Message("ID必须大于0")

	if val.HasErrors() {
		app.MarkErrors(val.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	articleService := article_service.Article{ID: id}
	exits, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_ARTICLE_FILE, nil)
	}
	if !exits {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)

}

// @Summary Get multiple articles
// @Produce json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failuer 500 {object} app.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagID := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagID = com.StrTo(arg).MustInt()
		valid.Min(tagID, 1, "tag_id")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	articleService := article_service.Article{
		TagID:    tagID,
		State:    state,
		PageNum:  utills.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddArticleFrom struct {
	TagID         int    `form:"tag_id"  valid:"Required;Min(1)"`
	Title         string `form:"title"  valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc"  valid:"Required;MaxSize(255)"`
	Content       string `form:"content"  valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by"  valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url"  valid:"Required;MaxSize(255)"`
	State         int    `form:"state"  valid:"Range(0,1)"`
}

// @Summary Add article
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	var (
		appG = app.Gin{c}
		form AddArticleFrom
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Dest:          form.Desc,
		Content:       form.Content,
		CreatedBy:     form.CreatedBy,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// @Summary Edit article
// @Produce json
// @Param id path int true "ID"
// @Param tag_id body int true "TagID"
// @Param modefied_by body string true "ModefiedBy
// @Param desc body string true "Dese"
// @Param tite body string true "Title
// @Param state body int true "State"
// @Param cover_image_url body string true "CoverImageUrl"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [put]
type EditArticleFrom struct {
	ID            int    `form:"id"  valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id"  valid:"Required;Min(1)"`
	Title         string `form:"title"  valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc"  valid:"Required;MaxSize(255)"`
	Content       string `form:"content"  valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by"  valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url"  valid:"Required;MaxSize(255)"`
	State         int    `form:"state"  valid:"Range(0,1)"`
}

func EditArticle(c *gin.Context) {

	var (
		appG = app.Gin{c}
		form EditArticleFrom
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Dest:          form.Desc,
		Content:       form.Content,
		ModifiedBy:    form.ModifiedBy,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	exist, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete article
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exist, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_ARTICLE_FILE, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ExportArticle(c *gin.Context) {
	// 初始化结构
	appG := app.Gin{c}

	// 获取请求
	tagID := com.StrTo(c.PostForm("tag_id")).MustInt()
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	// 结构体赋值
	articleService := article_service.Article{
		TagID: tagID,
		State: state,
	}

	fileName, err := articleService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROE_EXPORT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(fileName),
		"export_save_url": export.GetExcelFullPath() + fileName,
	})

}

func ImportArticle(c *gin.Context) {
	appG := app.Gin{c}

	// 获取http请求
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	defer file.Close()

	// 处理请求
	articleService := article_service.Article{}
	err = articleService.Import(file)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_IMPORT_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

const (
	QRCODE_URL = "https://github.com/yong-king/go-gin"
)

func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{c}
	// 创建二维码实例
	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	// 海报名称
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetFileName(qr.URL) + qr.GetFileExt()
	//fmt.Print(posterName)
	article := &article_service.Article{}
	// 创建海报实例
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)

	// 创建背景海报实例
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpeg",
		articlePoster,
		&article_service.Rect{X0: 0, Y0: 0, X1: 550, Y1: 700},
		&article_service.Pt{X: 125, Y: 298})
	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}
	//fmt.Println("filePath:", filePath)
	//fmt.Println("qrcode.GetQrCodeFullUrl:", qrcode.GetQrCodeFullUrl(filePath))
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(filePath) + posterName,
		"poster_save_url": filePath + posterName,
	})
}
