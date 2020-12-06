package view

import (
	"github.com/go-vgo/robotgo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/xiefeihong/crazyreply/utils"
	"os"
	"strconv"
)

var book *gtk.Notebook

func ShowApp() {
	utils.StartSettings()
	const appID = "top.xiefeihong.crazyreply"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		panic(err)
	}
	application.Connect("activate", func() {
		win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		win.SetIconFromFile("view/ui/icon.jpg")
		win.SetSizeRequest(450, 450)
		win.SetTitle("疯狂回复")
		book, _ = gtk.NotebookNew()
		for lab, msgs := range utils.Settings.Tags {
			bookPage := createBookPage(lab, msgs)
			bottonAspectFrame := createBottonAspectFrame(lab)
			bookPage.Add(bottonAspectFrame)
			label, _ := gtk.LabelNew(lab)
			book.AppendPage(bookPage, label)
		}
		win.Add(book)
		application.AddWindow(win)
		win.ShowAll()
	})
	go utils.KeyEvent(utils.Settings.EndKeys)
	os.Exit(application.Run(os.Args))
}

func createBookPage(pageLabel string, msgs []string) *gtk.Box {
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
	entrys := make([]*gtk.Entry, 0)
	for i := 0; i < utils.Settings.EditNum; i++ {
		inputLabel := strconv.FormatInt(int64(i + 1), 10) + ": "
		var msg string
		if i < len(msgs) {
			msg = msgs[i]
		} else {
			msg = ""
		}
		textBox, enrty := createInputBox(inputLabel, msg)
		textsBox.Add(textBox)
		entrys = append(entrys, enrty)
	}
	utils.Texts[pageLabel] = entrys
	topBox.Add(scrolledWindow)
	return topBox
}

func createInputBox(inputLabel string, messages string) (*gtk.Box, *gtk.Entry) {
	lineBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	label, _ := gtk.LabelNew(inputLabel)
	label.SetWidthChars(3)
	entry, _ := gtk.EntryNew()
	entry.SetHExpand(true)
	entry.SetWidthChars(50)
	entry.SetText(messages)
	lineBox.Add(label)
	lineBox.Add(entry)
	return lineBox, entry
}

func createBottonAspectFrame(label string) *gtk.AspectFrame {
	bottonAspectFrame, _ := gtk.AspectFrameNew("", 0.5, 0.5, 1, true)
	bottonAspectFrame.SetShadowType(gtk.SHADOW_NONE)
	bottonAspectFrame.SetMarginBottom(10)
	bottomBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	button1, _ := gtk.ButtonNew()
	button1.SetLabel("设置")
	button2, _ := gtk.ButtonNew()
	button2.SetLabel("开始")
	utils.BottonLabel = "开始"
	button1.Connect("clicked", func() {
		if utils.BottonLabel == "开始" {
			ShowSetting()
		}
	})
	button2.Connect("clicked", func() {
		utils.PageLabel = label
		utils.BottonLabel, _ = button2.GetLabel()
		if utils.BottonLabel == "开始" {
			msgs := utils.Settings.Tags[label]
			var validNum = 0
			for i := 0; i < utils.Settings.EditNum; i++ {
				str, _ := utils.Texts[label][i].GetText()
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
			utils.Settings.Tags[label] = msgs
			utils.SettingToFile()
			utils.BottonLabel = "结束"
			go utils.CarryReply(button2)
			robotgo.KeyTap("tab")
		} else if utils.BottonLabel == "结束" {
			utils.BottonLabel = "开始"
		}
		button2.SetLabel(utils.BottonLabel)
	})
	bottonAspectFrame.Add(bottomBox)
	bottomBox.Add(button1)
	bottomBox.Add(button2)
	return bottonAspectFrame
}