package main

import (
	"io/ioutil"
	"os"
	"log"
)

func setSystemHosts(content string) {
	stringBytes := []byte(content)
	if err := ioutil.WriteFile(hostsPath, stringBytes, os.ModePerm); err != nil {
		log.Panicln(err)
	}
}

func hItemsToSystemHosts(hosts []hostsItem) {
	ret := ""
	for _, h := range hosts {
		if h.Toggle {
			ret += h.Content + "\n"
		}
	}

	appendToFile("/tmp/yrt", ret)

	setSystemHosts(ret)
}