package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	Command "four-key/command"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

type Setting interface {
	Load() error
}

var settings Settings
var isLoaded = false

const TemplateConfig = `{"repositories":[]}`
const DefaultTeamName = "master"
const EnvironmentFileName = "four-key.json"
const DefaultRepositoryDirName = "repos"
const AllTeamsDefaultDirName = "allTeams"
const TeamBasedDefaultDirName = "teamBased"
const DefaultGeneratedFileOutputDirName = "metrics"
const DefaultDateFormat = "2006-01-02"

type Configuration struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Settings struct {
	Output           string
	Repositories     []Repository
	RepositoriesPath string
	commander        Command.ICommand
}

type Repository struct {
	CloneAddress      string   `json:"cloneAddress"`
	TeamName          string   `json:"teamName"`
	ReleaseTagPattern string   `json:"releaseTagPattern"`
	FixCommitPatterns []string `json:"fixCommitPatterns"`
}

func (r *Repository) Name() string {
	rx := regexp.MustCompile(`([^/]+)\.git$`)
	return strings.Replace(rx.FindString(r.CloneAddress), ".git", "", 1)
}

func (s *Settings) Load() error {
	cfg, err := ioutil.ReadFile(path.Join(s.commander.GetFourKeyPath(), EnvironmentFileName))

	if err != nil {
		fmt.Println(s.commander.Warn("Your configurations not found!"))
		fmt.Println(s.commander.Warn("Generating configuration file to -> ", path.Join(s.commander.GetFourKeyPath(), EnvironmentFileName)))

		f, err := os.Create(path.Join(s.commander.GetFourKeyPath(), EnvironmentFileName))
		if err != nil {
			fmt.Println(s.commander.Fatal("An error occurred while creating four-key.json to ", path.Join(s.commander.GetFourKeyPath(), EnvironmentFileName)))
			return err
		}

		_, err = f.WriteString(TemplateConfig)

		fmt.Println(s.commander.Good("Configuration file added."))
		fmt.Println(s.commander.Good("please add an repository and run command like -> ./four-key run -s 2018-01-13 -e 2021-01-13"))

		err = s.commander.Open(s.commander.GetFourKeyPath())
		if err != nil {
			fmt.Println(s.commander.Warn("Configuration file directory not opened", err.Error()))
		}

		cfg, err = ioutil.ReadFile(path.Join(s.commander.GetFourKeyPath(), EnvironmentFileName))
	}

	jsonErr := json.Unmarshal(cfg, &settings)

	if jsonErr != nil {
		return errors.New(fmt.Sprintf("An error occurred. Error: %v", jsonErr))
	}

	if settings.RepositoriesPath == "" {
		settings.RepositoriesPath = DefaultRepositoryDirName
	}

	if settings.Output == "" {
		settings.Output = s.commander.GetFourKeyPath()
	}

	return nil
}

func Initialize(cmd Command.ICommand) error {
	if isLoaded != true {
		settings.commander = cmd
		err := settings.Load()

		if err != nil {
			return err
		}

		isLoaded = true
	}

	return nil
}

func Get() (*Settings, error) {
	if isLoaded == true {
		return &settings, nil
	}

	return nil, errors.New("settings firstly must be initialized")
}
