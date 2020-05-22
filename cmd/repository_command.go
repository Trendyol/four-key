package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	Command "four-key/command"
	"four-key/helpers"
	"four-key/models"
	"four-key/settings"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"text/tabwriter"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add repository",
	Long:  "Add repository",
	Run:   onAddRepository,
}

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove repository",
	Long:  "Remove repository",
	Run:   onRemoveRepository,
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	Long:  "List repositories",
	Run:   onListRepositories,
}

func init() {
	var empty []string
	rootCmd.AddCommand(addCommand)
	settings.Initialize(Command.ACommander())
	addCommand.Flags().StringP("cloneAddress", "c", "", "Set your clone address")
	addCommand.Flags().StringP("team", "t", "", "Set your team of repository")
	addCommand.Flags().StringP("releaseTagPattern", "r", "", "Set your release tag pattern of repository")
	addCommand.Flags().StringArrayP("fixCommitPatterns", "f", empty, "Set your fix commit patterns of repository")

	rootCmd.AddCommand(removeCommand)
	removeCommand.Flags().StringP("repository", "r", "", "Set your repository name to remove from config")

	rootCmd.AddCommand(listCommand)
}

func writeFile(document models.Document) error {
	doc, err := json.Marshal(document)
	if err != nil {
		return errors.New("json convert error")
	}

	err = ioutil.WriteFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName), doc, os.FileMode(0644))
	if err != nil {
		return errors.New("file write error")
	}

	return nil
}

func onAddRepository(cmd *cobra.Command, args []string) {
	cloneAddress, err := cmd.Flags().GetString("cloneAddress")
	team, err := cmd.Flags().GetString("team")
	tagPattern, err := cmd.Flags().GetString("releaseTagPattern")
	commitPatterns, err := cmd.Flags().GetStringArray("fixCommitPatterns")

	if err != nil {
		fmt.Println(Command.ACommander().Fatal("an error occurred while adding repository, please check entered inputs."))
		return
	}

	document := models.Document{}

	fmt.Println(Command.ACommander().Good("Reading your configuration file..."))
	existFile, err := ioutil.ReadFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName))

	if err != nil {
		fmt.Println(Command.ACommander().Fatal("The file does not exist!"))

		existFile, err = ioutil.ReadFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName))

		if err != nil {
			fmt.Println(Command.ACommander().Fatal("The file does not exist!"), err)
			return
		}
	}

	err = json.Unmarshal(existFile, &document)
	if err != nil {
		fmt.Println(Command.ACommander().Fatal("json parse error."))
		return
	}

	document.Repositories = append(document.Repositories, &models.DocumentRepository{
		TeamName:          team,
		CloneAddress:      cloneAddress,
		ReleaseTagPattern: tagPattern,
		FixCommitPatterns: commitPatterns,
	})

	err = writeFile(document)
	if err != nil {
		fmt.Println(Command.ACommander().Fatal(err.Error()))
		return
	}

	fmt.Println(Command.ACommander().Good("successfully added your repository to config file."))

	s, err = settings.Get()

	if err != nil {
		fmt.Println(Command.ACommander().Fatal(err))
	}

	err = helpers.CloneRepository(cloneAddress, s.RepositoriesPath)

	if err != nil {
		fmt.Println(Command.ACommander().Fatal(err))
		return
	}
}

func onRemoveRepository(cmd *cobra.Command, args []string) {
	repository, err := cmd.Flags().GetString("repository")
	if err != nil {
		fmt.Println(Command.ACommander().Fatal("an error occurred while removing repository, please check entered inputs."))
		return
	}

	document := models.Document{}

	fmt.Println(Command.ACommander().Good("Reading your configuration file..."))
	existFile, err := ioutil.ReadFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName))

	if err != nil {
		fmt.Println(Command.ACommander().Fatal("The file does not exist!"))
		existFile, err = ioutil.ReadFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName))

		if err != nil {
			fmt.Println(Command.ACommander().Fatal("The file does not exist!"), err)
			return
		}
	}

	err = json.Unmarshal([]byte(string(existFile)), &document)
	if err != nil {
		fmt.Println(Command.ACommander().Fatal("json parse error."))
		return
	}

	isRemoved := false
	var newRepositories []*models.DocumentRepository
	for _, documentRepository := range document.Repositories {
		if repository != helpers.GetNameByRepositoryCloneUrl(documentRepository.CloneAddress) {
			newRepositories = append(newRepositories, documentRepository)
		} else {
			isRemoved = true
		}
	}

	document.Repositories = newRepositories

	err = writeFile(document)
	if err != nil {
		fmt.Println(Command.ACommander().Fatal(err.Error()))
		return
	}

	if isRemoved {
		fmt.Println(Command.ACommander().Good(fmt.Sprintf("successfully removed %s repository from the config file.", repository)))
	} else {
		fmt.Println(Command.ACommander().Fatal(fmt.Sprintf("The %s repository does not exist!", repository)))
	}
}

func onListRepositories(cmd *cobra.Command, args []string) {
	document := models.Document{}

	fmt.Println(Command.ACommander().Good("Reading your configuration file..."))
	existFile, err := ioutil.ReadFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName))

	if err != nil {
		fmt.Println(Command.ACommander().Fatal("The file does not exist!"))
		existFile, err = ioutil.ReadFile(path.Join(Command.ACommander().GetFourKeyPath(), settings.EnvironmentFileName))

		if err != nil {
			fmt.Println(Command.ACommander().Fatal("The file does not exist!"), err)
			return
		}
	}

	err = json.Unmarshal(existFile, &document)
	if err != nil {
		fmt.Println(Command.ACommander().Fatal("json parse error."))
		return
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 40, 10, 0, '\t', 0)

	fmt.Println(fmt.Sprintf("\nTotal %d repository/repositories has been found.", len(document.Repositories)))

	sort.Slice(document.Repositories, func(i, j int) bool {
		return document.Repositories[i].TeamName < document.Repositories[j].TeamName
	})

	for i, repository := range document.Repositories {
		if i == 0 {
			_, _ = fmt.Fprintln(w, fmt.Sprintf("\nTeam: %s", repository.TeamName))
		}

		_, _ = fmt.Fprintln(w, fmt.Sprintf("%d. %s\t%s", i+1, helpers.GetNameByRepositoryCloneUrl(repository.CloneAddress), repository.CloneAddress))
		if i < (len(document.Repositories)-1) && repository.TeamName != document.Repositories[i+1].TeamName {
			_, _ = fmt.Fprintln(w, fmt.Sprintf("\nTeam: %s", document.Repositories[i+1].TeamName))
		}
	}
	_ = w.Flush()
}
