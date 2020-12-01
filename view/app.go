package view

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var texts = make([]*gtk.Entry, 0)
var label string

func ShowApp() {
	StartSettings()
	const appID = "org.gtk.example"
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
			label, _ = button2.GetLabel()
			if label == "开始" {
				label = "结束"
				go carryReply(button2)
			} else if label == "结束" {
				label = "开始"
				fmt.Println("The end!")
			}
			button2.SetLabel(label)

			file, err := os.OpenFile("message.txt", os.O_WRONLY | os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			for i := 0; i<Settings.EditNum; i++ {
				str, _ := texts[i].GetText()
				file.WriteString(str + "\n")
			}
		})
		application.AddWindow(win)
		win.ShowAll()
	})
	go keyEvent(Settings.EndKeys)
	os.Exit(application.Run(os.Args))
}

func getMessages() []string {
	messages := make([]string, Settings.EditNum)
	for i:= 0; i<Settings.EditNum; i++ {
		messages[i], _ = texts[i].GetText()
	}
	return messages
}

func carryReply(button *gtk.Button) {
	messages := getMessages()
	for i:= 0; i<Settings.EditNum; i++ {
		texts[i].SetEditable(false)
	}
	if Settings.Random {
		for i:= 0; label == "结束" &&  i < Settings.ReplyNum * Settings.EditNum; i++ {
			reply(messages[rand.Intn(Settings.EditNum)])
		}
	} else {
		for i:= 0 ;label == "结束" && i < Settings.ReplyNum; i++ {
			for j:=0; j< Settings.EditNum; j++ {
				reply(messages[j])
			}
		}
	}
	for i:= 0; i<Settings.EditNum; i++ {
		texts[i].SetEditable(true)
	}
	label = "开始"
	button.SetLabel(label)
}

func reply(message string){
	startTime := time.Now().UnixNano()
	var date int
	if Settings.Persion {
		date = rand.Intn(100) * Settings.DateLimit
	} else {
		date = Settings.DateLimit
	}
	clipboard.WriteAll(message)
	time.Sleep(time.Duration(date) * time.Millisecond)
	robotgo.KeyTap("v", "ctrl")
	time.Sleep(time.Duration(date) * time.Millisecond)
	robotgo.KeyTap("enter")
	time.Sleep(time.Duration(date) * time.Millisecond)
	endTime := time.Now().UnixNano()
	fmt.Println("耗时：", (endTime - startTime)/1000000, "毫秒")
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

	for i := 0; i < Settings.EditNum; i++ {
		box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		label, _ := gtk.LabelNew(strconv.FormatInt(int64(i + 1), 10) + ": ")
		label.SetWidthChars(3)
		entry, _ := gtk.EntryNew()
		entry.SetHExpand(true)
		entry.SetWidthChars(50)

		str, e2 := reader.ReadString('\n')
		if e2 != io.EOF {
			entry.SetText(str[:len(str)-1])
		} else {
			fmt.Println("", e2)
		}
		box.Add(label)
		box.Add(entry)
		bigBox.Add(box)
		texts = append(texts, entry)
		fmt.Println(texts)
	}
}

func keyEvent(keys []string){
	l := len(keys)
	var ok bool
	if l == 2 {
		ok = robotgo.AddEvents(keys[0], keys[1])
	} else if l == 3 {
		ok = robotgo.AddEvents(keys[0], keys[1], keys[2])
	} else if l == 4 {
		ok = robotgo.AddEvents(keys[0], keys[1], keys[2], keys[3])
	} else {
		fmt.Println("不支持的组合键")
	}
	if ok {
		label = "开始"
		fmt.Println("监测到按下了退出组合键")
	}
}