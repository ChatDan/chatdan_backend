package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
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
	return Success(c, nil)
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
	return Success(c, nil)
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
	return Success(c, nil)
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
	return Success(c, nil)
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
	return Success(c, nil)
}
