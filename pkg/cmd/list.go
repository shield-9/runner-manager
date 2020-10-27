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
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List self-hosted runners for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			org, _ := cmd.Flags().GetString("org")
			runners, _, err := client.Actions.ListOrganizationRunners(context.TODO(), org, nil)
			if err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Name", "OS", "Status", "Busy"})
			for _, runner := range runners.Runners {
				t.AppendRow(table.Row{runner.GetID(), runner.GetName(), runner.GetOS(), runner.GetStatus(), *runner.Busy})
			}
			t.Render()

			fmt.Printf("Total: %d (Online: %d, Offline: %d)\n", len(runners.Runners),
				countRunnersByStatus(runners.Runners, "online"), countRunnersByStatus(runners.Runners, "offline"))

			return nil
		},
	}

	return cmd
}

func countRunnersByStatus(runners []*github.Runner, status string) (count uint) {
	count = 0

	for _, runner := range runners {
		if runner.GetStatus() == status {
			count++
		}
	}
	return count
}
