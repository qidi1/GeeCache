package GeeCache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"tom":  "630",
	"jack": "390",
	"ice":  "439",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("test_group", 2<<10, GetterFunc(func(key string) (value []byte, err error) {
		if v, ok := db[key]; ok {
			log.Println("[SlowDB] search key", key)
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%v not exist", key)
	}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %v", k)
		} // load from callback function
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
