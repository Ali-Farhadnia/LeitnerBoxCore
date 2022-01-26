package sqlitedatabase

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	database_interface "github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	_ "github.com/mattn/go-sqlite3"
)

/* Used to create a singleton object of sql.DB client.
Initialized and exposed through GetDB().*/
// nolint:gochecknoglobals
var clientInstance *sql.DB

// Used during creation of singleton client object in GetDB().
// nolint:gochecknoglobals
var errClientInstance error

// DBOnce used to execute client creation procedure only once.
// nolint:gochecknoglobals
var DBOnce sync.Once

// getdb create sqlite database client from ./database/mydb.db.
func getdb() (*sql.DB, error) {
	DBOnce.Do(func() {
		res, err := sql.Open("sqlite3", "./database/sqlitedatabase/mydb.db")
		if err != nil {
			errClientInstance = err
		}
		clientInstance = res
	})

	return clientInstance, errClientInstance
}

// DB is database type that implement database interface.
type DB struct {
	Client *sql.DB
}

// NewDB create new DB struct.
func NewDB() (*DB, error) {
	db, err := getdb()
	if err != nil {
		return nil, err
	}

	rdb := DB{Client: db}

	if err = rdb.CreateCardTable(); err != nil {
		return nil, fmt.Errorf("newdb error :%w", err)
	}

	return &rdb, nil
}

// CreateCardTable check if card table exist or not and if not creates one.
func (db *DB) CreateCardTable() error {
	e := fmt.Sprintf("table %s already exists", "card")
	sqlStatement := `
	CREATE TABLE card (
		id TEXT NOT NULL UNIQUE PRIMARY KEY,
		box INT NOT NULL,
		data TEXT NOT NULL,
		createtime datetime NOT NULL
	);
	`

	res, err := db.Client.Exec(sqlStatement)
	if err != nil {
		switch err.Error() {
		case e:
			// log.Println("books card created")
			return nil
		case "":
			log.Println("failed to create card table")

			return fmt.Errorf("createCardTable error :%w", err)
		default:
			if res == nil {
				return fmt.Errorf("createCardTable error :%w", err)
			}

			e, err := res.RowsAffected()
			if err != nil {
				return fmt.Errorf("createCardTable error :%w", err)
			}

			if e == 0 {
				return database_interface.ErrSomthingWentWrong
			}

			return nil
		}
	}

	return nil
}

// AddNewCard get card and add it to database.
func (db *DB) AddNewCard(card models.Card) error {
	sqlStatement := `INSERT INTO card (id, box, data, createtime) VALUES ($1,$2,$3,$4)`

	res, err := db.Client.Exec(sqlStatement, card.ID, card.Box, card.Data, card.CreateTime)
	if err != nil {
		return fmt.Errorf("AddNewCard error :%w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("AddNewCard error :%w", err)
	}

	if affected == 0 {
		return database_interface.ErrSomthingWentWrong
	}

	return nil
}

// GetCards query all cards from database and return them.
func (db *DB) GetCards() ([]models.Card, error) {
	sqlStatement := `SELECT * FROM card;`

	rows, err := db.Client.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("GetCards error :%w", err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetCards error :%w", err)
	}

	cards := make([]models.Card, 0)

	for rows.Next() {
		card := models.NewCard()

		err = rows.Scan(&card.ID, &card.Box, &card.Data, card.CreateTime)
		if err != nil {
			return nil, fmt.Errorf("GetCards error :%w", err)
		}

		cards = append(cards, *card)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetCards error :%w", err)
	}

	return cards, nil
}

// FindByID get card d and find it in the database and return it.
func (db *DB) FindByID(id string) (*models.Card, error) {
	var card models.Card

	sqlStatement := `SELECT * FROM card WHERE id=$1;`
	row := db.Client.QueryRow(sqlStatement, id)

	err := row.Err()
	if err != nil {
		return nil, fmt.Errorf("FindByID error :%w", err)
	}

	err = row.Scan(&card.ID, &card.Box, &card.Data, &card.CreateTime)
	if err != nil {
		return nil, fmt.Errorf("FindByID error :%w", err)
	}

	return &card, nil
}

// UpdateCard get card and update it in database.
func (db *DB) UpdateCard(card models.Card) error {
	stmt, err := db.Client.Prepare("UPDATE card set box = ?, data = ?, createtime =? where id = ?")
	if err != nil {
		return fmt.Errorf("UpdateCard error :%w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(card.Box, card.Data, card.CreateTime, card.ID)
	if err != nil {
		return fmt.Errorf("UpdateCard error :%w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateCard error :%w", err)
	}

	if affected == 0 {
		return database_interface.ErrSomthingWentWrong
	}

	return nil
}

// DeleteCard get card id and remove it from database.
func (db *DB) DeleteCard(id string) error {
	stmt, err := db.Client.Prepare("DELETE FROM card where id = ?")
	if err != nil {
		return fmt.Errorf("DeleteCard error :%w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("DeleteCard error :%w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteCard error :%w", err)
	}

	if affected == 0 {
		return database_interface.ErrSomthingWentWrong
	}

	return nil
}
