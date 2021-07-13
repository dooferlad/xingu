package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"

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
		ctx := cmd.Context()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}

		return ec2.List(ctx, cfg, name, filter)
	},
}

func init() {
	ec2Cmd.AddCommand(ec2ListCmd)
	ec2ListCmd.Flags().StringToStringVar(&filter, "filters", nil, "ec2 instance filters")
	ec2ListCmd.Flags().StringVarP(&name, "name", "n", "", "ec2 instance name")
}
