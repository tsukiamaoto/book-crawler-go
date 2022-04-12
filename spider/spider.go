package spider

import (
	"github.com/tsukiamaoto/book-crawler-go/configs"
	"github.com/tsukiamaoto/book-crawler-go/models"

	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

type Spider struct {
	Collector *colly.Collector
}

func (s *Spider) init() {
	s.Collector = colly.NewCollector(
		colly.MaxDepth(3),
		colly.Async(),
		colly.Debugger(&debug.LogDebugger{}),
	)

	extensions.RandomUserAgent(s.Collector)

	s.Collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
	})
}

func New() *Spider {
	s := &Spider{}
	s.init()

	return s
}

func (s *Spider) Run() {
	s.visitEveryBookOnList("https://www.books.com.tw/web/sys_bbotm/books/010101/")

}

func (s *Spider) visitEveryBookOnList(url string) {
	// before visiting web
	s.Collector.OnRequest(func(r *colly.Request) {
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		s.Collector.SetCookies("books.com.tw", cookies)
		delayTimeForVisiting()

		fmt.Println("Visiting", r.URL)
		// q.AddRequest(r)
	})

	reqCh := make(chan int, 1)
	reqCh <- 0

	// collect every book information on the book list
	s.Collector.OnHTML(".wrap > .item", func(e *colly.HTMLElement) {
		reqCount := <-reqCh + 1
		reqCh <- reqCount

		if reqCount%10 == 0 {
			fmt.Println("Please wait for 10 second to avoid for blocked IP.")
			time.Sleep(10 * 1000 * time.Millisecond)
		}

		productLink := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		
		clone := s.Collector.Clone() // create a new Collector thread, avoid to get resources from the origin thread
		product := getNewProduct(clone, productLink)
		fmt.Println(product)
	})

	// go to the next book list page
	s.Collector.OnHTML(".nxt", func(e *colly.HTMLElement) {
		nextPageLink := e.Attr("href")

		e.Request.Visit(nextPageLink)
	})

	// q.AddURL("https://www.books.com.tw/web/sys_bbotm/books/010101/")
	s.Collector.Visit(url)
	s.Collector.Wait()
	// Wait until threads are finished
	// q.Run(c)
}

func getNewProduct(c *colly.Collector, url string) models.Product {
	var product models.Product

	// before visiting web
	c.OnRequest(func(r *colly.Request) {
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		c.SetCookies("books.com.tw", cookies)
		delayTimeForVisiting()

		fmt.Println("Visiting", r.URL.String())
	})

	// get product description
	c.OnHTML(".bd > .content", func(h *colly.HTMLElement) {
		description, err := h.DOM.Html()
		if err != nil {
			log.Println(err)
		}
		product.Description = strings.TrimSpace(description)
	})

	// get product detail infomation
	c.OnHTML("html", func(e *colly.HTMLElement) {
		name := e.ChildText(".grid_10 > div > h1")
		product.Name = name

		editor := e.ChildText(".type02_p003 > ul > li:first-child > a:nth-child(2)")
		product.Editor = editor

		publisherChildIndex := findTagChildIndex(e, ".type02_p003 > ul > li", "出版社")
		queryPublisher := fmt.Sprintf(".type02_p003 > ul > li:nth-child(%d) > a:first-child", publisherChildIndex)
		publisher := e.ChildText(queryPublisher)
		product.Publisher = publisher

		publicationDateChildIndex := findTagChildIndex(e, ".type02_p003 > ul > li", "出版日期")
		querypublicationDate := fmt.Sprintf(".type02_p003 > ul > li:nth-child(%d)", publicationDateChildIndex)
		publicationDate := strings.Split(e.ChildText(querypublicationDate), "：")
		if len(publicationDate) > 0 && publicationDate[0] == "出版日期" {
			product.PublicationDate = publicationDate[1]
		}

		product.Categories = append(product.Categories, getCategory(e))

	})

	c.Visit(url)
	c.Wait()

	return product
}

func getCategory(e *colly.HTMLElement) models.Category {
	var Category models.Category

	e.ForEach(".type04_breadcrumb > li", func(_ int, e *colly.HTMLElement) {
		categoryType := e.ChildText("a")
		if categoryType != "" && categoryType != "博客來" {
			Category.Types = append(Category.Types, categoryType)
		}
	})

	mainImage := e.ChildAttr("#M201106_0_getTakelook_P00a400020052 > img", "src")
	Category.Images = append(Category.Images, mainImage)

	e.ForEach(".items > ul > li", func(_ int, e2 *colly.HTMLElement) {
		image := e2.ChildAttr("img", "src")
		Category.Images = append(Category.Images, image)
	})

	price := e.ChildText(".price > li:first-child > em")
	discountPrice := e.ChildText(".price01")
	if discountPrice != "" {
		intPrice, _ := strconv.Atoi(discountPrice)
		Category.Price = intPrice
	} else {
		intPrice, _ := strconv.Atoi(price)
		Category.Price = intPrice
	}

	return Category
}

func findTagChildIndex(e *colly.HTMLElement, query, tag string) int {
	var childIndex int

	e.ForEach(query, func(index int, e2 *colly.HTMLElement) {
		r, _ := regexp.Compile(tag)
		if r.MatchString(e2.Text) {
			childIndex = index + 1
		}
	})

	return childIndex
}

func delayTimeForVisiting() {
	var delayTime = []int{5, 7, 8}
	rand.Seed(time.Now().UnixNano())
	delayTimeIndex := rand.Intn(len(delayTime))
	time.Sleep(time.Duration(delayTime[delayTimeIndex]*1000) * time.Millisecond)
}
