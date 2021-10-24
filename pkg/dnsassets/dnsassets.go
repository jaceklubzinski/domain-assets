package dnsassets

type Inventory struct {
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

type Store struct {
	Storer
}

type Storer interface {
	CreateDomainsTable() error
	GetRow() error
	AddRow(na, rt, dz, rp, rr string) error
	CheckIfExist(name string) (bool, error)
	LastUpdate(name string) error
}

func NewDB(d Storer) *Store {
	return &Store{
		Storer: d,
	}
}
