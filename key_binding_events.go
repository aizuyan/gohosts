package main

import (
	"github.com/jroimartin/gocui"
)

// 绑定事件的自定义结构体
type KeyBindingItem struct {
	ViewName string
	Key interface{}
	Mod gocui.Modifier
	Handler func(g *gocui.Gui, v *gocui.View) error
}

func KeyBindingAction(g *gocui.Gui) error {
	bindingItems := []KeyBindingItem{
		{
			// 关闭软件，组合按键，Ctrl + C
			"",
			gocui.KeyCtrlC,
			gocui.ModNone,
			quit,
		},
		{
			// 关闭软件，组合按键，Ctrl + C
			"",
			gocui.KeyArrowUp,
			gocui.ModNone,
			arrowUpAction,
		},
		{
			// 关闭软件，组合按键，Ctrl + C
			"",
			gocui.KeyArrowDown,
			gocui.ModNone,
			arrowDownAction,
		},
	}

	for _, bindingItem := range bindingItems {
		if err := g.SetKeybinding(bindingItem.ViewName, bindingItem.Key, bindingItem.Mod, bindingItem.Handler); err != nil {
			return err
		}
	}

	return nil
}

func quit(g *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func arrowDownAction(g *gocui.Gui, v *gocui.View) error {
	// 到底了
	if slideCursorY + slideOriginY > len(hItems) - 2 {
		return nil
	}

	if slideCursorY < getSlideRowCount() {
		slideCursorY++
	} else {
		slideOriginY++
	}

	return nil
}

func arrowUpAction(g *gocui.Gui, v *gocui.View) error {
	if slideOriginY > 0 {
		slideOriginY--
	} else if slideCursorY > 0 {
		slideCursorY--
	}

	return nil
}