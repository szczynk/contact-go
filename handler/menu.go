package handler

import (
	"contact-go/helper"
	"contact-go/helper/input"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Menu(handler ContactHandler, input *input.InputReader) {
	err := helper.ClearTerminal()
	if err != nil {
		fmt.Println(err)
		return
	}

	helper.ShowMenuList()

	for {
		menuStr, _ := input.Scan()
		menu64, err := strconv.ParseInt(strings.TrimSpace(menuStr), 10, 32)
		if err != nil && !errors.Is(err, strconv.ErrSyntax) {
			fmt.Println(err)
			break
		}
		menu := int32(menu64)

		if menu == 6 {
			_ = helper.ClearTerminal()
			break
		}

		switch menu {
		default: // case 0 atau selain 0
			_ = helper.ClearTerminal()
			helper.ShowMenuList()
		case 1:
			handler.List()
		case 2:
			handler.Add()
		case 3:
			handler.Detail()
		case 4:
			handler.Update()
		case 5:
			handler.Delete()
		}
	}
}
