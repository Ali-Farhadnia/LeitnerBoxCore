package sqlitedatabase_test

import (
	"testing"
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database/sqlitedatabase"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateCardTable(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlStatement := "CREATE TABLE card"

	mock.ExpectExec(sqlStatement).WillReturnResult(sqlmock.NewResult(1, 1))

	database := sqlitedatabase.DB{Client: db}
	err = database.CreateCardTable()

	assert.NoError(t, err)
}

func TestAddNewCard(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlStatement := "INSERT INTO card"

	mock.ExpectExec(sqlStatement).WillReturnResult(sqlmock.NewResult(1, 1))

	database := sqlitedatabase.DB{Client: db}
	err = database.AddNewCard(*models.NewCard())

	assert.NoError(t, err)
}

func TestGetCards(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlStatement := "SELECT [*] FROM card"
	time := time.Now()
	rows := []string{"id", "box", "data", "createtime"}
	data := []byte("ali")
	mock.ExpectQuery(sqlStatement).WillReturnRows(sqlmock.NewRows(rows).AddRow("123", "1", data, time))

	database := sqlitedatabase.DB{Client: db}
	cards, err := database.GetCards()

	assert.NoError(t, err)
	assert.Equal(t, cards[0], models.Card{ID: "123", Box: 1, Data: data, CreateTime: &time})
}

func TestFindByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlStatement := "SELECT [*] FROM card"
	time := time.Now()
	rows := []string{"id", "box", "data", "createtime"}
	data := []byte("ali")
	mock.ExpectQuery(sqlStatement).WillReturnRows(sqlmock.NewRows(rows).AddRow("123", "1", data, time))

	database := sqlitedatabase.DB{Client: db}
	card, err := database.FindByID("123")

	assert.NoError(t, err)
	assert.Equal(t, *card, models.Card{ID: "123", Box: 1, Data: data, CreateTime: &time})
}

func TestUpdateCard(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlStatement := "UPDATE card"
	time := time.Now()
	data := []byte("ali")

	mock.ExpectPrepare(sqlStatement).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	card := models.Card{ID: "123", Box: 1, Data: data, CreateTime: &time}
	database := sqlitedatabase.DB{Client: db}
	err = database.UpdateCard(card)

	assert.NoError(t, err)
}

func TestDeleteCard(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlStatement := "DELETE FROM card"

	mock.ExpectPrepare(sqlStatement).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	database := sqlitedatabase.DB{Client: db}
	err = database.DeleteCard("123")

	assert.NoError(t, err)
}
