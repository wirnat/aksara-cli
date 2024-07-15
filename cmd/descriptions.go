package cmd

import "strings"

var generatorCmdDesc = `The '{cli_name} generate' command automatically generates design patterns based on a configuration file.

Usage:
  {cli_name} generate <config_file> --models <model_name>

Example:
  {cli_name} generate conf.yaml --models User

Arguments:
  <config_file>    The path to the configuration file.
  --models <model_name>    The name of the model(s) to generate.

Details:
Some folks say that Design Patterns are dead. How foolish.
The Design Patterns book is one of the most important books published in our industry.
The concepts within should be common rudimentary knowledge for all professional programmers.`

var getterCmdDesc = `The '{cli_name} get' command fetches a CLI item from a specified GitHub repository and stores it in a given destination directory.

Usage:
  {cli_name} get <repository_url> <destination_directory>

Example:
  {cli_name} get https://github.com/wirnat/axara-template-go-clean-architecture templates

Arguments:
  <repository_url>          URL of the GitHub repository to fetch.
  <destination_directory>   Directory where the fetched item will be stored.`

var newConfigDesc = `The '%s new' command creates a new Axara configuration file.

Usage:
  {cli_name} new <config_name>

Example:
  {cli_name} new wirnat_arch

Arguments:
  <config_name>   The name of the configuration file to be created (without the .yaml extension).

Details:
This command initializes a new Axara configuration file with the specified name, helping you to start your project with the necessary configuration setup.`

var migrateNewCmdDesc = `The '{cli_name} migrate new' command creates a new migration file.

Usage:
  {cli_name} migrate new <migration_name>

Example:
  {cli_name} migrate new create_users_table

Arguments:
  <migration_name>   The name of the new migration file to be created.`

var migrateUpCmdDesc = `The '{cli_name} migrate up' command applies all up migrations.

Usage:
  {cli_name} migrate up`

var migrateDownCmdDesc = `The '{cli_name} migrate down' command rolls back the last batch of migrations.

Usage:
  {cli_name} migrate down`

var migrateFreshCmdDesc = `The '{cli_name} migrate fresh' command drops all tables and re-runs all migrations.

Usage:
  {cli_name} migrate fresh`

var migrateFlushCmdDesc = `The '{cli_name} migrate flush' command resets the database by rolling back all migrations and then re-applying them.

Usage:
  {cli_name} migrate flush`

var migrateBuildCmdDesc = `The '{cli_name} migrate build' command builds migration plugins from migration files (*.go) and renames them to *.so files.

Usage:
  {cli_name} migrate build [flags] [--source-path <source_path>]... [--output-path <output_path>]

Example:
  {cli_name} migrate build -ldflags="-s -w" --source-path ./migrations1 --source-path https://github.com/user/repo --output-path ./plugins

Flags:
  -ldflags string   flags to pass to go build (default "")

Options:
  --source-path <source_path>   Directory or Git repository containing the migration files (*.go) to build into plugins. Can be specified multiple times.
  --output-path <output_path>   Directory to store the built plugin files (default: ./plugins).`

func replaceCLIName(template, cliName string) string {
	return strings.ReplaceAll(template, "{cli_name}", cliName)
}
