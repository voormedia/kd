package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func findGitDir(dir string) (string, error) {
	for {
		gitPath := filepath.Join(dir, ".git", "HEAD")
		if _, err := os.Stat(gitPath); err == nil {
			return gitPath, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	return "", fmt.Errorf("failed to find git directory")
}

func GetCurrentBranch(log *Logger, path string) (string, error) {
	log.Debug("Reading current git branch from .git/HEAD")

	gitHeadPath, err := findGitDir(path)
	if err != nil {
		return "", err
	}

	headContent, err := os.ReadFile(gitHeadPath)
	if err != nil {
		return "", fmt.Errorf("error reading .git/HEAD: %w", err)
	}

	// Check if HEAD points to a branch reference
	if strings.HasPrefix(string(headContent), "ref: ") {
		branchPath := strings.TrimSpace(strings.TrimPrefix(string(headContent), "ref: "))
		branchName := filepath.Base(branchPath)
		return branchName, nil
	}

	return "", fmt.Errorf("error parsing git ref")
}
