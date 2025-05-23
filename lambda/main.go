package main

import (
	"lambda_func/api"
	"lambda_func/app"
	"lambda_func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUser(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(api.ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
