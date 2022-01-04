package inputhandler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/dixonwille/wmenu"
)

func HandleFunc(db *database.DB, ops []wmenu.Opt) {

	reader := bufio.NewReader(os.Stdin)
	switch ops[0].Value {
	case 0:
		fmt.Print("data:")
		data, _ := reader.ReadString('\n')
		data = strings.TrimSuffix(data, "\n")
		data = strings.TrimSuffix(data, "\r")
		err := service.AddCard([]byte(data), db)
		if err != nil {
			log.Fatal(err)
		}
	case 1:
		carts, err := service.Review(db)
		if err != nil {
			log.Fatal(err)
		}
		if len(carts) == 0 {
			fmt.Println("nothings to review!!")
			fmt.Print("-------------------------\n\n")
			break
		}
		for _, cart := range carts {
			for {
				fmt.Println(string(cart.Data))
				fmt.Print("Are you remember(yes or no):")
				yn, _ := reader.ReadString('\n')
				yn = strings.TrimSuffix(yn, "\n")
				yn = strings.TrimSuffix(yn, "\r")
				if yn == "yes" {
					service.ConfirmTheCard(cart.ID, db)
					break
				} else if yn == "no" {
					service.RejectTheCard(cart.ID, db)
					break
				}

			}

		}

	}
}
