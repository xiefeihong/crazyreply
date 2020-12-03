package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/xiefeihong/crazyreply/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

func ShowSetting() {
	builder, err := gtk.BuilderNewFromFile("view/ui/setting.glade")
	if err != nil {
		fmt.Println(err)
	}
	winObj, _ := builder.GetObject("window1")
	win := winObj.(*gtk.Window)
	buttonObj1, _ := builder.GetObject("button1")
	buttonObj2, _ := builder.GetObject("button2")
	button1 := buttonObj1.(*gtk.Button)
	button2 := buttonObj2.(*gtk.Button)
	entryObj1, _ := builder.GetObject("setup1")
	entryObj2, _ := builder.GetObject("setup2")
	entryObj3, _ := builder.GetObject("setup3")
	entryObj4, _ := builder.GetObject("setup4")
	entryObj5, _ := builder.GetObject("setup5")
	entryObj6, _ := builder.GetObject("setup6")
	entry1 := entryObj1.(*gtk.Entry)
	entry2 := entryObj2.(*gtk.Entry)
	entry3 := entryObj3.(*gtk.Entry)
	entry4 := entryObj4.(*gtk.Entry)
	checkButton1 := entryObj5.(*gtk.CheckButton)
	checkButton2 := entryObj6.(*gtk.CheckButton)
	button1.Connect("clicked", func() {
		setSetting(entry1, entry2, entry3, entry4, checkButton1, checkButton2)
	})
	button2.Connect("clicked", func() {
		setSetting(entry1, entry2, entry3, entry4, checkButton1, checkButton2)
		win.Close()
	})
	settingsToUI(entry1, entry2, entry3, entry4, checkButton1, checkButton2)
	win.ShowAll()
}

func setSetting(entry1 *gtk.Entry, entry2 *gtk.Entry, entry3 *gtk.Entry, entry4 *gtk.Entry, checkButton1 *gtk.CheckButton, checkButton2 *gtk.CheckButton){
	text1, _ := entry1.GetText()
	text2, _ := entry2.GetText()
	text3, _ := entry3.GetText()
	text4, _ := entry4.GetText()
	dateLimit, _ := strconv.Atoi(text1)
	replyNum, _ := strconv.Atoi(text2)
	editNum, _ := strconv.Atoi(text3)
	endKeys := strings.Fields(text4)
	random := checkButton1.GetActive()
	printLog := checkButton2.GetActive()
	utils.Settings = utils.Setting{dateLimit, replyNum, editNum, endKeys, random, printLog}
	config, _ := json.Marshal(utils.Settings)
	file, err := os.OpenFile("config.json", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	var str bytes.Buffer
	json.Indent(&str, config,"", "\t")
	file.Write(str.Bytes())
}

func settingsToUI(entry1 *gtk.Entry, entry2 *gtk.Entry, entry3 *gtk.Entry, entry4 *gtk.Entry, checkButton1 *gtk.CheckButton, checkButton2 *gtk.CheckButton){
	entry1.SetText(strconv.FormatInt(int64(utils.Settings.DateLimit), 10))
	entry2.SetText(strconv.FormatInt(int64(utils.Settings.ReplyNum), 10))
	entry3.SetText(strconv.FormatInt(int64(utils.Settings.EditNum), 10))
	entry4.SetText(strings.Join(utils.Settings.EndKeys, "\t"))
	checkButton1.SetActive(utils.Settings.Random)
	checkButton2.SetActive(utils.Settings.Persion)
	fmt.Println(utils.Settings)
}