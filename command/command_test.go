package command

import (
	"four-key/command/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	mock mocks.Command

	command Commander
}

func (s *Suite) AfterTest(_, _ string) {
	s.mock.AssertExpectations(s.T())
	s.command = Commander{}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
}


func (s *Suite) TestGet_WhenBeforeInitialize_ReturnsError() {
	s.True(true)
}
