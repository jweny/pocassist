package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"pocassist/basic"
)

var GlobalDB *gorm.DB

func InitDB(dbname string) error {
	if dbname != "mysql" && dbname != "sqlite" {
		return errors.New("unsupported database kind. only 'sqlite' or 'mysql'")
	}
	var err error
	dbConfig := basic.GlobalConfig.DbConfig
	if dbname == "mysql" {
		// 配置mysql据源
		if dbConfig.Mysql.User == "" || dbConfig.Mysql.Password == "" || dbConfig.Mysql.Host == "" ||
			dbConfig.Mysql.Port == "" || dbConfig.Mysql.Database == "" {
			return errors.New("config.yaml mysql not set completed")
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
			dbConfig.Mysql.User,
			dbConfig.Mysql.Password,
			dbConfig.Mysql.Host,
			dbConfig.Mysql.Port,
			dbConfig.Mysql.Database,
			dbConfig.Mysql.Timeout)

		GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			basic.GlobalLogger.Error("[mysql db connect err ]", err)
			return err
		}
	}
	if dbname == "sqlite" {
		// 配置sqlite数据源
		if dbConfig.Sqlite == "" {
			return errors.New("config.yaml sqlite not set")
		}
		GlobalDB, err = gorm.Open(sqlite.Open(dbConfig.Sqlite), &gorm.Config{})
		if err != nil {
			basic.GlobalLogger.Error("[sqlite db connect err ]", err)
			return err
		}
	}
	if GlobalDB == nil {
		basic.GlobalLogger.Error("[db connect err ]", err)
		return errors.New("db connect err")
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

