package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	awsint "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/aws"
	"github.com/spf13/cobra"
)

func init() {
	execCmd.Args = cobra.MinimumNArgs(1)
	execCmd.DisableFlagParsing = true
	execCmd.SilenceUsage = true
	execCmd.SilenceErrors = true
	execCmd.RunE = runExec
}

func parseExecArgs(args []string) (profile string, cmdArgs []string, err error) {
	if len(args) == 0 {
		return "", nil, errors.New("missing profile")
	}

	profile = args[0]
	for i := 1; i < len(args); i++ {
		if args[i] == "--" {
			if i+1 >= len(args) {
				return "", nil, errors.New("missing command")
			}
			return profile, args[i+1:], nil
		}
	}

	return "", nil, errors.New("missing -- separator")
}

func runExec(cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return cmd.Help()
		}
	}

	profile, cmdArgs, err := parseExecArgs(args)
	if err != nil {
		fmt.Fprint(cmd.ErrOrStderr(), cmd.UsageString())
		return err
	}

	resolved, err := resolver.Resolve(profile)
	if err != nil {
		return err
	}

	creds, err := awsint.GetCredentials(context.Background(), resolved.AWSProfileName, resolved.Region)
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), err)
		return err
	}

	bin, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		return fmt.Errorf("command not found: %s", cmdArgs[0])
	}

	env := append(awsint.FilterAWSEnv(os.Environ()), awsint.EnvVars(creds)...)
	return syscall.Exec(bin, cmdArgs, env)
}
