package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/youngking/gin-blog/pkg/setting"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID         int    `gorm:"primary_key" json:"id"`
	CreatedOn  int    `json:"created_on"`
	ModifiedOn int    `json:"modified_on"`
	DeletedOn  *int64 `json:"deleted_on"`
}

func Setup() {
	//var (
	//	err                                               error
	//	dbType, password, dbName, host, user, tablePrefix string
	//)
	//
	//// 获取数据库配置
	//sec, err := setting.Cfg.GetSection("database")
	//if err != nil {
	//	log.Fatal(2, "Fail to get section 'database': %v", err)
	//}
	//dbType = sec.Key("TYPE").String()
	//dbName = sec.Key("NAME").String()
	//host = sec.Key("HOST").String()
	//user = sec.Key("USER").String()
	//password = sec.Key("PASSWORD").String()
	//tablePrefix = sec.Key("TABLE_PREFIX").String()

	var err error
	// 连接数据库
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	if err != nil {
		log.Println(err)
	}

	// 设置表前缀名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defultTableName
	}

	// 设置数据库配置
	db.SingularTable(true) // 默认 GORM 会自动给表名加 "s"（复数形式），这里关闭这个功能
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStreamForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStreamForUpdate)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetConnMaxIdleTime(10) // 最大空闲连接数
	db.DB().SetMaxOpenConns(100)   // 最大连接数

}

// 关闭数据库
func CloseDB() {
	defer db.Close()
}

// 定制gorm的callback
func updateTimeStreamForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeFiled, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeFiled.IsBlank {
				createTimeFiled.Set(nowTime)
			}
		}

		if modifyTimeFiled, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeFiled.IsBlank {
				modifyTimeFiled.Set(nowTime)
			}
		}
	}
}

func updateTimeStreamForUpdate(scope *gorm.Scope) {
	if !scope.HasError() {
		if _, ok := scope.Get("grom:update_column"); !ok {
			scope.SetColumn("ModifiedOn", time.Now().Unix())
		}
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_on"); ok {
			fmt.Sprint(str)
		}
		deletedOn, hasDeletedOnField := scope.FieldByName("DeletedOn")
		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOn.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption))).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.Quote(deletedOn.DBName)),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(s string) string {
	if s != "" {
		return " " + s
	}
	return ""
}
