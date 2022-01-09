package models_test

import (
	"testing"
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	card := models.NewCard()
	assert.Equal(t, "", card.ID)
	assert.Equal(t, []byte(nil), card.Data)
	assert.Equal(t, 0, card.Box)
	assert.Equal(t, &time.Time{}, card.CreateTime)
}
