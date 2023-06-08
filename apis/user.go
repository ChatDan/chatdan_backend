package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"time"
)

// ListUsers godoc
// @Summary 查询所有用户, admin only
// @Tags User Module
// @Produce json
// @Router /users [get]
// @Param body query UserListRequest true "page"
// @Success 200 {object} RespForSwagger{data=UserListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListUsers(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	if !user.IsAdmin {
		return Forbidden("只有管理员才能查看用户列表")
	}

	// get and validate query
	var query UserListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// load users from database
	var (
		users          []User
		version, total int
	)
	if version, total, err = PageLoad(DB, &users, "users", query.PageRequest); err != nil {
		return
	}

	// construct response
	var response UserListResponse
	if err = copier.Copy(&response, &users); err != nil {
		return
	}
	response.Version = version
	response.Total = total

	return Success(c, &response)
}

// GetUserMe godoc
// @Summary 获取当前用户信息
// @Tags User Module
// @Produce json
// @Router /user/me [get]
// @Success 200 {object} RespForSwagger{data=UserResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func GetUserMe(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// load user from database
	if err = LoadModel(DB, &user); err != nil {
		return
	}

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, &response)
}

// ModifyUserMe	godoc
// @Summary 修改当前用户信息
// @Tags User Module
// @Produce json
// @Router /user/me [put]
// @Param json body UserModifyRequest true "user"
// @Success 200 {object} RespForSwagger{data=UserResponse}
// @Failure 400 {object} RespForSwagger{data=ErrorDetail} "Invalid request body"
// @Failure 500 {object} RespForSwagger
func ModifyUserMe(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate request body
	var body UserModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	if err = DB.Transaction(func(tx *gorm.DB) (err error) {
		// load user from database
		if err = LoadModel(DB, &user); err != nil {
			return
		}

		// modify user
		if err = UpdateModel(DB, &user, body); err != nil {
			return
		}

		// update search
		if err = SearchAddOrReplace(user.ToSearchModel()); err != nil {
			return
		}

		return nil
	}); err != nil {
		return
	}

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, &response)
}

// DeleteUserMe godoc
// @Summary 注销当前用户
// @Tags User Module
// @Produce json
// @Router /user/me [delete]
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func DeleteUserMe(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// delete user
	if err = DB.Transaction(func(tx *gorm.DB) (err error) {
		// load user from database
		if err = DB.Clauses(LockClause).First(&user, user.ID).Error; err != nil {
			return
		}

		// change user's username
		if err = tx.Model(&user).Update("username", user.DeletedUsername()).Error; err != nil {
			return
		}

		// delete user
		if err = DeleteModel(tx, &user); err != nil {
			return
		}

		// delete search
		if err = SearchDelete[UserSearchModel](user.ID); err != nil {
			return
		}

		return nil
	}); err != nil {
		return
	}

	return Success(c, &EmptyStruct{})
}

// GetAUser godoc
// @Summary 获取用户信息
// @Tags User Module
// @Produce json
// @Router /user/{id} [get]
// @Param id path int true "user id"
// @Success 200 {object} RespForSwagger{data=UserResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func GetAUser(c *fiber.Ctx) (err error) {
	// get current user
	var currentUser User
	if err = GetCurrentUser(c, &currentUser); err != nil {
		return
	}

	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load user from database
	var user = User{ID: userID}
	if err = LoadModel(DB, &user); err != nil {
		return
	}

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, &response)
}

// ModifyAUser godoc
// @Summary 修改用户信息, admin only
// @Tags User Module
// @Produce json
// @Router /user/{id} [put]
// @Param id path int true "user id"
// @Param json body UserModifyRequest true "user"
// @Success 200 {object} RespForSwagger{data=UserResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ModifyAUser(c *fiber.Ctx) (err error) {
	// get current user
	var currentUser User
	if err = GetCurrentUser(c, &currentUser); err != nil {
		return
	}

	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// check permission
	if !currentUser.IsAdmin || currentUser.ID != userID {
		return Forbidden("只有管理员或自己才能修改用户信息")
	}

	// get and validate request body
	var body UserModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	var user User
	if err = DB.Transaction(func(tx *gorm.DB) (err error) {
		// load user from database
		if err = LoadModel(DB, &user); err != nil {
			return
		}

		// modify user
		if err = UpdateModel(DB, &user, body); err != nil {
			return
		}

		// update search
		if err = SearchAddOrReplace(user.ToSearchModel()); err != nil {
			return
		}

		return nil
	}); err != nil {
		return
	}

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, &response)
}

// DeleteAUser godoc
// @Summary 注销用户, admin only
// @Tags User Module
// @Produce json
// @Router /user/{id} [delete]
// @Param id path int true "user id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func DeleteAUser(c *fiber.Ctx) (err error) {
	// get current user
	var currentUser User
	if err = GetCurrentUser(c, &currentUser); err != nil {
		return
	}

	if !currentUser.IsAdmin {
		return Forbidden("只有管理员才能注销用户")
	}

	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	var user User
	if err = DB.Transaction(func(tx *gorm.DB) (err error) {
		// load user from database
		if err = DB.Clauses(LockClause).First(&user, userID).Error; err != nil {
			return
		}

		// change user's username
		if err = tx.Model(&user).Update("username", user.DeletedUsername()).Error; err != nil {
			return err
		}

		// delete user
		if err = DeleteModel(tx, &user); err != nil {
			return
		}

		// delete search
		if err = SearchDelete[UserSearchModel](user.ID); err != nil {
			return
		}

		return nil
	}); err != nil {
		return err
	}

	// 删除缓存
	Delete(CacheName(&user))

	return Success(c, &EmptyStruct{})
}

