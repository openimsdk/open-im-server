package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
)

func TestDeleteGroupMemberHash(t *testing.T) {
	mockGroupDB := new(controller.MockGroupDatabase)

	testGroupMemberHash := "testGroupMemberHash"

	err := mockGroupDB.DeleteGroupMemberHash(testGroupMemberHash)
	assert.Nil(t, err)

	nonExistentGroupMemberHash := "nonExistentGroupMemberHash"

	err = mockGroupDB.DeleteGroupMemberHash(nonExistentGroupMemberHash)
	assert.NotNil(t, err)
}
