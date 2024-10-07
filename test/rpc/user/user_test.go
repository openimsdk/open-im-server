package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) DoSomeThing(number int) (bool, error) {
	args := m.Called(number)
	return args.Bool(0), args.Error(1)
}

// var uServer = userServer{}

func TestGetDesignateUsers(t *testing.T) {
	testOjb := new(MyMockedObject)
	testOjb.On("DoSomeThing", mock.Anything).Return(true, nil)

	targetFuncThatDoesSomethingWithObj(*testOjb)

	testOjb.AssertExpectations(t)

	t.Logf("success")
}

func targetFuncThatDoesSomethingWithObj(testObj MyMockedObject) {
	fmt.Println(testObj.DoSomeThing(12))
}
