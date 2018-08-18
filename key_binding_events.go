package main

import (
	"github.com/jroimartin/gocui"
	"strings"
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
			// 切换view，tab键
			"",
			gocui.KeyTab,
			gocui.ModNone,
			changeTab,
		},
		{
			// 关闭软件，组合按键，Ctrl + C
			"slide",
			gocui.KeyArrowUp,
			gocui.ModNone,
			arrowUpAction,
		},
		{
			// 关闭软件，组合按键，Ctrl + C
			"slide",
			gocui.KeyArrowDown,
			gocui.ModNone,
			arrowDownAction,
		},
		{
			// 关闭hosts item 左箭头
			"slide",
			gocui.KeyArrowLeft,
			gocui.ModNone,
			arrowLeftAction,
		},
		{
			// 打开hosts item 右箭头
			"slide",
			gocui.KeyArrowRight,
			gocui.ModNone,
			arrowRightAction,
		},
		{
			// 添加新的hosts item, 组合键 Shift + a
			"slide",
			'A',
			gocui.ModNone,
			shiftAAction,
		},
		{
			// new-hosts-item-msg view 上面按下回车
			"new-hosts-item-msg",
			gocui.KeyEnter,
			gocui.ModNone,
			newHostsItemMsgEnterAction,
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

func changeTab(g *gocui.Gui, view *gocui.View) error {
	tabViewIndex++
	if tabViewIndex >= len(tabViews) {
		tabViewIndex = 0
	}
	return nil
}


func arrowDownAction(g *gocui.Gui, v *gocui.View) error {
	// 到底了
	if slideCursorY+slideOriginY > len(hItems)-2 {
		return nil
	}

	if slideCursorY < getSlideRowCount() {
		slideCursorY++
	} else {
		slideOriginY++
	}
	hItemCursorChanged = true
	mainContentChanged = true

	return nil
}

func arrowUpAction(g *gocui.Gui, v *gocui.View) error {
	if slideOriginY > 0 {
		slideOriginY--
	} else if slideCursorY > 0 {
		slideCursorY--
	}
	hItemCursorChanged = true
	mainContentChanged = true

	return nil
}

func arrowLeftAction(g *gocui.Gui, v *gocui.View) error {
	hostsItemIdx := getCurrentHostsItemIndex()
	hItems[hostsItemIdx].toggleHostsItemSwitch(false)
	hItemCursorChanged = true

	return nil
}

func arrowRightAction(g *gocui.Gui, v *gocui.View) error {
	hostsItemIdx := getCurrentHostsItemIndex()
	hItems[hostsItemIdx].toggleHostsItemSwitch(true)
	hItemCursorChanged = true

	return nil
}

func shiftAAction(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("new-hosts-item-msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "new hosts item name"
		v.Editable = true
		onMsgView = true
		if _, err := setViewOnTop(g, "new-hosts-item-msg"); err != nil {
			return err
		}
	}

	return nil
}

// 判断输入的名称是否合法
func checkHostsItemName(hostsName string) bool {
	if hostsName == "" {
		return  false
	}

	for _, host := range hItems {
		if host.HostsName == hostsName {
			return false
		}
	}

	return true
}

func newHostsItemMsgEnterAction(g *gocui.Gui, v *gocui.View) error {
	// 输入的名称
	hostsName := v.ViewBuffer()
	hostsName = strings.Trim(hostsName, "\n ")
	if !checkHostsItemName(hostsName) {
		// 提示不合法 TODO

	}
	hItems = append(hItems, hostsItem{
		hostsName,
		"# " + hostsName + " config\n",
		false,
	})


	if err := g.DeleteView("new-hosts-item-msg"); err != nil {
		return err
	}
	onMsgView = false
	jsonencodeHostsInfoToPath(dataPath, hItems)
	// 调整鼠标焦点
	slideCursorY = getSlideRowCount()
	if len(hItems) - 1 < slideCursorY {
		slideCursorY = len(hItems) - 1
	}
	slideOriginY = len(hItems) - getSlideRowCount() - 1
	hItemCursorChanged = true
	return nil
}