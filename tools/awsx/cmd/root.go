package cmd

import (
	"fmt"
	"os"

	"github.com/light-bringer/new-laptop-setup/tools/awsx/internal/config"
	"github.com/light-bringer/new-laptop-setup/tools/awsx/internal/profile"
	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	verbose  bool
	cfg      *config.Config
	resolver *profile.Resolver
)

var rootCmd = &cobra.Command{
	Use:     "awsx",
	Short:   "AWS profile switcher with exec, ECR, and CodeArtifact support",
	Version: "0.1.0",
}

var execCmd = &cobra.Command{
	Use:   "exec <profile> -- <command...>",
	Short: "Execute a command with AWS credentials injected",
}

var loginCmd = &cobra.Command{
	Use:   "login <profile>",
	Short: "Log in to AWS SSO for the given profile",
}

var envCmd = &cobra.Command{
	Use:   "env <profile>",
	Short: "Print AWS credentials as export statements",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami [profile]",
	Short: "Show current AWS identity",
}

var ecrLoginCmd = &cobra.Command{
	Use:   "ecr-login <profile>",
	Short: "Log in to Amazon ECR for the given profile",
}

var codeartifactCmd = &cobra.Command{
	Use:   "codeartifact <profile>",
	Short: "Get a CodeArtifact authorization token",
}

var ssoListCmd = &cobra.Command{
	Use:   "sso-list",
	Short: "List available SSO accounts and roles",
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync AWS SSO profiles to ~/.aws/config and aliases to ~/.awsx.yaml",
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	path := cfgFile
	if path == "" {
		path = config.DefaultPath()
	}
	var err error
	cfg, err = config.Load(path)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	resolver = profile.NewResolver(cfg)
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: ~/.awsx.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show detailed operation output")
	if versionFlag := rootCmd.Flags().Lookup("version"); versionFlag != nil {
		versionFlag.Shorthand = ""
	}
	rootCmd.PersistentPreRunE = persistentPreRunE
	rootCmd.AddCommand(execCmd, loginCmd, envCmd, listCmd, whoamiCmd, ecrLoginCmd, codeartifactCmd, ssoListCmd, syncCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
