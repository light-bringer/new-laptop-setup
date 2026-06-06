package sso

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssso "github.com/aws/aws-sdk-go-v2/service/sso"
)

// Session holds the parsed SSO session block from ~/.aws/config.
type Session struct {
	Name     string
	StartURL string
	Region   string
}

// AccountRole is a single (account, role) pair returned by the SSO API.
type AccountRole struct {
	AccountID   string
	AccountName string
	RoleName    string
}

// FindSession reads awsConfigPath and returns the first [sso-session] block
// that has both sso_start_url and sso_region configured.
func FindSession(awsConfigPath string) (Session, error) {
	f, err := os.Open(awsConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Session{}, fmt.Errorf("no SSO session found in %s. Run: aws configure sso", awsConfigPath)
		}
		return Session{}, err
	}
	defer f.Close()

	var current *Session
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[sso-session ") && strings.HasSuffix(line, "]") {
			if current != nil && current.StartURL != "" && current.Region != "" {
				return *current, nil
			}
			name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[sso-session "), "]"))
			current = &Session{Name: name}
			continue
		}

		if strings.HasPrefix(line, "[") {
			if current != nil && current.StartURL != "" && current.Region != "" {
				return *current, nil
			}
			current = nil
			continue
		}

		if current == nil {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "sso_start_url":
			current.StartURL = strings.TrimSpace(val)
		case "sso_region":
			current.Region = strings.TrimSpace(val)
		}
	}

	if err := scanner.Err(); err != nil {
		return Session{}, err
	}

	if current != nil && current.StartURL != "" && current.Region != "" {
		return *current, nil
	}

	return Session{}, fmt.Errorf("no SSO session found in %s. Run: aws configure sso", awsConfigPath)
}

type tokenCache struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
	StartURL    string `json:"startUrl"`
}

// FindAccessToken reads ~/.aws/sso/cache/*.json and returns a valid
// (non-expired) access token whose startUrl matches startURL.
func FindAccessToken(startURL string) (string, error) {
	cacheDir := filepath.Join(os.Getenv("HOME"), ".aws", "sso", "cache")
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("SSO cache directory not found: %s", cacheDir)
		}
		return "", fmt.Errorf("reading SSO cache directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(cacheDir, entry.Name()))
		if err != nil {
			continue
		}

		var tc tokenCache
		if err := json.Unmarshal(data, &tc); err != nil {
			continue
		}

		if tc.AccessToken == "" || tc.StartURL != startURL {
			continue
		}

		expiresAt, err := time.Parse(time.RFC3339, tc.ExpiresAt)
		if err != nil {
			expiresAt, err = time.Parse("2006-01-02T15:04:05UTC", tc.ExpiresAt)
			if err != nil {
				continue
			}
		}

		if time.Now().After(expiresAt) {
			continue
		}

		return tc.AccessToken, nil
	}

	return "", fmt.Errorf("no valid SSO token found for %s", startURL)
}

// ListAccountRoles enumerates all (account, role) pairs accessible to the
// current user by calling the SSO ListAccounts and ListAccountRoles APIs.
func ListAccountRoles(ctx context.Context, accessToken string, session Session) ([]AccountRole, error) {
	cfg := aws.Config{Region: session.Region}
	client := awssso.NewFromConfig(cfg)

	type accountEntry struct{ ID, Name string }
	var accounts []accountEntry
	var accountToken *string
	for {
		resp, err := client.ListAccounts(ctx, &awssso.ListAccountsInput{
			AccessToken: aws.String(accessToken),
			NextToken:   accountToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing SSO accounts: %w", err)
		}
		for _, a := range resp.AccountList {
			accounts = append(accounts, accountEntry{
				ID:   aws.ToString(a.AccountId),
				Name: aws.ToString(a.AccountName),
			})
		}
		if resp.NextToken == nil {
			break
		}
		accountToken = resp.NextToken
	}

	var result []AccountRole
	for _, account := range accounts {
		var roleToken *string
		for {
			resp, err := client.ListAccountRoles(ctx, &awssso.ListAccountRolesInput{
				AccessToken: aws.String(accessToken),
				AccountId:   aws.String(account.ID),
				NextToken:   roleToken,
			})
			if err != nil {
				return nil, fmt.Errorf("listing roles for account %s: %w", account.ID, err)
			}
			for _, r := range resp.RoleList {
				result = append(result, AccountRole{
					AccountID:   account.ID,
					AccountName: account.Name,
					RoleName:    aws.ToString(r.RoleName),
				})
			}
			if resp.NextToken == nil {
				break
			}
			roleToken = resp.NextToken
		}
	}

	return result, nil
}

var accountShortSubs = map[string]string{
	"nonproduction": "nonprod",
	"production":    "prod",
	"development":   "dev",
	"staging":       "stg",
	"sandbox":       "sandbox",
	"testing":       "test",
}

// nonAlphaNum matches any character that is not an ASCII letter or digit.
var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]`)

// ShortAccountName derives a short identifier from an AWS account name.
//
//	"Application Security - Common - NonProduction" → "nonprod"
//	"Application Security - Common - Production"    → "prod"
func ShortAccountName(accountName string) string {
	segment := accountName
	if idx := strings.LastIndex(accountName, " - "); idx >= 0 {
		segment = accountName[idx+3:]
	}
	stripped := nonAlphaNum.ReplaceAllString(strings.ToLower(strings.TrimSpace(segment)), "")
	if sub, ok := accountShortSubs[stripped]; ok {
		return sub
	}
	return stripped
}

var roleShortSubs = map[string]string{
	"administrator":       "admin",
	"administratoraccess": "admin",
	"readonly":            "ro",
	"readonlyaccess":      "ro",
}

// nonAlphaNumHyphen matches any character that is not an ASCII letter, digit, or hyphen.
var nonAlphaNumHyphen = regexp.MustCompile(`[^a-z0-9-]`)

// ShortRoleName derives a short identifier from an AWS role name.
//
//	"Administrator"       → "admin"
//	"ReadOnly"            → "ro"
//	"AdministratorAccess" → "admin"
func ShortRoleName(roleName string) string {
	lower := strings.ToLower(roleName)
	if sub, ok := roleShortSubs[nonAlphaNum.ReplaceAllString(lower, "")]; ok {
		return sub
	}
	return nonAlphaNumHyphen.ReplaceAllString(strings.ReplaceAll(lower, " ", "-"), "")
}

func ProfileName(accountName, roleName string) string {
	return ShortAccountName(accountName) + "-" + ShortRoleName(roleName)
}
