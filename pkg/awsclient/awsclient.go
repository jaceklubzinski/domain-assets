package awsclient

import (
	"context"
	"domain-assets/pkg/dnsassets"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	log "github.com/sirupsen/logrus"
)

type AWSClient struct {
	cfg *route53.Client
}

func NewAWSClient(cfg aws.Config) *AWSClient {
	c := route53.NewFromConfig(cfg)
	return &AWSClient{cfg: c}
}

func (a *AWSClient) ListZones() ([]types.HostedZone, error) {
	zones, err := a.cfg.ListHostedZones(context.TODO(), nil)

	if err != nil {
		return nil, err
	}
	return zones.HostedZones, nil
}

func (a *AWSClient) ListZoneDomains(zoneID string) ([]types.ResourceRecordSet, error) {
	var allDmains []types.ResourceRecordSet
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	}

	domains, err := a.cfg.ListResourceRecordSets(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	for domains.IsTruncated {
		allDmains = append(allDmains, domains.ResourceRecordSets...)
		params = &route53.ListResourceRecordSetsInput{
			HostedZoneId:          aws.String(zoneID),
			StartRecordName:       domains.NextRecordName,
			StartRecordType:       domains.NextRecordType,
			StartRecordIdentifier: domains.NextRecordIdentifier,
		}
		domains, err = a.cfg.ListResourceRecordSets(context.TODO(), params)
		if err != nil {
			return nil, err
		}
	}
	return allDmains, nil
}

func (a *AWSClient) ListDomains() ([]dnsassets.Inventory, error) {
	var (
		awsAsset []dnsassets.Inventory
		records  []string
	)
	zones, err := a.ListZones()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Debug("failed to list zones,")

		return awsAsset, err
	}

	for _, zone := range zones {
		domains, err := a.ListZoneDomains(*zone.Id)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Debug("failed to list domain,")

			return awsAsset, err
		}

		for _, domain := range domains {
			if string(domain.Type) == "A" || string(domain.Type) == "CNAME" {
				records = nil
				for _, record := range domain.ResourceRecords {
					records = append(records, *record.Value)
				}

				a := dnsassets.Inventory{
					Name:            *domain.Name,
					RecordType:      string(domain.Type),
					DNSZone:         *zone.Name,
					RecordProvider:  "AWS",
					ResourceRecords: strings.Join(records, ","),
				}
				awsAsset = append(awsAsset, a)
			}
		}
	}
	return awsAsset, nil
}
