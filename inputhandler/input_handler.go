package inputhandler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/database"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/dixonwille/wmenu"
	"gopkg.in/dixonwille/wlog.v2"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)
const (
	//add options
	set_data = 0
	add      = 1
	//review options
	yes    = 0
	no     = 1
	edit   = 2
	delete = 3
	next   = 4
	//edit options
	edit_data = 0
	edit_box  = 1
	save      = 2
	//comon options
	cancel = -2
)

var err_not_complete = errors.New("not complete yet")

func HandleFunc(db *database.DB, opts []wmenu.Opt) error {
	//review menu
	reviewmenu := wmenu.NewMenu("Select one:")
	reviewmenu.LoopOnInvalid()
	reviewmenu.AddColor(wlog.Cyan, wlog.BrightYellow, wlog.BrightCyan, wlog.BrightRed)
	reviewmenu.Option("I remember this card", yes, false, nil)
	reviewmenu.Option("I do not remember this card", no, false, nil)
	reviewmenu.Option("edit this card", edit, false, nil)
	reviewmenu.Option("delete this card", delete, false, nil)
	reviewmenu.Option("next card", next, false, nil)
	reviewmenu.Option("cancel", cancel, false, nil)

	//add card menu
	addmenu := wmenu.NewMenu("Please select:")
	addmenu.Option("set data", set_data, false, nil)
	addmenu.Option("add", add, false, nil)
	addmenu.Option("cancel", cancel, false, nil)
	addmenu.LoopOnInvalid()
	addmenu.AddColor(wlog.Cyan, wlog.BrightYellow, wlog.BrightCyan, wlog.BrightRed)

	//handle input
	switch opts[0].Value {
	case 0:
		fmt.Println("--------------------------------")
		defer fmt.Println("--------------------------------")
		newcard := models.Cart{}
		flag := true
		for flag {
			addmenu.Action(func(opts []wmenu.Opt) error {
				err := HandleAdd(opts[0], db, &newcard)
				if err == err_not_complete {
					flag = true
					return nil
				} else if err != nil {
					fmt.Println(string(ColorRed) + err.Error() + string(ColorReset))
					flag = false
					return err
				}
				flag = false
				return nil
			})
			PrintCard(newcard, false, false, true, false)
			addmenu.Run()
		}
	case 1:
		fmt.Println("--------------------------------")
		defer fmt.Println("--------------------------------")
		carts, err := service.Review(db)
		if err != nil {
			return err
		}
		if len(carts) == 0 {
			return errors.New("nothings to review")
		}
		for _, card := range carts {

			PrintCard(card, true, true, true, true)
			reviewmenu.Action(func(opts []wmenu.Opt) error {
				err := HandleReview(opts[0], db, &card)
				if err != nil {
					fmt.Println(string(ColorRed) + err.Error() + string(ColorReset))
					return err
				}
				return nil
			})

			reviewmenu.Run()

		}
	case -1:
		fmt.Println("Goodby")
		os.Exit(0)
	}
	return nil
}
func HandleReview(opt wmenu.Opt, db *database.DB, card *models.Cart) error {
	reader := bufio.NewReader(os.Stdin)

	//edit card menu
	editmenu := wmenu.NewMenu("select one:")
	editmenu.LoopOnInvalid()
	editmenu.AddColor(wlog.Cyan, wlog.BrightYellow, wlog.BrightCyan, wlog.BrightRed)
	editmenu.Option("edit data", edit_data, false, nil)
	editmenu.Option("edit box", edit_box, false, nil)
	editmenu.Option("save", save, false, nil)
	editmenu.Option("cancel", cancel, false, nil)
	switch opt.Value {
	case yes:
		if err := service.ConfirmTheCard(card.ID, db); err != nil {
			return err
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		return nil
	case no:
		if err := service.RejectTheCard(card.ID, db); err != nil {
			return err
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		return nil
	case edit:
		fmt.Println("--------------------------------")
		defer fmt.Println("--------------------------------")
		flag := true
		for flag {
			editmenu.Action(func(opts []wmenu.Opt) error {
				err := HandleEdit(opts[0], db, card)
				if err == err_not_complete {
					flag = true
					return nil
				} else if err != nil {
					fmt.Println(string(ColorRed) + err.Error() + string(ColorReset))
					flag = false
					return err
				}
				flag = false
				return nil
			})
			PrintCard(*card, false, true, true, false)
			editmenu.Run()
		}
		return nil
	case delete:
		input := ""
		var err error
		for {
			fmt.Print("Are you sure(yes or no)?")
			input, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			input = strings.TrimSuffix(input, "\n")
			input = strings.TrimSuffix(input, "\r")
			if input == "yes" {
				break
			} else if input == "no" {
				return nil
			} else {
				fmt.Println(string(ColorRed) + "unvalid input" + string(ColorReset))
			}
		}
		err = service.DeleteCart(card.ID, db)
		if err != nil {
			return err
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		return nil
	case next:
		return errors.New("not supported yet")
	case cancel:
		return nil
	}
	return nil
}
func HandleEdit(opt wmenu.Opt, db *database.DB, card *models.Cart) error {
	reader := bufio.NewReader(os.Stdin)

	switch opt.Value {
	case edit_data:
		data := ""
		var err error
		for {
			fmt.Print("Data :")
			data, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			data = strings.TrimSuffix(data, "\n")
			data = strings.TrimSuffix(data, "\r")
			if data == "" {
				fmt.Println(string(ColorRed) + "unvalid input" + string(ColorReset))
				continue
			} else {
				break
			}
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		card.Data = []byte(data)
		return err_not_complete

	case edit_box:
		sbox := ""
		ibox := 0
		var err error
		for {
			fmt.Print("Box :")
			sbox, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			sbox = strings.TrimSuffix(sbox, "\n")
			sbox = strings.TrimSuffix(sbox, "\r")
			i, err := strconv.Atoi(sbox)
			if err != nil {
				fmt.Println(string(ColorRed) + "unvalid input" + string(ColorReset))
				continue
			} else {
				ibox = i
				break
			}
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		card.Box = ibox
		return err_not_complete
	case save:
		err := service.UpdateCart(*card, db)
		if err != nil {
			return err
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		return nil
	case cancel:
		return nil
	}
	return nil
}
func HandleAdd(opt wmenu.Opt, db *database.DB, card *models.Cart) error {
	reader := bufio.NewReader(os.Stdin)
	switch opt.Value {
	case set_data:
		data := ""
		var err error
		for {
			fmt.Print("Data :")
			data, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			data = strings.TrimSuffix(data, "\n")
			data = strings.TrimSuffix(data, "\r")
			if data == "" {
				fmt.Println(string(ColorRed) + "unvalid input" + string(ColorReset))
				continue
			} else {
				break
			}
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		card.Data = []byte(data)
		return err_not_complete
	case add:
		err := service.AddCard(card.Data, db) //To Do return card or id
		if err != nil {
			return err
		}
		fmt.Println(string(ColorGreen) + "Successful" + string(ColorReset))
		return nil
	case cancel:
		return nil
	}
	return nil
}
func PrintCard(card models.Cart, show_id, show_box, show_data, show_create_time bool) {
	fmt.Println(string(ColorPurple) + "**************** card ****************")
	if show_id {
		fmt.Println("ID:", card.ID)
	}
	if show_box {
		fmt.Println("Box:", card.Box)
	}
	if show_data {
		fmt.Println("Data:", string(card.Data))
	}
	if show_create_time {
		fmt.Println("Create time:", card.CreateTime)
	}
	fmt.Println("\n**************************************", string(ColorReset))
}
