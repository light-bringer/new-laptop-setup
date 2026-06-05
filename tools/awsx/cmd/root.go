package cmd

import (
	"fmt"
	"os"

	_ "github.com/aws/aws-sdk-go-v2"
	_ "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/credentials"
	_ "github.com/aws/aws-sdk-go-v2/service/codeartifact"
	_ "github.com/aws/aws-sdk-go-v2/service/ecr"
	_ "github.com/aws/aws-sdk-go-v2/service/sts"
	_ "github.com/spf13/viper"
	_ "gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "awsx",
	Short: "AWS profile switcher with exec, ECR, and CodeArtifact support",
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
