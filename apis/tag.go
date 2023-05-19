package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
)

// ListTags godoc
// @Summary 查询标签，按照 热度 倒序 排序
// @Tags Tag Module
// @Produce json
// @Router /tags [get]
// @Param json query TagListRequest true "page"
// @Success 200 {object} Response{data=TagListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListTags(c *fiber.Ctx) (err error) {
	return Success(c, TagListResponse{})
}

// GetATag godoc
// @Summary 获取一个标签
// @Tags Tag Module
// @Produce json
// @Router /tag/{id} [get]
// @Param id path int true "tag id"
// @Success 200 {object} Response{data=TagCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func GetATag(c *fiber.Ctx) (err error) {
	return Success(c, TagCommonResponse{})
}

// CreateATag godoc
// @Summary 创建一个标签
// @Tags Tag Module
// @Accept json
// @Produce json
// @Router /tag [post]
// @Param json body TagCreateRequest true "tag"
// @Success 201 {object} Response{data=TagCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func CreateATag(c *fiber.Ctx) (err error) {
	return Created(c, TagCommonResponse{})
}

// ModifyATag godoc
// @Summary 修改一个标签，仅管理员可修改
// @Tags Tag Module
// @Accept json
// @Produce json
// @Router /tag/{id} [put]
// @Param json body TagModifyRequest true "tag"
// @Success 200 {object} Response{data=TagCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ModifyATag(c *fiber.Ctx) (err error) {
	return Success(c, TagCommonResponse{})
}

// DeleteATag godoc
// @Summary 删除一个标签，仅管理员可删除
// @Tags Tag Module
// @Produce json
// @Router /tag/{id} [delete]
// @Param id path int true "tag id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func DeleteATag(c *fiber.Ctx) (err error) {
	return Success(c, EmptyStruct{})
}
