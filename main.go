package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"github.com/fatih/color"
	"errors"
)

var (
	dataPath string = ""
	initFalgPath string = ""
	hostsPath string = ""

	// global terminal size，change on func layout
	maxX, maxY int

	// 全局变量存储gocui.Gui实例，操作更方便
	g *gocui.Gui

	// hosts item 实例列表，存储hosts的名称、内容、是否开启状态
	hItems []hostsItem

	// tab键切换支持的view列表
	tabViews []string = []string{
		"slide",
		"main",
	}
	// tab键切换view，当前停留的view索引
	tabViewIndex int = 0

	onMsgView bool = false

	// red color string factory
	red func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()

	// green color string factory
	green func(a ...interface{}) string = color.New(color.FgGreen).SprintFunc()

	// TODO resize cursor origin
	slideOriginX, slideOriginY int = 0, 0
	slideCursorX, slideCursorY int = 0, 0

	// hItems是否切换过
	hItemCursorChanged bool = false

	// mainContent 中的内容是否改变过；切换hosts改变
	mainContentChanged bool = true

	// hosts Item 中的内容是否改变；添加hosts item的时候改变
	hItemChanged bool = true
)

func main() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// 初始化变量
	gohostsInit()

	hItems = jsondecodeHostsInfoFromPath(dataPath)

	// 初始化
	hostsItemInit()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(layout)

	// 绑定⌨️事件
	if err := KeyBindingAction(g); err !=nil {
		log.Panicln(err)
	}

	// 校验语法
	go checkSyntax()

	// 循环开始
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
	if hItemChanged || hItemCursorChanged {
		renderStringOriginCursor(
			g, "slide", hostsNameToString(hItems),
			slideOriginX, slideOriginY, slideCursorX, slideCursorY)
	}

	// main
	if v, err := g.SetView("main",41, 0, maxX - 1, maxY - 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Editable = true
		v.Title = "Hosts Item Content"
	}

	if mainContentChanged {
		renderString(g, "main", getCurrentHostsItemContent())
		mainContentChanged = false
	}

	// 操作提示
	if v, err := g.SetView("footer", 0, maxY - 2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelFgColor = gocui.ColorMagenta
		v.Frame = false
		renderString(
			g, "footer",
			"help: `tab`: switch view; `↑` `↓` change hosts item; " +
				"`←` `→` toggle hosts item; `shift + a` add hosts item; `shift + q` cansle add action")
	}

	setCursorView(g)
	refreshEnd(g)
	return nil
}