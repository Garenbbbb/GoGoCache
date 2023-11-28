package mysqldb

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	dataSourceName := "root@tcp(localhost:3306)/gogocache"

	myDB, err := NewDB(dataSourceName)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer myDB.Close()

	// Example usage
	data, err := myDB.Get("John Doe")
	if err != nil {
		fmt.Println("Error getting data from the database:", err)
		return
	}
	if string(data) != "123-456-7890" {
		fmt.Printf("Error getting data from the database: data should be 123-456-7890, but ger %s", string(data))
	}
	fmt.Printf("Data: %s\n", data)
}
