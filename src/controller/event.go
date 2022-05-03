package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mickambar19/authorizer/src/model"
)

type AccountUsecase interface {
	CreateAccount(bool, int) (model.Account, []model.Violation)
}

type TransactionUsecase interface {
	CreateTransaction(int, string, time.Time) (model.Account, []model.Violation)
}

type Event struct {
	au AccountUsecase
	tu TransactionUsecase
}

func NewEvent(au AccountUsecase, tu TransactionUsecase) *Event {
	return &Event{
		au,
		tu,
	}
}

// HandleEvents reads
func (e *Event) HandleEvents() {
	var event model.Event
	decoder := json.NewDecoder(os.Stdin)
	var output string
	for {
		err := decoder.Decode(&event)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		switch event.Type {
		case "account":
			output = e.ProcessAccountEvent(event)
		case "transaction":
			output = e.ProcessTransactionEvent(event)
		}
		fmt.Println(output)
	}
}

func (e *Event) ProcessAccountEvent(event model.Event) string {
	activeCard, availableLimit := event.Data["active-card"].(bool), event.Data["available-limit"].(float64)
	account, violations := e.au.CreateAccount(activeCard, int(availableLimit))

	return e.GenerateEventOutput(account, violations)
}

func (e *Event) ProcessTransactionEvent(event model.Event) string {
	amount, merchant, dateString := event.Data["amount"].(float64), event.Data["merchant"].(string), event.Data["time"].(string)
	date, _ := time.Parse(time.RFC3339, dateString)
	account, violations := e.tu.CreateTransaction(int(amount), merchant, date)
	return e.GenerateEventOutput(account, violations)
}

func (e *Event) GenerateEventOutput(account model.Account, violations []model.Violation) string {

	result := map[string]interface{}{
		"violations": violations,
	}

	if (model.Account{}) != account {
		result["account"] = account
	} else {
		result["account"] = map[string]interface{}{}
	}

	data, _ := json.Marshal(result)
	text := string(data)
	text = strings.ReplaceAll(text, "\":", "\": ")
	text = strings.ReplaceAll(text, ",\"", ", \"")
	return text
}
