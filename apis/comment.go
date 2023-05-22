package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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

	return Success(c, EmptyStruct{})
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
	return Success(c, EmptyStruct{})
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
	return Success(c, EmptyStruct{})
}

// LikeOrDislikeAComment godoc
// @Summary 点赞或点踩一个评论，或者取消点赞或点踩
// @Tags Comment Module
// @Produce json
// @Router /comment/{id}/like [post]
// @Param id path int true "comment id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func LikeOrDislikeAComment(c *fiber.Ctx) (err error) {
	return Success(c, EmptyStruct{})
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
	return Success(c, EmptyStruct{})
}
