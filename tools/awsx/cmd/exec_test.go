package cmd

import (
	"testing"

	awsint "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/aws"
)

func TestExecArgParsing(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		profile   string
		cmdArgs   []string
		wantError bool
	}{
		{
			name:    "splits profile and command",
			args:    []string{"prod", "--", "echo", "hello"},
			profile: "prod",
			cmdArgs: []string{"echo", "hello"},
		},
		{
			name:    "help flag after -- is not intercepted",
			args:    []string{"dev", "--", "grep", "-h", "foo"},
			profile: "dev",
			cmdArgs: []string{"grep", "-h", "foo"},
		},
		{
			name:      "missing separator errors",
			args:      []string{"prod"},
			wantError: true,
		},
		{
			name:      "empty command errors",
			args:      []string{"prod", "--"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, cmdArgs, err := parseExecArgs(tt.args)
			if tt.wantError {
				if err == nil {
					t.Fatalf("parseExecArgs() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseExecArgs() error = %v", err)
			}
			if profile != tt.profile {
				t.Fatalf("parseExecArgs() profile = %q, want %q", profile, tt.profile)
			}
			if len(cmdArgs) != len(tt.cmdArgs) {
				t.Fatalf("parseExecArgs() cmdArgs len = %d, want %d", len(cmdArgs), len(tt.cmdArgs))
			}
			for i := range tt.cmdArgs {
				if cmdArgs[i] != tt.cmdArgs[i] {
					t.Fatalf("parseExecArgs() cmdArgs[%d] = %q, want %q", i, cmdArgs[i], tt.cmdArgs[i])
				}
			}
		})
	}
}

func TestExecEnvConstruction(t *testing.T) {
	baseEnv := []string{"PATH=/usr/bin", "HOME=/Users/test", "AWS_ACCESS_KEY_ID=STALE_KEY", "AWS_DEFAULT_REGION=eu-west-1"}
	creds := awsint.Credentials{
		AccessKeyID:     "AKIA123",
		SecretAccessKey: "secret123",
		SessionToken:    "token123",
		Region:          "us-west-2",
	}

	got := append(awsint.FilterAWSEnv(baseEnv), awsint.EnvVars(creds)...)
	want := map[string]bool{
		"AWS_ACCESS_KEY_ID=AKIA123":       false,
		"AWS_SECRET_ACCESS_KEY=secret123": false,
		"AWS_SESSION_TOKEN=token123":      false,
		"AWS_DEFAULT_REGION=us-west-2":    false,
	}

	for _, env := range got {
		if _, ok := want[env]; ok {
			want[env] = true
		}
	}

	for key, seen := range want {
		if !seen {
			t.Fatalf("missing env var %q in %v", key, got)
		}
	}

	for _, e := range got {
		if e == "AWS_ACCESS_KEY_ID=STALE_KEY" {
			t.Fatal("stale AWS_ACCESS_KEY_ID survived FilterAWSEnv")
		}
		if e == "AWS_DEFAULT_REGION=eu-west-1" {
			t.Fatal("stale AWS_DEFAULT_REGION survived FilterAWSEnv")
		}
	}
}
