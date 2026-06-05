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
