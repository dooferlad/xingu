package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/dooferlad/jat/shell"

	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/aws-sdk-go/aws/session"
)

func SSMConnect(sess *session.Session, name string, filter map[string]string) error {
	result, err := list(sess, name, filter)
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

			if *instance.State.Code == 16 { // Running
				return ssmConnect(sess, instance, name)
			}
		}
	}

	return nil
}

func ssmConnect(sess *session.Session, instance *ec2.Instance, name string) error {
	fmt.Printf("aws ssm start-session %s # %s\n", *instance.InstanceId, name)
	svc := ssm.New(sess, nil)
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	sessionConfig := &ssm.StartSessionInput{Target: instance.InstanceId}
	startSessionOutput, err := svc.StartSessionWithContext(subctx, sessionConfig)
	if err != nil {
		return err
	}

	sessJson, err := json.Marshal(startSessionOutput)
	if err != nil {
		return err
	}

	paramsJson, err := json.Marshal(sessionConfig)
	if err != nil {
		return err
	}

	shell.Shell(
		"session-manager-plugin",
		string(sessJson),
		*sess.Config.Region,
		"StartSession",
		os.Getenv("AWS_PROFILE"),
		string(paramsJson),
		*instance.InstanceId,
	)

	if _, err := svc.TerminateSessionWithContext(subctx, &ssm.TerminateSessionInput{SessionId: startSessionOutput.SessionId}); err != nil {
		return err
	}

	return nil
}
