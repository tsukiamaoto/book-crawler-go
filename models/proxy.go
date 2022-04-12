package models

type Proxy struct {
	Country  string `json:"country"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}
