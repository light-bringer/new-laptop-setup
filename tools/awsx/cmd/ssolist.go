package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	ssopkg "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/sso"
	"github.com/spf13/cobra"
)

func init() {
	ssoListCmd.Args = cobra.NoArgs
	ssoListCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error { return nil }
	ssoListCmd.RunE = runSSOList
}

func runSSOList(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = os.Getenv("HOME")
	}
	awsConfigPath := filepath.Join(home, ".aws", "config")

	session, err := ssopkg.FindSession(awsConfigPath)
	if err != nil {
		return err
	}

	token, err := ssopkg.FindAccessToken(session.StartURL)
	if err != nil {
		if loginErr := ssoLogin(session.Name); loginErr != nil {
			return loginErr
		}
		token, err = ssopkg.FindAccessToken(session.StartURL)
		if err != nil {
			return fmt.Errorf("no valid SSO token after login: %w", err)
		}
	}

	pairs, err := ssopkg.ListAccountRoles(context.Background(), token, session)
	if err != nil {
		return err
	}

	if len(pairs) == 0 {
		fmt.Println("No accounts or roles found.")
		return nil
	}

	fmt.Printf("%-16s %-48s %s\n", "ACCOUNT ID", "ACCOUNT NAME", "ROLE")
	for _, p := range pairs {
		fmt.Printf("%-16s %-48s %s\n", p.AccountID, p.AccountName, p.RoleName)
	}

	return nil
}

func ssoLogin(sessionName string) error {
	awsBin, err := exec.LookPath("aws")
	if err != nil {
		return fmt.Errorf("aws CLI not found in PATH. Install: brew install awscli")
	}

	c := exec.Command(awsBin, "sso", "login", "--sso-session", sessionName)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("aws sso login exited with code %d", exitErr.ExitCode())
		}
		return err
	}
	return nil
}
