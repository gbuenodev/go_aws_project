package app

import (
	"lambda_func/api"
	"lambda_func/store"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	userStore := store.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(userStore)

	return App{
		ApiHandler: apiHandler,
	}
}
