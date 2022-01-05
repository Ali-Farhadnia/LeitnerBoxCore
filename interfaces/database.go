package interfaces

import "github.com/Ali-Farhadnia/LeitnerBoxCore/models"

type Database interface {
	AddNewCard(models.Cart) error
	GetCards() ([]models.Cart, error)
	FindById(string) (models.Cart, error)
	UpdateCard(models.Cart) error
	DeleteCard(string) error
}
