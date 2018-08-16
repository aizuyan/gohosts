package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"github.com/fatih/color"
	"errors"
	"strconv"
)

var (
	dataPath string = "/tmp/hosts-data/hostsItems"
	hostsPath string = "/etc/hosts"

	// global terminal size，change on func layout
	maxX, maxY int

	// 全局变量存储gocui.Gui实例，操作更方便
	g *gocui.Gui

	// hosts item 实例列表，存储hosts的名称、内容、是否开启状态
	hItems []hostsItem
	// hosts item 当前聚焦的实例索引
	hItemIndex int = 0

	// tab键切换支持的view列表
	tabViews []string = []string{
		"slide",
		"main",
	}
	// tab键切换view，当前停留的view索引
	tabViewIndex int = 0

	// red color string factory
	red func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()

	// TODO resize cursor origin
	slideOriginX, slideOriginY int = 0, 0
	slideCursorX, slideCursorY int = 0, 0
)

func main() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	hItems = jsondecodeHostsInfoFromPath(dataPath)

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(layout)

	// 绑定⌨️事件
	if err := KeyBindingAction(g); err !=nil {
		log.Panicln(err)
	}

	if err := g.MainLoop();  err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY = g.Size()

	// 最大，最小限制
	if maxX < 100 || maxY < 30 {
		return errors.New("mininum width limit 100, mininum height limit 30")
	}

	// 调整坐标和偏移
	adjustCursorOrigin()

	// slide
	if v, err := g.SetView("slide",0, 0, 40, maxY - 2); err != nil{
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.FgColor = gocui.AttrBold
		v.SelFgColor = gocui.AttrBold
		v.SelBgColor = gocui.AttrReverse
		v.Title= "Hosts Items"
	}
	// 重新绘制slide中的内容，包括设置偏移和焦点
	renderStringOriginCursor(
		g, "slide", hostsNameToString(hItems),
		slideOriginX, slideOriginY, slideCursorX, slideCursorY)

	// main
	if v, err := g.SetView("main",41, 0, maxX - 1, maxY - 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = true
		v.Title = "Hosts Info"
	}
	renderString(g, "main", "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n" + strconv.Itoa(getSlideRowCount()))

	// 操作提示
	if v, err := g.SetView("footer", 0, maxY - 2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelFgColor = gocui.ColorMagenta
		v.Frame = false
	}
	renderString(
		g, "footer",
		"help: `tab`: switch view; `↑` `↓` change hosts item; " +
			"`←` `→` toggle hosts item; `shift + a` add hosts item; `shift + q` cansle add action")
	return nil
}