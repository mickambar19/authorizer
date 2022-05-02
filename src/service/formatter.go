package service

import (
	"encoding/json"
	"fmt"

	"github.com/mickambar19/authorizer/src/model"
)

type Formatter struct {
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) GenerateJSON(activeCard bool, availableLimit int, violations []model.Violation) string {

	account := map[string]interface{}{
		"active-card":     activeCard,
		"available-limit": availableLimit,
	}

	info := map[string]interface{}{
		"account":    account,
		"violations": violations,
	}

	data, _ := json.Marshal(info)
	fmt.Println(data)
	return string(data)
}
