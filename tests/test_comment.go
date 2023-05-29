package tests

import (
	"chatdan_backend/apis"
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"testing"
)

func testListComments(t *testing.T) {
	const url = "/api/comments"
	var response Response[apis.CommentListResponse]

	query1 := Map{
		"topic_id": 1,
		"order_by": "id",
	}

	query2 := Map{
		"topic_id": 2,
		"order_by": "id",
	}

	defaultTester.testGet(t, url, 401, query2, &response)
	userTester.testGet(t, url, 200, query2, &response)
	userTester.testGet(t, url, 404, query1, &response)

}

func testCreateComment(t *testing.T) {
	const url = "/api/comment"
	var response Response[apis.CommentCommonResponse]

	body := Map{
		"topic_id":     2,
		"content":      "Test Comment",
		"is_anonymous": false,
	}

	body2 := Map{
		"topic_id":     2,
		"content":      "Test Comment",
		"is_anonymous": true,
	}

	defaultTester.testPost(t, url, 401, body, &response)
	userTester.testPost(t, url, 200, body, &response)
	userTester.testPost(t, url, 200, body2, &response)

}
