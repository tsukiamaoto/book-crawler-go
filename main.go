package main

import (
	"github.com/tsukiamaoto/book-crawler-go/spider"
)

func main() {
	spider := spider.New()
	spider.Run()
}
