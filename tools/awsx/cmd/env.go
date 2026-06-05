package cmd

import (
	"context"
	"fmt"
	"os"

	awsint "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/aws"
	"github.com/spf13/cobra"
)

func init() {
	envCmd.Args = cobra.ExactArgs(1)
	envCmd.SilenceErrors = true
	envCmd.SilenceUsage = true
	envCmd.RunE = runEnv
}

func runEnv(cmd *cobra.Command, args []string) error {
	profile := args[0]

	resolved, err := resolver.Resolve(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	creds, err := awsint.GetCredentials(context.Background(), resolved.AWSProfileName, resolved.Region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(awsint.ExportStatements(creds))
	return nil
}
