package helpers

import (
	"errors"
	"fmt"
	Command "four-key/command"
	"four-key/settings"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"path"
	"regexp"
	"strings"
)

type RepositoryWrapper struct {
	Repository     *git.Repository
	Configurations settings.Repository
}

type IRepositoryHelper interface {
	GetRepositories(s *settings.Settings) ([]RepositoryWrapper, error)
}

func GetRepositories(s *settings.Settings) ([]RepositoryWrapper, error) {
	var repositoriesWrapper []RepositoryWrapper

	for _, repo := range s.Repositories {
		err := CheckDirectory(s.RepositoriesPath, repo.Name())

		if err != nil {
			err := CloneRepository(repo.CloneAddress, s.RepositoriesPath)
			if err != nil {
				return repositoriesWrapper, err
			}
		}

		repository, err := git.PlainOpen(path.Join(Command.ACommander().GetRepositoriesPath(s.RepositoriesPath), repo.Name()))

		if err != nil {
			fmt.Println(Command.ACommander().Warn(err))
			return nil, err
		}

		wrapper := RepositoryWrapper{
			Repository:     repository,
			Configurations: repo,
		}

		repositoriesWrapper = append(repositoriesWrapper, wrapper)
	}

	return repositoriesWrapper, nil
}

func GetRepositoryByName(s *settings.Settings, repositoryName string) (RepositoryWrapper, error) {
	var w RepositoryWrapper

	for _, repo := range s.Repositories {
		if repo.Name() == repositoryName {
			err := CheckDirectory(s.RepositoriesPath, repositoryName)

			if err != nil {
				err := CloneRepository(repo.CloneAddress, s.RepositoriesPath)
				if err != nil {
					return w, err
				}
			}
			repository, err := git.PlainOpen(path.Join(Command.ACommander().GetRepositoriesPath(s.RepositoriesPath), repo.Name()))

			if err != nil {
				fmt.Println(err)
				return w, err
			}

			w.Repository = repository
			w.Configurations = repo

			return w, nil
		}
	}

	return w, errors.New("repository not found with that given name -> " + repositoryName)
}

func RepoCheck(r *git.Repository) {
	_, err := r.Log(&git.LogOptions{
		From:     plumbing.Hash{},
		Order:    0,
		FileName: nil,
		All:      false,
	})

	w, err := r.Worktree()
	if err != nil {
		panic(fmt.Sprintf("An error occurred. Error: %v", err))
	}

	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Tags:       2,
		Force:      true,
	})
	err = w.Pull(&git.PullOptions{RemoteName: "origin", Force: true})

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		panic(err)
	}

	_, err = r.CommitObject(ref.Hash())

	if err != nil {
		panic(fmt.Sprintf("An error occurred. Error: %v", err))
	}

}

func CloneRepository(cloneLink string, p string) error {
	fmt.Println(Command.ACommander().Good("Cloning repository - clone address -> ", cloneLink))

	err := Command.ACommander().Command("git clone "+cloneLink+" --progress", Command.ACommander().GetRepositoriesPath(p))

	if err != nil {
		return err
	}

	fmt.Println(Command.ACommander().Good("Cloned repository ", GetNameByRepositoryCloneUrl(cloneLink)))
	return nil
}

func CreateDirectory(path string, name string) error {
	err := Command.ACommander().Command("mkdir "+name, path)

	if err != nil {
		return err
	}

	return nil
}

func CheckDirectory(args ...string) error {
	paths := append([]string{Command.ACommander().GetFourKeyPath()}, args...)
	return Command.ACommander().Command("pwd", path.Join(paths...))
}

func GetNameByRepositoryCloneUrl(cloneAddress string) string {
	re := regexp.MustCompile(`([^/]+)\.git$`)
	return strings.Replace(re.FindString(cloneAddress), ".git", "", 1)
}
