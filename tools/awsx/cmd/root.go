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
	cfg      *config.Config
	resolver *profile.Resolver
)

var rootCmd = &cobra.Command{
	Use:   "awsx",
	Short: "AWS profile switcher with exec, ECR, and CodeArtifact support",
	Version: "0.1.0",
}

var execCmd = &cobra.Command{
	Use:   "exec <profile> -- <command...>",
	Short: "Execute a command with AWS credentials injected",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
}

var loginCmd = &cobra.Command{
	Use:   "login <profile>",
	Short: "Log in to AWS SSO for the given profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
}

var envCmd = &cobra.Command{
	Use:   "env <profile>",
	Short: "Print AWS credentials as export statements",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami [profile]",
	Short: "Show current AWS identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
}

var ecrLoginCmd = &cobra.Command{
	Use:   "ecr-login <profile>",
	Short: "Log in to Amazon ECR for the given profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
}

var codeartifactCmd = &cobra.Command{
	Use:   "codeartifact <profile>",
	Short: "Get a CodeArtifact authorization token",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented")
	},
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
	rootCmd.PersistentPreRunE = persistentPreRunE
	rootCmd.AddCommand(execCmd, loginCmd, envCmd, listCmd, whoamiCmd, ecrLoginCmd, codeartifactCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
