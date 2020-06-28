package main

import (
	"fmt"
	"github.com/forensicanalysis/forensicstore"
	"github.com/patrickmn/go-cache"
	"time"
)

var queryCache *cache.Cache

func init() {
	queryCache = cache.New(5*time.Minute, 10*time.Minute)
}

func storequery(store *forensicstore.ForensicStore, q string) ([]forensicstore.JSONElement, error) {
	elems, found := queryCache.Get(q)
	if !found {
		fmt.Println(q)
		elems, err := store.Query(q)
		if err != nil {
			return nil, err
		}
		queryCache.Set(q, elems, cache.DefaultExpiration)
		return elems, nil
	}

	return elems.([]forensicstore.JSONElement), nil
}
