package view

import (
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
	end bool
)

func ShowSetting() {
	builder, err := gtk.BuilderNewFromFile(utils.Root + "/view/ui/setting.glade")
	if err != nil {
		panic(err)
	}
	winObj, _ := builder.GetObject("win")
	win := winObj.(*gtk.Window)
	dateLimitObj, _ := builder.GetObject("date_limit")
	replyNumObj, _ := builder.GetObject("reply_num")
	editNumObj, _ := builder.GetObject("edit_num")
	endKeysObj, _ := builder.GetObject("end_keys")
	randomObj, _ := builder.GetObject("random")
	withoutStopObj, _ := builder.GetObject("without_stop")
	persionObj, _ := builder.GetObject("persion")
	buttonDecreaseObj, _ := builder.GetObject("btn-")
	textLabsObj, _ := builder.GetObject("lab-texts")
	buttonIncreaseObj, _ := builder.GetObject("btn+")
	cancelObj, _ := builder.GetObject("btn-cancel")
	preserveObj, _ := builder.GetObject("btn-preserve")
	dateLimitSpinBtn := dateLimitObj.(*gtk.SpinButton)
	replyNumSpinBtn := replyNumObj.(*gtk.SpinButton)
	editNumSpinBtn := editNumObj.(*gtk.SpinButton)
	endKeysEntry := endKeysObj.(*gtk.Entry)
	randomSwitch := randomObj.(*gtk.Switch)
	withoutStopSwitch := withoutStopObj.(*gtk.Switch)
	persionSwitch := persionObj.(*gtk.Switch)
	buttonDecreaseBtn := buttonDecreaseObj.(*gtk.Button)
	textLabsBox := textLabsObj.(*gtk.Box)
	buttonIncrease := buttonIncreaseObj.(*gtk.Button)
	cancelBtn := cancelObj.(*gtk.Button)
	preserveBtn := preserveObj.(*gtk.Button)
	cancelBtn.Connect("clicked", func() {
		win.Close()
	})
	preserveBtn.Connect("clicked", func() {
		restart := setSetting(dateLimitSpinBtn, replyNumSpinBtn, editNumSpinBtn, labelEntrys, endKeysEntry, randomSwitch, withoutStopSwitch, persionSwitch)
		if restart {
			dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", "请重新打开此程序")
			dialog.Run()
			dialog.Destroy()
		}
		win.Close()
	})
	buttonIncrease.Connect("clicked", func() {
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
	endKeysEntry.Connect("focus_in_event", func() {
		end = false
	})
	endKeysEntry.Connect("focus_out_event", func() {
		end = true
	})
	win.Connect("destroy", func() {
		defer hook.End()
	})
	settingsToUI(dateLimitSpinBtn, replyNumSpinBtn, editNumSpinBtn, textLabsBox, endKeysEntry, randomSwitch, withoutStopSwitch, persionSwitch)
	win.ShowAll()
	end = true
	go keyUpEvent(endKeysEntry)
}

func keyUpEvent(endKeysEntry *gtk.Entry) {
	EvChan := hook.Start()
	for ev := range EvChan {
		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyUp {
			if !end {
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
					endKeysEntry.SetText(newKeys)
					vsHold = make([]uint16, 0)
					vsUp = make([]uint16, 0)
				}
			}
		}
	}
}

func setSetting(dateLimitSpinBtn *gtk.SpinButton, replyNumSpinBtn *gtk.SpinButton, editNumSpinBtn *gtk.SpinButton, textLabEntrys []*gtk.Entry,
		endKeysEntry *gtk.Entry, randomSwitch *gtk.Switch, withoutStopSwitch *gtk.Switch, persionSwitch *gtk.Switch) bool {
	dateLimitStr, _ := dateLimitSpinBtn.GetText()
	replyNumStr, _ := replyNumSpinBtn.GetText()
	editNumStr, _ := editNumSpinBtn.GetText()
	endKeyStr, _ := endKeysEntry.GetText()
	dateLimit, _ := strconv.Atoi(dateLimitStr[:len(dateLimitStr)-3])
	replyNum, _ := strconv.Atoi(replyNumStr[:len(replyNumStr)-3])
	editNum, _ := strconv.Atoi(editNumStr[:len(editNumStr)-3])
	endKeys := strings.Fields(endKeyStr)
	random := randomSwitch.GetActive()
	withoutStop := withoutStopSwitch.GetActive()
	persion := persionSwitch.GetActive()
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
				if  utils.Settings.Tags[pageIndex].Label != text {
					restart = true
				}
			} else {
				msgs = make([]string, 0)
			}
			tags = append(tags, utils.Tag{text, msgs})
		}
	}
	utils.Settings = utils.Setting{dateLimit, replyNum, editNum, tags, endKeys, random, withoutStop, persion}
	utils.SettingToFile()
	return restart
}

func settingsToUI(dateLimitSpinBtn *gtk.SpinButton, replyNumSpinBtn *gtk.SpinButton, editNumSpinBtn *gtk.SpinButton, textLabsBox *gtk.Box,
		endKeysEntry *gtk.Entry, randomSwitch *gtk.Switch, withoutStopSwitch *gtk.Switch, persionSwitch *gtk.Switch) {
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
	endKeysEntry.SetText(strings.Join(utils.Settings.EndKeys, "\t"))
	randomSwitch.SetActive(utils.Settings.Random)
	withoutStopSwitch.SetActive(utils.Settings.WithoutStop)
	persionSwitch.SetActive(utils.Settings.Persion)
}