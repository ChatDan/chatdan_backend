package tests

import "testing"

func TestAll(t *testing.T) {
	t.Run("TestAccountRegister", testAccountRegister)
	t.Run("TestAccountLogin", testAccountLogin)
	t.Run("TestListBoxes", testListBoxes)
	t.Run("TestCreateABox", testCreateABox)

	// chat
	t.Run("TestCreateMessage", testCreateMessage)
	t.Run("TestListMessages", testListMessages)
	t.Run("TestListChats", testListChats)

	//division
	t.Run("TestListDivisions", testListDivisions)
	//t.Run("TestCreateDivision", testCreateDivision)
	t.Run("TestGetADivision", testGetADivision)

	//topic
	t.Run("TestCreateTopic", testCreatATopic)
	t.Run("TestListTopics", testListTopics)
	t.Run("TestGetATopic", testGetATopic)
	t.Run("TestModifyTopic", testModifyTopic)
	t.Run("TestDeleteTopic", testDeleteTopic)
	t.Run("TestCreateTopic2", testCreatATopic)
	t.Run("TestGetATopic2", testGetATopic2)
	t.Run("TestCreateTopic3", testCreatATopic)

	t.Run("TestLikeOrDislikeATopic", testLikeOrDislikeATopic)

	//comment
	t.Run("TestCreateComment", testCreateComment)
	t.Run("TestListComments", testListComments)
}

func BenchmarkAll(b *testing.B) {
	registerOnce(b)

	b.Run("BenchAccountLogin", benchAccountLogin)
}
