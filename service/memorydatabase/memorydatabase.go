package memorydatabase

import (
	"errors"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
)

type MemoryDatabase struct {
	Cards []models.Card
}

func (db *MemoryDatabase) AddNewCard(card models.Card) error {
	if string(card.Data) == "error" {
		return errors.New("add card error")
	}
	db.Cards = append(db.Cards, card)
	return nil
}
func (db *MemoryDatabase) GetCards() ([]models.Card, error) {
	for _, card := range db.Cards {
		if string(card.Data) == "getCards error" {
			return nil, errors.New("getCards error")
		}
	}
	return db.Cards, nil
}
func (db *MemoryDatabase) FindById(id string) (models.Card, error) {
	for _, card := range db.Cards {
		if card.ID == id {
			return card, nil
		}
	}
	return models.Card{}, errors.New("nothings fount")
}
func (db *MemoryDatabase) UpdateCard(card2 models.Card) error {
	if string(card2.Data) == "updatecard error" {
		return errors.New("updatecard error")
	}
	for i, card1 := range db.Cards {
		if card1.ID == card2.ID {
			db.Cards[i] = card2
			return nil
		}
	}
	return errors.New("nothing updated")
}
func (db *MemoryDatabase) DeleteCard(id string) error {
	if id == "deletecard error" {
		return errors.New("deletecard error")
	}
	for i, card := range db.Cards {
		if card.ID == id {
			db.Cards = append(db.Cards[:i], db.Cards[i+1:]...)
			return nil
		}
	}
	return errors.New("nothing deleted")
}
