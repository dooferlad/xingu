package cmd

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dooferlad/xingu/ec2"
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
		cfg, err := config.LoadDefaultConfig(cmd.Context())
		if err != nil {
			return err
		}

		return ec2.SSH(cmd.Context(), cfg, name, filter)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd) // quick shortcut

	ec2Cmd.AddCommand(sshCmd)
	sshCmd.Flags().StringToStringVar(&filter, "filters", nil, "ec2 instance filters")
	sshCmd.Flags().StringVarP(&name, "name", "n", "", "ec2 instance name")
}
