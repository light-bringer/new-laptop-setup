package sso

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShortAccountName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Application Security - Common - NonProduction", "nonprod"},
		{"Application Security - Common - Production", "prod"},
		{"Network - Shared - Development", "dev"},
		{"Platform - Core - Staging", "stg"},
		{"Security - Common - Sandbox", "sandbox"},
		{"QA - Automation - Testing", "test"},
		{"SingleWordAccount", "singlewordaccount"},
		{"My-Account-Name", "myaccountname"},
		{"Special!@# Chars - Dev", "dev"},
		{"Already - - Short", "short"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := ShortAccountName(tc.input)
			if got != tc.want {
				t.Errorf("ShortAccountName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestShortRoleName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Administrator", "admin"},
		{"AdministratorAccess", "admin"},
		{"ReadOnly", "ro"},
		{"ReadOnlyAccess", "ro"},
		{"PowerUser", "poweruser"},
		{"DataScientist", "datascientist"},
		{"Custom Role Name", "custom-role-name"},
		{"S3FullAccess", "s3fullaccess"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := ShortRoleName(tc.input)
			if got != tc.want {
				t.Errorf("ShortRoleName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestProfileName(t *testing.T) {
	cases := []struct {
		account string
		role    string
		want    string
	}{
		{"Application Security - Common - NonProduction", "Administrator", "nonprod-admin"},
		{"Application Security - Common - NonProduction", "ReadOnly", "nonprod-ro"},
		{"Application Security - Common - Production", "ReadOnly", "prod-ro"},
		{"Application Security - Common - Production", "AdministratorAccess", "prod-admin"},
		{"SingleWordAccount", "PowerUser", "singlewordaccount-poweruser"},
	}

	for _, tc := range cases {
		t.Run(tc.account+"/"+tc.role, func(t *testing.T) {
			got := ProfileName(tc.account, tc.role)
			if got != tc.want {
				t.Errorf("ProfileName(%q, %q) = %q, want %q", tc.account, tc.role, got, tc.want)
			}
		})
	}
}

func TestFindSessionByName(t *testing.T) {
	config := `[profile default]
region = us-west-2

[sso-session first]
sso_start_url = https://first.awsapps.com/start
sso_region = us-east-1

[sso-session target]
sso_start_url = https://target.awsapps.com/start
sso_region = us-west-2
`

	path := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(path, []byte(config), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	session, err := FindSessionByName(path, "target")
	if err != nil {
		t.Fatalf("FindSessionByName() error = %v", err)
	}
	if session.Name != "target" {
		t.Fatalf("session.Name = %q, want %q", session.Name, "target")
	}
	if session.StartURL != "https://target.awsapps.com/start" {
		t.Fatalf("session.StartURL = %q", session.StartURL)
	}
	if session.Region != "us-west-2" {
		t.Fatalf("session.Region = %q, want %q", session.Region, "us-west-2")
	}

	if _, err := FindSessionByName(path, "missing"); err == nil {
		t.Fatal("FindSessionByName() error = nil, want error")
	}
}
