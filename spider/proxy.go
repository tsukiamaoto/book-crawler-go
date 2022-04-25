package spider

import (
	Config "github.com/tsukiamaoto/book-crawler-go/config"

	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Data []string `json: "data"`
}

func GetProxies() []string {
	config := Config.LoadConfig()
	api := "/api/v1/proxy"
	proxyUrl := fmt.Sprintf("http://%s%s", config.ProxyServerAddr, api)

	res, err := http.Get(proxyUrl)
	if err != nil {
		log.Error("無法請求Proxy Server的資料, 原因是:", err)
	}
	defer res.Body.Close()

	var result *Response
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Error("無法解析Proxy Server內容, 原因是:", err)
	}
	proxies := result.Data

	return proxies
}
