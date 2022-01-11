package database_interface

import (
	"errors"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
)

// Database use in service.
type Database interface {
	AddNewCard(models.Card) error
	GetCards() ([]models.Card, error)
	FindByID(string) (*models.Card, error)
	UpdateCard(models.Card) error
	DeleteCard(string) error
}

// nothing found error.
var ErrNothingFound = errors.New("nothing found")

// somthings went wrong error.
var ErrSomthingWentWrong = errors.New("somthings went wrong")
