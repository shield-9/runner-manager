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

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a self-hosted runner for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			org, _ := cmd.Flags().GetString("org")
			runnerID, _ := cmd.Flags().GetInt64("runner")
			runner, _, err := client.Actions.GetOrganizationRunner(context.TODO(), org, runnerID)
			if err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Name", "OS", "Status", "Busy"})
			t.AppendRow(table.Row{runner.GetID(), runner.GetName(), runner.GetOS(), runner.GetStatus(), *runner.Busy})
			t.Render()

			return nil
		},
	}

	cmd.PersistentFlags().Int64P("runner-id", "r", -1, "Unique identifier of the self-hosted runner.")
	cmd.MarkFlagRequired("runner-id")

	return cmd
}
