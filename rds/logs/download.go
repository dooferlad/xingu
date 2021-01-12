package logs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/dooferlad/xingu/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/rds"
)

func DownloadDays(days int, dbIdentifier string) error {
	sess, err := session.New()
	if err != nil {
		return err
	}

	svc := rds.New(sess)
	input := &rds.DescribeDBLogFilesInput{
		DBInstanceIdentifier: aws.String(dbIdentifier),
		FileLastWritten:      aws.Int64(time.Now().Add(-time.Hour*24*time.Duration(days)).Unix() * 1000),
	}

	result, err := svc.DescribeDBLogFiles(input)
	if err != nil {
		return err
	}

	for _, r := range result.DescribeDBLogFiles {
		if err := Download(*r.LogFileName, dbIdentifier); err != nil {
			return err
		}
	}

	return nil
}

func Download(fileName, dbIdentifier string) error {
	sess, err := session.New()
	if err != nil {
		return err
	}

	creds := defaults.CredChain(sess.Config, sess.Handlers)
	signer := v4.NewSigner(creds)

	region := *sess.Config.Region

	url := fmt.Sprintf(
		"https://rds.%s.amazonaws.com/v13/downloadCompleteLogFile/%s/%s",
		region,
		dbIdentifier,
		fileName,
	)

	request, _ := http.NewRequest("GET", url, nil)
	_, err = signer.Presign(request, nil, rds.ServiceName, region, 1*time.Hour, time.Now())
	if err != nil {
		return err
	}
	fmt.Println(request.URL)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error downloading log: %s", resp.Status)
	}

	defer resp.Body.Close()

	if out, err := os.Create(path.Base(fileName)); err != nil {
		return err
	} else {
		defer out.Close()
		io.Copy(out, resp.Body)
	}

	return nil
}
