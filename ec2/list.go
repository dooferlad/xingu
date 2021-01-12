package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
)

func List(sess *session.Session, name string, filter map[string]string) error {
	result, err := list(sess, name, filter)
	if err != nil {
		return err
	}

	svc := ec2.New(sess)

	addressesArray, err := svc.DescribeAddresses(nil)
	if err != nil {
		return err
	}

	addresses := map[string]*ec2.Address{}

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

func list(sess *session.Session, name string, filter map[string]string) (*ec2.DescribeInstancesOutput, error) {
	var filters []*ec2.Filter

	if name != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{&name},
		})
	}

	for k, v := range filter {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(k),
			Values: []*string{&v},
		})
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				fmt.Println(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}
	return result, nil
}
