package utils

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
)

func DelayTimeForVisiting() {
	var delayTime = []int{5, 7, 8}
	rand.Seed(time.Now().UnixNano())
	delayTimeIndex := rand.Intn(len(delayTime))
	time.Sleep(time.Duration(delayTime[delayTimeIndex]*1000) * time.Millisecond)
}

func FindTagChildIndex(e *colly.HTMLElement, query, tag string) int {
	var childIndex int

	e.ForEach(query, func(index int, e2 *colly.HTMLElement) {
		r, _ := regexp.Compile(tag)
		if r.MatchString(e2.Text) {
			childIndex = index + 1
		}
	})

	return childIndex
}
