package view

import (
	"github.com/go-vgo/robotgo"
	"github.com/gotk3/gotk3/gtk"
	hook "github.com/robotn/gohook"
	"github.com/xiefeihong/crazyreply/utils"
	"runtime"
	"strconv"
	"strings"
)

var (
	labelEntrys []*gtk.Entry
	vsHold []uint16
	vsUp []uint16
	replyNumBox *gtk.Box
	beforeKeysBox *gtk.Box
	withoutStopCheckButton *gtk.CheckButton
	beforeCheckButton *gtk.CheckButton
)

func ShowSetting() {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetIconFromFile(utils.Root + "/view/ui/icon.ico")
	win.SetTitle("疯狂回复")
	win.SetResizable(false)
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	topBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	topBox.SetMarginTop(10)
	topBox.SetMarginStart(10)
	topBox.SetMarginEnd(10)
	topBox.SetHomogeneous(true)
	topLeftBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	topLeft1Box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	dateLimitBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	dateLimitLabel, _ := gtk.LabelNew("时间间隔: ")
	dateLimitSpinBtn, _ := gtk.SpinButtonNewWithRange(0, 10000000, 100)
	dateLimitSpinBtn.SetHExpand(true)
	dateLimitBox.Add(dateLimitLabel)
	dateLimitBox.Add(dateLimitSpinBtn)
	topLeft1Box.Add(dateLimitBox)

	replyNumBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	replyNumLabel, _ := gtk.LabelNew("最大次数: ")
	replyNumSpinBtn, _ := gtk.SpinButtonNewWithRange(1, 1000, 10)
	replyNumSpinBtn.SetHExpand(true)
	replyNumBox.Add(replyNumLabel)
	replyNumBox.Add(replyNumSpinBtn)
	topLeft1Box.Add(replyNumBox)

	editNumBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	editNumLabel, _ := gtk.LabelNew("编辑条数: ")
	editNumSpinBtn, _ := gtk.SpinButtonNewWithRange(3, 30, 1)
	editNumSpinBtn.SetHExpand(true)
	editNumBox.Add(editNumLabel)
	editNumBox.Add(editNumSpinBtn)
	topLeft1Box.Add(editNumBox)

	strategyBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	strategyLabel, _ := gtk.LabelNew("发送方式: ")
	strategyComboBoxText, _ := gtk.ComboBoxTextNew()
	strategyComboBoxText.Append(string(utils.Key), "键盘")
	strategyComboBoxText.Append(string(utils.Clipboard), "粘贴板")
	strategyComboBoxText.SetHExpand(true)
	strategyBox.Add(strategyLabel)
	strategyBox.Add(strategyComboBoxText)
	topLeft1Box.Add(strategyBox)

	beforeKeysBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	beforeKeysLabel, _ := gtk.LabelNew("前置按键: ")
	beforeKeysEntry, _ := gtk.EntryNew()
	beforeKeysEntry.SetHExpand(true)
	beforeKeysBox.Add(beforeKeysLabel)
	beforeKeysBox.Add(beforeKeysEntry)
	topLeft1Box.Add(beforeKeysBox)

	endKeysBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	endKeysLabel, _ := gtk.LabelNew("快捷终止: ")
	endKeysEntry, _ := gtk.EntryNew()
	endKeysEntry.SetHExpand(true)
	endKeysBox.Add(endKeysLabel)
	endKeysBox.Add(endKeysEntry)
	topLeft1Box.Add(endKeysBox)
	topLeftBox.Add(topLeft1Box)

	topLeft2Box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	topLeft2Box.SetHomogeneous(true)
	checkButton1Box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	checkButton1Box.SetHAlign(gtk.ALIGN_CENTER)
	checkButton2Box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	checkButton2Box.SetHAlign(gtk.ALIGN_CENTER)
	randomCheckButton, _ := gtk.CheckButtonNewWithLabel("随机消息")
	withoutStopCheckButton, _ = gtk.CheckButtonNewWithLabel("连续不停")
	averageCheckButton, _ := gtk.CheckButtonNewWithLabel("间隔平均")
	beforeCheckButton, _ = gtk.CheckButtonNewWithLabel("前置按键")
	checkButton1Box.Add(randomCheckButton)
	checkButton1Box.Add(withoutStopCheckButton)
	checkButton2Box.Add(averageCheckButton)
	checkButton2Box.Add(beforeCheckButton)
	topLeft2Box.Add(checkButton1Box)
	topLeft2Box.Add(checkButton2Box)
	topLeftBox.Add(topLeft2Box)
	topBox.Add(topLeftBox)

	topRightBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	rightLabel, _ := gtk.LabelNew("标签分类")
	rightLabel.SetHExpand(true)
	rightLabel.SetMarginBottom(10)
	scrolledWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	scrolledWindow.SetShadowType(gtk.SHADOW_IN)
	scrolledWindow.SetVExpand(true)
	viewport, _ := gtk.ViewportNew(nil, nil)
	textLabsBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	viewport.Add(textLabsBox)
	scrolledWindow.Add(viewport)

	rightBottomBox, _:= gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	buttonDecreaseBtn, _ := gtk.ButtonNewWithLabel("-")
	buttonDecreaseBtn.SetHExpand(true)
	buttonIncreaseBtn, _ := gtk.ButtonNewWithLabel("+")
	buttonIncreaseBtn.SetHExpand(true)
	rightBottomBox.Add(buttonDecreaseBtn)
	rightBottomBox.Add(buttonIncreaseBtn)
	topRightBox.Add(rightLabel)
	topRightBox.Add(scrolledWindow)
	topRightBox.Add(rightBottomBox)
	topBox.Add(topRightBox)

	bottomBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	bottomBox.SetHAlign(gtk.ALIGN_CENTER)
	bottomBox.SetMarginBottom(10)
	cancelBtn, _ := gtk.ButtonNewWithLabel("取消")
	preserveBtn, _ := gtk.ButtonNewWithLabel("保存")
	bottomBox.Add(cancelBtn)
	bottomBox.Add(preserveBtn)

	box.Add(topBox)
	box.Add(bottomBox)
	win.Add(box)

	withoutStopCheckButton.Connect("toggled", func() {
		setReplyNumBox()
	})
	beforeCheckButton.Connect("toggled", func() {
		setBeforeBox()
	})
	cancelBtn.Connect("clicked", func() {
		win.Close()
	})
	preserveBtn.Connect("clicked", func() {
		restart := changeData(dateLimitSpinBtn, replyNumSpinBtn, editNumSpinBtn, strategyComboBoxText, beforeKeysEntry, endKeysEntry, randomCheckButton, averageCheckButton, beforeCheckButton, labelEntrys)
		if restart {
			dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", "请重新打开此程序")
			dialog.Run()
			dialog.Destroy()
		}
		win.Close()
	})
	buttonIncreaseBtn.Connect("clicked", func() {
		entry, _ := gtk.EntryNew()
		textLabsBox.Add(entry)
		labelEntrys = append(labelEntrys, entry)
		entry.ShowNow()
	})
	buttonDecreaseBtn.Connect("clicked", func() {
		removeIndex := len(labelEntrys)-1
		if removeIndex >= 0 {
			textLabsBox.Remove(labelEntrys[removeIndex])
			labelEntrys = append(labelEntrys[:removeIndex])
		}
	})
	beforeKeysEntry.Connect("focus_in_event", func() {
		go keyUpEvent(beforeKeysEntry)
	})
	beforeKeysEntry.Connect("focus_out_event", func() {
		robotgo.EventEnd()
	})
	endKeysEntry.Connect("focus_in_event", func() {
		go keyUpEvent(endKeysEntry)
	})
	endKeysEntry.Connect("focus_out_event", func() {
		robotgo.EventEnd()
	})
	changeUI(dateLimitSpinBtn, replyNumSpinBtn, editNumSpinBtn, strategyComboBoxText, beforeKeysEntry, endKeysEntry, randomCheckButton, averageCheckButton, beforeCheckButton, textLabsBox)
	win.ShowAll()
	setReplyNumBox()
	setBeforeBox()
}

