package main

import (
	"fmt"
	"log"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/inputhandler"
	"gopkg.in/dixonwille/wlog.v2"

	"github.com/dixonwille/wmenu"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	menu := wmenu.NewMenu("What would you like to do?")
	menu.Action(func(opts []wmenu.Opt) error { inputhandler.HandleFunc(db, opts); return nil })
	menu.Option("Add a new card", 0, false, nil)
	menu.Option("Review all cards", 1, false, nil)
	menu.LoopOnInvalid()
	menu.AddColor(wlog.BrightBlue, wlog.BrightYellow, wlog.BrightBlue, wlog.Red)
	for {
		fmt.Println("============================")
		err = menu.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
