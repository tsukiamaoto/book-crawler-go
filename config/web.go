package config

import (
	"net/http"

	"github.com/spf13/viper"
)

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	HttpOnly bool   `json:"httpOnly"`
}

func GetCookies() []*http.Cookie {
	config := viper.New()
	config.AddConfigPath(".")
	config.AddConfigPath("..")
	config.AddConfigPath("/config")
	config.SetConfigName("cookies")
	config.SetConfigType("json")

	config.AutomaticEnv()
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {
		panic("讀取設定檔出現錯誤，錯誤的原因為" + err.Error())
	}

	var cookies2HttpCookies = make([]*http.Cookie, 0)
	var cookies = make([]*Cookie, 0)
	viper.UnmarshalKey("cookies", cookies)

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
