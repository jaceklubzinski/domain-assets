package main

import (
	"domain-assets/pkg/dnsassets"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func appFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "sync-loop",
			Value:   "3600",
			Usage:   "How often DNS data should be synced in seconds",
			EnvVars: []string{"DNS_SYNC_LOOP"},
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "If set verbose logging is enabled",
		},
		&cli.BoolFlag{
			Name:  "trace",
			Usage: "Be even more verbose when logging stuff",
		},
	}
}

func initApp(c *cli.Context) {
	formatter := &log.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
	log.SetFormatter(formatter)
	if c.IsSet("trace") {
		log.SetLevel(log.TraceLevel)
	} else if c.IsSet("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
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
