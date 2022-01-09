package service_test

import (
	"testing"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service/memorydatabase"
	"github.com/stretchr/testify/assert"
)

func TestAddCard(t *testing.T) {
	t.Parallel()
	db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		_, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		_, err := service.AddCard([]byte("error"), &db)
		assert.EqualError(t, err, "add card error")
	})
}
func TestReview(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		_, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		cards, err := service.Review(&db)
		assert.NoError(t, err)
		assert.Equal(t, cards, db.Cards)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		_, err := service.AddCard([]byte("getCards error"), &db)
		assert.NoError(t, err)
		_, err = service.Review(&db)
		assert.EqualError(t, err, "getCards error")
	})
}
func TestConfirmTheCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card1, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		err = service.ConfirmTheCard(card1.ID, &db)
		card2, _ := db.FindByID(card1.ID)
		assert.NoError(t, err)
		assert.Equal(t, card1.Box+1, card2.Box)
	})

	t.Run("Wrong id", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		_, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		err = service.ConfirmTheCard("123", &db)
		assert.EqualError(t, err, "nothings fount")
	})
	t.Run("Updatecard error", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card, err := service.AddCard([]byte("updatecard error"), &db)
		assert.NoError(t, err)
		err = service.ConfirmTheCard(card.ID, &db)
		assert.EqualError(t, err, "updatecard error")
	})
}
func TestRejectTheCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card1, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		err = service.RejectTheCard(card1.ID, &db)
		card2, _ := db.FindByID(card1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, card2.Box)
	})

	t.Run("Wrong id", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		_, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		err = service.RejectTheCard("123", &db)
		assert.EqualError(t, err, "nothings fount")
	})
	t.Run("Updatecard error", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card, err := service.AddCard([]byte("updatecard error"), &db)
		assert.NoError(t, err)
		err = service.RejectTheCard(card.ID, &db)
		assert.EqualError(t, err, "updatecard error")
	})
}
func TestUpdateCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card1, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		card1.Data = []byte("hello")
		err = service.UpdateCard(card1, &db)
		assert.NoError(t, err)
		card2, err := db.FindByID(card1.ID)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello"), card2.Data)

	})
	t.Run("Updatecard error", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card, err := service.AddCard([]byte("updatecard error"), &db)
		assert.NoError(t, err)
		err = service.UpdateCard(card, &db)
		assert.EqualError(t, err, "updatecard error")
	})
}
func TestDeleteCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		card1, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		err = service.DeleteCard(card1.ID, &db)
		assert.NoError(t, err)
		_, err = db.FindByID(card1.ID)
		assert.EqualError(t, err, "nothings fount")
	})

	t.Run("deletecard error", func(t *testing.T) {
		t.Parallel()
		db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
		_, err := service.AddCard([]byte("hi"), &db)
		assert.NoError(t, err)
		err = service.DeleteCard("deletecard error", &db)
		assert.EqualError(t, err, "deletecard error")
	})
}
