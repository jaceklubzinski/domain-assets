package main

import (
	"context"
	"domain-assets/pkg/awsclient"
	"domain-assets/pkg/dnsassets"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func awsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "aws-access-keyi-id",
			Value:    "xxx",
			Usage:    "AWS CLI access ID",
			EnvVars:  []string{"AWS_ACCESS_KEY_ID"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "aws-secret-access-key",
			Value:    "xxx",
			Usage:    "AWS CLI secret access key",
			EnvVars:  []string{"AWS_SECRET_ACCESS_KEY"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "aws-default-region",
			Value:    "us-east-1",
			Usage:    "AWS region",
			EnvVars:  []string{"AWS_DEFAULT_REGION"},
			Required: true,
		},
	}
}

func runAWS(c *cli.Context) []dnsassets.Inventory {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(c.String("aws-default-region")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.String("aws-access-keyi-id"), c.String("aws-secret-access-key"), "")),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("unable to load SDK config")
	}
	awsClient := awsclient.NewAWSClient(cfg)

	awsProvider := dnsassets.NewDNSAsset(awsClient)
	awsAssets, err := awsProvider.Provider.ListDomains()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("unable to provide AWS assets")
	}

	return awsAssets
}
