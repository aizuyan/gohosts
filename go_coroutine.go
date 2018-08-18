package main

import (
	"time"
	"strings"
)

func checkSyntax() {
	old := ""
	for range time.Tick(time.Millisecond * 50) {
		mainView := g.CurrentView()
		if "main" == mainView.Name(){
			mainBuffer := mainView.Buffer()
			if old != mainBuffer {
				old = mainBuffer
				lines := mainView.ViewBufferLines()

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
		}
	}
}