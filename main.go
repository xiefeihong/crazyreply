package main

import (
	"github.com/xiefeihong/crazyreply/utils"
	"github.com/xiefeihong/crazyreply/view"
	"os"
	"path/filepath"
)

func main() {
	path, _ := os.Executable()
	utils.Root = filepath.Dir(path)
	utils.StartSettings()
	view.ShowApp()
}