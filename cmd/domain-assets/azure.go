package main

import (
	"domain-assets/pkg/azureclient"
	"domain-assets/pkg/dnsassets"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func azureFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "azure-subscription-id",
			Value:    "xxx",
			Usage:    "Azure CLI subscription ID",
			EnvVars:  []string{"AZURE_SUBSCRIPTION_ID"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "azure-resource-group",
			Value:    "xxx",
			Usage:    "Azure CLI resource group",
			EnvVars:  []string{"AZURE_RESOURCE_GROUP"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "azure-zone-name",
			Value:    "xxx.xx",
			Usage:    "azure DNS zone name",
			EnvVars:  []string{"AZURE_ZONE_NAME"},
			Required: true,
		},
	}
}

func runAzure(c *cli.Context) []dnsassets.Inventory {
	azureAuthorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Unable to load Azure SDK auth from CLI")
	}
	azureDNSClient := azureclient.NewAzureClient(azureAuthorizer, azureclient.Config{
		ResourceGroup:  c.String("azure-resource-group"),
		ZoneName:       c.String("azure-zone-name"),
		SubscriptionID: c.String("azure-subscription-id"),
	})

	azureProvider := dnsassets.NewDNSAsset(azureDNSClient)
	azureAssets, err := azureProvider.Provider.ListDomains()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Unable to provide Azure assets")
	}

	return azureAssets
}
