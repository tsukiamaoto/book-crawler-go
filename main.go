package main

import (
	"github.com/tsukiamaoto/book-crawler-go/spider"

	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func init() {
	// define logrus output format
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func main() {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(1).Hour().Do(task)
	scheduler.StartBlocking()
}

func task() {
	spider := spider.New()
	spider.Run()
}