func keyUpEvent(keysEntry *gtk.Entry) {
	evChan := robotgo.EventStart()
	for ev := range evChan {
		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyUp {
			if ev.Kind == hook.KeyUp {
				vsUp = append(vsUp, ev.Keycode)
			} else {
				vsHold = append(vsHold, ev.Keycode)
			}
			ks := make([]string, 0)
			for _, v := range vsUp {
				k := utils.KeyCode[v]
				ks = append(ks, k)
			}
			if len(vsHold) == len(vsUp) {
				utils.Reverse(ks)
				newKeys := strings.Join(ks, "    ")
				keysEntry.SetText(newKeys)
				vsHold = make([]uint16, 0)
				vsUp = make([]uint16, 0)
			}
		}
	}
}

func changeData(dateLimitSpinBtn *gtk.SpinButton, replyNumSpinBtn *gtk.SpinButton, editNumSpinBtn *gtk.SpinButton, strategyComboBoxText *gtk.ComboBoxText, beforeKeysEntry *gtk.Entry, endKeysEntry *gtk.Entry, randomCheckButton *gtk.CheckButton, averageCheckButton *gtk.CheckButton, beforeCheckButton *gtk.CheckButton, textLabEntrys []*gtk.Entry) bool {
	dateLimitStr, _ := dateLimitSpinBtn.GetText()
	replyNumStr, _ := replyNumSpinBtn.GetText()
	editNumStr, _ := editNumSpinBtn.GetText()
	beforeKeyStr, _ := beforeKeysEntry.GetText()
	endKeyStr, _ := endKeysEntry.GetText()
	dateLimit, _ := strconv.Atoi(dateLimitStr)
	replyNum, _ := strconv.Atoi(replyNumStr)
	editNum, _ := strconv.Atoi(editNumStr)
	strategy := strategyComboBoxText.GetActiveID()
	beforeKeys := strings.Fields(beforeKeyStr)
	endKeys := strings.Fields(endKeyStr)
	random := randomCheckButton.GetActive()
	withoutStop := withoutStopCheckButton.GetActive()
	average := averageCheckButton.GetActive()
	before := beforeCheckButton.GetActive()
	restart := false
	if editNum != utils.Settings.EditNum {
		restart = true
	}
	tagLen := len(utils.Settings.Tags)
	if tagLen != len(textLabEntrys) {
		restart = true
	}
	tags := make([]utils.Tag, 0)
	for pageIndex, textEntry := range textLabEntrys {
		text, _ := textEntry.GetText()
		if text != "" {
			var msgs []string
			if pageIndex + 1 < tagLen {
				msgs = utils.Settings.Tags[pageIndex].Msgs
				if utils.Settings.Tags[pageIndex].Label != text {
					restart = true
				}
			} else {
				msgs = make([]string, 0)
			}
			tags = append(tags, utils.Tag{text, msgs})
		}
	}
	utils.Settings = utils.Setting{dateLimit, replyNum, editNum, utils.Strategy(strategy), tags, beforeKeys, endKeys, random, withoutStop, average, before, runtime.GOOS}
	utils.SettingToFile()
	return restart
}

