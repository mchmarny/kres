package common

import (
	"time"
)

// SimpleStock is a placeholder for now, need to refactor to kapi module
type SimpleStock struct {
	ID        string                 `json:"id"`
	Symbol    string                 `json:"symbol"`
	RequestOn time.Time              `json:"requestOn"`
	Data      map[string]interface{} `json:"data"`
}
