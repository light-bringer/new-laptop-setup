package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	loginCmd.Args = cobra.ExactArgs(1)
	loginCmd.RunE = runLogin
}

func runLogin(cmd *cobra.Command, args []string) error {
	profile := args[0]

	resolved, err := resolver.Resolve(profile)
	if err != nil {
		return err
	}

	awsBin, err := exec.LookPath("aws")
	if err != nil {
		return fmt.Errorf("aws CLI not found in PATH. Install: brew install awscli")
	}

	awsCmd := exec.Command(awsBin, "sso", "login", "--profile", resolved.AWSProfileName)
	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	if err := awsCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// os.Exit is intentional here: Cobra's Execute() only exits with 1
			// on any error, losing the original exit code from aws sso login.
			// Preserving it matters for scripts that check $?.
			os.Exit(exitErr.ExitCode())
		}
		return err
	}

	return nil
}
