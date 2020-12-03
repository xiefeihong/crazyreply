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
		builder, err := gtk.BuilderNewFromFile("view/ui/app.glade")

		if err != nil {
			fmt.Println(err)
		}
		winObj, _ := builder.GetObject("window1")
		win := winObj.(*gtk.Window)
		setInputBox(builder)
		buttonObj1, _ := builder.GetObject("button1")
		button1 := buttonObj1.(*gtk.Button)
		button1.Connect("clicked", func() {
			ShowSetting()
		})
		buttonObj2, _ := builder.GetObject("button2")
		button2 := buttonObj2.(*gtk.Button)
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
		application.AddWindow(win)
		win.ShowAll()
	})
	go utils.KeyEvent(utils.Settings.EndKeys)
	os.Exit(application.Run(os.Args))
}

func setInputBox(builder *gtk.Builder){
	boxObj, _ := builder.GetObject("inputs")
	bigBox := boxObj.(*gtk.Box)

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