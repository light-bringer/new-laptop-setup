package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

var ecrRegion string

func init() {
	ecrLoginCmd.Args = cobra.ExactArgs(1)
	ecrLoginCmd.Flags().StringVar(&ecrRegion, "region", "", "AWS region (default: profile region or us-east-1)")
	ecrLoginCmd.RunE = runECRLogin
}

func runECRLogin(cmd *cobra.Command, args []string) error {
	profile := args[0]

	resolved, err := resolver.Resolve(profile)
	if err != nil {
		return err
	}

	region := ecrRegion
	if region == "" {
		region = resolved.Region
	}
	if region == "" {
		region = "us-east-1"
	}

	ctx := context.Background()
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithSharedConfigProfile(resolved.AWSProfileName),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return fmt.Errorf("loading AWS config for profile %q: %w", resolved.AWSProfileName, err)
	}

	accountID := resolved.AccountID
	if accountID == "" {
		stsClient := sts.NewFromConfig(awscfg)
		identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return fmt.Errorf("getting account ID for profile %q: %w", resolved.AWSProfileName, err)
		}
		accountID = strings.TrimSpace(*identity.Account)
	}

	if accountID == "" {
		return fmt.Errorf("unable to determine AWS account ID for profile %q", resolved.AWSProfileName)
	}

	registry := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", accountID, region)

	ecrClient := ecr.NewFromConfig(awscfg)
	resp, err := ecrClient.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return fmt.Errorf("getting ECR authorization token for profile %q: %w", resolved.AWSProfileName, err)
	}
	if len(resp.AuthorizationData) == 0 || resp.AuthorizationData[0].AuthorizationToken == nil {
		return fmt.Errorf("missing ECR authorization token for profile %q", resolved.AWSProfileName)
	}

	rawToken, err := base64.StdEncoding.DecodeString(*resp.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return fmt.Errorf("decoding ECR authorization token: %w", err)
	}

	parts := strings.SplitN(string(rawToken), ":", 2)
	if len(parts) != 2 || parts[0] != "AWS" {
		return fmt.Errorf("unexpected ECR authorization token format")
	}
	password := parts[1]

	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker not found in PATH")
	}

	dockerCmd := exec.Command(dockerBin, "login", "--username", "AWS", "--password-stdin", registry)
	dockerCmd.Stdin = strings.NewReader(password)
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	if err := dockerCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Logged in to %s\n", registry)
	return nil
}
