package inputhandler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/Ali-Farhadnia/LeitnerBoxCore/interfaces"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/models"
	"github.com/Ali-Farhadnia/LeitnerBoxCore/service"
	"github.com/dixonwille/wmenu"
	"gopkg.in/dixonwille/wlog.v2"
)

// Colors.
// nolint:gochecknoglobals
var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// init check if programm run on windows set color value to "".
// nolint:gochecknoinits
func init() {
	if runtime.GOOS == "windows" {
		ColorReset = ""
		ColorRed = ""
		ColorGreen = ""
		ColorYellow = ""
		ColorBlue = ""
		ColorPurple = ""
		ColorCyan = ""
		ColorWhite = ""
	}
}

const (
	// add options.
	setData = 0
	add     = 1
	// review options.
	yes        = 0
	no         = 1
	edit       = 2
	deletecard = 3
	next       = 4
	// edit options.
	editData = 0
	editBox  = 1
	save     = 2
	// comon options.
	cancel = -2
)

// errNotComplete is error that Used to announce that the program is not finished.
var errNotComplete = errors.New("not complete yet")

// HandleFunc main input handler.
func HandleFunc(db interfaces.Database, opts []wmenu.Opt) error {
	// review menu
	reviewmenu := wmenu.NewMenu("Select one:")

	reviewmenu.LoopOnInvalid()
	reviewmenu.AddColor(wlog.Cyan, wlog.BrightYellow, wlog.BrightCyan, wlog.BrightRed)
	reviewmenu.Option("I remember this card", yes, false, nil)
	reviewmenu.Option("I do not remember this card", no, false, nil)
	reviewmenu.Option("edit this card", edit, false, nil)
	reviewmenu.Option("delete this card", deletecard, false, nil)
	reviewmenu.Option("next card", next, false, nil)
	reviewmenu.Option("cancel", cancel, false, nil)

	// add card menu
	addmenu := wmenu.NewMenu("Please select:")

	addmenu.Option("set data", setData, false, nil)
	addmenu.Option("add", add, false, nil)
	addmenu.Option("cancel", cancel, false, nil)
	addmenu.LoopOnInvalid()
	addmenu.AddColor(wlog.Cyan, wlog.BrightYellow, wlog.BrightCyan, wlog.BrightRed)

	// handle input
	switch opts[0].Value {
	case 0:
		err := CoHandleAdd(db, *addmenu)
		if err != nil {
			return err
		}
	case 1:
		err := CoHandleReview(db, *reviewmenu)
		if err != nil {
			return err
		}
	case -1:
		fmt.Println("Goodby")
		os.Exit(0)
	}

	return nil
}

// CoHandleReview run HandleReview in proper way.
func CoHandleReview(db interfaces.Database, reviewmenu wmenu.Menu) error {
	fmt.Println("--------------------------------")
	defer fmt.Println("--------------------------------")

	carts, err := service.Review(db)
	if err != nil {
		fmt.Println(ColorRed + err.Error() + ColorReset)

		return err
	}

	if len(carts) == 0 {
		fmt.Println(ColorRed + "nothings to review" + ColorReset)

		return errors.New("nothings to review")
	}

	for _, cardd := range carts {
		card := cardd

		fmt.Println("in review loop")
		fmt.Println(card)
		PrintCard(card, true, true, true, true)
		reviewmenu.Action(func(opts []wmenu.Opt) error {
			err := HandleReview(opts[0], db, &card)
			if err != nil {
				fmt.Println(ColorRed + err.Error() + ColorReset)

				return err
			}

			return nil
		})

		err := reviewmenu.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleReview handle review menu.
func HandleReview(opt wmenu.Opt, db interfaces.Database, card *models.Card) error {
	// edit card menu
	editmenu := wmenu.NewMenu("select one:")

	editmenu.LoopOnInvalid()
	editmenu.AddColor(wlog.Cyan, wlog.BrightYellow, wlog.BrightCyan, wlog.BrightRed)
	editmenu.Option("edit data", editData, false, nil)
	editmenu.Option("edit box", editBox, false, nil)
	editmenu.Option("save", save, false, nil)
	editmenu.Option("cancel", cancel, false, nil)

	switch opt.Value {
	case yes:
		if err := service.ConfirmTheCard(card.ID, db); err != nil {
			return err
		}

		fmt.Println(ColorGreen + "Successful" + ColorReset)

		return nil
	case no:
		if err := service.RejectTheCard(card.ID, db); err != nil {
			return err
		}

		fmt.Println(ColorGreen + "Successful" + ColorReset)

		return nil
	case edit:
		err := CoHandleEdit(db, *editmenu, card)
		if err != nil {
			return err
		}
	case deletecard:
		err := DeleteCard(card.ID, db)
		if err != nil {
			return err
		}

		fmt.Println(ColorGreen + "Successful" + ColorReset)

		return nil
	case next:
		return nil
	case cancel:
		return nil
	}

	return nil
}

// CoHandleEdit run HandleEdit in proper way.
func CoHandleEdit(db interfaces.Database, editmenu wmenu.Menu, card *models.Card) error {
	fmt.Println("--------------------------------")
	defer fmt.Println("--------------------------------")

	flag := true

	for flag {
		editmenu.Action(func(opts []wmenu.Opt) error {
			err := HandleEdit(opts[0], db, card)
			if errors.Is(err, errNotComplete) {
				flag = true

				return nil
			} else if err != nil {
				fmt.Println(ColorRed + err.Error() + ColorReset)
				flag = false

				return err
			}
			flag = false

			return nil
		})
		PrintCard(*card, false, true, true, false)

		err := editmenu.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleEdit handle edit menu.
func HandleEdit(opt wmenu.Opt, db interfaces.Database, card *models.Card) error {
	reader := bufio.NewReader(os.Stdin)

	switch opt.Value {
	case editData:
		var data string

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
				fmt.Println(ColorRed + "unvalid input" + ColorReset)

				continue
			} else {
				break
			}
		}
		fmt.Println(ColorGreen + "Successful" + ColorReset)

		card.Data = []byte(data)

		return errNotComplete

	case editBox:
		var sbox string

		var ibox int

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
				fmt.Println(ColorRed + "unvalid input" + ColorReset)

				continue
			} else {
				ibox = i

				break
			}
		}
		fmt.Println(ColorGreen + "Successful" + ColorReset)

		card.Box = ibox

		return errNotComplete
	case save:
		err := service.UpdateCard(*card, db)
		if err != nil {
			return err
		}

		fmt.Println(ColorGreen + "Successful" + ColorReset)

		return nil
	case cancel:
		return nil
	}

	return nil
}

