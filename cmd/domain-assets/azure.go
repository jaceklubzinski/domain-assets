package main

import (
	"domain-assets/pkg/azureclient"
	"domain-assets/pkg/dnsassets"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type envAzure struct {
	AzureSubscriptionID string `required:"true" split_words:"true"`
	AzureResourceGroup  string `required:"true" split_words:"true"`
	AzureZoneName       string `required:"true" split_words:"true"`
}

func initAzure() *dnsassets.ProviderAsset {
	var env envAzure
	if err := envconfig.Process("", &env); err != nil {
		log.Fatal(err.Error())
	}
	azureAuthorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Unable to load Azure SDK auth from CLI")
	}
	azureDNSClient := azureclient.NewAzureClient(azureAuthorizer, azureclient.Config{
		ResourceGroup:  env.AzureResourceGroup,
		ZoneName:       env.AzureZoneName,
		SubscriptionID: env.AzureSubscriptionID,
	})
	azureProvider := dnsassets.NewDNSAsset(azureDNSClient)
	log.Debugln("Azure provider provisioned")
	return azureProvider
}
