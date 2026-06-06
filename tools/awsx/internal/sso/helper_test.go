package sso

import (
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
