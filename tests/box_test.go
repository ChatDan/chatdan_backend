package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	"chatdan_backend/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testListBoxes(t *testing.T) {
	const url = "/api/messageBoxes"

	var response utils.Response[apis.BoxListResponse]
	defaultTester.testGet(t, url, 401, nil, &response) // 401 Unauthorized
	userTester.testGet(t, url, 400, nil, &response)    // 没有分页

	data := Map{
		"page_num":  1,
		"page_size": 10,
	}
	userTester.testGet(t, url, 200, data, &response)
}

func testCreateABox(t *testing.T) {
	const url = "/api/messageBox"

	var response utils.Response[apis.BoxCommonResponse]
	defaultTester.testPost(t, url, 401, nil, &response) // 401 Unauthorized

	data := Map{
		"title": "123",
	}
	userTester.testPost(t, url, 201, data, &response)
	assert.EqualValues(t, "123", response.Data.Title)
}
