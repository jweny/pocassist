package db

import (
	"fmt"
	"github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/file"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"path/filepath"
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
		if dbConfig.EnableDefault {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatalf("db.Setup, fail to get current path: %v", err)
			}
			//配置文件路径 当前文件夹 + config.yaml
			defaultSqliteFile := path.Join(dir, "pocassist.db")
			// 检测 sqlite 文件是否存在
			if !file.Exists(defaultSqliteFile) {
				log.Fatalf("db.Setup err: pocassist.db not exist, download at https://github.com/jweny/pocassistdb/releases")
			}
		}

		GlobalDB, err = gorm.Open(sqlite.Open(dbConfig.Sqlite), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
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

