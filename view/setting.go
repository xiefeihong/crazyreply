package view

import (
	"github.com/go-vgo/robotgo"
	"github.com/gotk3/gotk3/gtk"
	hook "github.com/robotn/gohook"
	"github.com/xiefeihong/crazyreply/utils"
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
	builder, err := gtk.BuilderNewFromFile(utils.Root + "/view/ui/setting.glade")
	if err != nil {
		panic(err)
	}
	winObj, _ := builder.GetObject("win")
	win := winObj.(*gtk.Window)
	dateLimitObj, _ := builder.GetObject("date_limit")
	replyNumBoxObj, _ := builder.GetObject("reply_num_box")
	replyNumObj, _ := builder.GetObject("reply_num")
	editNumObj, _ := builder.GetObject("edit_num")
	beforeKeysBoxObj, _ := builder.GetObject("before_keys_box")
	beforeKeysObj, _ := builder.GetObject("before_keys")
	endKeysObj, _ := builder.GetObject("end_keys")
	randomObj, _ := builder.GetObject("random")
	withoutStopObj, _ := builder.GetObject("without_stop")
	averageObj, _ := builder.GetObject("average")
	beforeObj, _ := builder.GetObject("before")
	buttonDecreaseObj, _ := builder.GetObject("btn-")
	textLabsObj, _ := builder.GetObject("lab_texts")
	buttonIncreaseObj, _ := builder.GetObject("btn+")
	cancelObj, _ := builder.GetObject("btn_cancel")
	preserveObj, _ := builder.GetObject("btn_preserve")
	dateLimitSpinBtn := dateLimitObj.(*gtk.SpinButton)
	replyNumBox = replyNumBoxObj.(*gtk.Box)
	replyNumSpinBtn := replyNumObj.(*gtk.SpinButton)
	editNumSpinBtn := editNumObj.(*gtk.SpinButton)
	beforeKeysBox = beforeKeysBoxObj.(*gtk.Box)
	beforeKeysEntry := beforeKeysObj.(*gtk.Entry)
	endKeysEntry := endKeysObj.(*gtk.Entry)
	randomCheckButton := randomObj.(*gtk.CheckButton)
	withoutStopCheckButton = withoutStopObj.(*gtk.CheckButton)
	averageCheckButton := averageObj.(*gtk.CheckButton)
	beforeCheckButton = beforeObj.(*gtk.CheckButton)
	buttonDecreaseBtn := buttonDecreaseObj.(*gtk.Button)
	textLabsBox := textLabsObj.(*gtk.Box)
	buttonIncreaseBtn := buttonIncreaseObj.(*gtk.Button)
	cancelBtn := cancelObj.(*gtk.Button)
	preserveBtn := preserveObj.(*gtk.Button)
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
		restart := setSetting(dateLimitSpinBtn, replyNumSpinBtn, editNumSpinBtn, beforeKeysEntry, endKeysEntry, randomCheckButton, averageCheckButton, beforeCheckButton, labelEntrys)
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
	settingsToUI(dateLimitSpinBtn, replyNumSpinBtn, editNumSpinBtn, beforeKeysEntry, endKeysEntry, randomCheckButton, averageCheckButton, beforeCheckButton, textLabsBox)
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

func setSetting(dateLimitSpinBtn *gtk.SpinButton, replyNumSpinBtn *gtk.SpinButton, editNumSpinBtn *gtk.SpinButton, beforeKeysEntry *gtk.Entry, endKeysEntry *gtk.Entry, randomCheckButton *gtk.CheckButton, averageCheckButton *gtk.CheckButton, beforeCheckButton *gtk.CheckButton, textLabEntrys []*gtk.Entry) bool {
	dateLimitStr, _ := dateLimitSpinBtn.GetText()
	replyNumStr, _ := replyNumSpinBtn.GetText()
	editNumStr, _ := editNumSpinBtn.GetText()
	beforeKeyStr, _ := beforeKeysEntry.GetText()
	endKeyStr, _ := endKeysEntry.GetText()
	dateLimit, _ := strconv.Atoi(dateLimitStr)
	replyNum, _ := strconv.Atoi(replyNumStr)
	editNum, _ := strconv.Atoi(editNumStr)
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
	utils.Settings = utils.Setting{dateLimit, replyNum, editNum, tags, beforeKeys, endKeys, random, withoutStop, average, before}
	utils.SettingToFile()
	return restart
}

func settingsToUI(dateLimitSpinBtn *gtk.SpinButton, replyNumSpinBtn *gtk.SpinButton, editNumSpinBtn *gtk.SpinButton, beforeKeysEntry *gtk.Entry, endKeysEntry *gtk.Entry, randomCheckButton *gtk.CheckButton, averageCheckButton *gtk.CheckButton, beforeCheckButton *gtk.CheckButton, textLabsBox *gtk.Box) {
	dateLimitSpinBtn.SetValue(float64(utils.Settings.DateLimit))
	replyNumSpinBtn.SetValue(float64(utils.Settings.ReplyNum))
	editNumSpinBtn.SetValue(float64(utils.Settings.EditNum))
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