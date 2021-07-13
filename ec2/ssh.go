package ec2

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/spf13/viper"

	"github.com/dooferlad/jat/shell"
)

func SSH(ctx context.Context, cfg aws.Config, name string, filter map[string]string) error {
	userConfig := viper.GetStringMap(os.Getenv("AWS_PROFILE"))
	sshXinguConfig, ok := userConfig["ssh"]
	var sshArgs []string
	if ok {
		sshConfigFileMap, ok := sshXinguConfig.(map[string]interface{})["config"]
		if ok {
			sshArgs = append(sshArgs, "-F", sshConfigFileMap.(string))
		}
	}

	result, err := list(ctx, cfg, name, filter)
	if err != nil {
		return err
	}

	var ipAddress string

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}

			if instance.PublicIpAddress != nil {
				ipAddress = *instance.PublicIpAddress
			} else if instance.PrivateIpAddress != nil {
				ipAddress = *instance.PrivateIpAddress
			}

			if ipAddress != "" {
				sshArgs = append(sshArgs, ipAddress)
				fmt.Printf("ssh %s # %s\n", strings.Join(sshArgs, " "), name)
				return shell.Shell("ssh", sshArgs...)
			}
		}
	}

	return nil
}
