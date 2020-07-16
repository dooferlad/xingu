/*
Copyright © 2020 James Tunnicliffe <dooferlad@nanosheep.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/dooferlad/xingu/ec2"

	"github.com/spf13/cobra"
)

var name string
var filter map[string]string

func defaultArgFromPositionS(args []string, arg *string) error {
	if len(args) > 0 {
		if *arg == "" && len(args) == 1 {
			*arg = args[0]
			return nil
		}
		return fmt.Errorf("unknown arguments: %v", args)
	}
	return nil

}

// ec2ListCmd represents the ec2ListCmd command
var ec2ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ec2 instances",
	Args: func(cmd *cobra.Command, args []string) error {
		return defaultArgFromPositionS(args, &name)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return ec2.List(name, filter)
	},
}

func init() {
	ec2Cmd.AddCommand(ec2ListCmd)
	ec2ListCmd.Flags().StringToStringVar(&filter, "filters", nil, "ec2 instance filters")
	ec2ListCmd.Flags().StringVarP(&name, "name", "n", "", "ec2 instance name")
}
