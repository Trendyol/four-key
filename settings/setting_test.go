package settings

import (
	"errors"
	"four-key/command/mocks"
	"github.com/brianvoe/gofakeit/v5"
	_ "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"path"
	"path/filepath"
	"testing"
)

type Suite struct {
	suite.Suite
	mock mocks.Command

	settings Setting
}

func (s *Suite) AfterTest(_, _ string) {
	s.mock.AssertExpectations(s.T())
	settings = Settings{}
	isLoaded = false
	s.mock = mocks.Command{}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
}

func (s *Suite) TestGet_WhenBeforeInitialize_ReturnsError() {
	settings, err := Get()

	s.NotNil(err)
	s.Equal(err.Error(), "settings firstly must be initialized")
	s.Nil(settings)
}

func (s *Suite) TestGet_WhenAfterInitialize_ReturnsSettings() {
	isLoaded = true
	settings, err := Get()

	s.Nil(err)
	s.NotNil(settings)
}

func (s *Suite) TestInitialize_ReturnsSettings() {
	dir, err := os.Getwd()
	s.Nil(err)
	s.mock.On("GetFourKeyPath").Return(path.Join(dir, "/mock"))

	expectedRepoName := "four-key-metrics"
	expectedCloneAddress := "https://github.com/user/four-key-metrics.git"
	expectedReleaseTagPattern := "release-"
	expectedFixCommitPatterns := []string{"fix", "hot-fix", "hotfix"}
	expectedTeamName := "reform"

	err = Initialize(&s.mock)

	s.Nil(err)
	s.NotNil(settings)
	s.Equal(1, len(settings.Repositories))
	s.Equal(settings.Repositories[0].CloneAddress, expectedCloneAddress)
	s.Equal(settings.Repositories[0].FixCommitPatterns, expectedFixCommitPatterns)
	s.Equal(settings.Repositories[0].ReleaseTagPattern, expectedReleaseTagPattern)
	s.Equal(settings.Repositories[0].Name(), expectedRepoName)
	s.Equal(settings.Repositories[0].TeamName, expectedTeamName)
}

func (s *Suite) TestInitialize_WhenReaderReturnedError_CreatesNewConfigurationFileAndWritesDefaultTemplate() {
	dir, err := os.Getwd()
	s.Nil(err)

	fourKeyDir := gofakeit.BeerMalt()
	s.Nil(os.Mkdir(path.Join(dir, fourKeyDir), os.FileMode(0777)))

	s.mock.On("GetFourKeyPath").Return(path.Join(dir, fourKeyDir))
	s.mock.On("Open", path.Join(dir, fourKeyDir)).Return(nil)
	s.mock.On("Warn", "Your configurations not found!").Return("Your configurations not found!")
	s.mock.On("Warn", "Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName)).Return("Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName))
	s.mock.On("Good", "Configuration file added.").Return("Configuration file added.")
	s.mock.On("Good", "please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13").Return("please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13")

	err = Initialize(&s.mock)

	s.Nil(err)
	s.NotNil(settings)
	s.Equal(0, len(settings.Repositories))
	err = removeContents(path.Join(dir, fourKeyDir))
}

func (s *Suite) TestInitialize_ReturnsSettingsFromCache() {
	dir, err := os.Getwd()
	s.Nil(err)

	s.mock.On("GetFourKeyPath").Return(path.Join(dir, "/mock")).Twice()

	err = Initialize(&s.mock)
	err = Initialize(&s.mock)
	err = Initialize(&s.mock)
	err = Initialize(&s.mock)

	s.NotNil(settings)
	s.Equal(1, len(settings.Repositories))
}

func (s *Suite) TestInitialize_WhenFirstCreatingConfigurationFile_OpensFourKeyDirectory() {
	dir, err := os.Getwd()
	s.Nil(err)

	fourKeyDir := gofakeit.BeerName()
	s.Nil(os.Mkdir(path.Join(dir, fourKeyDir), os.FileMode(0777)))

	s.mock.On("GetFourKeyPath").Return(path.Join(dir, fourKeyDir))
	s.mock.On("Open", path.Join(dir, fourKeyDir)).Return(nil)
	s.mock.On("Warn", "Your configurations not found!").Return("Your configurations not found!")
	s.mock.On("Warn", "Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName)).Return("Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName))
	s.mock.On("Good", "Configuration file added.").Return("Configuration file added.")
	s.mock.On("Good", "please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13").Return("please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13")
	s.mock.On("Open", path.Join(dir, fourKeyDir)).Return(nil)

	err = Initialize(&s.mock)

	s.Nil(err)
	s.NotNil(settings)
	s.Equal(0, len(settings.Repositories))
	err = removeContents(path.Join(dir, fourKeyDir))
}

func (s *Suite) TestInitialize_WhenFirstCreatingConfigurationFileIfOpenReturnsError_LogsWarningOpeningDir() {
	openError := errors.New("open error")
	dir, err := os.Getwd()
	s.Nil(err)

	fourKeyDir := gofakeit.Adverb()
	s.Nil(os.Mkdir(path.Join(dir, fourKeyDir), os.FileMode(0777)))

	s.mock.On("GetFourKeyPath").Return(path.Join(dir, fourKeyDir))
	s.mock.On("Warn", "Your configurations not found!").Return("Your configurations not found!")
	s.mock.On("Warn", "Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName)).Return("Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName))
	s.mock.On("Warn", "Configuration file directory not opened", openError.Error()).Return("Configuration file directory not opened", openError.Error())
	s.mock.On("Good", "Configuration file added.").Return("Configuration file added.")
	s.mock.On("Good", "please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13").Return("please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13")
	s.mock.On("Open", path.Join(dir, fourKeyDir)).Return(openError)

	err = Initialize(&s.mock)

	s.Nil(err)
	s.NotNil(settings)
	s.Equal(0, len(settings.Repositories))
	err = removeContents(path.Join(dir, fourKeyDir))
}

func (s *Suite) TestInitialize_IfReturnsErrorWhenCreatingFourKeyDir_ReturnsError() {
	dir, err := os.Getwd()
	s.Nil(err)
	fourKeyDir := "\\//$+'"

	s.mock.On("GetFourKeyPath").Return(path.Join(dir, fourKeyDir))
	s.mock.On("Warn", "Your configurations not found!").Return("Your configurations not found!")
	s.mock.On("Warn", "Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName)).Return("Generating configuration file to -> ", path.Join(dir, fourKeyDir, EnvironmentFileName))
	s.mock.On("Fatal", "An error occurred while creating four-key.json to ", path.Join(dir, fourKeyDir, EnvironmentFileName)).Return("An error occurred while creating four-key.json to ", path.Join(dir, fourKeyDir, EnvironmentFileName))
	s.mock.On("Fatal", "Configurations not loaded").Return("Configurations not loaded")

	err = Initialize(&s.mock)

	s.NotNil(err)
	s.Equal(0, len(settings.Repositories))
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return os.Remove(dir)
}
