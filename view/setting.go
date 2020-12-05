package view

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/xiefeihong/crazyreply/utils"
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
	buttonObj, _ := builder.GetObject("button")
	button := buttonObj.(*gtk.Button)
	entryObj1, _ := builder.GetObject("setup1")
	entryObj2, _ := builder.GetObject("setup2")
	entryObj3, _ := builder.GetObject("setup3")
	entryObj4, _ := builder.GetObject("setup4")
	entryObj5, _ := builder.GetObject("setup5")
	entryObj6, _ := builder.GetObject("setup6")
	entryObj7, _ := builder.GetObject("setup7")
	entry1 := entryObj1.(*gtk.Entry)
	entry2 := entryObj2.(*gtk.Entry)
	entry3 := entryObj3.(*gtk.Entry)
	entry4 := entryObj4.(*gtk.Entry)
	entry5 := entryObj5.(*gtk.Entry)
	checkButton1 := entryObj6.(*gtk.CheckButton)
	checkButton2 := entryObj7.(*gtk.CheckButton)
	button.Connect("clicked", func() {
		setSetting(entry1, entry2, entry3, entry4, entry5, checkButton1, checkButton2)
		win.Close()
	})
	settingsToUI(entry1, entry2, entry3, entry4, entry5, checkButton1, checkButton2)
	win.ShowAll()
}

func setSetting(entry1 *gtk.Entry, entry2 *gtk.Entry, entry3 *gtk.Entry, entry4 *gtk.Entry, entry5 *gtk.Entry, checkButton1 *gtk.CheckButton, checkButton2 *gtk.CheckButton){
	text1, _ := entry1.GetText()
	text2, _ := entry2.GetText()
	text3, _ := entry3.GetText()
	text4, _ := entry4.GetText()
	text5, _ := entry5.GetText()
	dateLimit, _ := strconv.Atoi(text1)
	replyNum, _ := strconv.Atoi(text2)
	editNum, _ := strconv.Atoi(text3)
	tagLabs := strings.Split(text4, ",")
	endKeys := strings.Fields(text5)
	random := checkButton1.GetActive()
	printLog := checkButton2.GetActive()
	tags := make(map[string][]string, 0)
	for _, value := range tagLabs {
		messages := utils.Settings.Tags[value]
		tags[value] = messages
	}
	fmt.Println()
	utils.Settings = utils.Setting{dateLimit, replyNum, editNum, tags, endKeys, random, printLog}
	utils.SettingToFile()
}

func settingsToUI(entry1 *gtk.Entry, entry2 *gtk.Entry, entry3 *gtk.Entry, entry4 *gtk.Entry, entry5 *gtk.Entry, checkButton1 *gtk.CheckButton, checkButton2 *gtk.CheckButton){
	entry1.SetText(strconv.FormatInt(int64(utils.Settings.DateLimit), 10))
	entry2.SetText(strconv.FormatInt(int64(utils.Settings.ReplyNum), 10))
	entry3.SetText(strconv.FormatInt(int64(utils.Settings.EditNum), 10))
	tags := make([]string, len(utils.Settings.Tags))
	for label, _ := range utils.Settings.Tags {
		tags = append(tags, label)
		fmt.Print(len(tags))
	}
	entry4.SetText(strings.Join(tags, ","))
	entry5.SetText(strings.Join(utils.Settings.EndKeys, "\t"))
	checkButton1.SetActive(utils.Settings.Random)
	checkButton2.SetActive(utils.Settings.Persion)
	fmt.Println(utils.Settings)
}