// FollowAUser godoc
// @Summary 关注用户
// @Tags User Module
// @Produce json
// @Router /user/{id}/_follow [post]
// @Param id path int true "user id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func FollowAUser(c *fiber.Ctx) (err error) {
	// get current user
	var currentUser User
	if err = GetCurrentUser(c, &currentUser); err != nil {
		return
	}

	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load user from database
	var user User

	// transaction
	if err = DB.Transaction(func(tx *gorm.DB) (err error) {

		if err = tx.Clauses(LockClause).First(&user, userID).Error; err != nil {
			return
		}

		// follow user
		var userFollows = UserFollows{
			UserID:     userID,
			FollowerID: currentUser.ID,
			CreatedAt:  time.Now(),
		}
		result := tx.FirstOrCreate(&userFollows)
		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return BadRequest("已经关注过该用户")
		}

		// update count
		if err = tx.Model(&user).
			UpdateColumn("followers_count", gorm.Expr("followers_count + 1")).Error; err != nil {
			return
		}

		if err = tx.Model(&currentUser).
			UpdateColumn("following_users_count", gorm.Expr("following_users_count + 1")).Error; err != nil {
			return
		}

		return nil
	}); err != nil {
		return
	}

	// construct response
	var response EmptyStruct
	return Success(c, &response)
}

// UnfollowAUser godoc
// @Summary 取消关注用户
// @Tags User Module
// @Produce json
// @Router /user/{id}/_follow [delete]
// @Param id path int true "user id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func UnfollowAUser(c *fiber.Ctx) (err error) {
	// get current user
	var currentUser User
	if err = GetCurrentUser(c, &currentUser); err != nil {
		return
	}

	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	var user User

	// transaction
	if err = DB.Transaction(func(tx *gorm.DB) (err error) {

		if err = tx.Clauses(LockClause).First(&user, userID).Error; err != nil {
			return
		}

		// unfollow user
		var userFollows = UserFollows{
			UserID:     userID,
			FollowerID: currentUser.ID,
		}
		result := tx.Delete(&userFollows)
		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return BadRequest("未关注该用户")
		}

		// update count
		if err = tx.Model(&user).
			UpdateColumn("followers_count", gorm.Expr("followers_count - 1")).Error; err != nil {
			return
		}

		if err = tx.Model(&currentUser).
			UpdateColumn("following_users_count", gorm.Expr("following_users_count - 1")).Error; err != nil {
			return
		}
		return nil
	}); err != nil {
		return err
	}

	// construct response
	var response EmptyStruct
	return Success(c, &response)
}

// ListUserFollowers godoc
// @Summary 获取用户的粉丝列表
// @Tags User Module
// @Produce json
// @Router /users/{id}/_followers [get]
// @Param id path int true "user id"
// @Param page query UserListRequest true "page"
// @Success 200 {object} RespForSwagger{data=UserListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListUserFollowers(c *fiber.Ctx) (err error) {
	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// get and validate request query
	var query UserListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// load user from database
	var user User
	if err = DB.First(&user, userID).Error; err != nil {
		return
	}

	// get followers
	var followers []User
	if err = query.QuerySet(DB).Model(&user).
		Joins("inner join user_follows on user_follows.follower_id = user.id and user_follows.user_id = ?", user.ID).
		Order("user_follows.created_at desc").
		Find(&followers).Error; err != nil {
		return
	}

	// construct response
	var response UserListResponse
	if err = copier.Copy(&response.Users, &followers); err != nil {
		return
	}

	return Success(c, &response)
}

// ListUserFollowing godoc
// @Summary 获取用户的关注列表
// @Tags User Module
// @Produce json
// @Router /users/{id}/_following [get]
// @Param id path int true "user id"
// @Param page query UserListRequest true "page"
// @Success 200 {object} RespForSwagger{data=UserListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ListUserFollowing(c *fiber.Ctx) (err error) {
	// get userID
	var userID int
	if userID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// get and validate request query
	var query UserListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// load user from database
	var user User
	if err = DB.First(&user, userID).Error; err != nil {
		return
	}

	// get following
	var following []User
	if err = query.QuerySet(DB).Model(&user).
		Joins("inner join user_follows on user_follows.user_id = user.id and user_follows.follower_id = ?", user.ID).
		Order("user_follows.created_at desc").
		Find(&following).Error; err != nil {
		return
	}

	// construct response
	var response UserListResponse
	if err = copier.Copy(&response.Users, &following); err != nil {
		return
	}

	return Success(c, &response)
}

// SearchUsers godoc
// @Summary 搜索用户
// @Tags User Module
// @Produce json
// @Router /users/_search [get]
// @Param page query UserSearchRequest true "page"
// @Success 200 {object} RespForSwagger{data=UserListResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func SearchUsers(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return
	}

	// get and validate request query
	var query UserSearchRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// get users
	var users []User
	var total int
	if total, err = Search(query.QuerySet(DB), &users, query.Search, "", nil, "username", query.PageRequest); err != nil {
		return
	}

	// construct response
	var response UserListResponse
	if err = copier.Copy(&response.Users, &users); err != nil {
		return
	}
	response.Total = total

	return Success(c, &response)
}
