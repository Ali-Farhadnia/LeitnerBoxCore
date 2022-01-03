package service

import (
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/interfaces"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	uuid "github.com/satori/go.uuid"
)

func AddCard(data []byte, database interfaces.Database) error {
	id := uuid.NewV4().String()
	newcart := models.Cart{ID: id, Data: data, CreateTime: time.Now(), Box: 1}

	return database.AddNewCart(newcart)
}
func Review(database interfaces.Database) ([]*models.Cart, error) {
	allcarts, err := database.GetCarts()
	wantedcarts := make([]*models.Cart, 0)
	if err != nil {
		return nil, err
	}
	for _, cart := range allcarts {
		if cart.Box == 1 { //this must be some logic not just this it must be chainged
			wantedcarts = append(wantedcarts, cart)
		}
	}

	return wantedcarts, nil
}
func ConfirmTheCard(id string, database interfaces.Database) error {
	cart, err := database.FindById(id)
	if err != nil {
		return err
	}
	cart.Box += 1
	err = database.UpdateCart(cart)
	if err != nil {
		return err
	}

	return nil
}
func RejectTheCard(id string, database interfaces.Database) error {
	cart, err := database.FindById(id)
	if err != nil {
		return err
	}
	cart.Box = 1
	err = database.UpdateCart(cart)
	if err != nil {
		return err
	}

	return nil
}
