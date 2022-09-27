package mysql

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBCONFIG struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type DB struct {
	*gorm.DB
	NotFindError error
}

func Init(cfg *DBCONFIG) *DB {
	dsn := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ")/" + cfg.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}}); err != nil {
		log.Fatalln(err.Error())
		return nil
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalln(err.Error())
			return nil
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
		return &DB{db, gorm.ErrRecordNotFound}
	}
}
