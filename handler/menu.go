package handler

import (
	"contact-go/helper"
	"contact-go/helper/input"
	"fmt"
	"io"
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
		var menu int
		_, err = fmt.Sscan(menuStr, &menu)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				break
			}
		}

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
