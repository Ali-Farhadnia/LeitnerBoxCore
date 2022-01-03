package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do(add or review):")
		input, _ := reader.ReadString('\n')
		switch input {
		case "add\r\n":
			fmt.Print("data:")
			data, _ := reader.ReadString('\n')
			err := service.AddCart([]byte(data), db)
			if err != nil {
				log.Fatal(err)
			}
		case "review\r\n":
			carts, err := service.Review(db)
			if err != nil {
				log.Fatal(err)
			}
			for _, cart := range carts {
				for {
					fmt.Println(string(cart.Data))
					fmt.Print("Are you remember(yes or no):")
					yn, _ := reader.ReadString('\n')
					if yn == "yes\r\n" {
						service.ConfirmTheCard(cart.Id, db)
						break
					} else if yn == "no\r\n" {
						service.RejectTheCard(cart.Id, db)
						break
					}

				}

			}

		}
	}
}
