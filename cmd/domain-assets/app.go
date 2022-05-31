package main

import (
	"domain-assets/pkg/dnsassets"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initApp() {
	formatter := &log.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
	//only to verify env variables just after initial run
	var envAz envAzure
	if err := envconfig.Process("", &envAz); err != nil {
		log.Fatal(err.Error())
	}
	var envAWS envAWS
	if err := envconfig.Process("", &envAWS); err != nil {
		log.Fatal(err.Error())
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Cant initialize sqlite database")
	}
	err = db.AutoMigrate(&dnsassets.Inventory{})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Cant auto migrate database")
	}
	return db
}
