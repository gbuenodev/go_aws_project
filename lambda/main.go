package main

import (
	"lambda_func/app"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	myApp := app.NewApp()
	lambda.Start(myApp.ApiHandler.RegisterUserHandler)
}
