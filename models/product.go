package models

import (
	"github.com/lib/pq"
)

type Product struct {
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Categories      []Category `json:"categories"`
	Editor          string     `json:"editor"`
	Publisher       string     `json:"publisher"`
	PublicationDate string     `json:"publicaionDate"`
}

type Category struct {
	Types     pq.StringArray
	Images    pq.StringArray `json:"images"`
	Price     int            `json:"price"`
	Inventory int            `json:"inventory"`
}
