package queue

import (
	"log"
	"testing"
)

func TestToStock(t *testing.T) {

	json := `{"id":"id-03412ac5-2355-11e9-9944-0a580a1402df","symbol":"aaaa","requestOn":"2019-01-28T23:32:44.53083063Z"}`

	stock, err := toStock(json)

	if err != nil {
		t.Errorf("Error while converting string to stock: %v", err)
	}

	log.Printf("stock: %v", stock)

}
