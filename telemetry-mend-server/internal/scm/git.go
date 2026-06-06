package scm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitClient struct {
	BaseDir string
}

func NewGitClient(baseDir string) (*GitClient, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &GitClient{BaseDir: baseDir}, nil
}

func (c *GitClient) CloneOrFetch(ctx context.Context, repoURL string, appID int64) (string, error) {
	repoPath := filepath.Join(c.BaseDir, fmt.Sprintf("app-%d", appID))

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		_, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return "", err
		}
	} else {
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			return "", err
		}
		err = repo.FetchContext(ctx, &git.FetchOptions{
			Progress: os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return "", err
		}
	}

	return repoPath, nil
}

func (c *GitClient) GetFileContent(repoPath, commitHash, filePath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	hash := plumbing.NewHash(commitHash)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	file, err := tree.File(filePath)
	if err != nil {
		return "", err
	}

	reader, err := file.Reader()
	if err != nil {
		return "", err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
