package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	_ "github.com/mattn/go-sqlite3"
)

/* Used to create a singleton object of sql.DB client.
Initialized and exposed through GetDB().*/
var clientInstance *sql.DB

//Used during creation of singleton client object in GetDB().
var clientInstanceError error

//Used to execute client creation procedure only once.
var DbOnce sync.Once

func getdb() (*sql.DB, error) {
	DbOnce.Do(func() {
		res, err := sql.Open("sqlite3", "./database/mydb.db")
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = res

	})

	return clientInstance, clientInstanceError
}

type DB struct {
	client *sql.DB
}

func NewDB() (*DB, error) {
	db, err := getdb()
	if err != nil {
		return nil, err
	}
	rdb := DB{client: db}
	err = rdb.createbooktable()
	if err != nil {
		return nil, err
	}
	return &rdb, nil
}

func (db *DB) checkconnection() error {
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

func (db *DB) createbooktable() error {
	if err := db.checkconnection(); err != nil {
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

func (db *DB) AddNewCard(card models.Card) error {
	if err := db.checkconnection(); err != nil {
		return err
	}
	sqlStatement := `INSERT INTO card (id, box, data, createtime) VALUES ($1,$2,$3,$4)`
	res, err := db.client.Exec(sqlStatement, card.ID, card.Box, card.Data, card.CreateTime)
	if err != nil {
		return err
	}
	/*
		fmt.Println("---------------------------")
		fmt.Println("id:", card.ID)
		fmt.Println("Box:", card.Box)
		fmt.Println("Data:", card.Data)
		fmt.Println("CreateTime:", card.CreateTime)
		fmt.Println("---------------------------")
	*/
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("nothing updated")
	}

	return nil
}

func (db *DB) GetCards() ([]models.Card, error) {
	if err := db.checkconnection(); err != nil {
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
		card := models.Card{CreateTime: &time.Time{}}
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
		cards = append(cards, card)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (db *DB) FindById(id string) (models.Card, error) {
	if err := db.checkconnection(); err != nil {
		return models.Card{}, err
	}
	var card models.Card
	sqlStatement := `SELECT * FROM card WHERE id=$1;`
	row := db.client.QueryRow(sqlStatement, id)
	err := row.Err()
	if err != nil {
		return models.Card{}, err
	}
	var t string

	err = row.Scan(&card.ID, &card.Box, &card.Data, &t)
	if err != nil {
		return models.Card{}, err
	}
	ti, err := time.Parse("2006-01-02 03:04:05+06:00", t)
	if err != nil {
		return models.Card{}, err
	}
	card.CreateTime = &ti
	return card, nil
}

func (db *DB) UpdateCard(card models.Card) error {
	//fmt.Println("--------in update")
	//defer fmt.Println("--------in update")
	if err := db.checkconnection(); err != nil {
		return err
	}
	stmt, err := db.client.Prepare("UPDATE card SET box = ?, data = ?, createtime = ? where id = ?")
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
func (db *DB) DeleteCard(id string) error {
	if err := db.checkconnection(); err != nil {
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
