package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/spf13/cobra"
)

var (
	caDomain      string
	caDomainOwner string
	caRegion      string
)

func init() {
	codeartifactCmd.Args = cobra.ExactArgs(1)
	codeartifactCmd.Flags().StringVar(&caDomain, "domain", "", "CodeArtifact domain name (required)")
	codeartifactCmd.Flags().StringVar(&caDomainOwner, "domain-owner", "", "Domain owner account ID (default: profile account_id)")
	codeartifactCmd.Flags().StringVar(&caRegion, "region", "", "AWS region (default: profile region or us-east-1)")
	if err := codeartifactCmd.MarkFlagRequired("domain"); err != nil {
		panic(err)
	}
	codeartifactCmd.RunE = runCodeArtifact
}

func runCodeArtifact(cmd *cobra.Command, args []string) error {
	profile := args[0]

	resolved, err := resolver.Resolve(profile)
	if err != nil {
		return err
	}

	region := caRegion
	if region == "" {
		region = resolved.Region
	}
	if region == "" {
		region = "us-east-1"
	}

	domainOwner := caDomainOwner
	if domainOwner == "" {
		domainOwner = resolved.AccountID
	}

	ctx := context.Background()
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithSharedConfigProfile(resolved.AWSProfileName),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return fmt.Errorf("loading AWS config for profile %q: %w", resolved.AWSProfileName, err)
	}

	result, err := fetchCodeArtifactToken(ctx, awscfg, resolved.AWSProfileName, domainOwner)
	if err != nil {
		if isExpiredOrSSOError(err) {
			if loginErr := awsSSOLogin(resolved.AWSProfileName); loginErr != nil {
				return loginErr
			}
			result, err = fetchCodeArtifactToken(ctx, awscfg, resolved.AWSProfileName, domainOwner)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if result.AuthorizationToken == nil {
		return fmt.Errorf("missing CodeArtifact authorization token for profile %q", resolved.AWSProfileName)
	}

	fmt.Printf("export CODEARTIFACT_AUTH_TOKEN=%s\n", *result.AuthorizationToken)
	return nil
}

func fetchCodeArtifactToken(ctx context.Context, awscfg aws.Config, profile, domainOwner string) (*codeartifact.GetAuthorizationTokenOutput, error) {
	client := codeartifact.NewFromConfig(awscfg)

	input := &codeartifact.GetAuthorizationTokenInput{Domain: &caDomain}
	if domainOwner != "" {
		input.DomainOwner = &domainOwner
	}

	result, err := client.GetAuthorizationToken(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("getting CodeArtifact authorization token for profile %q: %w", profile, err)
	}
	return result, nil
}

func awsSSOLogin(profile string) error {
	awsBin, err := exec.LookPath("aws")
	if err != nil {
		return fmt.Errorf("aws CLI not found in PATH. Install: brew install awscli")
	}

	awsCmd := exec.Command(awsBin, "sso", "login", "--profile", profile)
	awsCmd.Stdin = os.Stdin
	awsCmd.Stdout = os.Stdout
	awsCmd.Stderr = os.Stderr

	if err := awsCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("aws sso login failed with exit code %d", exitErr.ExitCode())
		}
		return err
	}

	return nil
}

func isExpiredOrSSOError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "expired") || strings.Contains(msg, "sso")
}
