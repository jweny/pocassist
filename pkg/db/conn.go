package db

import (
	"fmt"
	"github.com/jweny/pocassist/pkg/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var GlobalDB *gorm.DB

func Setup() {
	var err error
	dbConfig := conf.GlobalConfig.DbConfig
	if conf.GlobalConfig.DbConfig.Sqlite == "" {
		// 配置mysql数据源
		if dbConfig.Mysql.User == "" ||
			dbConfig.Mysql.Password == "" ||
			dbConfig.Mysql.Host == "" ||
			dbConfig.Mysql.Port == "" ||
			dbConfig.Mysql.Database == "" {
			log.Fatalf("db.Setup err: config.yaml mysql config not set")
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
			log.Fatalf("db.Setup err: %v", err)
		}
	} else {
		// 配置sqlite数据源
		if dbConfig.Sqlite == "" {
			log.Fatalf("db.Setup err: config.yaml sqlite config not set")
		}
		GlobalDB, err = gorm.Open(sqlite.Open(dbConfig.Sqlite), &gorm.Config{})
		if err != nil {
			log.Fatalf("db.Setup err: %v", err)
		}
	}

	if GlobalDB == nil {
		log.Fatalf("db.Setup err: db connect failed")
	}

	err = GlobalDB.AutoMigrate(&Auth{}, &Vulnerability{}, &Webapp{}, &Plugin{}, &Task{}, &Result{})

	if err != nil {
		log.Fatalf("db.Setup err: %v", err)
	}

	if conf.GlobalConfig.ServerConfig.RunMode == "release" {
		// release下
		GlobalDB.Logger = logger.Default.LogMode(logger.Silent)
	}
}

