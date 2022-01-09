package service

import (
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/interfaces"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	uuid "github.com/satori/go.uuid"
)

// AddCard() create card with given data and time.now() and add it to database and return it.
func AddCard(data []byte, database interfaces.Database) (models.Card, error) {
	id := uuid.NewV4().String()
	now := time.Now()
	newcard := models.NewCard()
	newcard.ID = id
	newcard.Data = data
	newcard.Box = 1
	newcard.CreateTime = &now
	err := database.AddNewCard(*newcard)
	if err != nil {
		return *models.NewCard(), err
	}
	return *newcard, nil
}

// Review() return all card with box==1.
func Review(database interfaces.Database) ([]models.Card, error) {
	allcarts, err := database.GetCards()
	if err != nil {
		return nil, err
	}
	wantedcarts := make([]models.Card, 0)

	for _, cart := range allcarts {
		if cart.Box == 1 { //this must be some logic not just this it must be chainged.
			wantedcarts = append(wantedcarts, cart)
		}
	}

	return wantedcarts, nil
}

// ConfirmTheCard() get card id and increases it by one unit.
func ConfirmTheCard(id string, database interfaces.Database) error {
	cart, err := database.FindByID(id)

	if err != nil {
		return err
	}
	cart.Box++
	err = database.UpdateCard(cart)
	if err != nil {
		return err
	}

	return nil
}

// RejectTheCard() get card id and set it box to one.
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

// UpdateCard() get one card and update it in database.
func UpdateCard(card models.Card, database interfaces.Database) error {
	err := database.UpdateCard(card)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCard() get one card and remove it from database.
func DeleteCard(id string, database interfaces.Database) error {
	if err := database.DeleteCard(id); err != nil {
		return err
	}

	return nil
}
