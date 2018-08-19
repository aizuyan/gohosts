package main

import (
	"time"
	"strings"
)

func checkSyntax() {
	old := ""
	for range time.Tick(time.Millisecond * 50) {
		view := g.CurrentView()
		viewName := view.Name()
		if "main" == viewName{
			mainBuffer := view.Buffer()
			if old != mainBuffer {
				old = mainBuffer
				lines := view.ViewBufferLines()

				for idx, line := range lines {
					if false == checkHostsItemLine(line) {
						lines[idx] = red(line)
					}
				}
				newBuffer := "\n" + strings.Join(lines, "\n")
				renderString(g, "main", newBuffer)
				if newBuffer[len(newBuffer) - 1] != '\n' {
					newBuffer += "\n"
				}
				// 写入数据库
				setCurrentHostsItemContent(newBuffer)
				jsonencodeHostsInfoToPath(dataPath, hItems)

				// 提现到系统hosts中
				hItemsToSystemHosts(hItems)
			}
		} else if "new-hosts-item-msg" == viewName {
			hostsName := view.ViewBuffer()
			hostsName = strings.Trim(hostsName, "\n ")
			if !checkHostsItemName(hostsName) {
				renderString(g, "new-hosts-item-msg", red(hostsName))
			} else {
				renderString(g, "new-hosts-item-msg", green(hostsName))
			}
		}
	}
}