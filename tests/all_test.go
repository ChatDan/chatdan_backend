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
}

func BenchmarkAll(b *testing.B) {
	registerOnce(b)

	b.Run("BenchAccountLogin", benchAccountLogin)
}
