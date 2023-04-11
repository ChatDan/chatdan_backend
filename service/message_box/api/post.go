package api

import "github.com/gofiber/fiber/v2"

// CreateAPost godoc
// @Summary Create a post
// @Description Create a post
// @Tags Post
// @Accept json
// @Produce json
// @Param post query PostCreateRequest true "post"
// @Success 200 {object} common.Response{data=PostCreateResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /post [post]
func CreateAPost(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// ListPosts godoc
// @Summary List posts
// @Description List posts
// @Tags Post
// @Accept json
// @Produce json
// @Param body query PostListRequest true "page"
// @Success 200 {object} common.Response{data=PostListResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /posts [get]
func ListPosts(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// GetAPost godoc
// @Summary Get a post
// @Description Get a post
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} common.Response{data=PostGetResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /post/{id} [get]
func GetAPost(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// ModifyAPost godoc
// @Summary Modify a post
// @Description Modify a post, owner only
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param post body PostModifyRequest true "post"
// @Success 200 {object} common.Response{data=PostModifyResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /post/{id} [put]
func ModifyAPost(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// DeleteAPost godoc
// @Summary Delete a post
// @Description Delete a post, owner only
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} common.Response{data=PostDeleteResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /post/{id} [delete]
func DeleteAPost(c *fiber.Ctx) error {
	return c.JSON(nil)
}
