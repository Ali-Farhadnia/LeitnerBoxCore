package interfaces

import "github.com/Ali-Farhadnia/LeitnerBoxCore/models"

type Database interface {
	AddNewCart(models.Cart) error
	GetCarts() ([]*models.Cart, error)
	FindById(string) (models.Cart, error)
	UpdateCart(models.Cart) error
}
