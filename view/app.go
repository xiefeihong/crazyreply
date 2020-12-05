package view

import (
	"bufio"
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/xiefeihong/crazyreply/utils"
	"os"
	"strconv"
)

func ShowApp() {
	utils.StartSettings()
	const appID = "top.xiefeihong.crazyreply"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		fmt.Println("Could not create application.", err)
	}
	application.Connect("activate", func() {
		win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
		win.SetSizeRequest(450, 450)
		win.SetTitle("疯狂回复")
		book, _ := gtk.NotebookNew()
		setTags(book)
		win.Add(book)
		application.AddWindow(win)
		win.ShowAll()
	})
	go utils.KeyEvent(utils.Settings.EndKeys)
	os.Exit(application.Run(os.Args))
}
func setTags(book *gtk.Notebook){
	bigBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	hdjustment, _ := gtk.AdjustmentNew(100, -1, -1, -1, -1, -1)
	vdjustment, _ := gtk.AdjustmentNew(100, -1, -1, -1, -1, -1)
	scrolledWindow, _ := gtk.ScrolledWindowNew(hdjustment, vdjustment)
	scrolledWindow.SetVExpand(true)
	viewport, _ := gtk.ViewportNew(hdjustment, vdjustment)
	aspectFrame, _ := gtk.AspectFrameNew("", 0.5, 0, 1, true)
	aspectFrame.SetShadowType(gtk.SHADOW_NONE)
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	bigBox.Add(scrolledWindow)
	scrolledWindow.Add(viewport)
	scrolledWindow.SetMarginTop(10)
	scrolledWindow.SetMarginStart(10)
	scrolledWindow.SetMarginEnd(10)
	viewport.Add(aspectFrame)
	aspectFrame.Add(box)
	setInputBox(box)
	bigBox.Add(createBottonAspectFrame())
	book.Add(bigBox)
	label, _ := gtk.LabelNew("label1")
	label.SetProperty("tab-fill", false)
	book.SetTabLabel(bigBox, label)
}

func setInputBox(bigBox *gtk.Box){
	file, e := os.Open("message.txt")
	if e != nil {
		fmt.Printf("打开文件失败：%v\n", e)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	for i := 0; i < utils.Settings.EditNum; i++ {
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		label, _ := gtk.LabelNew(strconv.FormatInt(int64(i + 1), 10) + ": ")
		label.SetWidthChars(3)
		entry, _ := gtk.EntryNew()
		entry.SetHExpand(true)
		entry.SetWidthChars(50)

		str, _ := reader.ReadString('\n')
		l := len(str)
		if l>0 {
			entry.SetText(str[:l-1])
		}
		box.Add(label)
		box.Add(entry)
		bigBox.Add(box)
		utils.Texts = append(utils.Texts, entry)
		fmt.Print(len(utils.Texts))
	}
	fmt.Println()
}

func createBottonAspectFrame() *gtk.AspectFrame {
	bottonAspectFrame, _ := gtk.AspectFrameNew("", 0.5, 0.5, 1, true)
	bottonAspectFrame.SetShadowType(gtk.SHADOW_NONE)
	bottonAspectFrame.SetMarginBottom(10)
	bottomBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	button1, _ := gtk.ButtonNew()
	button1.SetLabel("设置")
	button2, _ := gtk.ButtonNew()
	button2.SetLabel("开始")
	button1.Connect("clicked", func() {
		ShowSetting()
	})
	button2.Connect("clicked", func() {
		utils.Label, _ = button2.GetLabel()
		if utils.Label == "开始" {
			utils.Label = "结束"
			go utils.CarryReply(button2)
		} else if utils.Label == "结束" {
			utils.Label = "开始"
			fmt.Println("The end!")
		}
		button2.SetLabel(utils.Label)

		file, err := os.OpenFile("message.txt", os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		for i := 0; i<utils.Settings.EditNum; i++ {
			str, _ := utils.Texts[i].GetText()
			file.WriteString(str + "\n")
		}
	})
	bottonAspectFrame.Add(bottomBox)
	bottomBox.Add(button1)
	bottomBox.Add(button2)
	return bottonAspectFrame
}