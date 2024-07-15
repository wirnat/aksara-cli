package migrate_v1

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *Migration) Build() {
	//if len(args) < 1 {
	//	logrus.Fatal("require argument migration_files_directory, example: axara migrate build ./migrations ./migrations_output")
	//}
	outputPath := m.OutputPath
	sourcePaths := m.SourcePath
	buildCommand := m.BuildCommand
	rebuild := m.Rebuild

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

	// Track processed source paths to skip duplicates
	processedPaths := make(map[string]bool)
	// Track processed file names to detect duplicates
	processedFiles := make(map[string]bool)

	// Iterate through each source path
	for _, sourcePath := range sourcePaths {
		// Skip if the source path has been processed before
		if processedPaths[sourcePath] {
			fmt.Printf("Skipping duplicate source path: %s\n", sourcePath)
			continue
		}

		processedPaths[sourcePath] = true
		var goFiles []string
		var sourceDir string

		// Check if sourcePath is a Git repository URL
		if strings.HasPrefix(sourcePath, "https://") || strings.HasPrefix(sourcePath, "http://") || strings.HasPrefix(sourcePath, "git://") {
			sourceDir, err = m.downloadGitRepository(sourcePath)
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

			// Check for duplicate file names
			if processedFiles[baseName] {
				logrus.Fatalf("Duplicate file name detected: %s. Ensure file names are unique across all source paths.", baseName)
			}

			processedFiles[baseName] = true

			outputName := strings.TrimSuffix(baseName, ".go") + ".so"
			outputFilePath, err := filepath.Abs(filepath.Join(outputPath, outputName))
			if err != nil {
				logrus.Fatalf("Error creating output file %s: %v", outputFilePath, err)
			}

			fileMigrate, err := filepath.Abs(file)
			if err != nil {
				logrus.Fatalf("Error finding migration file %s: %v", fileMigrate, err)
			}

			if _, err := os.Stat(outputFilePath); os.IsNotExist(err) || rebuild == true {

				buildCmd := []string{"build", "-buildmode=plugin", "-o", outputFilePath}
				buildCmd = append(buildCmd, strings.Fields(buildCommand)...)
				buildCmd = append(buildCmd, fileMigrate)

				xcmd := exec.Command(goPath, buildCmd...)
				xcmd.Dir = sourceDir

				output, err := xcmd.CombinedOutput()
				if err != nil {
					logrus.Fatalf("Error building migration from %s: %v\n%s", file, err, output)
				}

				fmt.Printf("Built migration: %s\n", outputName)
			} else {
				fmt.Printf("migration plugin already exists, skipping build: %s\n", outputName)
			}
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
}
