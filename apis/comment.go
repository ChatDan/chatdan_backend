package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"time"
)

// ListComments godoc
// @Summary 查询评论，按照 id 升序 或 点赞数倒序 排序
// @Tags Comment Module
// @Produce json
// @Router /comments [get]
// @Param json query CommentListRequest true "page"
// @Success 200 {object} Response{data=CommentListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
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

	tx := DB.Where("topic_id = ?", query.TopicID)
	if query.OrderBy == "id" {
		tx = tx.Order(query.OrderBy + "asc")
	} else {
		tx = tx.Order(query.OrderBy + "desc")
	}

	tx = tx.Limit(query.PageSize).Offset(query.PageNum * query.PageSize)

	var comments []Comment

	result := tx.Find(&comments)

	if result.Error != nil {
		return result.Error
	}

	var response CommentListRequest

	if err = copier.CopyWithOption(&response, &comments, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	return Success(c, response)
}

// GetAComment godoc
// @Summary 获取一个评论
// @Tags Comment Module
// @Produce json
// @Router /comment/{id} [get]
// @Param id path int true "comment id"
// @Success 200 {object} Response{data=CommentCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
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

	result := DB.First(&comment, id)
	if result.Error != nil {
		return result.Error
	}

	var response CommentCommonResponse

	if err = copier.CopyWithOption(&response, &comment, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	return Success(c, response)
}

// CreateAComment godoc
// @Summary 创建一个评论
// @Tags Comment Module
// @Accept json
// @Produce json
// @Router /comment [post]
// @Param json body CommentCreateRequest true "comment"
// @Success 201 {object} Response{data=CommentCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
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

	comment := Comment{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Content:   body.Content,
		ReplyToID: body.ReplyToID,
	}

	result := DB.Create(&comment)
	if result.RowsAffected == 0 {
		return BadRequest()
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	return Success(c, response)
}

// ModifyAComment godoc
// @Summary 修改一个评论
// @Tags Comment Module
// @Accept json
// @Produce json
// @Router /comment/{id} [put]
// @Param json body CommentModifyRequest true "comment"
// @Success 200 {object} Response{data=CommentCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
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
	if body.IsEmpty() {
		return BadRequest()
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
		// load division with lock
		if err = tx.Clauses(LockClause).First(&comment, id).Error; err != nil {
			return err
		}

		// copy body to division
		if err = copier.CopyWithOption(&comment, &body, copier.Option{IgnoreEmpty: true}); err != nil {
			return err
		}

		// update division
		return tx.Model(&comment).Updates(&comment).Error
	}); err != nil {
		return err
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	return Success(c, response)
}

// DeleteAComment godoc
// @Summary 删除一个评论，仅作者或管理员可删除
// @Tags Comment Module
// @Produce json
// @Router /comment/{id} [delete]
// @Param id path int true "comment id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
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

	return Success(c, EmptyStruct{})
}

// LikeOrDislikeAComment godoc
// @Summary 点赞或点踩一个评论，或者取消点赞或点踩
// @Tags Comment Module
// @Produce json
// @Router /comment/{id}/like/{like_data} [post]
// @Param id path int true "comment id"
// @Param like_data path int true "1: like, -1: dislike, 0: reset" Enums(1, -1, 0)
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
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
	result := DB.First(&comment, id)
	if result.Error != nil {
		return result.Error
	}

	var commentUserLikes CommentUserLikes
	result = DB.Model(&commentUserLikes).Where("uesr_id = ? AND comment_id = ?", user.ID, id).First(&commentUserLikes)

	if result.Error != nil {
		comment.LikeCount = comment.LikeCount - commentUserLikes.LikeData + likeData
		commentUserLikes.LikeData = likeData
		result = DB.Model(&commentUserLikes).Updates(commentUserLikes)
		if result.Error != nil {
			return result.Error
		}
	} else {
		comment.LikeCount = comment.LikeCount + likeData
		commentUserLikes.CommentID = id
		commentUserLikes.UserID = user.ID
		commentUserLikes.CreatedAt = time.Now()
		commentUserLikes.LikeData = likeData

		result = DB.Create(&commentUserLikes)
		if result.RowsAffected == 0 {
			return BadRequest()
		}
	}

	result = DB.Model(&comment).Updates(comment)
	if result.RowsAffected == 0 {
		return BadRequest()
	}

	var response CommentCommonResponse
	if err = copier.CopyWithOption(&response, &comment, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	return Success(c, response)
}

// ListCommentsByUser godoc
// @Summary 查询某个用户的评论，按照 id 升序 或 点赞数倒序 排序
// @Tags Comment Module
// @Produce json
// @Router /comments/_user/{user_id} [get]
// @Param user_id path int true "user id"
// @Param json query CommentListRequest true "page"
// @Success 200 {object} Response{data=CommentListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListCommentsByUser(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}

	uid, err := c.ParamsInt("id")

	var query CommentListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	tx := DB.Where("topic_id = ? AND poster_id = ?", query.TopicID, uid)
	if query.OrderBy == "id" {
		tx = tx.Order(query.OrderBy + "asc")
	} else {
		tx = tx.Order(query.OrderBy + "desc")
	}

	tx = tx.Limit(query.PageSize).Offset(query.PageNum * query.PageSize)

	var comments []Comment

	result := tx.Find(&comments)

	if result.Error != nil {
		return result.Error
	}

	var response CommentListRequest

	if err = copier.CopyWithOption(&response, &comments, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	return Success(c, response)
}
