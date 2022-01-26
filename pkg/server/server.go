package server

import (
	"domain-assets/pkg/dnsassets"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	storage *gorm.DB
}

func (s *Server) Start() {
	router := gin.Default()
	router.GET("/dns", s.getAll)
	router.GET("/dns/name/:name", s.getDNSbyName)
	router.GET("/dns/resource/:name", s.getDNSbyResource)
	router.GET("/dns/lastadded", s.getLastAdded)
	router.GET("/dns/inactive", s.getInactive)
	router.Run("localhost:8080")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("Unable to load Azure SDK auth from CLI")
	}
}

func (s *Server) getAll(c *gin.Context) {
	var i []dnsassets.Inventory
	r := s.storage.Find(&i)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"error": r.Error,
		}).Fatalln("unable to check DNS asset")
	}
	c.IndentedJSON(http.StatusOK, i)
}

func (s *Server) getDNSbyName(c *gin.Context) {
	name := c.Param("name")
	var i dnsassets.Inventory
	r := s.storage.Where("name = ?", name).First(&i)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"error": r.Error,
		}).Fatalln("unable to check DNS asset")
	}
	c.IndentedJSON(http.StatusOK, i)
}

func (s *Server) getDNSbyResource(c *gin.Context) {
	name := c.Param("name")
	var i []dnsassets.Inventory
	r := s.storage.Where("resource_records = ?", name).Find(&i)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"error": r.Error,
		}).Fatalln("unable to check DNS asset")
	}
	c.IndentedJSON(http.StatusOK, i)
}

func (s *Server) getLastAdded(c *gin.Context) {
	last := time.Now().AddDate(0, 0, -1)
	var i []dnsassets.Inventory
	r := s.storage.Where("added_at > ?", last).Find(&i)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"error": r.Error,
		}).Fatalln("unable to check DNS asset")
	}
	c.IndentedJSON(http.StatusOK, i)
}

func (s *Server) getInactive(c *gin.Context) {
	var i []dnsassets.Inventory
	r := s.storage.Where("status > ?", "Inactive").Find(&i)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"error": r.Error,
		}).Fatalln("unable to check DNS asset")
	}
	c.IndentedJSON(http.StatusOK, i)
}

func NewServer(db *gorm.DB) *Server {
	return &Server{
		storage: db,
	}
}
