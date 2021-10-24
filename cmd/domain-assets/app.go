package main

import (
	"domain-assets/pkg/dnsassets"
	"domain-assets/pkg/sqlite"

	log "github.com/sirupsen/logrus"
)

func initApp() *dnsassets.Store {
	formatter := &log.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)

	sqliteDB, err := sqlite.NewsqliteDB()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Cant initialize sqlite database")
	}
	db := dnsassets.NewDB(sqliteDB)
	err = db.CreateDomainsTable()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Cant create table")
	}
	return db
}
