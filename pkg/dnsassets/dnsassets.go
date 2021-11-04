package dnsassets

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	Name            string
	RecordType      string
	Description     string
	DNSZone         string
	RecordProvider  string
	ResourceRecords string
	AddetAt         string
	LastUpdate      string
	Status          string
}

type ProviderAsset struct {
	Provider ProviderAsseter
}

type ProviderAsseter interface {
	ListDomains() ([]Inventory, error)
}

func NewDNSAsset(p ProviderAsseter) *ProviderAsset {
	return &ProviderAsset{
		Provider: p,
	}
}
