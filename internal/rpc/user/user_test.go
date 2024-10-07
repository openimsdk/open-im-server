package user

import (
	"context"
	"fmt"
	"testing"

	pbuser "github.com/OpenIMSDK/protocol/user"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/user/mocks"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/user/service"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockUserDatabase = &mocks.UserDatabase{}

// var uServer userServer

// func inti() {

// }

func TestGetDesignateUsers(t *testing.T) {

	uServer := userServer{UserDatabase: mockUserDatabase}

	mockUserDatabase.On("FindWithError", mock.Anything, mock.Anything).Return([]*relation.UserModel{{UserID: "99"}}, nil)
	ctx := context.Background()

	res, err := uServer.GetDesignateUsers(ctx, &pbuser.GetDesignateUsersReq{UserIDs: []string{}})

	assert.NoError(t, err)
	assert.Equal(t, "99", res.UsersInfo[0].UserID)

	t.Logf("success")
}

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) DoSomeThing(number int) (bool, error) {
	args := m.Called(number)
	return args.Bool(0), args.Error(1)
}

func TestDataProvider(t *testing.T) {

	testOjb := new(MyMockedObject)
	testOjb.On("DoSomeThing", mock.Anything).Return(true, nil)

	targetFuncThatDoesSomethingWithObj(*testOjb)

	testOjb.AssertExpectations(t)

	///////

	mock := &mocks.DataProvider{}

	// Set the expected return value for the GetRandomNumber method
	mock.On("GetRandomNumber", 5).Return(3, nil).Once()
	// Call ConsumeData that using the mocked DataProvider
	result, err := ConsumeData(mock, 5)
	// Assert that the result and error are as expected
	assert.Equal(t, "Odd", result)
	assert.NoError(t, err)
	// Assert that the GetRandomNumber method was called with the expected input
	mock.AssertExpectations(t)
}

func targetFuncThatDoesSomethingWithObj(testObj MyMockedObject) {
	fmt.Println(testObj.DoSomeThing(12))
}

// ConsumeData is a function that uses a DataProvider to fetch and process data.
func ConsumeData(provider service.DataProvider, id int) (string, error) {
	// Use GetRandomNumber to get a random number between 0 and id
	randomNumber, err := provider.GetRandomNumber(id)
	if err != nil {
		return "", err
	}
	// Check whether the value is even or odd
	result := checkEvenOrOdd(randomNumber)
	// Return the result
	return result, nil
}

// checkEvenOrOdd checks whether the given value is even or odd.
func checkEvenOrOdd(value int) string {
	if value%2 == 0 {
		return "Even"
	}
	return "Odd"
}

//mockery --dir ../../../pkg/common/db/controller --name UserDatabase
//mockery --all --output path/to/output
//mockery --all --recursive
//mockery  --output ./mocks --dir ./service --all
