package main

import (
	"encoding/csv"
	"os"
)

type Store struct {
	AreaCode  string
	StoreName string
	StoreID   string
}

var storeMasterMap map[string]Store

func LoadStoreMaster(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	storeMasterMap = make(map[string]Store)
	for i, rec := range records {
		if i == 0 {
			continue
		}
		if len(rec) < 3 {
			continue
		}
		store := Store{
			AreaCode:  rec[0],
			StoreName: rec[1],
			StoreID:   rec[2],
		}
		storeMasterMap[store.StoreID] = store
	}
	return nil
}

func StoreExists(storeID string) bool {
	_, exists := storeMasterMap[storeID]
	return exists
}
