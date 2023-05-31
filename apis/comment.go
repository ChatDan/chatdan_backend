package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListComments godoc
// @Summary 查询评论，按照 id 升序 或 点赞数倒序 排序
// @Tags Comment Module
// @Produce json
// @Router /comments [get]
// @Param json query CommentListRequest true "page"
// @Success 200 {object} RespForSwagger{data=CommentListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListComments(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	var query CommentListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	var topic Topic
	err = DB.First(&topic, query.TopicID).Error
	if err != nil {
		return err
	}

	tx := query.QuerySet(DB).Where("topic_id = ?", query.TopicID)
	if query.OrderBy == "id" {
		tx = tx.Order(query.OrderBy + " asc")
	} else {
		tx = tx.Order(query.OrderBy + " desc")
	}

	var comments []Comment
	result := tx.Preload("Poster").Find(&comments)
	if result.Error != nil {
		return result.Error
	}

	var response CommentListResponse
	if err = copier.CopyWithOption(&response.Comments, &comments, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// GetAComment godoc
// @Summary 获取一个评论
// @Tags Comment Module
// @Produce json
// @Router /comment/{id} [get]
// @Param id path int true "comment id"
// @Success 200 {object} RespForSwagger{data=CommentCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func GetAComment(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var comment Comment

	result := DB.Preload("Poster").First(&comment, id)
	if result.Error != nil {
		return result.Error
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, CopyOption); err != nil {
		return err
	}
	return Success(c, &response)
}

// CreateAComment godoc
// @Summary 创建一个评论
// @Tags Comment Module
// @Accept json
// @Produce json
// @Router /comment [post]
// @Param json body CommentCreateRequest true "comment"
// @Success 201 {object} RespForSwagger{data=CommentCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func CreateAComment(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	var body CommentCreateRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var topic Topic
	err = DB.First(&topic, body.TopicID).Error
	if err != nil {
		return NotFound()
	}

	var comment Comment
	err = DB.Transaction(func(tx *gorm.DB) error {

		comment = Comment{
			Content:     body.Content,
			ReplyToID:   body.ReplyToID,
			PosterID:    user.ID,
			TopicID:     body.TopicID,
			IsAnonymous: body.IsAnonymous,
		}

		if body.IsAnonymous {
			var anonyname string
			anonyname, err = FindOrGenerateAnonyname(tx, body.TopicID, user.ID)
			if err != nil {
				return err
			}
			comment.Anonyname = &anonyname
		}

		result := tx.Create(&comment)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return BadRequest()
		}

		// update topic
		result = tx.Model(&topic).Update("comment_count", gorm.Expr("comment_count + 1"))
		if result.Error != nil {
			return result.Error
		}

		// update user comment count
		result = tx.Model(&user).Update("comment_count", gorm.Expr("comment_count + 1"))
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = SearchAddOrReplace(comment.ToSearchModel())
	if err != nil {
		return err
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, CopyOption); err != nil {
		return err
	}
	return Created(c, &response)
}

// ModifyAComment godoc
// @Summary 修改一个评论
// @Tags Comment Module
// @Accept json
// @Produce json
// @Router /comment/{id} [put]
// @Param id path int true "comment id"
// @Param json body CommentModifyRequest true "comment"
// @Success 200 {object} RespForSwagger{data=CommentCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ModifyAComment(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var body CommentModifyRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var comment Comment
	result := DB.First(&comment, id)
	if result.Error != nil {
		return NotFound()
	}
	if !user.IsAdmin && user.ID != comment.PosterID {
		return Forbidden()
	}
	if body.IsHidden != nil {
		if !user.IsAdmin {
			return Forbidden()
		}
	}
	if err = DB.Transaction(func(tx *gorm.DB) error {
		// load comment with lock
		if err = tx.Clauses(LockClause).First(&comment, id).Error; err != nil {
			return err
		}

		// copy body to comment
		if err = copier.CopyWithOption(&comment, &body, CopyOption); err != nil {
			return err
		}

		// update comment
		return tx.Model(&comment).Select("Content", "IsHidden").Updates(&comment).Error
	}); err != nil {
		return err
	}

	err = SearchAddOrReplace(comment.ToSearchModel())
	if err != nil {
		return err
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// DeleteAComment godoc
// @Summary 删除一个评论，仅作者或管理员可删除
// @Tags Comment Module
// @Produce json
// @Router /comment/{id} [delete]
// @Param id path int true "comment id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func DeleteAComment(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var comment Comment
	result := DB.First(&comment, id)
	if result.Error != nil {
		return result.Error
	}

	if !user.IsAdmin && user.ID != comment.PosterID {
		return Forbidden()
	}

	result = DB.Delete(&comment)
	if result.RowsAffected == 0 {
		return BadRequest()
	}

	return Success(c, &EmptyStruct{})
}

