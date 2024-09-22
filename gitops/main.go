package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
)

func main() {
	timeSec := 5 * time.Second
	gitopsRepo := "https://github.com/txuna/argocd-tutorial.git"
	localPath := "/tmp/argocd-tutorial"
	pathToApply := "ch01"

	for {
		fmt.Println("start repo sync")
		err := syncRepo(gitopsRepo, localPath)
		if err != nil {
			fmt.Printf("repo sync error: %s\n", err)
			return
		}

		fmt.Println("start manifest apply")
		err = applyManifestsClient(path.Join(localPath, pathToApply))
		if err != nil {
			fmt.Printf("manifest apply error: %s\n", err)
		}

		syncTimer := time.NewTimer(timeSec)
		fmt.Printf("\n next sync in %s \n", timeSec)
		<-syncTimer.C
	}
}

func syncRepo(repoUrl, localPath string) error {
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})

	if err == git.ErrRepositoryAlreadyExists {
		repo, err := git.PlainOpen(localPath)
		if err != nil {
			return err
		}

		w, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		})

		if err == git.NoErrAlreadyUpToDate {
			return nil
		}

		return err
	}

	return err
}

func applyManifestsClient(localPath string) error {
	cmd := exec.Command("kubectl", "apply", "-f", localPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
