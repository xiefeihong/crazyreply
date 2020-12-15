package view

import (
	"github.com/go-vgo/robotgo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/xiefeihong/crazyreply/utils"
	"os"
	"strconv"
	"time"
)

var book *gtk.Notebook

func ShowApp() {
	const appID = "top.xiefeihong.crazyreply"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		panic(err)
	}
	application.Connect("activate", func() {
		win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		win.SetIconFromFile(utils.Root + "/view/ui/icon.ico")
		win.SetSizeRequest(450, 450)
		win.SetTitle("疯狂回复")
		book, _ = gtk.NotebookNew()
		for tagIndex, tag := range utils.Settings.Tags {
			bookPage := createBookPage(tagIndex, tag.Msgs)
			bottonAspectFrame := createBottonAspectFrame(tagIndex, tag.Label)
			bookPage.Add(bottonAspectFrame)
			label, _ := gtk.LabelNew(tag.Label)
			book.AppendPage(bookPage, label)
		}
		win.Add(book)
		application.AddWindow(win)
		win.ShowAll()
	})
	os.Exit(application.Run(os.Args))
}

func createBookPage(pageIndex int, msgs []string) *gtk.Box {
	topBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	hdjustment, _ := gtk.AdjustmentNew(-1, -1, -1, -1, -1, -1)
	vdjustment, _ := gtk.AdjustmentNew(-1, -1, -1, -1, -1, -1)
	scrolledWindow, _ := gtk.ScrolledWindowNew(hdjustment, vdjustment)
	scrolledWindow.SetVExpand(true)
	viewport, _ := gtk.ViewportNew(hdjustment, vdjustment)
	aspectFrame, _ := gtk.AspectFrameNew("", 0.5, 0, 1, true)
	aspectFrame.SetShadowType(gtk.SHADOW_NONE)
	textsBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	scrolledWindow.Add(viewport)
	scrolledWindow.SetMarginTop(10)
	scrolledWindow.SetMarginStart(10)
	scrolledWindow.SetMarginEnd(10)
	viewport.Add(aspectFrame)
	aspectFrame.Add(textsBox)
	for i := 0; i < utils.Settings.EditNum; i++ {
		inputLabel := strconv.FormatInt(int64(i + 1), 10) + ": "
		var msg string
		if i < len(msgs) {
			msg = msgs[i]
		} else {
			msg = ""
		}
		textBox, textEnrty := createInputBox(inputLabel, msg)
		textsBox.Add(textBox)
		utils.Texts[pageIndex] = append(utils.Texts[pageIndex], textEnrty)
	}
	topBox.Add(scrolledWindow)
	return topBox
}

func createInputBox(inputLabel string, messages string) (*gtk.Box, *gtk.Entry) {
	lineBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	label, _ := gtk.LabelNew(inputLabel)
	label.SetWidthChars(3)
	inputEntry, _ := gtk.EntryNew()
	inputEntry.SetHExpand(true)
	inputEntry.SetWidthChars(50)
	inputEntry.SetText(messages)
	lineBox.Add(label)
	lineBox.Add(inputEntry)
	return lineBox, inputEntry
}

func createBottonAspectFrame(tagIndex int, label string) *gtk.AspectFrame {
	bottonAspectFrame, _ := gtk.AspectFrameNew("", 0.5, 0.5, 1, true)
	bottonAspectFrame.SetShadowType(gtk.SHADOW_NONE)
	bottonAspectFrame.SetMarginBottom(10)
	bottomBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
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
		utils.PageIndex = tagIndex
		utils.BottonLabel, _ = startBtn.GetLabel()
		if utils.BottonLabel == "开始" {
			msgs := utils.Settings.Tags[tagIndex].Msgs
			var validNum = 0
			for i := 0; i < utils.Settings.EditNum; i++ {
				str, _ := utils.Texts[tagIndex][i].GetText()
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
			utils.Settings.Tags[tagIndex].Msgs = msgs
			utils.SettingToFile()
			utils.BottonLabel = "结束"
			robotgo.KeyTap("tab")
			time.Sleep(50 * time.Millisecond)
			go utils.CarryReply(startBtn)
			go utils.KeyDownEvent(utils.Settings.EndKeys)
		} else if utils.BottonLabel == "结束" {
			utils.BottonLabel = "开始"
		}
		startBtn.SetLabel(utils.BottonLabel)
	})
	bottonAspectFrame.Add(bottomBox)
	bottomBox.Add(settingBtn)
	bottomBox.Add(startBtn)
	return bottonAspectFrame
}