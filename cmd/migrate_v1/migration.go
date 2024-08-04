package migrate_v1

import (
	"context"
	"embed"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"plugin"
	"time"
)

//go:embed template.tmpl

var migrateTemplate embed.FS

type MigrateInterface interface {
	MigrateUp(ctx context.Context, db *gorm.DB) error
	MigrateDown(ctx context.Context, db *gorm.DB) error
}

type Migration struct {
	SourcePath   []string `json:"source"`
	OutputPath   string   `json:"output"`
	Rebuild      bool
	BuildCommand string
}

type MigrateTable struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Batch     int8      `json:"batch"`
}

func (t MigrateTable) TableName() string {
	return "migration"
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

func findMigration(file string) (res MigrateInterface, err error) {
	fn, err := plugin.Open(file)
	if err != nil {
		return nil, err
	}

	symbol, err := fn.Lookup("Migrate")
	if err != nil {
		return nil, err
	}

	res, ok := symbol.(MigrateInterface)
	if !ok {
		return nil, fmt.Errorf("invalid migration interface")
	}
	return
}
