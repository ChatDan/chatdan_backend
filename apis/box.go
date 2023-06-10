package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// ListBoxes godoc
// @Summary 查询提问箱
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param body query BoxListRequest true "page"
// @Success 200 {object} RespForSwagger{data=BoxListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
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

	// 从缓存或数据库中读取数据
	var (
		boxes   []Box
		version int
		total   int
		key     = "boxes"
	)

	if query.Title != "" {
		// 使用 Meilisearch 模糊搜索
		var filter string
		if query.Owner != 0 {
			filter = "owner_id = " + strconv.Itoa(query.Owner)
		}
		if total, err = Search(
			DB, &boxes, query.Title,
			filter, []string{query.OrderBy}, "title", query.PageRequest,
		); err != nil {
			return
		}
	} else {
		tx := DB.Session(&gorm.Session{NewDB: true}).Model(&Box{}).Order(query.OrderBy)
		if query.Owner != 0 {
			tx = tx.Where("owner_id = ?", query.Owner)
			key = "boxes:" + strconv.Itoa(query.Owner)
		}
		key = key + ":" + strings.Replace(query.OrderBy, " ", "_", -1)
		if version, total, err = PageLoad(tx, &boxes, key, query.PageRequest); err != nil {
			return
		}
	}

	// 构建响应
	var response BoxListResponse
	if err = copier.CopyWithOption(&response.MessageBoxes, &boxes, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}
	response.Version = version
	response.Total = total

	return Success(c, &response)
}

// GetABox godoc
// @Summary 获取提问箱信息
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} RespForSwagger{data=BoxGetResponse}
// @Failure 400 {object} RespForSwagger "Bad Request"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
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

	// 从缓存说数据库中读取
	var box = Box{ID: boxID}
	if err = LoadModel(DB, &box); err != nil {
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

	return Success(c, &response)
}

// CreateABox godoc
// @Summary 创建提问箱
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param box body BoxCreateRequest true "box"
// @Success 201 {object} RespForSwagger{data=BoxCommonResponse}
// @Failure 400 {object} RespForSwagger{data=ErrorDetail} "Bad Request"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
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

	// 创建提问箱
	if err = DB.Create(&box).Error; err != nil {
		return err
	}
	if err = SearchAddOrReplace(box.ToBoxSearchModel()); err != nil {
		return
	}

	// 删除缓存
	go DeleteInBatch(
		fmt.Sprintf("boxes:%d:id_asc:latest", user.ID),
		fmt.Sprintf("boxes:%d:updated_at_desc:latest", user.ID),
		"boxes:id_asc:latest",
		"boxes:updated_at_desc:latest",
	)

	var response BoxCommonResponse
	if err = copier.CopyWithOption(&response, &box, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}

	return Created(c, &response)
}

// ModifyABox godoc
// @Summary 修改提问箱信息
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param json body BoxModifyRequest true "box"
// @Success 200 {object} RespForSwagger{data=BoxCommonResponse}
// @Failure 400 {object} RespForSwagger{data=ErrorDetail} "Bad Request"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
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
	if err = copier.CopyWithOption(&box, &body, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	if err = DB.Model(&box).Select("title").Updates(&box).Error; err != nil {
		return
	}
	if err = SearchAddOrReplace(box.ToBoxSearchModel()); err != nil {
		return
	}

	// 删除缓存
	go DeleteInBatch(
		fmt.Sprintf("boxes:%d:updated_at_desc:latest", user.ID),
		"boxes:updated_at_desc:latest",
	)

	var response BoxCommonResponse
	if err = copier.CopyWithOption(&response, &box, copier.Option{IgnoreEmpty: true}); err != nil {
		return
	}

	return Success(c, &response)
}

// DeleteABox godoc
// @Summary 删除提问箱
// @Tags MessageBox Module
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger "Bad Request"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
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
	if err = SearchDelete[BoxSearchModel](box.ID); err != nil {
		return
	}

	// 删除缓存
	go DeleteInBatch(
		fmt.Sprintf("boxes:%d:id_asc:latest", user.ID),
		fmt.Sprintf("boxes:%d:updated_at_desc:latest", user.ID),
		"boxes:id_asc:latest",
		"boxes:updated_at_desc:latest",
		CacheName(&box),
	)

	return Success(c, &EmptyStruct{})
}
