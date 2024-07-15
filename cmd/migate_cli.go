package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration commands",
}

// migrateNew represents the migrate new command
var migrateNew = &cobra.Command{
	Use:   "new <migration_name>",
	Short: "Create a new migration file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("❌ require argument migration_name, example: axara migrate new create_users_table")
			return
		}
		// Add logic to create new migration file here
		fmt.Printf("Creating new migration file: %s\n", args[0])
	},
}

// migrateUp represents the migrate up command
var migrateUp = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// Add logic to apply all up migrations here
		fmt.Println("Applying all up migrations")
	},
}

// migrateDown represents the migrate down command
var migrateDown = &cobra.Command{
	Use:   "down",
	Short: "Roll back the last batch of migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// Add logic to roll back the last batch of migrations here
		fmt.Println("Rolling back the last batch of migrations")
	},
}

// migrateFresh represents the migrate fresh command
var migrateFresh = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// Add logic to drop all tables and re-run all migrations here
		fmt.Println("Dropping all tables and re-running all migrations")
	},
}

// migrateFlush represents the migrate flush command
var migrateFlush = &cobra.Command{
	Use:   "flush",
	Short: "Reset the database by rolling back and re-applying all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// Add logic to reset the database here
		fmt.Println("Resetting the database by rolling back and re-applying all migrations")
	},
}

// migrateBuild represents the migrate build command
var migrateBuild = &cobra.Command{
	Use:   "build [flags] [--source-path <source_path>]... [--output-path <output_path>]",
	Short: "Build a plugin from migration files (*.go)",
	Run: func(cmd *cobra.Command, args []string) {
		//if len(args) < 1 {
		//	logrus.Fatal("require argument migration_files_directory, example: axara migrate build ./migrations ./migrations_output")
		//}

		// Check if 'go' command is available
		goPath, err := exec.LookPath("go")
		if err != nil {
			logrus.Fatal("'go' command not found in PATH. Please ensure Go is installed and added to PATH.")
		}

		// Validate source paths
		if len(sourcePaths) == 0 {
			logrus.Fatal("at least one --source-path must be provided")
		}

		// Determine output directory
		if outputPath == "" {
			outputPath = "./migrate-out" // Default output directory
		}

		// Create the output directory if it doesn't exist'
		err = os.MkdirAll(outputPath, 0755)
		if err != nil {
			logrus.Fatalf("Error creating output directory %s: %v", outputPath, err)
		}

		// Iterate through each source path
		for _, sourcePath := range sourcePaths {
			var goFiles []string
			var sourceDir string

			// Check if sourcePath is a Git repository URL
			if strings.HasPrefix(sourcePath, "https://") || strings.HasPrefix(sourcePath, "http://") || strings.HasPrefix(sourcePath, "git://") {
				sourceDir, err = downloadGitRepository(sourcePath)
				if err != nil {
					logrus.Fatalf("Error downloading Git repository %s: %v", sourcePath, err)
				}
			} else {
				sourceDir = sourcePath
			}

			// Get all .go files in the source directory
			goFiles, err = filepath.Glob(filepath.Join(sourceDir, "*.go"))
			if err != nil {
				logrus.Fatalf("Error finding .go files in directory %s: %v", sourceDir, err)
			}

			if len(goFiles) == 0 {
				logrus.Fatalf("No .go files found in directory: %s", sourceDir)
			}

			// Build each .go file into a .so plugin
			for _, file := range goFiles {
				baseName := filepath.Base(file)
				outputName := strings.TrimSuffix(baseName, ".go") + ".so"
				outputFilePath := filepath.Join(outputPath, outputName)

				buildCmd := []string{"build", "-buildmode=plugin", "-o", outputFilePath}
				buildCmd = append(buildCmd, strings.Fields(buildCommand)...)
				buildCmd = append(buildCmd, file)

				xcmd := exec.Command(goPath, buildCmd...)
				xcmd.Dir = sourceDir

				output, err := xcmd.CombinedOutput()
				if err != nil {
					logrus.Fatalf("Error building migration from %s: %v\n%s", file, err, output)
				}

				fmt.Printf("Built migration: %s\n", outputPath)
			}

			// Clean up downloaded Git repository if applicable
			if sourcePath != sourceDir {
				err := os.RemoveAll(sourceDir)
				if err != nil {
					logrus.Warnf("Failed to clean up temporary directory %s: %v", sourceDir, err)
				}
			}
		}

		fmt.Println("Plugins built successfully.")
	},
}

func downloadGitRepository(repoURL string) (string, error) {
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
