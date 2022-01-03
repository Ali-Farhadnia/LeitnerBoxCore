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
		id VARCHAR ( 50 ) NOT NULL UNIQUE PRIMARY KEY,
		box INT ( 50 ) NOT NULL,
		data VARBINARY(500)  NOT NULL,
		createtime VARCHAR (50) NOT NULL
	);
	`
	res, err := db.client.Exec(sqlStatement)
	if err != nil {
		switch err.Error() {
		case e:
			log.Println("books card created")
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
			log.Println("books card created")
			return nil
		}
	}
	return nil
}

func (db *DB) AddNewCart(card models.Cart) error {
	if err := db.checkconnection(); err != nil {
		fmt.Println("---------------in check")
		return err
	}
	fmt.Println("---------------1------------")
	sqlStatement := `INSERT INTO card (id, box, data, createtime) VALUES ($1,$2,$3,$4)`
	res, err := db.client.Exec(sqlStatement, card.Id, card.Box, card.Data, card.CreateTime)
	if err != nil {
		return err
	}
	//stmt, _ := db.client.Prepare("INSERT INTO card (id, box, data, createtime) VALUES (?, ?, ?, ?)")
	fmt.Println("---------------2------------")
	fmt.Println("id:", card.Id)
	fmt.Println("Box:", card.Box)
	fmt.Println("Data:", card.Data)
	fmt.Println("CreateTime:", card.CreateTime)
	//stmt.Exec(card.Id, card.Box, card.Data, card.CreateTime)
	fmt.Println("---------------3------------")
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("nothing updated")
	}
	//defer stmt.Close()
	fmt.Println("---------------4------------")

	return nil
	//fmt.Printf("Added %v %v \n", newPerson.first_name, newPerson.last_name
}

func (db *DB) GetCarts() ([]*models.Cart, error) {
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
	cards := make([]*models.Cart, 0)
	for rows.Next() {
		card := models.Cart{}
		var t string
		err = rows.Scan(&card.Id, &card.Box, &card.Data, &t)
		if err != nil {
			return nil, err
		}
		time.Now()
		ti, err := time.Parse("2006-01-02 03:04:05.999999999Z07:00", t)
		if err != nil {
			return nil, err
		}
		card.CreateTime = ti
		if err != nil {
			return nil, err
		}
		cards = append(cards, &card)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return cards, nil
}
func (db *DB) FindById(id string) (models.Cart, error) {
	if err := db.checkconnection(); err != nil {
		return models.Cart{}, err
	}
	var card models.Cart
	sqlStatement := `SELECT * FROM card WHERE id=$1;`
	row := db.client.QueryRow(sqlStatement, id)
	err := row.Err()
	if err != nil {
		return models.Cart{}, err
	}
	err = row.Scan(&card.Id, &card.Box, &card.Data, &card.CreateTime)
	if err != nil {
		return models.Cart{}, err
	}
	return card, nil
}
func (db *DB) UpdateCart(card models.Cart) error {
	if err := db.checkconnection(); err != nil {
		return err
	}
	sqlStatement := `
	UPDATE card
	SET box = $2, data = $3,createtime=$4
	WHERE id = $1;`
	res, err := db.client.Exec(sqlStatement, card.Id, card.Box, card.Data, card.CreateTime)
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
