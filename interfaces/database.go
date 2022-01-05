package interfaces

import "github.com/Ali-Farhadnia/LeitnerBoxCore/models"

type Database interface {
	AddNewCard(models.Card) error
	GetCards() ([]models.Card, error)
	FindById(string) (models.Card, error)
	UpdateCard(models.Card) error
	DeleteCard(string) error
}
