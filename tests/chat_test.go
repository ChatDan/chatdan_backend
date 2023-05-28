package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	"chatdan_backend/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testListChats(t *testing.T) {
	const url = "/api/chats"
	user0 := otherTester[0]
	user1 := otherTester[1]

	var response = utils.Response[apis.ChatListResponse]{}
	user0.testGet(t, url, 200, nil, &response)
	assert.EqualValues(t, 1, len(response.Data.Chats))
	assert.EqualValues(t, user1.ID, response.Data.Chats[0].AnotherUser.ID)
}

func testListMessages(t *testing.T) {
	const url = "/api/messages"
	user0 := otherTester[0]
	user1 := otherTester[1]

	data := Map{
		"to_user_id": user1.ID,
	}

	var response = utils.Response[apis.MessageListResponse]{}
	user0.testGet(t, url, 200, data, &response) // get in time reversed order
	assert.EqualValues(t, 2, len(response.Data.Messages))
	assert.EqualValues(t, "hello", response.Data.Messages[1].Content)
	assert.EqualValues(t, "world", response.Data.Messages[0].Content)
	assert.EqualValues(t, true, response.Data.Messages[1].IsOwner)
}

func testCreateMessage(t *testing.T) {
	const url = "/api/messages"
	user0 := otherTester[0]
	user1 := otherTester[1]

	data := Map{
		"content":    "hello",
		"to_user_id": user1.ID,
	}
	var response = utils.Response[apis.MessageCommonResponse]{}
	user0.testPost(t, url, 201, data, &response)
	assert.EqualValues(t, "hello", response.Data.Content)

	data = Map{
		"content":    "world",
		"to_user_id": user0.ID,
	}
	user1.testPost(t, url, 201, data, &response)
	assert.EqualValues(t, "world", response.Data.Content)
}
