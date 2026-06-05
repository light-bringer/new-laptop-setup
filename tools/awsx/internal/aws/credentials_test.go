package aws

import (
	"strings"
	"testing"
)

func TestEnvVars(t *testing.T) {
	creds := Credentials{
		AccessKeyID:     "AKIA123",
		SecretAccessKey: "secret123",
		SessionToken:    "token123",
		Region:          "us-west-2",
	}

	got := EnvVars(creds)
	want := []string{
		"AWS_ACCESS_KEY_ID=AKIA123",
		"AWS_SECRET_ACCESS_KEY=secret123",
		"AWS_SESSION_TOKEN=token123",
		"AWS_DEFAULT_REGION=us-west-2",
	}

	if len(got) != len(want) {
		t.Fatalf("EnvVars() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("EnvVars()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestEnvVars_EmptySessionToken(t *testing.T) {
	creds := Credentials{
		AccessKeyID:     "AKIA123",
		SecretAccessKey: "secret123",
		SessionToken:    "",
		Region:          "us-west-2",
	}

	got := EnvVars(creds)
	if got[2] != "AWS_SESSION_TOKEN=" {
		t.Fatalf("EnvVars()[2] = %q, want %q", got[2], "AWS_SESSION_TOKEN=")
	}
}

func TestFilterAWSEnv(t *testing.T) {
	input := []string{
		"PATH=/usr/bin",
		"HOME=/Users/test",
		"AWS_ACCESS_KEY_ID=OLD_KEY",
		"AWS_SECRET_ACCESS_KEY=OLD_SECRET",
		"AWS_SESSION_TOKEN=OLD_TOKEN",
		"AWS_DEFAULT_REGION=us-west-2",
		"AWS_REGION=us-west-2",
		"AWS_PROFILE=default",
		"GOPATH=/Users/test/go",
	}

	got := FilterAWSEnv(input)

	kept := map[string]bool{"PATH=/usr/bin": false, "HOME=/Users/test": false, "GOPATH=/Users/test/go": false}
	for _, e := range got {
		if _, ok := kept[e]; ok {
			kept[e] = true
		}
	}
	for k, seen := range kept {
		if !seen {
			t.Errorf("FilterAWSEnv() dropped non-AWS var %q", k)
		}
	}

	stripped := []string{"AWS_ACCESS_KEY_ID=OLD_KEY", "AWS_SECRET_ACCESS_KEY=OLD_SECRET",
		"AWS_SESSION_TOKEN=OLD_TOKEN", "AWS_DEFAULT_REGION=us-west-2",
		"AWS_REGION=us-west-2", "AWS_PROFILE=default"}
	for _, e := range got {
		for _, bad := range stripped {
			if e == bad {
				t.Errorf("FilterAWSEnv() kept AWS var %q", e)
			}
		}
	}
}

func TestExportStatements(t *testing.T) {
	creds := Credentials{
		AccessKeyID:     "AKIA123",
		SecretAccessKey: "secret123",
		SessionToken:    "token123",
		Region:          "us-west-2",
	}

	got := ExportStatements(creds)
	if !strings.HasPrefix(got, "export AWS_ACCESS_KEY_ID=") {
		t.Fatalf("ExportStatements() prefix mismatch: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("ExportStatements() must end with newline: %q", got)
	}
}

func TestExportStatements_EvalSafe(t *testing.T) {
	creds := Credentials{
		AccessKeyID:     "AKIA123",
		SecretAccessKey: "secret123",
		SessionToken:    "token123",
		Region:          "us-west-2",
	}

	got := ExportStatements(creds)
	if strings.Contains(got, "`") {
		t.Fatalf("ExportStatements() contains backticks: %q", got)
	}
	if strings.Contains(got, "$(") {
		t.Fatalf("ExportStatements() contains command substitution: %q", got)
	}
}
