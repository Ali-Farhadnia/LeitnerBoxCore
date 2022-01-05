package service

import (
	"fmt"
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/interfaces"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	uuid "github.com/satori/go.uuid"
)

func AddCard(data []byte, database interfaces.Database) error {
	id := uuid.NewV4().String()
	now := time.Now().Format("2006-01-02 03")
	t, err := time.Parse("2006-01-02 03", now)
	if err != nil {
		return err
	}
	newcart := models.Cart{ID: id, Data: data, CreateTime: t, Box: 1}

	return database.AddNewCard(newcart)
}
func Review(database interfaces.Database) ([]models.Cart, error) {
	allcarts, err := database.GetCards()
	wantedcarts := make([]models.Cart, 0)
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
	//fmt.Println("=========in confirm")
	//defer fmt.Println("=========in confirm")
	cart, err := database.FindById(id)

	if err != nil {
		return err
	}
	cart.Box += 1
	err = database.UpdateCard(cart)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
func RejectTheCard(id string, database interfaces.Database) error {
	cart, err := database.FindById(id)
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
func UpdateCart(card models.Cart, database interfaces.Database) error {
	err := database.UpdateCard(card)
	if err != nil {
		return err
	}
	return nil
}
func DeleteCart(id string, database interfaces.Database) error {
	err := database.DeleteCard(id)
	if err != nil {
		return err
	}
	return nil
}
