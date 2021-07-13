package cmd

import (
	"github.com/spf13/cobra"
)

// rdsCmd represents the rds command
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Interact with Amazon RDS",
	Long:  `Interact with Amazon RDS. See subcommands for help`,
}

func init() {
	rootCmd.AddCommand(rdsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rdsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rdsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
