package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/dooferlad/jat/shell"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func SSMConnect(ctx context.Context, cfg aws.Config, name string, filter map[string]string) error {
	result, err := list(ctx, cfg, name, filter)
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
				return ssmConnect(ctx, cfg, instance, name)
			}
		}
	}

	return nil
}

func ssmConnect(ctx context.Context, cfg aws.Config, instance types.Instance, name string) error {
	fmt.Printf("aws ssm start-session %s # %s\n", *instance.InstanceId, name)
	svc := ssm.NewFromConfig(cfg)
	subctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	sessionConfig := &ssm.StartSessionInput{Target: instance.InstanceId}
	startSessionOutput, err := svc.StartSession(subctx, sessionConfig)
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
		cfg.Region,
		"StartSession",
		os.Getenv("AWS_PROFILE"),
		string(paramsJson),
		*instance.InstanceId,
	)

	if _, err := svc.TerminateSession(subctx, &ssm.TerminateSessionInput{SessionId: startSessionOutput.SessionId}); err != nil {
		return err
	}

	return nil
}
