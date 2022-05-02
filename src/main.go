package main

import (
	"github.com/mickambar19/authorizer/src/controller"
	"github.com/mickambar19/authorizer/src/service"
	"github.com/mickambar19/authorizer/src/usecase"
)

func main() {

	accountManager := service.NewAccountManager()
	accountUsecase := usecase.NewAccount(accountManager)

	transactionsManager := service.NewTransactionsManager()
	transactionsUsecase := usecase.NewTransaction(transactionsManager, accountManager)

	eventController := controller.NewEvent(accountUsecase, transactionsUsecase)

	eventController.HandleEvents()
}
