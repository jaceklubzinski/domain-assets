package azureclient

import (
	"context"
	"domain-assets/pkg/dnsassets"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2017-09-01/dns"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

type AzureClient struct {
	ZoneClient       dns.ZonesClient
	RecordSetsClient dns.RecordSetsClient
	Config
}

type Config struct {
	ResourceGroup  string
	ZoneName       string
	SubscriptionID string
}

func NewAzureClient(a autorest.Authorizer, c Config) *AzureClient {
	newZoneClient := dns.NewZonesClient(c.SubscriptionID)
	newZoneClient.Authorizer = a
	newRecordSetsClient := dns.NewRecordSetsClient(c.SubscriptionID)
	newRecordSetsClient.Authorizer = a
	return &AzureClient{
		ZoneClient:       newZoneClient,
		RecordSetsClient: newRecordSetsClient,
		Config:           c,
	}
}

func (a *AzureClient) ListDomainsByType(rt dns.RecordType) ([]dns.RecordSet, error) {
	var top int32 = 100

	azureRecords, err := a.RecordSetsClient.ListByType(context.TODO(), a.ResourceGroup, a.ZoneName, rt, &top, "")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Debug("failed to list dns records,")

		return azureRecords.Values(), err
	}

	return azureRecords.Values(), err
}

func (a *AzureClient) ListZOnes() (dns.ZoneListResultPage, error) {
	var (
		top int32 = 100
	)
	azureZones, err := a.ZoneClient.List(context.TODO(), &top)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Debug("failed to list zones,")

		return azureZones, err
	}
	return azureZones, nil
}

func (a *AzureClient) ListDomains() ([]dnsassets.Inventory, error) {
	var (
		azureAsset []dnsassets.Inventory
		records    []string
	)

	recordTypes := []dns.RecordType{"A", "CNAME"}
	for _, rt := range recordTypes {
		azureRecords, err := a.ListDomainsByType(rt)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Debug("failed to list dns records,")

			return azureAsset, err
		}
		for _, domain := range azureRecords {
			records = nil
			if string(rt) == "A" {
				for _, record := range *domain.ARecords {
					records = append(records, *record.Ipv4Address)
				}
			} else if string(rt) == "CNAME" {
				record := *domain.CnameRecord.Cname
				records = append(records, record)
			}

			a := dnsassets.Inventory{
				Name:            *domain.Fqdn,
				RecordType:      string(rt),
				DNSZone:         a.ZoneName,
				RecordProvider:  "Azure",
				ResourceRecords: strings.Join(records, ","),
			}
			azureAsset = append(azureAsset, a)
		}
	}

	return azureAsset, nil
}
