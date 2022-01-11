package service_test

import (
	"fmt"
	"testing"

	database_interface "github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/database/mockeddatabase"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddCard(t *testing.T) {
	t.Parallel()

	db := mockeddatabase.NewMockedDatabase()

	db.On("AddNewCard", mock.AnythingOfType("models.Card")).Return(nil)
	_, err := service.AddCard([]byte("hi"), db)

	assert.NoError(t, err)
}

func TestReview(t *testing.T) {
	t.Parallel()

	db := mockeddatabase.NewMockedDatabase()
	expected := make([]models.Card, 0)
	card1 := models.NewCard()

	card1.Box = 1
	expected = append(expected, *card1)
	db.On("GetCards").Return(expected, nil)

	cards, err := service.Review(db)

	assert.NoError(t, err)
	assert.Equal(t, cards, expected)
}

func TestConfirmTheCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockeddatabase.NewMockedDatabase()
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
		db := mockeddatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()

		card.ID = id
		card.Box = 1
		db.On("FindByID", card.ID).Return(card, database_interface.ErrNothingFound)
		card2 := models.NewCard()

		card2.ID = id
		card2.Box = 2
		db.On("UpdateCard", *card2).Return(nil)
		err := service.ConfirmTheCard(card.ID, db)

		assert.EqualError(t, err, fmt.Errorf("ConfirmTheCard error :%w", database_interface.ErrNothingFound).Error())
	})
}

func TestRejectTheCard(t *testing.T) {
	t.Parallel()
	t.Run("Successful", func(t *testing.T) {
		t.Parallel()
		db := mockeddatabase.NewMockedDatabase()
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
		db := mockeddatabase.NewMockedDatabase()
		id := "123"
		card := models.NewCard()

		card.ID = id
		card.Box = 2
		db.On("FindByID", card.ID).Return(card, database_interface.ErrNothingFound)
		card2 := models.NewCard()

		card2.ID = id
		card2.Box = 1
		db.On("UpdateCard", *card2).Return(nil)
		err := service.RejectTheCard(card.ID, db)

		assert.EqualError(t, err, fmt.Errorf("RejectTheCard error :%w", database_interface.ErrNothingFound).Error())
	})
}

func TestUpdateCard(t *testing.T) {
	t.Parallel()

	db := mockeddatabase.NewMockedDatabase()
	card := models.NewCard()
	card.Box = 1
	db.On("UpdateCard", *card).Return(nil)
	err := service.UpdateCard(*card, db)

	assert.NoError(t, err)
}

func TestDeleteCard(t *testing.T) {
	t.Parallel()

	db := mockeddatabase.NewMockedDatabase()
	id := "123"

	db.On("DeleteCard", id).Return(nil)

	err := service.DeleteCard(id, db)

	assert.NoError(t, err)
}
