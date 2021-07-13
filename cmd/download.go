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

	"github.com/dooferlad/xingu/rds/logs"
	"github.com/spf13/cobra"
)

var FileName, Database string
var Days int

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an RDS log file",
	Long: `Download an RDS log file:

  xingu rds logs download --filename errors/something --database production`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if FileName != "" {
			err = logs.Download(cmd.Context(), FileName, Database)
		} else if Days != 0 {
			err = logs.DownloadDays(cmd.Context(), Days, Database)
		} else {
			err = fmt.Errorf("days or filename must be specified")
		}

		if err != nil {
			fmt.Printf("Error downloading logs: %v\n", err)
		}

		return err
	},
}

func init() {
	logsCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&FileName, "filename", "f", "", "File name to download")
	downloadCmd.Flags().IntVarP(&Days, "days", "", 0, "Days of logs to download")

	downloadCmd.Flags().StringVarP(&Database, "database", "d", "", "RDS database name")
	downloadCmd.MarkFlagRequired("database")

	logsCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&Database, "database", "d", "", "RDS database name")
	listCmd.MarkFlagRequired("database")
}
