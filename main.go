package main

import (
	"log"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/inputhandler"

	"github.com/dixonwille/wmenu"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	menu := wmenu.NewMenu("What would you like to do?")
	menu.Action(func(opts []wmenu.Opt) error { inputhandler.HandleFunc(db, opts); return nil })
	menu.Option("Add a new card", 0, true, nil)
	menu.Option("Review all cards", 1, true, nil)
	menu.LoopOnInvalid()
	for {
		err = menu.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
