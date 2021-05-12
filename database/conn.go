package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"pocassist/basic"
)

var GlobalDB *gorm.DB

func InitDB() error {
	var err error
	dbConfig := basic.GlobalConfig.DbConfig

	// 配置数据源
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.Timeout)

	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		basic.GlobalLogger.Error("[db connect err ]", err)
		return err
	}

	err = GlobalDB.AutoMigrate(&Auth{}, &Vulnerability{}, &Webapp{},&Plugin{},)

	if err != nil {
		basic.GlobalLogger.Error("[db migrate err ]", err)
		return err
	}
	GlobalDB.Logger = logger.Default.LogMode(logger.Silent)
	basic.GlobalLogger.Debug("[db connect success ]")
	return nil
}

