package main

import (
	"fmt"
	"log"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database/database_sqlite"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/inputhandler"
	"github.com/dixonwille/wmenu"
	"gopkg.in/dixonwille/wlog.v2"
)

const exit = -1

func main() {
	db, err := database_sqlite.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	menu := wmenu.NewMenu("What would you like to do?")

	menu.Action(func(opts []wmenu.Opt) error {
		err := inputhandler.HandleFunc(db, opts)
		if err != nil {
			log.Println(err)
		}

		return nil
	})
	menu.Option("Add a new card", 0, false, nil)
	menu.Option("Review all cards", 1, false, nil)
	menu.Option("exit", exit, false, nil)
	menu.LoopOnInvalid()
	menu.AddColor(wlog.BrightBlue, wlog.BrightYellow, wlog.BrightBlue, wlog.Red)

	for {
		fmt.Println("=========================================================")

		err = menu.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
