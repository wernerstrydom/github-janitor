package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	token        string
	organization string
	client       *github.Client
	ctx          context.Context
)

// getGitHubCLIToken retrieves the token from the GitHub CLI
func getGitHubCLIToken() (string, error) {
	// Check if gh is installed
	_, err := exec.LookPath("gh")
	if err != nil {
		return "", fmt.Errorf("gh command not found: %v", err)
	}

	cmd := exec.Command("gh", "auth", "token")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute gh auth token: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func initAuth() error {
	var err error

	// Get token from config or flag
	token = viper.GetString("token")
	if token == "" {
		// Try to get the token from GitHub CLI
		token, err = getGitHubCLIToken()
		if err != nil {
			return fmt.Errorf("GitHub access token is required. Use --token flag, set it in the config file, or authenticate with the GitHub CLI")
		}
	}

	// Get organization from config or flag
	organization = viper.GetString("organization")
	if organization == "" {
		return fmt.Errorf("GitHub organization name is required. Use --organization flag or set it in the config file")
	}

	// Initialize GitHub client
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-janitor",
	Short: "A tool for housekeeping activities on GitHub repositories",
	Long: `github-janitor is a CLI tool to perform housekeeping tasks
on GitHub repositories, such as identifying empty repositories.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration
		if err := initAuth(); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent Flags, which will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-janitor.yaml)")
	rootCmd.PersistentFlags().StringP("token", "t", "", "GitHub Access Token")
	rootCmd.PersistentFlags().StringP("organization", "o", "", "GitHub organization name")

	// Bind flags to Viper
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("organization", rootCmd.PersistentFlags().Lookup("organization"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".github-janitor" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".github-janitor")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