func changeUI(dateLimitSpinBtn *gtk.SpinButton, replyNumSpinBtn *gtk.SpinButton, editNumSpinBtn *gtk.SpinButton, strategyComboBoxText *gtk.ComboBoxText, beforeKeysEntry *gtk.Entry, endKeysEntry *gtk.Entry, randomCheckButton *gtk.CheckButton, averageCheckButton *gtk.CheckButton, beforeCheckButton *gtk.CheckButton, textLabsBox *gtk.Box) {
	dateLimitSpinBtn.SetValue(float64(utils.Settings.DateLimit))
	replyNumSpinBtn.SetValue(float64(utils.Settings.ReplyNum))
	editNumSpinBtn.SetValue(float64(utils.Settings.EditNum))
	strategyComboBoxText.SetActiveID(string(utils.Settings.Strategy))
	labelEntrys = make([]*gtk.Entry, 0)
	for _, tag := range utils.Settings.Tags {
		textLabEntry, _ := gtk.EntryNew()
		textLabEntry.SetText(tag.Label)
		textLabsBox.Add(textLabEntry)
		labelEntrys = append(labelEntrys, textLabEntry)
	}
	beforeKeysEntry.SetText(strings.Join(utils.Settings.BeforeKeys, "\t"))
	endKeysEntry.SetText(strings.Join(utils.Settings.EndKeys, "\t"))
	randomCheckButton.SetActive(utils.Settings.Random)
	withoutStopCheckButton.SetActive(utils.Settings.WithoutStop)
	averageCheckButton.SetActive(utils.Settings.Average)
	beforeCheckButton.SetActive(utils.Settings.Before)
}

func setReplyNumBox() {
	if withoutStopCheckButton.GetActive() {
		replyNumBox.Hide()
	} else {
		replyNumBox.Show()
	}
}

func setBeforeBox() {
	if beforeCheckButton.GetActive() {
		beforeKeysBox.Show()
	} else {
		beforeKeysBox.Hide()
	}
}