package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

/* Used to create a singleton object of sql.DB client.
Initialized and exposed through GetDB().*/
var clientInstance *sql.DB

//Used during creation of singleton client object in GetDB().
var errClientInstance error

//DBOnce used to execute client creation procedure only once.
var DBOnce sync.Once

// getdb create sqlite database client from ./database/mydb.db.
func getdb() (*sql.DB, error) {
	DBOnce.Do(func() {
		res, err := sql.Open("sqlite3", "./database/mydb.db")
		if err != nil {
			errClientInstance = err
		}
		clientInstance = res

	})

	return clientInstance, errClientInstance
}

// DB is database type that implement database interface.
type DB struct {
	client *sql.DB
}

// NewDB create new DB struct.
func NewDB() (*DB, error) {
	db, err := getdb()
	if err != nil {
		return nil, err
	}
	rdb := DB{client: db}
	err = rdb.createCardTable()
	if err != nil {
		return nil, err
	}
	return &rdb, nil
}

// checkConnection ping database client and if not response trys to build new one.
func (db *DB) checkConnection() error {
	err := db.client.Ping()
	if err != nil {
		res, err := getdb()
		if err != nil {
			return err
		}
		db.client = res
	}
	return nil
}

// createCardTable check if card table exist or not and if not creates one.
func (db *DB) createCardTable() error {
	if err := db.checkConnection(); err != nil {
		return err
	}

	e := fmt.Sprintf("table %s already exists", "card")
	sqlStatement := `
	CREATE TABLE card (
		id TEXT NOT NULL UNIQUE PRIMARY KEY,
		box INT NOT NULL,
		data TEXT NOT NULL,
		createtime datetime NOT NULL
	);
	`
	res, err := db.client.Exec(sqlStatement)
	if err != nil {
		switch err.Error() {
		case e:
			//log.Println("books card created")
			return nil
		case "":
			log.Println("failed to create card table")
			return err
		default:
			if res == nil {
				return err
			}
			e, err := res.RowsAffected()

			if err != nil {

				return err
			}
			if e == 0 {
				return errors.New("somthing went wrong in create table")
			}
			//log.Println("books card created")
			return nil
		}
	}
	return nil
}

// AddNewCard get card and add it to database.
func (db *DB) AddNewCard(card models.Card) error {
	if err := db.checkConnection(); err != nil {
		return err
	}
	sqlStatement := `INSERT INTO card (id, box, data, createtime) VALUES ($1,$2,$3,$4)`
	res, err := db.client.Exec(sqlStatement, card.ID, card.Box, card.Data, card.CreateTime)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("nothing updated")
	}

	return nil
}

// GetCards query all cards from database and return them.
func (db *DB) GetCards() ([]models.Card, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	sqlStatement := `SELECT * FROM card;`
	rows, err := db.client.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	cards := make([]models.Card, 0)
	for rows.Next() {
		card := models.NewCard()
		//var t time.Time
		err = rows.Scan(&card.ID, &card.Box, &card.Data, card.CreateTime)
		if err != nil {
			return nil, err
		}

		/*
			ti, err := time.Parse("2006-01-02 15:04:05 -0700 MST", t)
			if err != nil {
				return nil, err
			}
			card.CreateTime = ti.UTC().Local()
		*/
		cards = append(cards, *card)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// FindByID get card d and find it in the database and return it.
func (db *DB) FindByID(id string) (*models.Card, error) {
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	var card models.Card
	sqlStatement := `SELECT * FROM card WHERE id=$1;`
	row := db.client.QueryRow(sqlStatement, id)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	err = row.Scan(&card.ID, &card.Box, &card.Data, &card.CreateTime)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

// UpdateCard get card and update it in database.
func (db *DB) UpdateCard(card models.Card) error {
	//fmt.Println("--------in update")
	//defer fmt.Println("--------in update")
	if err := db.checkConnection(); err != nil {
		return err
	}
	stmt, err := db.client.Prepare("UPDATE card set box = ?, data = ?, createtime = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(card.Box, card.Data, card.CreateTime, card.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("nothing updated")
	}
	return nil
}

// DeleteCard get card id and remove it from database.
func (db *DB) DeleteCard(id string) error {
	if err := db.checkConnection(); err != nil {
		return err
	}
	stmt, err := db.client.Prepare("DELETE FROM card where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("nothing deleted")
	}
	return nil
}
