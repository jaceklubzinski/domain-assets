package main

import (
	"domain-assets/pkg/dnsassets"

	log "github.com/sirupsen/logrus"
)

func manageAssets(db *dnsassets.Store, assets []dnsassets.Inventory) {
	for _, records := range assets {
		exist, err := db.CheckIfExist(records.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatalln("unable to check Azure asset")
		}
		if exist {
			err := db.LastUpdate(records.Name)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatalln("unable to update Azure asset")
			}
		} else {
			err := db.AddRow(records.Name, records.RecordType, records.RecordProvider, records.DNSZone, records.ResourceRecords)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatalln("unable to save Azure asset")
			}
		}
	}
}
