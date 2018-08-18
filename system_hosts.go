package main

import (
	"io/ioutil"
	"os"
	"log"
	"runtime"
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

	setSystemHosts(ret)
}

func getUserHome() string {
	home := ""
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	} else {
		home = os.Getenv("HOME")
	}

	return home
}

func getWinSystemDir() string {
	dir := ""
	if runtime.GOOS == "windows" {
		dir = os.Getenv("windir")
	}

	return dir
}