// CoHandleAdd run HandleAdd in proper way.
func CoHandleAdd(db interfaces.Database, addmenu wmenu.Menu) error {
	fmt.Println("--------------------------------")
	defer fmt.Println("--------------------------------")

	newcard := models.NewCard()
	flag := true

	for flag {
		addmenu.Action(func(opts []wmenu.Opt) error {
			err := HandleAdd(opts[0], db, newcard)
			if errors.Is(err, errNotComplete) {
				flag = true

				return nil
			} else if err != nil {
				fmt.Println(ColorRed + err.Error() + ColorReset)

				flag = false

				return err
			}

			flag = false

			return nil
		})
		PrintCard(*newcard, false, false, true, false)

		err := addmenu.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleAdd handle add menu.
func HandleAdd(opt wmenu.Opt, db interfaces.Database, card *models.Card) error {
	reader := bufio.NewReader(os.Stdin)

	switch opt.Value {
	case setData:
		var data string

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
				fmt.Println(ColorRed + "unvalid input" + ColorReset)

				continue
			} else {
				break
			}
		}

		fmt.Println(ColorGreen + "Successful" + ColorReset)

		card.Data = []byte(data)

		return errNotComplete
	case add:
		card, err := service.AddCard(card.Data, db)
		if err != nil {
			return err
		}

		PrintCard(*card, true, true, true, true)
		fmt.Println(ColorGreen + "Successful" + ColorReset)

		return nil
	case cancel:
		return nil
	}

	return nil
}

// DeleteCard handle delete option.
func DeleteCard(cardid string, db interfaces.Database) error {
	reader := bufio.NewReader(os.Stdin)

	var input string

	var err error

myloop:
	for {
		fmt.Print("Are you sure(yes or no)?")

		input, err = reader.ReadString('\n')
		if err != nil {
			return err
		}

		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		switch input {
		case "yes":
			break myloop
		case "no":
			return nil
		default:
			fmt.Println(ColorRed + "unvalid input" + ColorReset)
		}
	}

	err = service.DeleteCard(cardid, db)
	if err != nil {
		return err
	}

	return nil
}

// PrintCard print card to stdout.
func PrintCard(card models.Card, showID, showBox, showData, showCreateTime bool) {
	fmt.Println(ColorPurple + "**************** card ****************")

	if showID {
		fmt.Println("ID:", card.ID)
	}

	if showBox {
		fmt.Println("Box:", card.Box)
	}

	if showData {
		fmt.Println("Data:", string(card.Data))
	}

	if showCreateTime {
		fmt.Println("Create time:", card.CreateTime)
	}

	fmt.Println("\n**************************************", ColorReset)
}
