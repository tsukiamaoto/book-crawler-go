package spider

import (
	configs "github.com/tsukiamaoto/book-crawler-go/config"
	"github.com/tsukiamaoto/book-crawler-go/db"
	models "github.com/tsukiamaoto/book-crawler-go/model"
	"github.com/tsukiamaoto/book-crawler-go/redis"
	"github.com/tsukiamaoto/book-crawler-go/repository"
	"github.com/tsukiamaoto/book-crawler-go/service"
	"github.com/tsukiamaoto/book-crawler-go/utils"

	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	log "github.com/sirupsen/logrus"
)

type Spider struct {
	Collector *colly.Collector
	Redis     *redis.Redis
	DB        *db.DB
}

func (s *Spider) init() {
	s.Collector = colly.NewCollector(
		// colly.MaxDepth(3),
		colly.Async(),
		colly.Debugger(&debug.LogDebugger{}),
	)

	extensions.RandomUserAgent(s.Collector)

	s.Collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		Delay:       5 * time.Second,
		RandomDelay: 3 * time.Second,
	})

	// setProxy(s.Collector)
}

func New() *Spider {
	s := &Spider{}
	s.init()

	// connect to redis server
	s.Redis = redis.New()
	s.Redis.ConnectRDB()

	var hasCreatedDB bool
	// connect to database
	s.DB, hasCreatedDB = db.DbConnect()
	if !hasCreatedDB {
		db.AutoMigrate(s.DB)
	}

	return s
}

// get book information from 博客來
func (s *Spider) Run() {
	var wg sync.WaitGroup

	productCh := make(chan *models.Product, 10)
	urlOfMainPageBookListCh := make(chan string, 10)
	urlOfAllBookListCh := make(chan string, 10)

	homeUrl := "https://www.books.com.tw/"
	go s.getUrlOfMainPageBookListCh(urlOfMainPageBookListCh, urlOfAllBookListCh, homeUrl)
	go s.getNextPageOfBottomBookList(urlOfMainPageBookListCh, urlOfAllBookListCh)
	go s.visitEveryBookOnBottomBookList(urlOfAllBookListCh, productCh)

	// create a repository instance
	repos := repository.NewRepositories(s.DB)

	// create a service instance
	services := service.NewServices(service.Repos{
		Repos: repos,
	})

	for product := range productCh {
		if product.Name != "" {
			key := product.Name

			wg.Add(1)
			go func(product *models.Product) {
				// save product to redis for caching
				jsonProduct, _ := json.Marshal(product)
				s.Redis.Set(key, jsonProduct)
				// if product did't exist in the database, save product to database
				if p, _ := services.Products.GetProductByName(key); p.Name == "" {
					result, err := services.Products.AddProduct(product)
					if err != nil {
						log.Error("Faild to add product,the reason is ", err)
					}

					if err := services.Products.AddTypeByCategoryId(result.Categories[0].ID, result.Categories[0].Types); err != nil {
						log.Error("Failed to add type by category id, the reason is ", err)
					}
				}

				defer wg.Done()
			}(product)
			wg.Wait()
		}
	}
}

func (s *Spider) visitEveryBookOnBottomBookList(urlCh chan string, productCh chan *models.Product) {
	wg := sync.WaitGroup{}
	// before visiting web
	s.Collector.OnRequest(func(r *colly.Request) {
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		s.Collector.SetCookies("books.com.tw", cookies)
	})

	// collect every book information on the book list
	reqTimes := 0
	s.Collector.OnHTML(".wrap > .item", func(e *colly.HTMLElement) {
		productName := e.ChildText(".msg > h4 > a")
		if !s.Redis.Exists(productName) {
			reqTimes++
			if reqTimes%10 == 0 {
				setProxy(s.Collector)
			}

			productLink := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))

			wg.Add(1)
			go func(url string) {
				product := s.getNewProduct(url)
				productCh <- product
				defer wg.Done()
			}(productLink)
		}
	})

	for url := range urlCh {
		s.Collector.Visit(url)
	}
	s.Collector.Wait()
	wg.Wait()
	close(productCh)
}

