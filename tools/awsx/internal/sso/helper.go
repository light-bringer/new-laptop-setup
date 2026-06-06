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

type Session struct {
	Name     string
	StartURL string
	Region   string
}

type AccountRole struct {
	AccountID   string
	AccountName string
	RoleName    string
}

func FindSession(awsConfigPath string) (Session, error) {
	return findSession(awsConfigPath, "")
}

func FindSessionByName(awsConfigPath, name string) (Session, error) {
	return findSession(awsConfigPath, name)
}

func findSession(awsConfigPath, targetName string) (Session, error) {
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

		if sessionName, ok := strings.CutPrefix(line, "[sso-session "); ok && strings.HasSuffix(sessionName, "]") {
			if current != nil && current.StartURL != "" && current.Region != "" && (targetName == "" || current.Name == targetName) {
				return *current, nil
			}
			current = &Session{Name: strings.TrimSpace(strings.TrimSuffix(sessionName, "]"))}
			continue
		}

		if strings.HasPrefix(line, "[") {
			if current != nil && current.StartURL != "" && current.Region != "" && (targetName == "" || current.Name == targetName) {
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

	if current != nil && current.StartURL != "" && current.Region != "" && (targetName == "" || current.Name == targetName) {
		return *current, nil
	}

	return Session{}, fmt.Errorf("no SSO session found in %s. Run: aws configure sso", awsConfigPath)
}

type tokenCache struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
	StartURL    string `json:"startUrl"`
}

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

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]`)

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

var nonAlphaNumHyphen = regexp.MustCompile(`[^a-z0-9-]`)

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
