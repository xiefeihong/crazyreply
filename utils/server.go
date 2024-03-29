package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
	"github.com/gotk3/gotk3/gtk"
	hook "github.com/robotn/gohook"
	"math/rand"
	"runtime"
)

type Tag struct {
	Label string `json:"label"`
	Msgs []string `json:"msgs"`
}

type Strategy string
const(
	Key Strategy = "key"
	Clipboard Strategy = "clipboard"
)

type Setting struct {
	DateLimit int `json:"date_limit"`
	ReplyNum int `json:"reply_num"`
	EditNum int `json:"edit_num"`
	Strategy Strategy `json:"strategy"`
	Tags []Tag `json:"tags"`
	BeforeKeys []string `json:"before_keys"`
	EndKeys []string `json:"end_keys"`
	Random bool `json:"random"`
	WithoutStop bool `json:"without_stop"`
	Average bool `json:"average"`
	Before bool `json:"before"`
	OS string `json:"os"`
}

var (
	Root string
	Texts = make([][]*gtk.Entry, 0)
	Book *gtk.Notebook
	BottonLabel string
	Settings Setting
	ReplyIndex int
	MsgIndex int
	KeyCode = map[uint16]string{41:"`",2:"1",3:"2",4:"3",5:"4",6:"5",7:"6",8:"7",9:"8",10:"9",11:"0",12:"-",13:"+",
		16:"q",17:"w",18:"e",19:"r",20:"t",21:"y",22:"u",23:"i",24:"o",25:"p",26:"[",27:"]",43:"\\",30:"a",31:"s",32:"d",33:"f",34:"g",35:"h",36:"j",37:"k",38:"l",39:";",40:"'",44:"z",45:"x",46:"c",47:"v",48:"b",49:"n",50:"m",
		51:",",52:".",53:"/",59:"f1",60:"f2",61:"f3",62:"f4",63:"f5",64:"f6",65:"f7",66:"f8",67:"f9",68:"f10",69:"f11",70:"f12",
		1:"esc",14:"delete",15:"tab",29:"control",56:"alt",57:"space",42:"shift",54:"rshift",28:"enter",3675:"command",3676:"rcmd",3640:"ralt",57416:"up",57424:"down",57419:"left",57421:"right"}
)

func CarryReply(button *gtk.Button) {
	setTextState(Texts, false)
	msgs := Settings.Tags[Book.GetCurrentPage()].Msgs
	if Settings.WithoutStop {
		for ;BottonLabel == "结束"; {
			replys(msgs)
		}
	} else {
		replys(msgs)
	}
	robotgo.EventEnd()
	setTextState(Texts, true)
	BottonLabel = "开始"
	button.SetLabel(BottonLabel)
}

func setTextState(texts [][]*gtk.Entry, disable bool) {
	for _, msgs := range texts {
		for _, entry := range msgs {
			entry.SetEditable(disable)
		}
	}
}

func replys(msgs []string)  {
	if Settings.Random {
		for ReplyIndex= 0; BottonLabel == "结束" &&  ReplyIndex < Settings.ReplyNum * len(msgs); ReplyIndex++ {
			msgs = Settings.Tags[Book.GetCurrentPage()].Msgs
			if len(msgs) > 0 {
				reply(msgs[rand.Intn(len(msgs))])
			}
		}
	} else {
		for ReplyIndex= 0 ; BottonLabel == "结束" && ReplyIndex < Settings.ReplyNum; ReplyIndex++ {
			for MsgIndex=0; MsgIndex< len(msgs); MsgIndex++ {
				msgs = Settings.Tags[Book.GetCurrentPage()].Msgs
				if len(msgs) > 0 {
					reply(msgs[MsgIndex])
				}
			}
		}
	}
}

func reply(message string){
	date := Settings.DateLimit
	space := 30
	if Settings.Average {
		date = rand.Intn(Settings.DateLimit) + Settings.DateLimit / 2
	}
	if Settings.Before {
		keys := Settings.BeforeKeys
		keyLen := len(keys)
		switch keyLen {
		case 0:
			break
		case 1:
			robotgo.KeyTap(keys[0])
			break
		default:
			robotgo.KeyTap(keys[keyLen - 1], keys[:keyLen - 1])
			break
		}
		robotgo.MilliSleep(space)
	}
	if Settings.Strategy == Key {
		robotgo.TypeStr(message)
	} else if Settings.Strategy == Clipboard {
		clipboard.WriteAll(message)
		robotgo.MilliSleep(space)
		robotgo.KeyTap("v", "ctrl")
	}
	robotgo.MilliSleep(space)
	robotgo.KeyTap("enter")
	robotgo.MilliSleep(date)
}

func KeyDownEvent(keys []string){
	robotgo.EventHook(hook.KeyDown, keys, func(e hook.Event) {
		BottonLabel = "开始"
	})
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func StartSettings() {
	defer func(){
		if err := recover(); err != nil {
			tags := []Tag{{"LOL", make([]string, 0)}, {"弹幕", make([]string, 0)}, {"网友", make([]string, 0)}}
			Settings = Setting{50, 10, 10, Clipboard, tags, nil, []string{"control", "t"}, true, false, false, false, runtime.GOOS}
		}
	}()
	out := ReadBytesToFile(Root + "/config.json")
	e := json.Unmarshal(out, &Settings)
	if e != nil {
		panic(e)
	}
}

func SettingToFile(){
	config, _ := json.Marshal(Settings)
	var str bytes.Buffer
	json.Indent(&str, config,"", "\t")
	WriteStringToFile(Root + "/config.json", str.String())
}