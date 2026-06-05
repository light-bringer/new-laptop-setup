package cmd

import (
	"context"
	"fmt"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsint "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/aws"
	"github.com/spf13/cobra"
)

func init() {
	whoamiCmd.Args = cobra.MaximumNArgs(1)
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
		resolved, _ := resolver.Resolve(args[0])
		profileName = resolved.AWSProfileName
		region = resolved.Region
	default:
		return fmt.Errorf("unexpected arguments")
	}

	if _, err := awsint.GetCredentials(ctx, profileName, region); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return err
	}

	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithSharedConfigProfile(profileName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Credentials expired. Run: awsx login %s\n", profileName)
		return err
	}

	stsClient := sts.NewFromConfig(awscfg)
	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Credentials expired. Run: awsx login %s\n", profileName)
		return err
	}

	fmt.Printf("Account: %s\n", *result.Account)
	fmt.Printf("Arn:     %s\n", *result.Arn)
	fmt.Printf("UserId:  %s\n", *result.UserId)

	return nil
}
