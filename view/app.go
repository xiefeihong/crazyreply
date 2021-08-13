package view

import (
	"github.com/go-vgo/robotgo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/xiefeihong/crazyreply/utils"
	"os"
	"strconv"
)

func ShowApp() {
	const appID = "top.xiefeihong.crazyreply"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		panic(err)
	}
	application.Connect("activate", func() {
		onActivate(application)
	})
	os.Exit(application.Run(os.Args))
}

func onActivate(application *gtk.Application) {
	win, _ := gtk.ApplicationWindowNew(application)
	win.SetIconFromFile(utils.Root + "/view/ui/icon.ico")
	win.SetSizeRequest(450, 450)
	win.SetTitle("疯狂回复")
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	utils.Book, _ = gtk.NotebookNew()
	for _, tag := range utils.Settings.Tags {
		bookPage := createBookPage(tag.Msgs)
		label, _ := gtk.LabelNew(tag.Label)
		utils.Book.AppendPage(bookPage, label)
		utils.Book.Connect("switch-page", func() {
			utils.ReplyIndex = 0
			utils.MsgIndex = 0
		})
	}
	box.Add(utils.Book)
	bottonBox := createBottonBox()
	box.Add(bottonBox)
	win.Add(box)
	win.ShowAll()
}

func createBookPage(msgs []string) *gtk.Box {
	topBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	scrolledWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	scrolledWindow.SetVExpand(true)
	viewport, _ := gtk.ViewportNew(nil, nil)
	textsBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	textsBox.SetMarginTop(10)
	textsBox.SetMarginStart(10)
	textsBox.SetMarginEnd(10)
	textsBox.SetMarginBottom(10)
	ts := make([]*gtk.Entry, 0)
	for i := 0; i < utils.Settings.EditNum; i++ {
		inputLabel := strconv.FormatInt(int64(i + 1), 10) + ": "
		var msg string
		if i < len(msgs) {
			msg = msgs[i]
		} else {
			msg = ""
		}
		textBox, inputEntry := createInputBox(inputLabel, msg)
		textsBox.Add(textBox)
		ts = append(ts, inputEntry)
	}
	utils.Texts = append(utils.Texts, ts)
	viewport.Add(textsBox)
	scrolledWindow.Add(viewport)
	topBox.Add(scrolledWindow)
	return topBox
}

func createInputBox(inputLabel string, messages string) (*gtk.Box, *gtk.Entry) {
	lineBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	label, _ := gtk.LabelNew(inputLabel)
	label.SetWidthChars(3)
	lineBox.Add(label)
	inputEntry, _ := gtk.EntryNew()
	inputEntry.SetHExpand(true)
	inputEntry.SetText(messages)
	lineBox.Add(inputEntry)
	return lineBox, inputEntry
}

func createBottonBox() *gtk.Box {
	bottomBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	bottomBox.SetHAlign(gtk.ALIGN_CENTER)
	bottomBox.SetMarginBottom(10)
	settingBtn, _ := gtk.ButtonNew()
	settingBtn.SetLabel("设置")
	startBtn, _ := gtk.ButtonNew()
	startBtn.SetLabel("开始")
	utils.BottonLabel = "开始"
	settingBtn.Connect("clicked", func() {
		if utils.BottonLabel == "开始" {
			ShowSetting()
		}
	})
	startBtn.Connect("clicked", func() {
		pageIndex := utils.Book.GetCurrentPage()
		utils.BottonLabel, _ = startBtn.GetLabel()
		if utils.BottonLabel == "开始" {
			msgs := utils.Settings.Tags[pageIndex].Msgs
			var validNum = 0
			for i := 0; i < utils.Settings.EditNum; i++ {
				str, _ := utils.Texts[pageIndex][i].GetText()
				if str != "" {
					if i < len(msgs) {
						msgs[i] = str
					} else {
						msgs = append(msgs, str)
					}
					validNum ++
				}
			}
			if validNum < len(msgs) {
				msgs = msgs[:validNum]
			}
			utils.Settings.Tags[pageIndex].Msgs = msgs
			utils.SettingToFile()
			if len(msgs) != 0 {
				utils.BottonLabel = "结束"
				robotgo.KeyTap("tab")
				go utils.CarryReply(startBtn)
				go utils.KeyDownEvent(utils.Settings.EndKeys)
			}
		} else if utils.BottonLabel == "结束" {
			utils.BottonLabel = "开始"
			robotgo.KeyTap("tab")
		}
		startBtn.SetLabel(utils.BottonLabel)
	})
	bottomBox.Add(settingBtn)
	bottomBox.Add(startBtn)
	return bottomBox
}