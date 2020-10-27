// Copyright 2020 Daisuke Takahashi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var client *github.Client

func NewRunnerCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "runner-manager",
		Short:             "GitHub Actions Self-hosted runner management",
		PersistentPreRunE: defaultPersistentPreRunE,
	}

	rootCmd.AddCommand(newDownloadsCommand())
	rootCmd.AddCommand(newRegistrationTokenCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newGetCommand())
	rootCmd.AddCommand(newRemoveTokenCommand())
	rootCmd.AddCommand(newRemoveCommand())

	rootCmd.PersistentFlags().Int64P("app-id", "a", 62182, "GitHub Apps ID")
	rootCmd.PersistentFlags().Int64P("installation-id", "i", 8315538, "GitHub Apps installation ID")
	rootCmd.PersistentFlags().StringP("private-key-file", "p", "", "Path to GitHub Apps private key file.")
	rootCmd.PersistentFlags().StringP("org", "o", "", "GitHub Organization")
	// rootCmd.MarkPersistentFlagRequired("app-id")
	// rootCmd.MarkPersistentFlagRequired("installation-id")
	rootCmd.MarkPersistentFlagRequired("org")

	cobra.OnInitialize(func() { initializeConfig(rootCmd, "RUNNER") })
	return rootCmd
}

func initializeConfig(cmd *cobra.Command, envPrefix string) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			viper.BindEnv(
				f.Name,
				fmt.Sprintf("%s_%s", envPrefix, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))),
			)
		}

		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func defaultPersistentPreRunE(cmd *cobra.Command, args []string) error {
	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport

	// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
	appID, _ := cmd.Flags().GetInt64("app-id")
	installationID, _ := cmd.Flags().GetInt64("installation-id")
	privateKeyFile, _ := cmd.Flags().GetString("private-key-file")
	itr, err := ghinstallation.NewKeyFromFile(tr, appID, installationID, privateKeyFile)
	if err != nil {
		return err
	}

	// Use installation transport with github.com/google/go-github
	client = github.NewClient(&http.Client{Transport: itr})

	return nil
}

func inheritParentPersistentPreRun(cmd *cobra.Command, args []string) error {
	if cmd.HasParent() {
		parent := cmd.Parent()
		if parent.PersistentPreRun != nil {
			parent.PersistentPreRun(parent, args)
			return nil
		}
		if parent.PersistentPreRunE != nil {
			return parent.PersistentPreRunE(parent, args)
		}
	}

	return nil
}
