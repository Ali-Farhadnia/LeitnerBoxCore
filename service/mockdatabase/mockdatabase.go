package mockdatabase

import (
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

//nolint:exhaustivestruct
func NewMockedDatabase() *MockDatabase {
	return &MockDatabase{}
}

func (m *MockDatabase) AddNewCard(card models.Card) error {
	arg := m.Called(card)

	return arg.Error(0)
}

func (m *MockDatabase) GetCards() ([]models.Card, error) {
	args := m.Called()

	if v := args.Get(0); v != nil {
		return v.([]models.Card), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockDatabase) FindByID(id string) (*models.Card, error) {
	args := m.Called(id)

	if v := args.Get(0); v != nil {
		return v.(*models.Card), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockDatabase) UpdateCard(card models.Card) error {
	args := m.Called(card)

	return args.Error(0)
}

func (m *MockDatabase) DeleteCard(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
