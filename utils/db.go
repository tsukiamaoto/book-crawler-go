package utils

import "github.com/tsukiamaoto/book-crawler-go/model"

func RelationMap(data []string) map[string]string {
	relations := make(map[string]string)
	for index := range data {
		if index == 0 {
			relations["root"] = data[index]
		} else {
			child, parent := data[index], data[index-1]
			relations[parent] = child
		}
	}
	return relations
}

func BuildTypes(keys []string, relations map[string]string) []*model.Type {
	Types := make([]*model.Type, 0)
	for _, key := range keys {
		if name, ok := relations[key]; ok {
			Type := new(model.Type)
			Type.Name = name

			Types = append(Types, Type)
		}
	}
	return Types
}
