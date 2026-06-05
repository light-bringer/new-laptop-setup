package cmd

import (
	"context"
	"fmt"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

func init() {
	whoamiCmd.Args = cobra.MaximumNArgs(1)
	whoamiCmd.SilenceErrors = true
	whoamiCmd.SilenceUsage = true
	whoamiCmd.RunE = runWhoami
}

func runWhoami(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	profileName := ""
	region := ""

	switch len(args) {
	case 0:
		// default credential chain
	case 1:
		resolved, err := resolver.Resolve(args[0])
		if err != nil {
			return fmt.Errorf("resolving profile: %w", err)
		}
		profileName = resolved.AWSProfileName
		region = resolved.Region
	default:
		return fmt.Errorf("unexpected arguments")
	}

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithSharedConfigProfile(profileName),
	}
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}

	awscfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		if profileName != "" {
			fmt.Fprintf(os.Stderr, "Credentials expired or unavailable. Run: awsx login %s\n", profileName)
		} else {
			fmt.Fprintln(os.Stderr, "AWS credentials not available")
		}
		return err
	}

	stsClient := sts.NewFromConfig(awscfg)
	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		if profileName != "" {
			fmt.Fprintf(os.Stderr, "Credentials expired or unavailable. Run: awsx login %s\n", profileName)
		} else {
			fmt.Fprintln(os.Stderr, "AWS credentials not available")
		}
		return err
	}

	fmt.Printf("Account: %s\n", *result.Account)
	fmt.Printf("Arn:     %s\n", *result.Arn)
	fmt.Printf("UserId:  %s\n", *result.UserId)

	return nil
}
