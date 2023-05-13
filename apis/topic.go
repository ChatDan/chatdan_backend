package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
)

// ListTopics godoc
// @Summary 查询话题，按照最近创建或最近回复排序
// @Tags Topic Module
// @Produce json
// @Router /topics [get]
// @Param json query TopicListRequest true "page"
// @Success 200 {object} Response{data=TopicListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListTopics(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// GetATopic godoc
// @Summary 获取一个话题
// @Tags Topic Module
// @Produce json
// @Router /topic/{id} [get]
// @Param id path int true "topic id"
// @Success 200 {object} Response{data=TopicCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func GetATopic(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// CreateATopic godoc
// @Summary 创建一个话题
// @Tags Topic Module
// @Accept json
// @Produce json
// @Router /topic [post]
// @Param json body TopicCreateRequest true "topic"
// @Success 201 {object} Response{data=TopicCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func CreateATopic(c *fiber.Ctx) (err error) {
	return Created(c, nil)
}

// ModifyATopic godoc
// @Summary 修改一个话题
// @Description 管理员可修改标题、内容、标签、是否隐藏，用户可修改标题、内容、标签、是否匿名
// @Tags Topic Module
// @Accept json
// @Produce json
// @Router /topic/{id} [put]
// @Param json body TopicModifyRequest true "topic"
// @Success 200 {object} Response{data=TopicCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ModifyATopic(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// DeleteATopic godoc
// @Summary 删除一个话题，仅管理员
// @Tags Topic Module
// @Produce json
// @Router /topic/{id} [delete]
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func DeleteATopic(c *fiber.Ctx) (err error) {
	return Success(c, EmptyStruct{})
}

// LikeOrDislikeATopic godoc
// @Summary 点赞或点踩一个话题，或者重置点赞点踩数据
// @Description 1: like, -1: dislike, 0: reset，点赞点踩二选一
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_like/{like_data} [put]
// @Param id path int true "topic id"
// @Param like_data path int true "1: like, -1: dislike, 0: reset" Enums(1, -1, 0)
// @Success 200 {object} Response{data=TopicCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func LikeOrDislikeATopic(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// ViewATopic godoc
// @Summary 浏览一个话题，浏览数 +1
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_view [post]
// @Param id path int true "topic id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ViewATopic(c *fiber.Ctx) (err error) {
	return Success(c, EmptyStruct{})
}

// FavorATopic godoc
// @Summary 收藏一个话题，收藏数 +1
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_favor [post]
// @Param id path int true "topic id"
// @Success 200 {object} Response{data=TopicCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func FavorATopic(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// ListFavoriteTopics godoc
// @Summary 查询收藏的话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_favor [get]
// @Param json query TopicListRequest true "page"
// @Success 200 {object} Response{data=TopicListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListFavoriteTopics(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// UnfavorATopic godoc
// @Summary 取消收藏一个话题
// @Tags Topic Module
// @Produce json
// @Router /topic/{id}/_favor [delete]
// @Param id path int true "topic id"
// @Success 200 {object} Response{data=TopicCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func UnfavorATopic(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// ListTopicsByUser godoc
// @Summary 查询用户发布的话题
// @Description 用户查询自己发布的话题，以及点进其他用户主页时查询用户发布的话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_user/{user_id} [get]
// @Param user_id path int true "user id"
// @Param json query TopicListRequest true "page"
// @Success 200 {object} Response{data=TopicListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListTopicsByUser(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// ListTopicsByTag godoc
// @Summary 查询标签下的话题
// @Tags Topic Module
// @Produce json
// @Router /topics/_tag/{tag_id} [get]
// @Param tag_id path int true "tag id"
// @Param json query TopicListRequest true "page"
// @Success 200 {object} Response{data=TopicListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListTopicsByTag(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}
