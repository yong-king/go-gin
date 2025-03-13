package article_service

import (
	"encoding/json"
	"github.com/unknwon/com"
	"github.com/xuri/excelize/v2"
	"github.com/youngking/gin-blog/models"
	"github.com/youngking/gin-blog/pkg/export"
	"github.com/youngking/gin-blog/pkg/gredis"
	"github.com/youngking/gin-blog/pkg/logging"
	"github.com/youngking/gin-blog/service/cache_service"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Dest          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Get() (*models.Article, error) {
	var articleModel *models.Article
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()

	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &articleModel)
			return articleModel, nil
		}
	}
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"title":           a.Title,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"desc":            a.Dest,
		"created_by":      a.CreatedBy,
		"state":           a.State,
		"tag_id":          a.TagID,
	}

	err := models.AddArticle(article)
	if err != nil {
		return err
	}
	return nil
}

func (a *Article) Edit() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"desc":            a.Dest,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	}

	return models.EditArticle(a.ID, article)
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var articles, cacheArticles []*models.Article
	cache := cache_service.Article{
		ID:    a.ID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			if len(cacheArticles) > 0 {
				return cacheArticles, nil // 直接返回缓存数据，避免后续查询失败
			}
		}
	}

	// 确保分页参数有效，防止 0 值影响查询
	pageNum := a.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}

	pageSize := a.PageSize
	if pageSize <= 0 {
		pageSize = math.MaxInt32 // 查询所有数据
	}

	articles, err := models.GetArticles(pageNum, pageSize, a.getMaps())
	if err != nil {
		logging.Info(err)
		return nil, err // 直接返回错误，避免后续代码执行
	}
	//logging.Info("数据库查询结果:", articles) // ✅ 打印数据库返回的数据
	if len(articles) > 0 {
		gredis.Set(key, articles, 3600)
	}
	return articles, nil
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}

func (a *Article) Export() (string, error) {
	// 获取所有文章
	articles, err := a.GetAll()
	//if err != nil {
	//	return "", err
	//}
	//fmt.Printf("数据库查询结果: %+v\n", articles)

	if err != nil {
		return "", err
	}

	// 新建一个表用于接受数据
	file := excelize.NewFile()
	defer file.Close()
	// 新建工作表
	indexName := "文章信息"
	index, err := file.NewSheet(indexName)
	if err != nil {
		return "", err
	}

	// 表头
	title := []string{"ID", "tag_id", "Title", "Dest", "Content", "ModifiedBy", "ModifiedOn", "CreatedBy", "CreatedOn"}
	for i, value := range title {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(indexName, cell, value)
	}

	for rowIndex, v := range articles {
		rowIndex += 2
		values := []string{
			strconv.Itoa(v.ID),
			strconv.Itoa(v.TagID),
			v.Title,
			v.Desc,
			v.Content,
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
		}
		for colIndex, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
			file.SetCellValue(indexName, cell, value)
		}
	}
	// 设置工作表为默认工作表
	file.SetActiveSheet(index)

	// 保存文件
	time := strconv.Itoa(int(time.Now().Unix()))
	fileName := "atricles-" + time + ".xlsx"
	fullPath := export.GetExcelFullPath() + "/" + fileName
	err = file.SaveAs(fullPath)
	if err != nil {
		return "", err
	}
	return fileName, nil

}

func (a *Article) Import(r io.Reader) error {
	// 读取excle文件
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	// 获得表单
	index := "文章信息"
	rows, err := xlsx.GetRows(index)
	if err != nil {
		return err
	}

	// 维护一个map用来保证不重复
	existArticles := make(map[string]bool)
	articls, err := a.GetAll()
	if err != nil {
		return err
	}
	for _, article := range articls {
		normalizedTitle := strings.TrimSpace(strings.ToLower(article.Title)) // 统一格式
		existArticles[normalizedTitle] = true
	}

	for i, row := range rows {
		// 跳过表头
		if i == 0 {
			continue
		}

		// 内容不完整
		if len(row) < 7 {
			continue
		}
		tag_id, title, state, created_by := row[1], row[2], row[4], row[5]
		content, desc := row[3], row[6]

		// 统一 `title` 进行检查，避免大小写和空格问题
		normalizedTitle := strings.TrimSpace(strings.ToLower(title))

		// 标题不能为空
		if normalizedTitle == "" {
			logging.Warn("Excel 第 %d 行标题为空，跳过", i+1)
			continue
		}

		// 检查是否存在
		if existArticles[normalizedTitle] {
			logging.Info("文章《%s》已存在，跳过导入", title)
			continue
		}

		// `state` 和 `tag_id` 转换
		tagID := com.StrTo(tag_id).MustInt()
		stateVal := com.StrTo(state).MustInt()

		// 避免 `nil` 赋值导致 `panic`
		if content == "" {
			content = "默认内容"
		}
		if desc == "" {
			desc = "暂无描述"
		}

		// **数据库中再次检查**
		exists, err := models.ArticleExistsByTitle(normalizedTitle)
		if err != nil {
			logging.Error("检查文章《%s》是否存在时出错: %v", title, err)
			continue
		}
		if exists {
			logging.Error("数据库中已存在文章《%s》，跳过导入", title)
			continue
		}

		// 添加文章
		err = models.AddArticle(map[string]interface{}{
			"title":      title,
			"tag_id":     tagID,
			"created_by": created_by,
			"content":    content,
			"desc":       desc,
			"state":      stateVal,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
