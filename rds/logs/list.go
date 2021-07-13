package logs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func List(ctx context.Context, dbIdentifier string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	svc := rds.NewFromConfig(cfg)
	input := &rds.DescribeDBLogFilesInput{
		DBInstanceIdentifier: aws.String(dbIdentifier),
		FileLastWritten:      time.Now().Add(-time.Hour*48).Unix() * 1000,
	}

	result, err := svc.DescribeDBLogFiles(ctx, input)
	if err != nil {
		var dbNotFound rdsTypes.DBInstanceNotFoundFault
		if errors.As(err, &dbNotFound) {
			fmt.Println(dbNotFound.ErrorCode(), dbNotFound.ErrorMessage())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}

	fmt.Println(result)

	return nil
}
