package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// ListTopics godoc
// @Summary 查询话题，按照最近创建或最近回复排序
// @Tags Topic Module
// @Produce json
// @Router /topics [get]
// @Param json query TopicListRequest true "page"
// @Success 200 {object} RespForSwagger{data=TopicListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListTopics(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	var query TopicListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}
	nowTime := time.Now()
	if query.StartTime == nil {
		query.StartTime = &nowTime
	}

	var topics []Topic
	querySet := DB.Order(query.OrderBy+" desc").Limit(query.PageSize).
		Where("? < ?", clause.Column{Name: query.OrderBy}, query.StartTime)
	if query.DivisionID != nil {
		querySet = querySet.Where("division_id = ?", *query.DivisionID)
	}
	result := querySet.Preload("Tags").Preload("Poster").Find(&topics)
	if result.Error != nil {
		return result.Error
	}

	var response TopicListResponse
	if err = copier.CopyWithOption(&response.Topics, &topics, CopyOption); err != nil {
		return err
	}
	return Success(c, &response)
}

// GetATopic godoc
// @Summary 获取一个话题
// @Tags Topic Module
// @Produce json
// @Router /topic/{id} [get]
// @Param id path int true "topic id"
// @Success 200 {object} RespForSwagger{data=TopicCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func GetATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var topic Topic
	result := DB.Preload("Tags").Preload("Poster").First(&topic, id)
	if result.Error != nil {
		return result.Error
	}

	var response TopicCommonResponse
	if err = copier.CopyWithOption(&response, &topic, CopyOption); err != nil {
		return err
	}
	return Success(c, &response)
}

// CreateATopic godoc
// @Summary 创建一个话题
// @Tags Topic Module
// @Accept json
// @Produce json
// @Router /topic [post]
// @Param json body TopicCreateRequest true "topic"
// @Success 201 {object} RespForSwagger{data=TopicCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func CreateATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	var body TopicCreateRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var topic Topic
	if err = copier.CopyWithOption(&topic, &body, CopyOption); err != nil {
		return err
	}
	topic.PosterID = user.ID
	err = topic.FindOrCreateTags(DB, ToTagNames(body.Tags))
	if err != nil {
		return err
	}

	if topic.IsAnonymous {
		newAnonyname := GenerateName([]string{})
		topic.Anonyname = &newAnonyname
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		// Create topic
		err = tx.Omit(clause.Associations).Create(&topic).Error
		if err != nil {
			return err
		}
		// Create topic_tags association only
		err = tx.Omit("Tags.*", "UpdatedAt").Select("Tags").Save(&topic).Error
		if err != nil {
			return err
		}
		// Update tag temperature
		err = tx.Model(&topic.Tags).Update("temperature", gorm.Expr("temperature + 1")).Error
		if err != nil {
			return err
		}

		if topic.IsAnonymous {
			// Create topic_anonyname_mapping
			err = tx.Create(&TopicAnonynameMapping{
				TopicID:   topic.ID,
				UserID:    user.ID,
				Anonyname: *topic.Anonyname,
			}).Error
			if err != nil {
				return err
			}
		}

		result := tx.Model(&user).Update("topic_count", gorm.Expr("topic_count + 1"))
		if result.Error != nil {
			return result.Error
		}

		// save to meilisearch
		err = SearchAddOrReplace(topic.ToSearchModel())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	var response TopicCommonResponse
	if err = copier.CopyWithOption(&response, &topic, CopyOption); err != nil {
		return err
	}
	response.IsOwner = true

	return Created(c, &response)
}

// ModifyATopic godoc
// @Summary 修改一个话题
// @Description 管理员可修改标题、内容、标签、是否隐藏，用户可修改标题、内容、标签
// @Tags Topic Module
// @Accept json
// @Produce json
// @Router /topic/{id} [put]
// @Param json body TopicModifyRequest true "topic"
// @Success 200 {object} RespForSwagger{data=TopicCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ModifyATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var body TopicModifyRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if body.IsHidden != nil || body.DivisionID != nil {
		if !user.IsAdmin {
			return Forbidden()
		}
	}

	var topic Topic
	err = DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(LockClause).First(&topic, id)
		if result.Error != nil {
			return NotFound()
		}

		if !user.IsAdmin && user.ID != topic.PosterID {
			return Forbidden()
		}

		err = copier.CopyWithOption(&topic, &body, CopyOption)
		if err != nil {
			return err
		}

		err = tx.Model(&topic).Select(body.Fields()).UpdateColumns(&topic).Error
		if err != nil {
			return err
		}

		if body.Tags != nil {
			// clear associations
			err = tx.Model(&topic).Association("Tags").Clear()

			err = topic.FindOrCreateTags(tx, ToTagNames(body.Tags))
			if err != nil {
				return err
			}

			// Create topic_tags association only
			err = tx.Omit("Tags.*", "UpdatedAt").Select("Tags").Save(&topic).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// save to meilisearch
	err = SearchAddOrReplace(topic.ToSearchModel())
	if err != nil {
		return err
	}

	var response TopicCommonResponse
	if err = copier.CopyWithOption(&response, &topic, CopyOption); err != nil {
		return err
	}
	return Success(c, &response)
}

