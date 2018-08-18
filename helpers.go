package main

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"fmt"
	"github.com/jroimartin/gocui"
	"regexp"
	"strconv"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func hostsItemInit() error {
	//
	exists, err := fileExists(initFalgPath)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if backup, err := ioutil.ReadFile(hostsPath); err != nil {
		return err
	} else {
		hItems = append(hItems, hostsItem{
			"backup",
			string(backup),
			true,
		})
		jsonencodeHostsInfoToPath(dataPath, hItems)
	}

	ioutil.WriteFile(initFalgPath, []byte{'y', 'e', 's'}, os.ModePerm)

	return nil
}

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

// 从文件中解析出来hostsItem列表
func jsondecodeHostsInfoFromPath(path string) []hostsItem {
	var items []hostsItem
	b := getHostsItems(path)
	json.Unmarshal(b, &items)

	return items
}

func jsonencodeHostsInfoToPath(path string, hItems []hostsItem) error {
	if b, err := json.Marshal(hItems); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}

	return nil
}

func marshal() []byte {
	b, err := json.Marshal(hItems)
	if err != nil {
		return []byte{}
	}
	return b
}


// 从数据文件中读取配置信息
func getHostsItems(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}

	return b
}

func setCursorView(g *gocui.Gui) {
	if !onMsgView {
		setViewOnTop(g, getCurrentTabViewName())
	}
}

func refreshEnd(g *gocui.Gui) {
	if hItemCursorChanged {
		setCursor(g, "main", 0, 0)
		setOrigin(g, "main", 0, 0)
		hItemCursorChanged = !hItemCursorChanged
	}


	if hItemChanged {
		hItemChanged = !hItemChanged
	}
}

// 设置当前view聚焦
func setViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}


// 获取当前所聚焦的tabView名称
func getCurrentTabViewName() string {
	viewName := tabViews[tabViewIndex]

	return viewName
}

// 渲染一个view中的内容 Editable 有bug，更新之后，不能edit
func renderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		// just in case the view disappeared as this function was called, we'll
		// silently return if it's not found
		if err != nil {
			return err
		}
		v.Wrap = true
		v.Clear()
		fmt.Fprint(v, s)
		return nil
	})
	return nil
}

// 渲染一个view中的内容
func renderStringOriginCursor(g *gocui.Gui, viewName, s string, originX, originY, cursorX, cursorY int) error {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(viewName)
		// just in case the view disappeared as this function was called, we'll
		// silently return if it's not found
		if err != nil {
			return nil
		}
		v.Clear()
		fmt.Fprint(v, s)
		v.SetOrigin(originX, originY)
		v.SetCursor(cursorX, cursorY)
		v.Wrap = true
		return nil
	})
	return nil
}

func setCursor(g *gocui.Gui, viewName string, cursorX, cursorY int) error {
	v, err := g.View(viewName)
	if err != nil {
		return err
	}
	if err := v.SetCursor(cursorX, cursorY); err != nil {
		return err
	}
	return nil
}
func setOrigin(g *gocui.Gui, viewName string, originX, originY int) error {
	v, err := g.View(viewName)
	if err != nil {
		return err
	}
	if err := v.SetOrigin(originX, originY); err != nil {
		return err
	}
	return nil
}


func getSlideRowCount() int {
	ret := 0
	ret = maxY - 4
	return ret
}

func adjustCursorOrigin() error {
	gap := getSlideRowCount() - slideCursorY
	if gap > 0 {
		if slideOriginY > 0 {
			slideCursorY += gap
			slideOriginY -= gap

			if slideOriginY < 0 {
				slideOriginY = 0
			}
		}
	} else if gap < 0 {
		slideCursorY += gap
		slideOriginY -= gap
	}

	if slideCursorY < 0 {
		slideCursorY = 0
	}

	if slideOriginY < 0 {
		slideOriginY = 0
	}

	return nil
}

func setCurrentHostsItemContent(content string) {
	idx := getCurrentHostsItemIndex()
	hItems[idx].Content = content
}

// 获取真实的索引
func getCurrentHostsItemIndex() int {
	ret := 0

	ret = slideCursorY + slideOriginY
	appendToFile("/tmp/yrt", strconv.Itoa(slideCursorY) + "," + strconv.Itoa(slideOriginY) + "\n")

	return ret
}

func getCurrentHostsItemContent() string {
	ret := ""
	ret = hItems[getCurrentHostsItemIndex()].Content
	return ret
}



func checkHostsItemLine(line string) bool {
	if ret, err := regexp.Match("^#.*", []byte(line)); err != nil {
		return false
	} else {
		if ret {
			return true
		}
	}


	if ret, err := regexp.Match("^[    ]*$", []byte(line)); err != nil {
		return false
	} else {
		if ret {
			return true
		}
	}

	ret, err := regexp.Match("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}(?:[  ]{1,})(?:[ 0-9a-zA-Z._-]{1,}){1,}$", []byte(line))
	if err != nil {
		return false
	}

	return ret
}