// LikeOrDislikeAComment godoc
// @Summary 点赞或点踩一个评论，或者取消点赞或点踩
// @Tags Comment Module
// @Produce json
// @Router /comment/{id}/like/{like_data} [post]
// @Param id path int true "comment id"
// @Param like_data path int true "1: like, -1: dislike, 0: reset" Enums(1, -1, 0)
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func LikeOrDislikeAComment(c *fiber.Ctx) (err error) {

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

	var comment Comment
	err = DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(LockClause).First(&comment, id)
		if result.Error != nil {
			return result.Error
		}

		var commentUserLikes = CommentUserLikes{
			CommentID: id,
			UserID:    user.ID,
			LikeData:  likeData,
		}
		result = tx.Save(&commentUserLikes)
		if result.Error != nil {
			return result.Error
		}

		var likeCount int64
		result = tx.Model(CommentUserLikes{}).Where("topic_id = ? and like_data = 1", id).Count(&likeCount)
		if result.Error != nil {
			return result.Error
		}

		var dislikeCount int64
		result = tx.Model(CommentUserLikes{}).Where("topic_id = ? and like_data = -1", id).Count(&dislikeCount)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&comment).UpdateColumns(Map{
			"like_count":    likeCount,
			"dislike_count": dislikeCount,
		})

		return result.Error
	})
	if err != nil {
		return err
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, CopyOption); err != nil {
		return err
	}
	return Success(c, &response)
}

// ListCommentsByUser godoc
// @Summary 查询某个用户的评论，按照 id 升序 或 点赞数倒序 排序
// @Tags Comment Module
// @Produce json
// @Router /comments/_user/{user_id} [get]
// @Param user_id path int true "user id"
// @Param json query CommentListByUserRequest true "page"
// @Success 200 {object} RespForSwagger{data=CommentListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListCommentsByUser(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	uid, err := c.ParamsInt("id")

	var query CommentListByUserRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	tx := DB.Where("poster_id = ? and is_anonymous", uid, false)
	if query.OrderBy == "id" {
		tx = tx.Order(query.OrderBy + " asc")
	} else {
		tx = tx.Order(query.OrderBy + " desc")
	}

	tx = query.QuerySet(tx)

	var comments []Comment
	result := tx.Preload("Poster").Find(&comments)
	if result.Error != nil {
		return result.Error
	}

	var response CommentListResponse
	if err = copier.CopyWithOption(&response.Comments, &comments, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}

// SearchComments godoc
// @Summary 搜索评论
// @Tags Comment Module
// @Produce json
// @Router /comments/_search [get]
// @Param json query CommentSearchRequest true "page"
// @Success 200 {object} RespForSwagger{data=CommentListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func SearchComments(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	var query CommentSearchRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return
	}

	var comments []Comment
	_, err = Search(DB, &comments, query.Search, "", []string{"id desc"}, "", query.PageRequest)
	if err != nil {
		return
	}

	var response CommentListResponse
	if err = copier.CopyWithOption(&response.Comments, &comments, CopyOption); err != nil {
		return err
	}

	return Success(c, &response)
}
