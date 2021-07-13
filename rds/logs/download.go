package logs

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func DownloadDays(ctx context.Context, days int, dbIdentifier string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	svc := rds.NewFromConfig(cfg)
	input := &rds.DescribeDBLogFilesInput{
		DBInstanceIdentifier: aws.String(dbIdentifier),
		FileLastWritten:      time.Now().Add(-time.Hour*24*time.Duration(days)).Unix() * 1000,
	}

	result, err := svc.DescribeDBLogFiles(ctx, input)
	if err != nil {
		return err
	}

	for _, r := range result.DescribeDBLogFiles {
		func(name string) {
			fmt.Printf("Downloading: %s\n", name)
			if err := Download(ctx, name, dbIdentifier); err != nil {
				fmt.Printf("Error downloading %s: %s\n", name, err.Error())
			}

		}(*r.LogFileName)
	}

	return nil
}

func Download(ctx context.Context, fileName, dbIdentifier string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	svc := rds.NewFromConfig(cfg)

	downloadDBLogFilePortionInput := &rds.DownloadDBLogFilePortionInput{
		DBInstanceIdentifier: &dbIdentifier,
		LogFileName:          &fileName,
	}

	out, err := os.Create(path.Base(fileName))
	if err != nil {
		return err
	} else {
		defer out.Close()
	}

	lfp := rds.NewDownloadDBLogFilePortionPaginator(svc, downloadDBLogFilePortionInput)

	for {
		result, err := lfp.NextPage(ctx)
		if err != nil {
			return err
		}

		fmt.Println(fileName, result.AdditionalDataPending, lfp.HasMorePages())
		if result.LogFileData != nil {
			out.WriteString(*result.LogFileData)
		}
		if !result.AdditionalDataPending {
			break
		}

	}

	return nil
}
