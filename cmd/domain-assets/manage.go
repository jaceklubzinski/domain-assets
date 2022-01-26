package main

import (
	"domain-assets/pkg/dnsassets"
	"domain-assets/pkg/server"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func manageAssets(assets []dnsassets.Inventory) {
	db := initDB()

	for _, records := range assets {
		var i dnsassets.Inventory
		r := db.Where("name = ?", records.Name).First(&i)
		if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
			log.WithFields(log.Fields{
				"error": r.Error,
			}).Fatalln("unable to check DNS asset")
		}
		if r.RowsAffected > 0 {
			i.Status = "Active"
			i.LastUpdate = time.Now()
			r := db.Save(&i)
			if r.Error != nil {
				log.WithFields(log.Fields{
					"error": r.Error,
				}).Fatalln("unable to update DNS asset")
			}
		} else {
			r := db.Create(&dnsassets.Inventory{
				Name:            records.Name,
				RecordType:      records.RecordType,
				Description:     records.Description,
				DNSZone:         records.DNSZone,
				RecordProvider:  records.RecordProvider,
				ResourceRecords: records.ResourceRecords,
				Status:          "Active",
				AddedAt:         time.Now(),
			})
			if r.Error != nil {
				log.WithFields(log.Fields{
					"error": r.Error,
				}).Fatalln("unable to save DNS asset")
			}
		}
	}
	HTTPServer := server.NewServer(db)
	HTTPServer.Start()
}
