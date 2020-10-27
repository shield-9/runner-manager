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
	// "os"

	// "github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

func newRegistrationTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registration-token",
		Short: "Create a registration token for an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			org, _ := cmd.Flags().GetString("org")
			token, _, err := client.Actions.CreateOrganizationRegistrationToken(context.TODO(), org)
			if err != nil {
				return err
			}

			fmt.Println(token.GetToken())
			fmt.Printf("Valid until %s\n", token.GetExpiresAt())

			/*
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Token", "Expires At"})
				t.AppendRow(table.Row{token.GetToken(), token.GetExpiresAt()})
				t.Render()
			*/

			return nil
		},
	}

	return cmd
}
