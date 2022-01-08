package service

import (
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/interfaces"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	uuid "github.com/satori/go.uuid"
)

func AddCard(data []byte, database interfaces.Database) (models.Card, error) {
	id := uuid.NewV4().String()
	now := time.Now() //.UTC()
	newcart := models.Card{ID: id, Data: data, CreateTime: &now, Box: 1}
	err := database.AddNewCard(newcart)
	if err != nil {
		return models.Card{}, err
	}
	return newcart, nil
}
func Review(database interfaces.Database) ([]models.Card, error) {
	allcarts, err := database.GetCards()
	if err != nil {
		return nil, err
	}
	wantedcarts := make([]models.Card, 0)

	for _, cart := range allcarts {
		if cart.Box == 1 { //this must be some logic not just this it must be chainged
			wantedcarts = append(wantedcarts, cart)
		}
	}

	return wantedcarts, nil
}
func ConfirmTheCard(id string, database interfaces.Database) error {
	cart, err := database.FindByID(id)

	if err != nil {
		return err
	}
	cart.Box += 1
	err = database.UpdateCard(cart)
	if err != nil {
		return err
	}

	return nil
}
func RejectTheCard(id string, database interfaces.Database) error {
	cart, err := database.FindByID(id)
	if err != nil {
		return err
	}
	cart.Box = 1
	err = database.UpdateCard(cart)
	if err != nil {
		return err
	}

	return nil
}
func UpdateCard(card models.Card, database interfaces.Database) error {
	err := database.UpdateCard(card)
	if err != nil {
		return err
	}
	return nil
}
func DeleteCard(id string, database interfaces.Database) error {
	if err := database.DeleteCard(id); err != nil {
		return err
	}

	return nil
}
