package tag_service

import (
	"encoding/json"
	"github.com/xuri/excelize/v2"
	"github.com/youngking/gin-blog/models"
	"github.com/youngking/gin-blog/pkg/export"
	"github.com/youngking/gin-blog/pkg/gredis"
	"github.com/youngking/gin-blog/pkg/logging"
	"github.com/youngking/gin-blog/service/cache_service"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tags{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()
	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tags, 3600)
	//fmt.Printf("service %v", tags)
	return tags, nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	data["modified_by"] = t.ModifiedBy
	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Export() (string, error) {
	tag, err := t.GetAll()
	if err != nil {
		return "", err
	}

	// 创建新的excel文件
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			logging.Error(err)
		}
	}()

	// 创建新的工作表
	sheetName := "标签信息"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		logging.Warn(err)
		return "", err
	}

	// 设置表头
	title := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	for i, value := range title {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // 从 A1 开始
		file.SetCellValue(sheetName, cell, value)
	}

	// 数据填充
	for rowIndex, v := range tag {
		rowIndex += 2 // 从第二行开始填充数据
		valuse := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}
		for colIndex, value := range valuse {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
			file.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置工作表为默认工作表
	file.SetActiveSheet(index)

	// 保存文件
	time := strconv.Itoa(int(time.Now().Unix()))
	fileName := "tags-" + time + ".xlsx"
	fullPath := export.GetExcelFullPath() + fileName
	err = file.SaveAs(fullPath)
	if err != nil {
		logging.Warn(err)
		return "", err
	}
	return fileName, nil

	//fmt.Printf("Tags retrieved: %+v\n", tag)

	//file := xlsx.NewFile()
	//sheet, err := file.AddSheet("标签信息")
	//if err != nil {
	//	return "", err
	//}
	//
	//title := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	//row := sheet.AddRow()
	//
	//var cell *xlsx.Cell
	//for _, value := range title {
	//	cell = row.AddCell()
	//	cell.Value = value
	//}
	//
	//for _, v := range tag {
	//	values := []string{
	//		strconv.Itoa(v.ID),
	//		v.Name,
	//		v.CreatedBy,
	//		strconv.Itoa(v.CreatedOn),
	//		v.ModifiedBy,
	//		strconv.Itoa(v.ModifiedOn),
	//	}
	//
	//	row := sheet.AddRow()
	//	for _, value := range values {
	//		cell := row.AddCell()
	//		cell.Value = value
	//	}
	//}
	//
	//time := strconv.Itoa(int(time.Now().Unix()))
	//fileName := "tags-" + time + ".xlsx"
	//
	//fullPath := export.GetExcelFullPath() + fileName
	//err = file.Save(fullPath)
	//if err != nil {
	//	return "", err
	//}
	//return fileName, nil
}

func (t *Tag) Import(r io.Reader) error {
	// 读取xlsx文件
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	// 获得表单
	rows, err := xlsx.GetRows("标签信息")
	if err != nil {
		return err
	}

	// // 维护一个去重 Map，避免重复插入
	existingTags := make(map[string]bool)
	tags, err := t.GetAll()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		existingTags[tag.Name] = true
	}

	// 遍历excle记录（跳过表头）
	for i, row := range rows {
		// 跳过表头
		if i == 0 {
			continue
		}

		if len(row) < 3 {
			logging.Info("跳过不完整数据!")
			continue
		}

		name, createdBy := row[1], row[2]
		// 检查是否已存在
		if existingTags[name] {
			logging.Info("跳过重复标签: %s", name)
			continue
		}

		// 添加标签
		if err := models.AddTag(name, 1, createdBy); err != nil {
			logging.Warn("插入标签失败: %s, 错误: %v", name, err)
			continue
		}

		// 记录插入的标签，避免重复插入
		existingTags[name] = true

	}
	return nil
}
