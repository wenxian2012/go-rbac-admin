package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wenxian2012/go-rbac-admin/pkg/gormkit"
	"github.com/wenxian2012/go-rbac-admin/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID        int                `gorm:"primary_key" json:"id"`
	CreatedAt gormkit.LocalTime  `json:"created_at"`
	UpdatedAt gormkit.LocalTime  `json:"updated_at"`
	DeletedAt *gormkit.LocalTime `json:"deleted_at"`
}

func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", setting.DatabaseSetting.TablePrefix, name)
}

func init() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		tableName := defaultTableName
		if !strings.HasPrefix(defaultTableName, setting.DatabaseSetting.TablePrefix) {
			tableName = fmt.Sprintf("%s%s", setting.DatabaseSetting.TablePrefix, defaultTableName)
		}
		return tableName
	}
	if setting.DatabaseSetting.Debug {
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}
