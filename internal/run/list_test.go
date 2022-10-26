package run

import (
	"github.com/golang/mock/gomock"
	"github.com/mongodb-forks/drone-helm3/internal/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type ListTestSuite struct {
	suite.Suite
}

func TestLIstTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}

func (suite *ListTestSuite) TestNewList() {
	cfg := env.Config{
		Command: "everybody dance NOW!!",
	}
	list := NewList(cfg)
	suite.Require().NotNil(list)
	suite.Equal("everybody dance NOW!!", list.helmCommand)
}

func (suite *ListTestSuite) TestPrepare() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mCmd := NewMockcmd(ctrl)
	originalCommand := command

	command = func(path string, args ...string) cmd {
		assert.Equal(suite.T(), helmBin, path)
		assert.Equal(suite.T(), []string{"list", "--output", "json"}, args)
		return mCmd
	}
	defer func() { command = originalCommand }()

	stdout := strings.Builder{}
	stderr := strings.Builder{}

	mCmd.EXPECT().
		Stdout(&stdout)
	mCmd.EXPECT().
		Stderr(&stderr)

	cfg := env.Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	h := NewList(cfg)
	err := h.Prepare()
	suite.NoError(err)
}

func (suite *ListTestSuite) TestExecute() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mCmd := NewMockcmd(ctrl)

	mCmd.EXPECT().
		Run().
		Times(2)

	list := NewList(env.Config{Command: "list"})
	list.cmd = mCmd
	suite.NoError(list.Execute())

	list.helmCommand = "get down on friday"
	suite.EqualError(list.Execute(), "unknown command 'get down on friday'")
}
