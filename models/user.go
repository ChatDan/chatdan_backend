package models

import (
	"chatdan_backend/config"
	"chatdan_backend/utils"
	"fmt"
	"github.com/juju/errors"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"time"
)

type User struct {
	// 元数据
	ID             int            `json:"id"`
	Username       string         `json:"username" gorm:"index,size:100"`
	Email          *string        `json:"email" gorm:"index"` // 邮箱，可选登录
	HashedPassword string         `json:"hashed_password" gorm:"size:256"`
	LoginTime      time.Time      `json:"login_time" gorm:"autoUpdateTime"`
	RegisterTime   time.Time      `json:"register_time" gorm:"autoCreateTime"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Banned         bool           `json:"banned"`
	IsAdmin        bool           `json:"is_admin"`
	Avatar         *string        `json:"avatar" gorm:"size:256"`       // 头像链接
	Introduction   *string        `json:"introduction" gorm:"size:256"` // 个人简介/个性签名

	// 关联数据
	UserJwtSecret  *UserJwtSecret `json:"-" gorm:"foreignKey:UserID"`
	ViewedTopics   []*Topic       `json:"viewed_topics" gorm:"many2many:topic_user_views"`       // 浏览过的话题
	FavoriteTopics []*Topic       `json:"favorite_topics" gorm:"many2many:topic_user_favorites"` // 收藏的话题
	Followers      []*User        `json:"followed_users" gorm:"many2many:user_followers"`        // 关注的用户

	// 统计数据
	TopicCount          int `json:"topic_count" gorm:"not null;default:0"`           // 发表的话题数
	CommentCount        int `json:"comment_count" gorm:"not null;default:0"`         // 发表的评论数
	FavoriteTopicsCount int `json:"favorite_topics_count" gorm:"not null;default:0"` // 收藏的话题数
	FollowersCount      int `json:"followers_count" gorm:"not null;default:0"`       // 被关注数
	FollowingUsersCount int `json:"following_users_count" gorm:"not null;default:0"` // 关注数
}

func (user User) GetID() int {
	return user.ID
}

func (User) TableName() string {
	return "user"
}

func (user User) DeletedUsername() string {
	if user.DeletedAt.Valid {
		return user.Username
	} else {
		return fmt.Sprintf("%s_d_%d", user.Username, time.Now().Unix())
	}
}

type UserFollows struct {
	UserID     int       `json:"user_id" gorm:"primaryKey"`
	FollowerID int       `json:"follower_id" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserJwtSecret struct {
	UserID int    `json:"id" gorm:"primaryKey"`
	Secret string `json:"secret" gorm:"size:256"`
}

func (u UserJwtSecret) GetID() int {
	return u.UserID
}

func (u UserJwtSecret) TableName() string {
	return "user_jwt_secret"
}

func CreateJwtToken(user *User) (string, error) {
	if config.Config.Standalone {
		// no gateway, store jwt secret in database
		var userJwtSecret UserJwtSecret
		err := DB.Take(&userJwtSecret, user.ID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userJwtSecret = UserJwtSecret{
					UserID: user.ID,
					Secret: randstr.Base62(32),
				}
				err = DB.Create(&userJwtSecret).Error
				if err != nil {
					return "", err
				}
			} else {
				return "", err
			}
		}

		userClaims := utils.UserClaims{
			UserID:  user.ID,
			IsAdmin: user.IsAdmin,
			Key:     "",
		}

		return utils.CreateJwtTokenStandalone(userClaims, []byte(userJwtSecret.Secret))
	} else {
		// gateway, store jwt secret in etcd
		userClaims := utils.UserClaims{
			UserID:  user.ID,
			IsAdmin: user.IsAdmin,
		}

		return utils.CreateJwtTokenFromApisix(userClaims)
	}
}

func DeleteJwtToken(user *User) error {
	if config.Config.Standalone {
		// no gateway, delete jwt secret from database
		return DB.Delete(&UserJwtSecret{}, user.ID).Error
	} else {
		// gateway, delete jwt secret from etcd
		return utils.DeleteJwtTokenFromApisix(user.ID)
	}
}
