package ec2

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func List(ctx context.Context, cfg aws.Config, name string, filter map[string]string) error {
	result, err := list(ctx, cfg, name, filter)
	if err != nil {
		return err
	}

	svc := ec2.NewFromConfig(cfg)

	addressesArray, err := svc.DescribeAddresses(ctx, nil)
	if err != nil {
		return err
	}

	addresses := map[string]types.Address{}

	for _, a := range addressesArray.Addresses {
		if a.InstanceId != nil {
			addresses[*a.InstanceId] = a
		}
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			var privateIPAddress, publicIPAddress string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
			}
			if instance.PrivateIpAddress != nil {
				privateIPAddress = *instance.PrivateIpAddress
			}

			if a, ok := addresses[*instance.InstanceId]; ok {
				publicIPAddress = *a.PublicIp
			}

			fmt.Printf("%s %30s %15s %15s\n", *instance.InstanceId, name, privateIPAddress, publicIPAddress)
		}
	}

	return nil
}

func list(ctx context.Context, cfg aws.Config, name string, filter map[string]string) (*ec2.DescribeInstancesOutput, error) {
	var filters []types.Filter

	if name != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:Name"),
			Values: []string{name},
		})
	}

	for k, v := range filter {
		filters = append(filters, types.Filter{
			Name:   aws.String(k),
			Values: []string{v},
		})
	}

	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	result, err := svc.DescribeInstances(ctx, input)
	if err != nil {
		var dbNotFound rdsTypes.DBInstanceNotFoundFault
		if errors.As(err, &dbNotFound) {
			fmt.Println(dbNotFound.ErrorCode(), dbNotFound.ErrorMessage())
		} else {
			fmt.Println(err.Error())
		}
		return nil, err
	}
	return result, nil
}