func (s *Spider) getNewProduct(url string) *models.Product {
	var product = new(models.Product)
	productCollector := s.Collector.Clone() // add a thread to avoid for blocked

	// before visiting web
	productCollector.OnRequest(func(r *colly.Request) {
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		productCollector.SetCookies(".books.com.tw", cookies)

	})

	// get product description
	productCollector.OnHTML(".grid_19", func(e *colly.HTMLElement) {
		var (
			description string
			err         error
		)
		descriptionIndex := utils.FindTagChildIndex(e, "div", "內容簡介")
		queryDescription := fmt.Sprintf("div:nth-child(%d) > .bd > .content", descriptionIndex)

		e.ForEach(queryDescription, func(_ int, e2 *colly.HTMLElement) {
			description, err = e2.DOM.Html()
			if err != nil {
				log.Error(err)
			}
		})

		product.Description = strings.TrimSpace(description)
	})

	// get product detail infomation
	productCollector.OnHTML("html", func(e *colly.HTMLElement) {
		name := e.ChildText(".grid_10 > div > h1")
		product.Name = name

		editor := e.ChildText(".type02_p003 > ul > li:first-child > a:nth-child(2)")
		product.Editor = editor

		publisherChildIndex := utils.FindTagChildIndex(e, ".type02_p003 > ul > li", "出版社")
		queryPublisher := fmt.Sprintf(".type02_p003 > ul > li:nth-child(%d) > a:first-child", publisherChildIndex)
		publisher := e.ChildText(queryPublisher)
		product.Publisher = publisher

		publicationDateChildIndex := utils.FindTagChildIndex(e, ".type02_p003 > ul > li", "出版日期")
		querypublicationDate := fmt.Sprintf(".type02_p003 > ul > li:nth-child(%d)", publicationDateChildIndex)
		publicationDate := strings.Split(e.ChildText(querypublicationDate), "：")
		if len(publicationDate) > 0 && publicationDate[0] == "出版日期" {
			product.PublicationDate = publicationDate[1]
		}

		product.Categories = append(product.Categories, getCategory(e))

	})

	productCollector.Visit(url)
	productCollector.Wait()

	return product
}

func getCategory(e *colly.HTMLElement) models.Category {
	var Category = new(models.Category)

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

	return *Category
}

func (s *Spider) getUrlOfMainPageBookListCh(urlOfMainPageBookListCh, urlOfAllBookListCh chan string, url string) {
	// copy collector to avoid for blocked the same resource
	mainPageCollector := s.Collector.Clone()
	subPageCollector := s.Collector.Clone()

	// before visiting web
	mainPageReqTimes := 0
	mainPageCollector.OnRequest(func(r *colly.Request) {
		mainPageReqTimes++
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		mainPageCollector.SetCookies("books.com.tw", cookies)

		if mainPageReqTimes%10 == 0 {
			setProxy(mainPageCollector)
		}
	})

	subPageReqTimes := 0
	subPageCollector.OnRequest(func(r *colly.Request) {
		subPageReqTimes++
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		subPageCollector.SetCookies("books.com.tw", cookies)

		if subPageReqTimes%10 == 0 {
			setProxy(subPageCollector)
		}
	})

	// visit url of main book type
	mainPageCollector.OnHTML(".menu > li[data-key] > h3 > a", func(e *colly.HTMLElement) {
		name := e.Text
		url := e.Request.AbsoluteURL(e.Attr("href"))

		hasVisited, err := e.Request.HasVisited(url)
		if err != nil {
			log.Error(err)
		}

		if (name == "中文書" ||
			name == "簡體書") && !hasVisited {
			subPageCollector.Visit(url)
		}
	})

	// visit url of sub book type
	subPageCollector.OnHTML(".type02_l001-1 > ul > li > span > a", func(e *colly.HTMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))

		hasVisited, err := e.Request.HasVisited(url)
		if err != nil {
			log.Error(err)
		}

		if !hasVisited {
			e.Request.Visit(url)
		}
	})

	// visit url of middle level of sub book type
	subPageCollector.OnHTML(".sub > li > span > a", func(e *colly.HTMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))

		hasVisited, err := e.Request.HasVisited(url)
		if err != nil {
			log.Error(err)
		}

		// is in the page of bottom level of sub book type
		r, _ := regexp.Compile("sys_b{1,2}otm")
		URL2string := url
		if r.MatchString(URL2string) {
			urlOfMainPageBookListCh <- URL2string
			urlOfAllBookListCh <- URL2string
		}

		if !hasVisited {
			e.Request.Visit(url)
		}
	})

	mainPageCollector.Visit(url)
	mainPageCollector.Wait()
	subPageCollector.Wait()
	close(urlOfMainPageBookListCh)
}

func (s *Spider) getNextPageOfBottomBookList(urlOfMainPageBookListCh, urlOfAllBookListCh chan string) {
	// before visiting web
	reqTimes := 0
	s.Collector.OnRequest(func(r *colly.Request) {
		reqTimes++
		// Set header set
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")

		cookies := configs.GetCookies()
		s.Collector.SetCookies("books.com.tw", cookies)

		if reqTimes%10 == 0 {
			setProxy(s.Collector)
		}
	})

	// go to the next book list page
	s.Collector.OnHTML(".nxt", func(e *colly.HTMLElement) {
		nextPageUrl := e.Request.AbsoluteURL(e.Attr("href"))
		urlOfAllBookListCh <- nextPageUrl

		e.Request.Visit(nextPageUrl)
	})

	for url := range urlOfMainPageBookListCh {
		s.Collector.Visit(url)
	}
	s.Collector.Wait()
	close(urlOfAllBookListCh)
}

func setProxy(c *colly.Collector) {
	proxies := GetProxies()
	rp, err := proxy.RoundRobinProxySwitcher(proxies...)
	if err != nil {
		log.Error(err)
	}

	c.SetProxyFunc(rp)
}
