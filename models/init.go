package models

import (
	"ChatDanBackend/config"
	"ChatDanBackend/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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
	switch config.Config.DbType {
	case "mysql":
		if config.Config.DbUrl == "" {
			panic("mysql db url required")
		}
		DB, err = gorm.Open(mysql.Open(config.Config.DbUrl), gormConfig)
	case "sqlite":
		if config.Config.DbUrl == "" {
			config.Config.DbUrl = "data.db"
		}
		DB, err = gorm.Open(sqlite.Open(config.Config.DbUrl), gormConfig)
	case "memory":
		DB, err = gorm.Open(mysql.Open(":memory:"), gormConfig)
	default:
		panic("unknown db type")
	}
	if err != nil {
		panic(err)
	}

	if config.Config.Debug {
		DB = DB.Debug()
	}

	if err = DB.SetupJoinTable(User{}, "Followers", &UserFollows{}); err != nil {
		panic(err)
	}
	if err = DB.SetupJoinTable(Topic{}, "LikedUsers", &TopicUserLikes{}); err != nil {
		panic(err)
	}
	if err = DB.SetupJoinTable(Topic{}, "FavoredUsers", &TopicUserFavorites{}); err != nil {
		panic(err)
	}
	if err = DB.SetupJoinTable(Topic{}, "ViewedUsers", &TopicUserViews{}); err != nil {
		panic(err)
	}
	if err = DB.SetupJoinTable(Comment{}, "LikedUsers", &CommentUserLikes{}); err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(
		User{},
		Box{},
		Post{},
		Channel{},
		Wall{},
		Division{},
		Topic{},
		Comment{},
		Tag{},
	)
	if err != nil {
		panic(err)
	}

	if config.Config.Standalone {
		err = DB.AutoMigrate(UserJwtSecret{})
	}
	if err != nil {
		panic(err)
	}

	utils.Logger.Info("database connected")
}
