package aws

import (
	"context"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
}

// awsEnvKeys is the set of AWS credential/region env var prefixes to strip
// before injecting fresh credentials so stale values don't shadow the new ones.
var awsEnvKeys = map[string]struct{}{
	"AWS_ACCESS_KEY_ID":     {},
	"AWS_SECRET_ACCESS_KEY": {},
	"AWS_SESSION_TOKEN":     {},
	"AWS_DEFAULT_REGION":    {},
	"AWS_REGION":            {},
	"AWS_PROFILE":           {},
}

func GetCredentials(ctx context.Context, profileName string, region string) (Credentials, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithSharedConfigProfile(profileName),
	}
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return Credentials{}, err
	}

	v, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		if isExpiredSSOError(err) {
			return Credentials{}, fmt.Errorf("credentials expired or unavailable for profile %q - run: awsx login %s\n%w", profileName, profileName, err)
		}
		return Credentials{}, err
	}

	return Credentials{
		AccessKeyID:     v.AccessKeyID,
		SecretAccessKey: v.SecretAccessKey,
		SessionToken:    v.SessionToken,
		Region:          cfg.Region,
	}, nil
}

func EnvVars(creds Credentials) []string {
	return []string{
		"AWS_ACCESS_KEY_ID=" + creds.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY=" + creds.SecretAccessKey,
		"AWS_SESSION_TOKEN=" + creds.SessionToken,
		"AWS_DEFAULT_REGION=" + creds.Region,
		"AWS_REGION=" + creds.Region,
	}
}

// FilterAWSEnv returns a copy of env with all AWS credential/region variables
// removed. Use this before appending fresh credentials from EnvVars().
func FilterAWSEnv(env []string) []string {
	out := make([]string, 0, len(env))
	for _, e := range env {
		key, _, _ := strings.Cut(e, "=")
		if _, skip := awsEnvKeys[key]; !skip {
			out = append(out, e)
		}
	}
	return out
}

func ExportStatements(creds Credentials) string {
	return strings.Join([]string{
		"export AWS_ACCESS_KEY_ID=" + creds.AccessKeyID,
		"export AWS_SECRET_ACCESS_KEY=" + creds.SecretAccessKey,
		"export AWS_SESSION_TOKEN=" + creds.SessionToken,
		"export AWS_DEFAULT_REGION=" + creds.Region,
		"export AWS_REGION=" + creds.Region,
		"",
	}, "\n")
}

func isExpiredSSOError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "sso") && (strings.Contains(msg, "expired") || strings.Contains(msg, "token has expired") || strings.Contains(msg, "token expired"))
}
