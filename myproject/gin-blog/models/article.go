package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/youngking/gin-blog/pkg/logging"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Desc       string `json:"desc"`
	State      int    `json:"stae"`
}

// 根据ID判断是否存在
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound { //允许 找不到记录，但如果是其他错误（如数据库连接错误），就返回错误信息。
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 获取单篇文章
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Preload("Tag").Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

// 获取多篇文章
func GetArticles(pageNum int, pageSize int, maps map[string]interface{}) ([]*Article, error) {
	var articles []*Article

	// 计算 offset，并设置默认分页大小
	offset := (pageNum - 1) * pageSize

	err := db.Where(maps).Offset(offset).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error("查询文章失败:", err)
		return nil, err
	}

	logging.Info("查询到的文章数量:", len(articles))
	return articles, nil
}

// 多篇文章的数量
func GetArticleTotal(maps map[string]interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil

}

// 新增文章
func AddArticle(data map[string]interface{}) error {
	content := ""
	if val, ok := data["content"].(string); ok {
		content = val
	}
	state := -1
	if arg, ok := data["state"].(int); ok {
		state = arg
	}
	desc := ""
	if arg, ok := data["desc"].(string); ok {
		desc = arg
	}
	err := db.Create(&Article{
		Title:     data["title"].(string),
		Content:   content,
		CreatedBy: data["created_by"].(string),
		TagID:     data["tag_id"].(int),
		State:     state,
		Desc:      desc,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// 修改文章
func EditArticle(id int, data map[string]interface{}) error {
	err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// 删除文章
func DeleteArticle(id int) error {
	err := db.Where("id = ?", id).Delete(&Article{}).Error
	if err != nil {
		return err
	}
	return nil
}

//// 创建时间
//func (article *Article) BerofeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//// 修改时间
//func (article *Article) BerofeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//	return nil
//}

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})
	return true
}

// 检查文章标题是否存在
func ArticleExistsByTitle(title string) (bool, error) {
	var count int
	err := db.Model(&Article{}).Where("LOWER(TRIM(title)) = LOWER(TRIM(?))", title).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, fmt.Errorf("查询文章是否存在出错: %v", err)
	}
	return count > 0, nil
}
