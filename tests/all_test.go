package tests

import "testing"

func TestAll(t *testing.T) {
	t.Run("TestAccountRegister", testAccountRegister)
	t.Run("TestAccountLogin", testAccountLogin)
	t.Run("TestListBoxes", testListBoxes)
	t.Run("TestCreateABox", testCreateABox)
}
