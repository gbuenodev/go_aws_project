package api

import (
	"fmt"
	"lambda_func/database"
	"lambda_func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}
	// does a user with this username already exist?
	userExists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("checking if user exists %w", err)
	}
	if userExists {
		return fmt.Errorf("a user with this username already exists")
	}

	// here we know that user does not exist
	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("error registering the user %w", err)
	}

	return nil
}
