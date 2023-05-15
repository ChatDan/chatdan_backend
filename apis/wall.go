package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
)

// ListWalls
// @Summary 获取今日表白墙
// @Description 获取今日表白墙，即昨日发送的表白墙
// @Tags Wall Module
// @Router /wall [get]
// @Produce json
// @Param json query WallListRequest true "query"
// @Success 200 {object} Response{data=WallListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response "服务器错误"
func ListWalls(c *fiber.Ctx) (err error) {
	// get and validate query
	var query WallListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return err
	}

	// construct querySet and load walls from database（昨日发送的表白墙）
	queryDate := carbon.Date{Carbon: carbon.Yesterday()}
	if query.Date != nil {
		queryDate = *query.Date
	}
	if queryDate.Carbon.Gt(carbon.Now()) {
		return BadRequest("不允许查询未来的表白墙")
	}

	var walls []Wall
	if err = query.QuerySet(DB).Where(
		"created_at between ? and ?",
		queryDate.StartOfDay().ToDateTimeString(),
		queryDate.EndOfDay().ToDateTimeString(),
	).Find(&walls).Error; err != nil {
		return err
	}

	// construct response
	var response WallListResponse
	if err = copier.Copy(&response.Posts, &walls); err != nil {
		return err
	}
	response.Date = queryDate

	return Success(c, response)
}

// GetAWall
// @Summary 获取表白墙信息
// @Tags Wall Module
// @Router /wall/{id} [get]
// @Produce json
// @Param id path int true "wall id"
// @Success 200 {object} Response{data=WallCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response "服务器错误"
func GetAWall(c *fiber.Ctx) (err error) {
	// get wall id
	var wallID int
	if wallID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load wall from database
	var wall Wall
	if err = DB.First(&wall, wallID).Error; err != nil {
		return
	}

	// construct response
	var response WallCommonResponse
	if err = copier.Copy(&response, &wall); err != nil {
		return
	}

	return Success(c, response)
}

// CreateAWall
// @Summary 创建表白墙
// @Tags Wall Module
// @Router /wall [post]
// @Produce json
// @Param json body WallCreateRequest true "json"
// @Success 201 {object} Response{data=WallCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response "服务器错误"
func CreateAWall(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate request
	var body WallCreateRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	// create wall
	var wall Wall
	if err = copier.CopyWithOption(&wall, &body, CopyOption); err != nil {
		return
	}
	wall.PosterID = user.ID
	if err = DB.Create(&wall).Error; err != nil {
		return
	}

	// construct response
	var response WallCommonResponse
	if err = copier.Copy(&response, &wall); err != nil {
		return
	}
	response.IsShown = carbon.FromStdTime(wall.CreatedAt).Lte(carbon.Yesterday().EndOfDay()) // 创建时间在昨天结束之前的才会显示

	return Created(c, response)
}
