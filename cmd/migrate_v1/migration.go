package migrate_v1

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
)

type Migration struct {
	SourcePath   []string `json:"source"`
	OutputPath   string   `json:"output"`
	Rebuild      bool
	BuildCommand string
}

func NewMigration(sourcePath []string, outputPath string, rebuild bool, buildCommand string) *Migration {
	return &Migration{SourcePath: sourcePath, OutputPath: outputPath, Rebuild: rebuild, BuildCommand: buildCommand}
}

func (m *Migration) downloadGitRepository(repoURL string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "axara-git-*")
	if err != nil {
		return "", err
	}

	repoDir := filepath.Join(tmpDir, "repository")
	_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return "", err
	}

	return repoDir, nil
}
