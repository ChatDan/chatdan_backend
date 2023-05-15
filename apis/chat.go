package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListChats godoc
// @Summary 查询所有聊天记录，按照 time_updated 倒序排序
// @Tags Chat Module
// @Produce json
// @Router /chats [get]
// @Success 200 {object} Response{data=ChatListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListChats(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// load chats from database
	var chats []Chat
	if err = DB.Where("one_user_id = @user_id or another_user_id = @user_id", sql.Named("user_id", user.ID)).
		Order("updated_at desc").Find(&chats).Error; err != nil {
		return
	}

	// construct response
	var response ChatListResponse
	if err = copier.Copy(&response.Chats, &chats); err != nil {
		return
	}

	return Success(c, response)
}

// ListMessages godoc
// @Summary 查询所有聊天记录，按照 time_created 或 id 倒序排序
// @Tags Chat Module
// @Produce json
// @Router /messages [get]
// @Param body query MessageListRequest true "page"
// @Success 200 {object} Response{data=MessageListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListMessages(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate query
	var query MessageListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// load messages from database
	OneUserID := user.ID
	AnotherUserID := query.ToUserID
	if OneUserID == AnotherUserID {
		return BadRequest("不能和自己聊天 :-)")
	}
	if OneUserID > AnotherUserID {
		OneUserID, AnotherUserID = AnotherUserID, OneUserID
	}
	querySet := DB.Limit(query.PageSize).Where("one_user_id = ? and another_user_id = ?", OneUserID, AnotherUserID)
	if query.StartTime != nil {
		querySet = querySet.Where("created_at < ?", query.StartTime)
	}
	var messages []ChatMessage
	if err = querySet.Order("id desc").Find(&messages).Error; err != nil {
		return
	}

	// construct response
	var response MessageListResponse
	if err = copier.Copy(&response.Messages, &messages); err != nil {
		return
	}
	for i := range response.Messages {
		response.Messages[i].IsOwner = response.Messages[i].FromUserID == user.ID
	}

	return Success(c, response)
}

// CreateMessage godoc
// @Summary 发送消息
// @Tags Chat Module
// @Produce json
// @Router /messages [post]
// @Param json body MessageCreateRequest true "message"
// @Success 201 {object} Response{data=MessageCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func CreateMessage(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate request body
	var query MessageCreateRequest
	if err = ValidateBody(c, &query); err != nil {
		return
	}

	message := ChatMessage{
		FromUserID: user.ID,
		ToUserID:   query.ToUserID,
		Content:    query.Content,
	}

	// load and execute transaction
	if err = DB.Transaction(func(tx *gorm.DB) (err error) {
		// load another user, and check if it exists
		var anotherUser User
		if err = tx.First(&anotherUser, query.ToUserID).Error; err != nil {
			return
		}

		// load chat_id or create a new one
		OneUserID := user.ID
		AnotherUserID := query.ToUserID
		if OneUserID == AnotherUserID {
			return BadRequest("不能和自己聊天 :-)")
		}

		if OneUserID > AnotherUserID {
			OneUserID, AnotherUserID = AnotherUserID, OneUserID
		}

		var chat Chat
		if err = tx.Clauses(LockClause).
			Where("one_user_id = ? and another_user_id = ?", OneUserID, AnotherUserID).First(&chat).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				chat = Chat{
					OneUserID:     OneUserID,
					AnotherUserID: AnotherUserID,
				}
				if err = tx.Create(&chat).Error; err != nil {
					return
				}
			} else {
				return
			}
		}

		// create message
		message.ChatID = chat.ID
		if err = tx.Create(&message).Error; err != nil {
			return
		}

		// update chat message_count and last_message
		if err = tx.Model(&chat).Updates(Map{
			"message_count":        gorm.Expr("message_count + ?", 1),
			"last_message_content": message.Content,
			"last_message_id":      message.ID,
		}).Error; err != nil {
			return
		}

		return

	}); err != nil {
		return
	}

	// construct response
	var response MessageCommonResponse
	if err = copier.Copy(&response, &message); err != nil {
		return
	}
	response.IsOwner = true

	return Created(c, response)
}
