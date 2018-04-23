package executer_test

import (
	"errors"
	"testing"

	"github.com/segurosfalabella/imperium-worker/executer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCmd struct {
	mock.Mock
}

func (mockCmd *MockCmd) Run() error {
	args := mockCmd.Called()
	return args.Error(0)
}

func TestShouldExecuteFailWhenSpawnlingCommandFail(t *testing.T) {
	job := new(executer.Job)
	oldCreateCommand := executer.CreateCommand
	defer func() { executer.CreateCommand = oldCreateCommand }()
	mock := new(MockCmd)
	mock.On("Run").Return(errors.New("527c090d-4102-4671-9033-b3363f78b343"))
	executer.CreateCommand = func(name string, arg ...string) executer.Commander {
		return mock
	}

	err := job.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, "527c090d-4102-4671-9033-b3363f78b343", err.Error())
}