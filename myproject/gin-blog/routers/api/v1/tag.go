package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/youngking/gin-blog/pkg/app"
	"github.com/youngking/gin-blog/pkg/e"
	"github.com/youngking/gin-blog/pkg/export"
	"github.com/youngking/gin-blog/pkg/logging"
	"github.com/youngking/gin-blog/pkg/setting"
	utills "github.com/youngking/gin-blog/pkg/util"
	"github.com/youngking/gin-blog/service/tag_service"
	"net/http"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {

	appG := app.Gin{c}
	// 获取查询条件
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  utills.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
	}
	data := make(map[string]interface{})
	data["lists"] = tags
	data["total"] = count

	appG.Response(http.StatusOK, e.SUCCESS, data)

}

type AddTagForm struct {
	Name      string `from:"name" validate:"Required;MaxSize(100)"`
	State     int    `from:"state" validate:"Range(0,1)"`
	CreatedBy string `from:"created_by" valid:"Required;MaxSize(100)"`
}

// @Summary Add article tag
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {

	var (
		appG = app.Gin{c}
		form AddTagForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		State:     form.State,
		CreatedBy: form.CreatedBy,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `from:"id" validate:"Required;MinSize(1)"`
	Name       string `from:"name" validate:"MaxSize(100)"`
	State      int    `from:"state" validate:"Range(0,1)"`
	ModifiedBy string `from:"modified_by" validate:"MaxSize(100)"`
}

// @Summary 修改文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{c}
		form EditTagForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		State:      form.State,
		ModifiedBy: form.ModifiedBy,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete article tag
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	var appG = app.Gin{c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}

	tagService := tag_service.Tag{
		ID: id,
	}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ExportTag(c *gin.Context) {
	appG := app.Gin{c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}
	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_Save_url": export.GetExcelPath() + filename,
	})
}

func ImportTag(c *gin.Context) {
	// 上下文结构体赋值
	appG := app.Gin{c}

	// 获取http请求中的file上传请求
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn("文件上传失败", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	// 处理请求
	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn("导入标签失败", err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}
	//logging.Info("Excel 导入成功")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
