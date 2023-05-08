package handler

import (
	"contact-go/helper"
	"contact-go/helper/input"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Menu struct {
	h            ContactHandler
	i            *input.InputReader
	clear        func() error
	showMenuList func()
}

func NewMenu(handler ContactHandler, input *input.InputReader, clear func() error, showMenuList func()) *Menu {
	menu := new(Menu)
	menu.h = handler
	menu.i = input
	menu.clear = clear
	menu.showMenuList = showMenuList

	return menu
}

func (m *Menu) ShowMenu() error {
	err := helper.ClearTerminal()
	if err != nil {
		return err
	}

	m.showMenuList()

	for {
		menuStr, _ := m.i.Scan()
		menu64, err := strconv.ParseInt(strings.TrimSpace(menuStr), 10, 32)
		if err != nil && !errors.Is(err, strconv.ErrSyntax) {
			fmt.Println(err)
			break
		}
		menu := int32(menu64)

		if menu == 6 {
			_ = m.clear
			break
		}

		switch menu {
		default: // case 0 atau selain 0
			_ = m.clear
			m.showMenuList()
		case 1:
			fmt.Println("Contact list")
			m.h.List()
		case 2:
			fmt.Println("Add a new contact")
			m.h.Add()
		case 3:
			fmt.Println("Contact detail")
			m.h.Detail()
		case 4:
			fmt.Println("Update a contact")
			m.h.Update()
		case 5:
			fmt.Println("Delete a contact")
			m.h.Delete()
		}
	}
	return nil
}
