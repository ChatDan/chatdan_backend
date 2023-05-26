package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	"chatdan_backend/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
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

	// test duplicate username
	defaultTester.testPost(t, "/api/user/register", 400, data, &response)

	// 注册10个用户，从 user0 到 user9
	for i := 0; i < 10; i++ {
		data["username"] = "user" + strconv.Itoa(i)
		defaultTester.testPost(t, "/api/user/register", 200, data, &response)
		assert.EqualValues(t, "user"+strconv.Itoa(i), response.Data.Username)
		otherTester[i] = tester{Token: response.Data.AccessToken, ID: response.Data.ID}
	}
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
	userTester.ID = response.Data.ID
}
