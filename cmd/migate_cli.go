package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wirnat/axara/cmd/migrate_v1"
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
		m := migrate_v1.NewMigration(sourcePaths, outputPath, rebuild, buildCommand)
		m.Build()
	},
}
