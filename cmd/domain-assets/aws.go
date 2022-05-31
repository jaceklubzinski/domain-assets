package main

import (
	"context"
	"domain-assets/pkg/awsclient"
	"domain-assets/pkg/dnsassets"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

//envAWS env variable configuration for AWS
type envAWS struct {
	AWSSecretAccessKey string `required:"true" split_words:"true"`
	AWSAccessKeyID     string `required:"true" split_words:"true"`
	AWSDefaultRegion   string `required:"true" split_words:"true"`
}

func initAWS() *dnsassets.ProviderAsset {
	var env envAWS
	if err := envconfig.Process("", &env); err != nil {
		log.Fatal(err.Error())
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(env.AWSDefaultRegion))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("unable to load SDK config")
	}
	awsClient := awsclient.NewAWSClient(cfg)
	awsProvider := dnsassets.NewDNSAsset(awsClient)
	return awsProvider
}
