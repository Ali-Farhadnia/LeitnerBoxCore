package service_test

import (
	"errors"
	"testing"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/stretchr/testify/assert"
)

type FakeDatabase struct {
	cards []models.Card
}

func (db *FakeDatabase) AddNewCard(card models.Card) error {
	if string(card.Data) == "error" {
		return errors.New("add card error")
	}
	db.cards = append(db.cards, card)
	return nil
}
func (db *FakeDatabase) GetCards() ([]models.Card, error) {
	for _, card := range db.cards {
		if string(card.Data) == "getCards error" {
			return nil, errors.New("getCards error")
		}
	}
	return db.cards, nil
}
func (db *FakeDatabase) FindById(id string) (models.Card, error) {
	for _, card := range db.cards {
		if card.ID == id {
			return card, nil
		}
	}
	return models.Card{}, errors.New("nothings fount")
}
func (db *FakeDatabase) UpdateCard(card2 models.Card) error {
	if string(card2.Data) == "updatecard error" {
		return errors.New("updatecard error")
	}
	for i, card1 := range db.cards {
		if card1.ID == card2.ID {
			db.cards[i] = card2
			return nil
		}
	}
	return errors.New("nothing updated")
}
func (db *FakeDatabase) DeleteCard(id string) error {
	if id == "deletecard error" {
		return errors.New("deletecard error")
	}
	for i, card := range db.cards {
		if card.ID == id {
			db.cards = append(db.cards[:i], db.cards[i+1:]...)
			return nil
		}
	}
	return errors.New("nothing deleted")
}
func TestAddCard(t *testing.T) {
	db := FakeDatabase{make([]models.Card, 0)}
	t.Run("Successful", func(t *testing.T) {
		_, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		_, err := service.AddCard([]byte("error"), &db)
		assert.EqualError(t, err, "add card error")
	})

}
func TestReview(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		service.AddCard([]byte("hi"), &db)
		cards, err := service.Review(&db)
		assert.NoError(t, err)
		assert.Equal(t, cards, db.cards)
	})

	t.Run("Error", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		service.AddCard([]byte("getCards error"), &db)
		_, err := service.Review(&db)
		assert.EqualError(t, err, "getCards error")
	})
}

func TestConfirmTheCard(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card1, _ := service.AddCard([]byte("hi"), &db)
		err := service.ConfirmTheCard(card1.ID, &db)
		card2, _ := db.FindById(card1.ID)
		assert.NoError(t, err)
		assert.Equal(t, card1.Box+1, card2.Box)
	})

	t.Run("Wrong id", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		service.AddCard([]byte("hi"), &db)
		err := service.ConfirmTheCard("123", &db)
		assert.EqualError(t, err, "nothings fount")
	})
	t.Run("Updatecard error", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card, _ := service.AddCard([]byte("updatecard error"), &db)
		err := service.ConfirmTheCard(card.ID, &db)
		assert.EqualError(t, err, "updatecard error")
	})

}

func TestRejectTheCard(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card1, _ := service.AddCard([]byte("hi"), &db)
		err := service.RejectTheCard(card1.ID, &db)
		card2, _ := db.FindById(card1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, card2.Box)
	})

	t.Run("Wrong id", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		service.AddCard([]byte("hi"), &db)
		err := service.RejectTheCard("123", &db)
		assert.EqualError(t, err, "nothings fount")
	})
	t.Run("Updatecard error", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card, _ := service.AddCard([]byte("updatecard error"), &db)
		err := service.RejectTheCard(card.ID, &db)
		assert.EqualError(t, err, "updatecard error")
	})

}
func TestUpdateCard(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card1, _ := service.AddCard([]byte("hi"), &db)
		card1.Data = []byte("hello")
		err := service.UpdateCard(card1, &db)
		card2, _ := db.FindById(card1.ID)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello"), card2.Data)

	})
	t.Run("Updatecard error", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card, _ := service.AddCard([]byte("updatecard error"), &db)
		err := service.UpdateCard(card, &db)
		assert.EqualError(t, err, "updatecard error")
	})
}
func TestDeleteCard(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		card1, _ := service.AddCard([]byte("hi"), &db)
		err := service.DeleteCard(card1.ID, &db)
		assert.NoError(t, err)
		_, err = db.FindById(card1.ID)
		assert.EqualError(t, err, "nothings fount")
	})

	t.Run("deletecard error", func(t *testing.T) {
		db := FakeDatabase{make([]models.Card, 0)}
		service.AddCard([]byte("hi"), &db)
		err := service.DeleteCard("deletecard error", &db)
		assert.EqualError(t, err, "deletecard error")
	})

}
