package main

import (
	"os"
	"strings"

	"github.com/Mr-Spect3r/FAD/core"
	"github.com/Mr-Spect3r/FAD/web"
)

func main() {
	if isWebDownloaderMode() {
	
		lib.RunWebDownloader()
		return
	}
	
	core.RunFAD()
}

func isWebDownloaderMode() bool {
	for i, arg := range os.Args[1:] {
		if arg == "-wd" || arg == "--web-downloader" {
			return true
		}
		if strings.HasPrefix(arg, "-wd=") || strings.HasPrefix(arg, "--web-downloader=") {
			val := strings.Split(arg, "=")[1]
			if val == "true" || val == "1" {
				return true
			}
		}

		if i > 0 && (os.Args[i] == "-wd" || os.Args[i] == "--web-downloader") {
			return true
		}
	}
	
	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		if firstArg == "web" || firstArg == "wd" {
			newArgs := []string{os.Args[0]}
			newArgs = append(newArgs, os.Args[2:]...)
			os.Args = newArgs
			return true
		}
	}
	
	return false
}
