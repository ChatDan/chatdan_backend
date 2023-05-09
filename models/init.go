package models

import (
	"ChatDanBackend/config"
	"ChatDanBackend/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var DB *gorm.DB

var LockClause = clause.Locking{Strength: "UPDATE"}

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	},
	Logger: logger.New(
		utils.StdOutLogger,
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.Config.DbUrl), gormConfig)
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(
		User{},
		Box{},
		Post{},
		Channel{},
		Wall{},
	)
	if err != nil {
		panic(err)
	}

	if config.Config.Debug {
		DB = DB.Debug()
	}

	utils.Logger.Info("database connected")
}
