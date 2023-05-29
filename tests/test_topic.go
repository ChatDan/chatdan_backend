package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"testing"
)

func testListTopics(t *testing.T) {
	const url = "/api/topics"
	var response Response[apis.TopicListResponse]
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url, 200, nil, &response)

}

func testGetATopic(t *testing.T) {
	const url = "/api/topic/1"
	const url1 = "/api/topic/2"
	var response Response[apis.TopicCommonResponse]
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url, 200, nil, &response)
	userTester.testGet(t, url1, 404, nil, &response)

}

func testCreatATopic(t *testing.T) {
	const url = "/api/topic"
	var response Response[apis.TopicCommonResponse]
	body := Map{
		"title":        "TestTile",
		"content":      "TestContent",
		"division_id":  1,
		"is_anonymous": false,
		"tags":         []string{"aaaa"},
	}
	defaultTester.testPost(t, url, 401, nil, &response)
	userTester.testPost(t, url, 201, body, &response)

}
