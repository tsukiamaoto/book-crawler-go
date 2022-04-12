package configs

import (
	"github.com/tsukiamaoto/book-crawler-go/models"
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func GetCookies() []*http.Cookie {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fileName := "cookies.json"
	filePath := rootDir + "\\configs\\" + fileName
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cookies = make([]*models.Cookie, 0)
	if err := json.Unmarshal(byteValue, &cookies); err != nil {
		log.Println(err)
	}

	var cookies2HttpCookies = make([]*http.Cookie, 0)
	for _, cookie := range cookies {
		httpCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			HttpOnly: cookie.HttpOnly,
		}

		cookies2HttpCookies = append(cookies2HttpCookies, httpCookie)
	}

	return cookies2HttpCookies
}

func GetProxies() []string {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fileName := "proxy.json"
	filePath := rootDir + "\\configs\\" + fileName
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var proxies = make([]models.Proxy, 0)
	if err := json.Unmarshal(byteValue, &proxies); err != nil {
		log.Println(err)
	}

	var proxyUrls []string
	for _, proxy := range proxies {
		proxyUrl := fmt.Sprintf("%s://%s:%s", proxy.Protocol, proxy.IP, proxy.Port)
		proxyUrls = append(proxyUrls, proxyUrl)
	}

	defer jsonFile.Close()
	return proxyUrls
}
