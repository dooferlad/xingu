package ec2

import (
	"fmt"

	"github.com/dooferlad/jat/shell"
)

func SSH(name string, filter map[string]string) error {
	result, err := list(name, filter)
	if err != nil {
		return err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}

			if instance.PublicIpAddress != nil {
				fmt.Printf("ssh %s # %s\n", *instance.PublicIpAddress, name)
				return shell.Shell("ssh", *instance.PublicIpAddress)
			}

			if instance.PrivateIpAddress != nil {
				fmt.Printf("ssh %s # %s\n", *instance.PrivateIpAddress, name)
				return shell.Shell("ssh", *instance.PrivateIpAddress)
			}
		}
	}

	return nil
}
