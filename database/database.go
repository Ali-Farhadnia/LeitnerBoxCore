package database

import (
	"errors"
	"fmt"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
)

type DB struct {
	Carts []*models.Cart
	Len   int
}

func NewDb() *DB {
	return &DB{Carts: make([]*models.Cart, 0)}
}

func (db *DB) AddNewCart(cart models.Cart) error {
	db.Carts = append(db.Carts, &cart)
	db.Len++
	return nil
}
func (db *DB) GetCarts() ([]*models.Cart, error) {
	return db.Carts, nil
}
func (db *DB) FindById(id string) (models.Cart, error) {
	for _, cart := range db.Carts {
		if cart.Id == id {
			return *cart, nil
		}
	}
	return models.Cart{}, errors.New("no cart with this id")
}
func (db *DB) UpdateCart(cart models.Cart) error {
	fmt.Println("in update")
	defer fmt.Println("in update")
	for i := 0; i < db.Len; i++ {
		if db.Carts[i].Id == cart.Id {
			db.Carts[i] = &cart
			return nil
		}
	}

	return errors.New("no cart with this id")
}
