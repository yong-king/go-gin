package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/youngking/gin-blog/middleware/jwt"
	"github.com/youngking/gin-blog/pkg/export"
	"github.com/youngking/gin-blog/pkg/setting"
	"github.com/youngking/gin-blog/pkg/upload"
	"github.com/youngking/gin-blog/routers/api"
	v1 "github.com/youngking/gin-blog/routers/api/v1"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServeSetting.RunMode)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//fmt.Println("图片存储路径:", upload.GetImageFullPath())
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/auth", api.GetAuth)
	r.POST("upload", api.UploadImage)
	r.POST("/tags/export", v1.ExportTag)
	r.POST("article/export", v1.ExportArticle)

	//router.GET("/test", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "test",
	//	})
	//})

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		// 新增标签列表
		apiv1.POST("/tags", v1.AddTag)
		// 更新标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		// 删除标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 获取单篇文章
		apiv1.GET("/articles/id", v1.GetArticle)
		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 更新文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		// 删除文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		// 导入标签
		r.POST("/tags/import", v1.ImportTag)
		r.POST("/article/import", v1.ImportArticle)

		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	return r
}
