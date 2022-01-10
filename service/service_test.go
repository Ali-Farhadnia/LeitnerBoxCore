package service_test

import (
	"errors"
	"testing"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service/mockdatabase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddCard(t *testing.T) {
	t.Parallel()
	//db := memorydatabase.MemoryDatabase{Cards: make([]models.Card, 0)}
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		db.On("AddNewCard", mock.AnythingOfType("models.Card")).Return(nil)
		_, err := service.AddCard([]byte("hi"), db)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		db.On("AddNewCard", mock.AnythingOfType("models.Card")).Return(errors.New("add card error"))

		_, err := service.AddCard([]byte("error"), db)

		assert.EqualError(t, err, "add card error")
	})
}

func TestReview(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		expected := make([]models.Card, 0)
		card1 := models.NewCard()
		card1.Box = 1
		expected = append(expected, *card1)
		db.On("GetCards").Return(expected, nil)

		cards, err := service.Review(db)
		assert.NoError(t, err)
		assert.Equal(t, cards, expected)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		db.On("GetCards").Return(nil, errors.New("getCards error"))

		_, err := service.Review(db)
		assert.EqualError(t, err, "getCards error")
	})
}

func TestConfirmTheCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()
		card.ID = id
		card.Box = 1
		db.On("FindByID", card.ID).Return(card, nil)
		card2 := models.NewCard()
		card2.ID = id
		card2.Box = 2
		db.On("UpdateCard", *card2).Return(nil)
		err := service.ConfirmTheCard(card.ID, db)

		assert.NoError(t, err)
	})

	t.Run("Wrong id", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()
		card.ID = id
		card.Box = 1
		db.On("FindByID", card.ID).Return(card, errors.New("nothings fount"))
		card2 := models.NewCard()
		card2.ID = id
		card2.Box = 2
		db.On("UpdateCard", *card2).Return(nil)
		err := service.ConfirmTheCard(card.ID, db)

		assert.EqualError(t, err, "nothings fount")
	})
	t.Run("Updatecard error", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()
		card.ID = id
		card.Box = 1
		db.On("FindByID", card.ID).Return(card, nil)
		card2 := models.NewCard()
		card2.ID = id
		card2.Box = 2
		db.On("UpdateCard", *card2).Return(errors.New("updatecard error"))
		err := service.ConfirmTheCard(card.ID, db)

		assert.EqualError(t, err, "updatecard error")
	})
}

func TestRejectTheCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()
		card.ID = id
		card.Box = 2
		db.On("FindByID", card.ID).Return(card, nil)
		card2 := models.NewCard()
		card2.ID = id
		card2.Box = 1
		db.On("UpdateCard", *card2).Return(nil)
		err := service.RejectTheCard(card.ID, db)

		assert.NoError(t, err)
	})

	t.Run("Wrong id", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()
		card.ID = id
		card.Box = 2
		db.On("FindByID", card.ID).Return(card, errors.New("nothings fount"))
		card2 := models.NewCard()
		card2.ID = id
		card2.Box = 1
		db.On("UpdateCard", *card2).Return(nil)
		err := service.RejectTheCard(card.ID, db)

		assert.EqualError(t, err, "nothings fount")
	})
	t.Run("Updatecard error", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()
		card.ID = id
		card.Box = 2
		db.On("FindByID", card.ID).Return(card, nil)
		card2 := models.NewCard()
		card2.ID = id
		card2.Box = 1
		db.On("UpdateCard", *card2).Return(errors.New("updatecard error"))
		err := service.RejectTheCard(card.ID, db)

		assert.EqualError(t, err, "updatecard error")
	})
}

func TestUpdateCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		card := models.NewCard()
		card.Box = 1
		db.On("UpdateCard", *card).Return(nil)
		err := service.UpdateCard(*card, db)

		assert.NoError(t, err)
	})
	t.Run("Updatecard error", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		card := models.NewCard()
		card.Box = 1
		db.On("UpdateCard", *card).Return(errors.New("updatecard error"))
		err := service.UpdateCard(*card, db)

		assert.EqualError(t, err, "updatecard error")
	})
}
func TestDeleteCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		db.On("DeleteCard", id).Return(nil)
		err := service.DeleteCard(id, db)

		assert.NoError(t, err)
	})

	t.Run("deletecard error", func(t *testing.T) {
		t.Parallel()
		db := mockdatabase.NewMockedDatabase()
		id := "123"
		db.On("DeleteCard", id).Return(errors.New("deletecard error"))
		err := service.DeleteCard(id, db)

		assert.EqualError(t, err, "deletecard error")
	})
}
