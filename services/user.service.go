package services

import "github.com/Fog3211/learn-golang-mongodb-api/models"

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
}
