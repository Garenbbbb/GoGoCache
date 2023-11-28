package cache

import (
	"fmt"
	"gogocache/mysqldb"
	"testing"
)

var db = map[string]string{
	"John Doe":      "123-456-7890",
	"Jane Smith":    "987-654-3210",
	"Alice Johnson": "555-123-4567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	dataSourceName := "root@tcp(localhost:3306)/gogocache"
	myDB, err := mysqldb.NewDB(dataSourceName)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer myDB.Close()

	gee := NewGroup("contacts", 2<<10, myDB)
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.Value() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := gee.Get("unknown"); err != nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
