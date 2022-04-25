package main

import (
	"github.com/tsukiamaoto/book-crawler-go/spider"

	log "github.com/sirupsen/logrus"
)

func init() {
	// define logrus output format
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func main() {
	spider := spider.New()
	spider.Run()
}
