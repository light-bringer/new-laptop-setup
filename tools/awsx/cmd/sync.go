package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/light-bringer/new-laptop-setup/tools/awsx/internal/config"
	ssopkg "github.com/light-bringer/new-laptop-setup/tools/awsx/internal/sso"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	syncSSOSession string
	syncForce      bool
	syncDryRun     bool
	syncRegion     string
)

func init() {
	syncCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error { return nil }
	syncCmd.Args = cobra.NoArgs
	syncCmd.RunE = runSync

	syncCmd.Flags().StringVar(&syncSSOSession, "sso-session", "", "override SSO session name (default: auto-detect from ~/.aws/config)")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "overwrite existing profiles/aliases")
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "print proposed changes without writing any files")
	syncCmd.Flags().StringVar(&syncRegion, "region", "us-east-1", "default region for all profiles")
}

func runSync(cmd *cobra.Command, args []string) error {
	awsConfigPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")

	session, err := ssopkg.FindSession(awsConfigPath)
	if err != nil {
		return err
	}
	if syncSSOSession != "" {
		session, err = ssopkg.FindSessionByName(awsConfigPath, syncSSOSession)
		if err != nil {
			return fmt.Errorf("SSO session %q not found in ~/.aws/config: %w", syncSSOSession, err)
		}
	}

	token, err := ssopkg.FindAccessToken(session.StartURL)
	if err != nil {
		if loginErr := ssoLogin(session.Name); loginErr != nil {
			return loginErr
		}
		token, err = ssopkg.FindAccessToken(session.StartURL)
		if err != nil {
			return fmt.Errorf("no valid SSO token after login: %w", err)
		}
	}

	pairs, err := ssopkg.ListAccountRoles(context.Background(), token, session)
	if err != nil {
		return err
	}

	awsCfgCount, err := syncAWSConfig(awsConfigPath, pairs, session.Name, syncRegion)
	if err != nil {
		return err
	}

	awsxPath := config.DefaultPath()
	awsxCount, err := syncAWSXConfig(awsxPath, pairs)
	if err != nil {
		return err
	}

	if !syncDryRun {
		fmt.Printf("Synced %d profiles to ~/.aws/config, %d aliases to ~/.awsx.yaml\n", awsCfgCount, awsxCount)
	}
	return nil
}

func syncAWSConfig(path string, pairs []ssopkg.AccountRole, sessionName, region string) (int, error) {
	existingContent := ""
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return 0, fmt.Errorf("reading %s: %w", path, err)
	}
	if err == nil {
		existingContent = string(data)
	}

	added := 0
	var newBlocks strings.Builder

	for _, p := range pairs {
		name := ssopkg.ProfileName(p.AccountName, p.RoleName)
		beginMarker := fmt.Sprintf("# BEGIN ANSIBLE MANAGED BLOCK - AWS Profile %s", name)
		endMarker := "# END ANSIBLE MANAGED BLOCK"

		alreadyExists := strings.Contains(existingContent, beginMarker)

		if alreadyExists && !syncForce {
			if verbose {
				fmt.Printf("  skip   %s (already exists)\n", name)
			}
			continue
		}

		block := fmt.Sprintf(`%s
[profile %s]
sso_session = %s
sso_account_id = %s
sso_role_name = %s
region = %s
output = json
%s
`, beginMarker, name, sessionName, p.AccountID, p.RoleName, region, endMarker)

		action := "add"
		if alreadyExists {
			action = "update"
		}
		if verbose {
			fmt.Printf("  %-6s %s → %s (%s %s)\n", action, name, name, p.AccountID, p.RoleName)
		}

		if !syncDryRun {
			newBlocks.WriteString(block)
		}
		added++
	}

	if syncDryRun || added == 0 {
		return added, nil
	}

	if syncForce {
		existingContent = removeExistingManagedBlocks(existingContent, pairs)
		content := existingContent
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += newBlocks.String()
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return 0, fmt.Errorf("writing %s: %w", path, err)
		}
		return added, nil
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return 0, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "\n%s", newBlocks.String()); err != nil {
		return 0, fmt.Errorf("writing %s: %w", path, err)
	}

	return added, nil
}

func removeExistingManagedBlocks(content string, pairs []ssopkg.AccountRole) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var out strings.Builder
	inBlock := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# BEGIN ANSIBLE MANAGED BLOCK - AWS Profile ") {
			profileName := strings.TrimPrefix(line, "# BEGIN ANSIBLE MANAGED BLOCK - AWS Profile ")
			for _, p := range pairs {
				if ssopkg.ProfileName(p.AccountName, p.RoleName) == profileName {
					inBlock = true
					break
				}
			}
			if inBlock {
				continue
			}
		}
		if inBlock && strings.TrimSpace(line) == "# END ANSIBLE MANAGED BLOCK" {
			inBlock = false
			continue
		}
		if !inBlock {
			out.WriteString(line + "\n")
		}
	}
	return out.String()
}

func syncAWSXConfig(path string, pairs []ssopkg.AccountRole) (int, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return 0, fmt.Errorf("loading %s: %w", path, err)
	}

	added := 0
	for _, p := range pairs {
		name := ssopkg.ProfileName(p.AccountName, p.RoleName)
		if _, exists := cfg.Profiles[name]; exists && !syncForce {
			if verbose {
				fmt.Printf("  skip   %s (already exists)\n", name)
			}
			continue
		}
		if verbose {
			fmt.Printf("  add    %s → %s (%s %s)\n", name, name, p.AccountID, p.RoleName)
		}
		if !syncDryRun {
			cfg.Profiles[name] = config.ProfileAlias{
				AWSProfile:  name,
				AccountID:   p.AccountID,
				Description: p.AccountName + " " + p.RoleName,
			}
		}
		added++
	}

	if syncDryRun || added == 0 {
		return added, nil
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return 0, fmt.Errorf("marshalling config: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return 0, fmt.Errorf("writing %s: %w", path, err)
	}

	return added, nil
}
