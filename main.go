package main

import (
	"github.com/tsukiamaoto/book-crawler-go/spider"

	log "github.com/sirupsen/logrus"
)

func init() {
	// define logrus out format
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	spider := spider.New()
	spider.Run()
}
