package logs

import (
	"fmt"
	"time"

	"github.com/dooferlad/xingu/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
)

func List(dbIdentifier string) error {
	sess, err := session.New()
	if err != nil {
		return err
	}

	svc := rds.New(sess)
	input := &rds.DescribeDBLogFilesInput{
		DBInstanceIdentifier: aws.String(dbIdentifier),
		FileLastWritten:      aws.Int64(time.Now().Add(-time.Hour*48).Unix() * 1000),
	}

	result, err := svc.DescribeDBLogFiles(input)
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
		return err
	}

	fmt.Println(result)

	return nil
}
