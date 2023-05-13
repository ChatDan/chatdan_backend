package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// ListBoxes godoc
// @Summary 查询提问箱
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param body query BoxListRequest true "page"
// @Success 200 {object} Response{data=BoxListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /messageBoxes [get]
func ListBoxes(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate query
	var query BoxListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// construct querySet
	querySet := query.QuerySet(DB)
	if query.Title != "" {
		querySet = querySet.Where("title=?", query.Title) // TODO: fuzzy search
	}
	if query.Owner != 0 {
		querySet = querySet.Where("owner_id=?", query.Owner)
	}

	// load boxes from database
	var boxes []Box
	if err = querySet.Find(&boxes).Error; err != nil {
		return
	}

	// construct response
	var response BoxListResponse
	if err = copier.CopyWithOption(&response.MessageBoxes, &boxes, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}

	return Success(c, response)
}

// GetABox godoc
// @Summary 获取提问箱信息
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=BoxGetResponse}
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /messageBox/{id} [get]
func GetABox(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get box id
	var boxID int
	if boxID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load box from database
	var box Box
	if err = DB.Take(&box, boxID).Error; err != nil {
		return
	}

	// load post content from database by box_id
	var postsContent []string
	if err = DB.Model(&Post{}).Where("box_id=?", boxID).Pluck("content", &postsContent).Error; err != nil {
		return
	}

	// construct response
	var response BoxGetResponse
	if err = copier.CopyWithOption(&response, &box, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}
	response.PostsContent = postsContent

	return Success(c, response)
}

// CreateABox godoc
// @Summary 创建提问箱
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param box body BoxCreateRequest true "box"
// @Success 201 {object} Response{data=BoxCommonResponse}
// @Failure 400 {object} Response{data=ErrorDetail} "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /messageBox [post]
func CreateABox(c *fiber.Ctx) (err error) {
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	var body BoxCreateRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	box := Box{
		OwnerID: user.ID,
		Title:   body.Title,
	}

	if err = DB.Create(&box).Error; err != nil {
		return err
	}

	var response BoxCommonResponse
	if err = copier.CopyWithOption(&response, &box, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}

	return Created(c, response)
}

// ModifyABox godoc
// @Summary 修改提问箱信息
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param json body BoxModifyRequest true "box"
// @Success 200 {object} Response{data=BoxCommonResponse}
// @Failure 400 {object} Response{data=ErrorDetail} "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /messageBox/{id} [put]
func ModifyABox(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get box id
	var boxID int
	if boxID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// get and validate body
	var body BoxModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}
	if body.IsEmpty() {
		return BadRequest("empty body")
	}

	// load box from database
	var box Box
	if err = DB.Take(&box, boxID).Error; err != nil {
		return
	}

	// check if current user is the owner of the box
	if box.OwnerID != user.ID {
		return Forbidden()
	}

	// update box
	if err = DB.Model(&box).Select("title").Updates(&box).Error; err != nil {
		return
	}

	var response BoxCommonResponse
	if err = copier.CopyWithOption(&response, &box, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}

	return Success(c, response)
}

// DeleteABox godoc
// @Summary 删除提问箱
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /messageBox/{id} [delete]
func DeleteABox(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}

	// get box id
	var boxID int
	if boxID, err = c.ParamsInt("id"); err != nil {
		return err
	}

	// load box from database
	var box Box
	if err = DB.Take(&box, boxID).Error; err != nil {
		return err
	}

	// check if current user is the owner of the box
	if box.OwnerID != user.ID {
		return Forbidden()
	}

	// delete box
	if err = DB.Delete(&box).Error; err != nil {
		return err
	}

	return Success(c, EmptyStruct{})
}
