package models

import (
	"github.com/jinzhu/gorm"
	"github.com/youngking/gin-blog/pkg/logging"
)

type Tag struct {
	Model
	Name       string `json:name`
	CreatedBy  string `json:created_by`
	ModifiedBy string `json:modified_by`
	State      int    `json:state`
	//DeletedOn  *int64 `json:"deleted_on"`
}

// 获取符合条件的文章标签
func GetTags(pageNum int, pageSize int, maps map[string]interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&tags).Offset(pageNum - 1).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

// 获取到的tag的数量
func GetTagTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 判断是否存在
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name=?", name).First(&tag).Error
	if err != nil {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 新增标签
func AddTag(name string, state int, createdBy string) error {
	deletedOn := int64(0) // 创建一个值为 0 的 int64
	tag := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
		Model:     Model{DeletedOn: &deletedOn},
	})
	err := tag.Error
	if err != nil {
		logging.Warn("数据库插入失败: ", err)
		return err
	}
	//logging.Info("数据插入成功: ", tag.Value)
	return nil
}

//// 添加之前添加创建时间
//func (t *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreaterOn", time.Now().Unix())
//	return nil
//}

//// 添加更新时间
//func (t *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModeifiedOn", time.Now().Unix())
//	return nil
//}

// 修改
func EditTag(id int, data map[string]interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// 删除
func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}
	return nil
}

// id是否存在
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id=?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	return tag.ID > 0, nil
}

func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})
	return true
}
