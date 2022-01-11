package service

import (
	"fmt"
	"time"

	database_interface "github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	uuid "github.com/satori/go.uuid"
)

// AddCard create card with given data and time.now and add it to database and return it.
func AddCard(data []byte, database database_interface.Database) (*models.Card, error) {
	id := uuid.NewV4().String()
	now := time.Now()
	newcard := models.NewCard()

	newcard.ID = id
	newcard.Data = data
	newcard.Box = 1
	newcard.CreateTime = &now

	if err := database.AddNewCard(*newcard); err != nil {
		return nil, fmt.Errorf("AddCard error :%w", err)
	}

	return newcard, nil
}

// Review return all card with box==1.
func Review(database database_interface.Database) ([]models.Card, error) {
	allcarts, err := database.GetCards()
	if err != nil {
		return nil, fmt.Errorf("Review error :%w", err)
	}

	wantedcarts := make([]models.Card, 0)

	for _, cart := range allcarts {
		// this must be some logic not just this it must be chainged.
		if cart.Box == 1 {
			wantedcarts = append(wantedcarts, cart)
		}
	}

	return wantedcarts, nil
}

// ConfirmTheCard get card id and increases it by one unit.
func ConfirmTheCard(id string, database database_interface.Database) error {
	cart, err := database.FindByID(id)
	if err != nil {
		return fmt.Errorf("ConfirmTheCard error :%w", err)
	}
	cart.Box++

	err = database.UpdateCard(*cart)
	if err != nil {
		return fmt.Errorf("ConfirmTheCard error :%w", err)
	}

	return nil
}

// RejectTheCard get card id and set it box to one.
func RejectTheCard(id string, database database_interface.Database) error {
	cart, err := database.FindByID(id)
	if err != nil {
		return fmt.Errorf("RejectTheCard error :%w", err)
	}

	cart.Box = 1

	err = database.UpdateCard(*cart)
	if err != nil {
		return fmt.Errorf("RejectTheCard error :%w", err)
	}

	return nil
}

// UpdateCard get one card and update it in database.
func UpdateCard(card models.Card, database database_interface.Database) error {
	err := database.UpdateCard(card)
	if err != nil {
		return fmt.Errorf("UpdateCard error :%w", err)
	}

	return nil
}

// DeleteCard get one card and remove it from database.
func DeleteCard(id string, database database_interface.Database) error {
	if err := database.DeleteCard(id); err != nil {
		return fmt.Errorf("DeleteCard error :%w", err)
	}

	return nil
}
