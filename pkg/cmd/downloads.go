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
	"context"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

func newDownloadsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "downloads",
		Short: "List runner applications for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			org, _ := cmd.Flags().GetString("org")
			apps, _, err := client.Actions.ListOrganizationRunnerApplicationDownloads(context.TODO(), org)
			if err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"OS", "Architecture", "Download URL", "Filename"})
			for _, app := range apps {
				t.AppendRow(table.Row{app.GetOS(), app.GetArchitecture(), app.GetDownloadURL(), app.GetFilename()})
			}
			t.Render()

			return nil
		},
	}

	return cmd
}
