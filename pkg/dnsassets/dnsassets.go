package dnsassets

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	Name            string    `json:"name"`
	RecordType      string    `json:"type"`
	Description     string    `json:"description"`
	DNSZone         string    `json:"zone"`
	RecordProvider  string    `json:"provider"`
	ResourceRecords string    `json:"resource"`
	AddedAt         time.Time `json:"added_at"`
	LastUpdate      time.Time `json:"last_update"`
	Status          string    `json:"status"`
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