// DeleteATopic godoc
// @Summary 删除一个话题，仅管理员
// @Tags Topic Module
// @Produce json
// @Router /topic/{id} [delete]
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func DeleteATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var topic Topic
	result := DB.First(&topic, id)
	if result.Error != nil {
		return NotFound()
	}
	if !user.IsAdmin && topic.PosterID != user.ID {
		return Forbidden()
	}
	result = DB.Where("id = ?", id).Delete(&topic)
	if result.Error != nil {
		return result.Error
	}

	// delete from meilisearch
	err = SearchDelete[TagSearchModel](topic.ID)
	if err != nil {
		return err
	}

	return Success(c, &EmptyStruct{})
}

// LikeOrDislikeATopic godoc
// @Summary 点赞或点踩一个话题，或者重置点赞点踩数据
// @Description 1: like, -1: dislike, 0: reset，点赞点踩二选一
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_like/{like_data} [put]
// @Param id path int true "topic id"
// @Param like_data path int true "1: like, -1: dislike, 0: reset" Enums(1, -1, 0)
// @Success 200 {object} RespForSwagger{data=TopicCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func LikeOrDislikeATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	likeData, err := c.ParamsInt("data")
	if err != nil {
		return err
	}

	var topic Topic
	err = DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(LockClause).Preload("Tags").Preload("Poster").First(&topic, id)
		if result.Error != nil {
			return result.Error
		}

		var topicUserLikes = TopicUserLikes{
			UserID:   user.ID,
			TopicID:  id,
			LikeData: likeData,
		}
		result = tx.Save(&topicUserLikes)
		if result.Error != nil {
			return result.Error
		}

		var likeCount int64
		result = tx.Model(TopicUserLikes{}).Where("topic_id = ? and like_data = 1", id).Count(&likeCount)
		if result.Error != nil {
			return result.Error
		}

		var dislikeCount int64
		result = tx.Model(TopicUserLikes{}).Where("topic_id = ? and like_data = -1", id).Count(&dislikeCount)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&topic).UpdateColumns(Map{
			"like_count":    likeCount,
			"dislike_count": dislikeCount,
		})

		return result.Error
	})
	if err != nil {
		return err
	}

	var response TopicCommonResponse
	if err = copier.CopyWithOption(&response, &topic, CopyOption); err != nil {
		return err
	}
	return Success(c, &response)
}

// ViewATopic godoc
// @Summary 浏览一个话题，浏览数 +1
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_view [put]
// @Param id path int true "topic id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ViewATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var topic Topic
	result := DB.Preload("Tags").Preload("Poster").First(&topic, id)
	if result.Error != nil {
		return NotFound()
	}

	var topicUser = TopicUserViews{
		TopicID: id,
		UserID:  user.ID,
		Count:   1,
	}
	err = DB.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "topic_id"}},
			DoUpdates: clause.Assignments(Map{
				"updated_at": time.Now(),
				"count":      gorm.Expr("count + 1"),
			}),
		}).Create(&topicUser).Error
		if err != nil {
			return err
		}

		return tx.Model(&topic).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
	})
	if err != nil {
		return err
	}
	return Success(c, &EmptyStruct{})
}

// FavorATopic godoc
// @Summary 收藏一个话题，收藏数 +1
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_favor [put]
// @Param id path int true "topic id"
// @Success 200 {object} RespForSwagger{data=TopicCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func FavorATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var topic Topic
	result := DB.Preload("Tags").Preload("Poster").First(&topic, id)
	if result.Error != nil {
		return NotFound()
	}
	var topicUserFavorites TopicUserFavorites
	result = DB.Where("user_id = ? AND topic_id = ?", user.ID, topic.ID).First(&topicUserFavorites)
	if result.Error != nil {
		topicUserFavorites.UserID = user.ID
		topicUserFavorites.TopicID = topic.ID
		result = DB.Model(&topicUserFavorites).Create(&topicUserFavorites)
		if result.RowsAffected == 0 {
			return BadRequest()
		}
		topic.FavorCount++
		result = DB.Model(&topic).UpdateColumn("favor_count", gorm.Expr("favor_count + 1"))
		if result.RowsAffected == 0 {
			return BadRequest()
		}

		result = DB.Model(&user).UpdateColumn("favorite_topics_count", gorm.Expr("favorite_topics_count + 1"))
		if result.Error != nil {
			return result.Error
		}
	}

	var response TopicCommonResponse
	if err = copier.CopyWithOption(&response, &topic, CopyOption); err != nil {
		return err
	}

	return Success(c, &TopicCommonResponse{})
}

