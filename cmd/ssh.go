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
	"github.com/dooferlad/xingu/ec2"
	"github.com/dooferlad/xingu/session"
	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ssh into ec2 instance",
	Args: func(cmd *cobra.Command, args []string) error {
		return defaultArgFromPositionS(args, &name)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sess, err := session.New()
		if err != nil {
			return err
		}

		defer sess.SaveCreds()
		return ec2.SSH(sess.Session, name, filter)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd) // quick shortcut

	ec2Cmd.AddCommand(sshCmd)
	sshCmd.Flags().StringToStringVar(&filter, "filters", nil, "ec2 instance filters")
	sshCmd.Flags().StringVarP(&name, "name", "n", "", "ec2 instance name")
}
