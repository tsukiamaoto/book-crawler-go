package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	HttpOnly bool   `json:"httpOnly"`
}

func GetCookies() []*http.Cookie {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}

	fileName := "cookies.json"
	filePath := rootDir + "\\config\\" + fileName
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cookies = make([]*Cookie, 0)
	if err := json.Unmarshal(byteValue, &cookies); err != nil {
		log.Error(err)
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