// ListFavoriteTopics godoc
// @Summary 查询收藏的话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_favor [get]
// @Param json query TopicListRequest true "page"
// @Success 200 {object} RespForSwagger{data=TopicListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListFavoriteTopics(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	var query TopicListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}
	nowTime := time.Now()
	if query.StartTime == nil {
		query.StartTime = &nowTime
	}
	var topicUserFavorites []TopicUserFavorites
	result := DB.Find(&topicUserFavorites, "user_id = ?", user.ID)
	if result.Error != nil {
		return NotFound()
	}

	var orderColumn = clause.Column{
		Name:  query.OrderBy,
		Table: "topic_user_favorites",
	}
	if query.OrderBy == "updated_at" {
		orderColumn.Table = "topic"
	}

	var topics []Topic
	tx := DB.
		Where("? < ?", orderColumn, query.StartTime).
		Order(
			clause.OrderByColumn{
				Column: orderColumn,
			}).Limit(query.PageSize)
	if query.DivisionID != nil {
		tx = tx.Where(&Topic{DivisionID: *query.DivisionID})
	}
	result = tx.Joins("inner join topic_user_favorites on topic_user_favorites.topic_id = topic.id and topic_user_favorites.user_id = ?", user.ID).
		Preload("Tags").Preload("Poster").Find(&topics)
	if result.Error != nil {
		return result.Error
	}
	var response TopicListResponse
	if err = copier.CopyWithOption(&response.Topics, &topics, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// UnfavorATopic godoc
// @Summary 取消收藏一个话题
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_favor [delete]
// @Param id path int true "topic id"
// @Success 200 {object} RespForSwagger{data=TopicCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func UnfavorATopic(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var topic Topic
	result := DB.First(&topic, id)
	if result.Error != nil {
		return result.Error
	}

	var topicUserFavorites = TopicUserFavorites{UserID: user.ID, TopicID: topic.ID}
	result = DB.Delete(&topicUserFavorites)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 1 {
		topic.FavorCount--
		result = DB.Model(&topic).UpdateColumn("favor_count", gorm.Expr("favor_count - 1"))
		if result.Error != nil {
			return result.Error
		}
	}

	var response TopicCommonResponse
	if err = copier.CopyWithOption(&response, &topic, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// ListTopicsByUser godoc
// @Summary 查询用户发布的话题
// @Description 用户查询自己发布的话题，以及点进其他用户主页时查询用户发布的话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_user/{user_id} [get]
// @Param user_id path int true "user id"
// @Param json query TopicListRequest true "page"
// @Success 200 {object} RespForSwagger{data=TopicListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListTopicsByUser(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	uid, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var query TopicListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	nowTime := time.Now()
	if query.StartTime == nil {
		query.StartTime = &nowTime
	}
	var topics []Topic
	querySet := DB.Order(query.OrderBy+" desc").Limit(query.PageSize).
		Where("? < ?", clause.Column{Name: query.OrderBy}, query.StartTime).
		Where("poster_id = ? and is_anonymous = false", uid)
	if query.DivisionID != nil {
		querySet = querySet.Where("division_id = ?", *query.DivisionID)
	}
	result := querySet.Preload("Tags").Preload("Poster").Find(&topics)
	if result.Error != nil {
		return result.Error
	}

	var response TopicListResponse
	if err = copier.CopyWithOption(&response.Topics, &topics, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// ListTopicsByTag godoc
// @Summary 查询标签下的话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_tag/{tag_id} [get]
// @Param tag_id path int true "tag id"
// @Param json query TopicListRequest true "page"
// @Success 200 {object} RespForSwagger{data=TopicListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListTopicsByTag(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	tagID, err := c.ParamsInt("tag_id")
	if err != nil {
		return err
	}

	tx := DB.Model(&Topic{}).Joins("inner join topic_tags on topic_tags.topic_id = topic.id")

	var topics []Topic
	result := tx.Where("tag_id = ?", tagID).Preload("Tags").Preload("Poster").Find(&topics)
	if result.Error != nil {
		return result.Error
	}

	var response TopicListResponse
	if err = copier.CopyWithOption(&response.Topics, &topics, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// SearchTopics godoc
// @Summary 搜索话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_search [get]
// @Param json query TopicSearchRequest true "page"
// @Success 200 {object} RespForSwagger{data=TopicListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func SearchTopics(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	var query TopicSearchRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	var topics []Topic
	_, err = Search(DB, &topics, query.Search, "", []string{"id desc"}, "title", query.PageRequest)
	if err != nil {
		return
	}

	var response TopicListResponse
	if err = copier.CopyWithOption(&response.Topics, &topics, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}
