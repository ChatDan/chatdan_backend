package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"testing"
)

func testListDivisions(t *testing.T) {
	const url = "/api/divisions"

	var response Response[apis.DivisionListResponse]
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url, 200, nil, &response)

	//DB.Model(User{ID: 1}).Update("is_admin", true)
}

//func testCreateDivision(t *testing.T) {
//	const url = "/api/division"
//	var response Response[apis.DivisionCreateRequest]
//	defaultTester.testPost(t, url, 401, nil, &response)
//	userTester.testPost(t, url, 403, nil, &response)
//
//	DB.Model(User{ID: 1}).Update("is_admin", true)
//
//	userTester.testPost(t, url, 201, nil, &response)
//}

func testGetADivision(t *testing.T) {
	const url = "/api/division/1"
	var response Response[apis.DivisionCreateRequest]
	division := Division{Name: "a"}
	DB.FirstOrCreate(&division)
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url, 200, nil, &response)

}
