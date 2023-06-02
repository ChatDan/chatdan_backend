package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func testListTopics(t *testing.T) {
	const url = "/api/topics"
	var response Response[apis.TopicListResponse]
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url, 200, nil, &response)
	fmt.Printf("%v", response)
}

func testGetATopic(t *testing.T) {
	const url = "/api/topic/1"
	const url1 = "/api/topic/2"
	var response Response[apis.TopicCommonResponse]
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url, 200, nil, &response)
	userTester.testGet(t, url1, 404, nil, &response)

}

func testGetATopic2(t *testing.T) {
	const url = "/api/topic/1"
	const url1 = "/api/topic/2"
	var response Response[apis.TopicCommonResponse]
	defaultTester.testGet(t, url, 401, nil, &response)
	userTester.testGet(t, url1, 200, nil, &response)
	userTester.testGet(t, url, 404, nil, &response)

}

func testCreatATopic(t *testing.T) {
	const url = "/api/topic"
	var response Response[apis.TopicCommonResponse]
	body := Map{
		"title":        "TestTitle",
		"content":      "TestContent",
		"division_id":  1,
		"is_anonymous": false,
		"tags": []Map{
			{"name": "testTag1"},
			{"name": "testTag2"},
		},
	}
	defaultTester.testPost(t, url, 401, nil, &response)
	userTester.testPost(t, url, 201, body, &response)
	assert.EqualValues(t, "TestTitle", response.Data.Title)
}

func testModifyTopic(t *testing.T) {
	const url = "/api/topic/1"
	var response Response[apis.TopicCommonResponse]
	body := Map{
		"title":   "newTestTitle",
		"content": "newTestContent",
		"tags": []Map{
			{"name": "newTestTag1"},
			{"name": "newTestTag2"},
		},
	}
	defaultTester.testPut(t, url, 401, nil, &response)
	userTester.testPut(t, url, 200, body, &response)
	assert.EqualValues(t, "newTestTitle", response.Data.Title)
}

func testDeleteTopic(t *testing.T) {
	const url = "/api/topic/1"
	defaultTester.testDelete(t, url, 401, nil, nil)
	userTester.testDelete(t, url, 200, nil, nil)

}

func testLikeOrDislikeATopic(t *testing.T) {
	const url = "/api/topic/2/_like/1"
	const url2 = "/api/topic/3/_like/-1"
	var response Response[apis.TopicCommonResponse]

	defaultTester.testPut(t, url, 401, nil, &response)
	userTester.testPut(t, url, 200, nil, &response)
	userTester.testPut(t, url2, 200, nil, &response)
	log.Printf("%+v", response.Data)
}
