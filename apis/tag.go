package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"strings"
)

// ListTags godoc
// @Summary 查询标签，按照 热度 倒序 排序
// @Tags Tag Module
// @Produce json
// @Router /tags [get]
// @Param json query TagListRequest true "page"
// @Success 200 {object} RespForSwagger{data=TagListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListTags(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	var query TagListRequest
	err = ValidateQuery(c, &query)
	if err != nil {
		return
	}

	var (
		tags    []Tag
		version int
		total   int
		key     = "tags"
	)
	if query.Search != "" {
		// using meilisearch
		var filter string
		if total, err = Search(
			DB, &tags, query.Search,
			filter, []string{query.OrderBy}, "name", query.PageRequest,
		); err != nil {
			return
		}
	} else {
		// load from database
		tx := DB.Session(&gorm.Session{NewDB: true}).Model(&Tag{}).Order(query.OrderBy)
		key = key + ":" + strings.Replace(query.OrderBy, " ", "_", -1)
		if version, total, err = PageLoad(tx, &tags, key, query.PageRequest); err != nil {
			return
		}
	}

	// copy to response
	var response TagListResponse
	err = copier.Copy(&response.Tags, &tags)
	if err != nil {
		return
	}
	response.Total = total
	response.Version = version

	return Success(c, &response)
}

// GetATag godoc
// @Summary 获取一个标签
// @Tags Tag Module
// @Produce json
// @Router /tag/{id} [get]
// @Param id path int true "tag id"
// @Success 200 {object} RespForSwagger{data=TagCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func GetATag(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	tagID, err := c.ParamsInt("id")
	if err != nil {
		return
	}

	var tag = Tag{ID: tagID}
	err = LoadModel(DB, &tag)
	if err != nil {
		return
	}

	var response TagCommonResponse
	err = copier.Copy(&response, &tag)
	if err != nil {
		return
	}
	return Success(c, &response)
}

// CreateATag godoc
// @Summary 创建一个标签
// @Tags Tag Module
// @Accept json
// @Produce json
// @Router /tag [post]
// @Param json body TagCreateRequest true "tag"
// @Success 201 {object} RespForSwagger{data=TagCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func CreateATag(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	var request TagCreateRequest
	err = ValidateBody(c, &request)
	if err != nil {
		return
	}

	var tag = Tag{Name: request.Name}

	err = CreateModel(DB, &tag)
	if err != nil {
		return
	}

	var response TagCommonResponse
	err = copier.Copy(&response, &tag)
	if err != nil {
		return
	}

	return Created(c, &response)
}

// ModifyATag godoc
// @Summary 修改一个标签，仅管理员可修改
// @Tags Tag Module
// @Accept json
// @Produce json
// @Router /tag/{id} [put]
// @Param json body TagModifyRequest true "tag"
// @Success 200 {object} RespForSwagger{data=TagCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ModifyATag(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	if !user.IsAdmin {
		return Forbidden("非管理员无法修改标签")
	}

	var request TagModifyRequest
	err = ValidateBody(c, &request)
	if err != nil {
		return
	}

	tagID, err := c.ParamsInt("id")
	if err != nil {
		return
	}

	var tag = Tag{ID: tagID}
	err = LoadModel(DB, &tag)
	if err != nil {
		return
	}

	err = UpdateModel(DB, &tag, request)
	if err != nil {
		return
	}

	var response TagCommonResponse
	err = copier.Copy(&response, &tag)
	if err != nil {
		return
	}

	return Success(c, &response)
}

// DeleteATag godoc
// @Summary 删除一个标签，仅管理员可删除
// @Tags Tag Module
// @Produce json
// @Router /tag/{id} [delete]
// @Param id path int true "tag id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func DeleteATag(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	if !user.IsAdmin {
		return Forbidden("非管理员无法删除标签")
	}

	tagID, err := c.ParamsInt("id")
	if err != nil {
		return
	}

	var tag = Tag{ID: tagID}
	err = LoadModel(DB, &tag)
	if err != nil {
		return
	}

	err = DeleteModel(DB, &tag)
	if err != nil {
		return
	}
	return Success(c, &EmptyStruct{})
}
