package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	listCmd.Args = cobra.NoArgs
	listCmd.RunE = runList
}

func runList(cmd *cobra.Command, args []string) error {
	printed := false
	covered := map[string]struct{}{}

	if cfg != nil && len(cfg.Profiles) > 0 {
		fmt.Println("Aliases:")
		for alias, profile := range cfg.Profiles {
			fmt.Printf("%-20s → %-45s %s\n", alias, profile.AWSProfile, profile.Description)
			covered[profile.AWSProfile] = struct{}{}
		}
		printed = true
	}

	profiles, err := rawAWSProfiles()
	if err != nil {
		return err
	}
	if len(profiles) > 0 {
		if printed {
			fmt.Println()
		}
		fmt.Println("AWS Profiles:")
		for _, profile := range profiles {
			if _, ok := covered[profile]; ok {
				continue
			}
			fmt.Println(profile)
			printed = true
		}
	}

	if !printed {
		fmt.Println("No profiles configured. Create ~/.awsx.yaml or run: aws configure sso")
	}

	return nil
}

func rawAWSProfiles() ([]string, error) {
	cmd := exec.Command("aws", "configure", "list-profiles")
	out, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			return nil, nil
		}
		if strings.Contains(err.Error(), "executable file not found") {
			return nil, nil
		}
		return nil, err
	}

	var profiles []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			profiles = append(profiles, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return profiles, nil
}
