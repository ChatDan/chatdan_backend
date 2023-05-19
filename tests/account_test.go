package tests

import (
	"ChatDanBackend/apis"
	. "ChatDanBackend/models"
	"ChatDanBackend/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testAccountRegister(t *testing.T) {
	var data = Map{
		"username": "test",
		"password": "test123456",
	}
	var response utils.Response[apis.LoginResponse]
	defaultTester.testPost(t, "/api/user/register", 200, data, &response)
	assert.EqualValues(t, "test", response.Data.Username)
}

func testAccountLogin(t *testing.T) {
	var data = Map{
		"username": "test",
		"password": "test123456",
	}
	var response utils.Response[apis.LoginResponse]
	defaultTester.testPost(t, "/api/user/login", 200, data, &response)
	assert.EqualValues(t, "test", response.Data.Username)
	userTester.Token = response.Data.AccessToken
}
