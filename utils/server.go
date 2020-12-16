package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
	"github.com/gotk3/gotk3/gtk"
	hook "github.com/robotn/gohook"
	"math/rand"
	"time"
)

type Tag struct {
	Label string `json:"label"`
	Msgs []string `json:"msgs"`
}

type Setting struct {
	DateLimit int `json:"date_limit"`
	ReplyNum int `json:"reply_num"`
	EditNum int `json:"edit_num"`
	Tags []Tag `json:"tags"`
	EndKeys []string `json:"end_keys"`
	Random bool `json:"random"`
	WithoutStop bool `json:"without_stop"`
	Persion bool `json:"persion"`
}

var (
	Root string
	Texts = make([][]*gtk.Entry, 0)
	PageIndex int
	BottonLabel string
	Settings Setting
	KeyCode = map[uint16]string{41:"`",2:"1",3:"2",4:"3",5:"4",6:"5",7:"6",8:"7",9:"8",10:"9",11:"0",12:"-",13:"+",
		16:"q",17:"w",18:"e",19:"r",20:"t",21:"y",22:"u",23:"i",24:"o",25:"p",26:"[",27:"]",43:"\\",30:"a",31:"s",32:"d",33:"f",34:"g",35:"h",36:"j",37:"k",38:"l",39:";",40:"'",44:"z",45:"x",46:"c",47:"v",48:"b",49:"n",50:"m",
		51:",",52:".",53:"/",59:"f1",60:"f2",61:"f3",62:"f4",63:"f5",64:"f6",65:"f7",66:"f8",67:"f9",68:"f10",69:"f11",70:"f12",
		1:"esc",14:"delete",15:"tab",29:"control",56:"alt",57:"space",42:"shift",54:"rshift",28:"enter",3675:"command",3676:"rcmd",3640:"ralt",57416:"up",57424:"down",57419:"left",57421:"right"}
)

func CarryReply(button *gtk.Button) {
	msgs := Settings.Tags[PageIndex].Msgs
	msgLen := len(msgs)
	setText(Texts, false)
	for ;BottonLabel == "结束"; {
		if Settings.Random {
			for i:= 0; BottonLabel == "结束" &&  i < Settings.ReplyNum * msgLen; i++ {
				reply(msgs[rand.Intn(msgLen)])
			}
		} else {
			for i:= 0 ; BottonLabel == "结束" && i < Settings.ReplyNum; i++ {
				for j:=0; j< msgLen; j++ {
					reply(msgs[j])
				}
			}
		}
		if !Settings.WithoutStop {
			break
		}
	}
	defer hook.End()
	setText(Texts, true)
	BottonLabel = "开始"
	button.SetLabel(BottonLabel)
}

func setText(texts [][]*gtk.Entry, disable bool) {
	for _, msgs := range texts {
		for _, entry := range msgs {
			entry.SetEditable(disable)
		}
	}
}

func reply(message string){
	var date int
	if !Settings.Persion {
		date = Settings.DateLimit
	} else {
		date = rand.Intn(100) * Settings.DateLimit
	}
	clipboard.WriteAll(message)
	t := time.Duration(date)
	time.Sleep(t * time.Millisecond)
	robotgo.KeyTap("v", "ctrl")
	time.Sleep(t * time.Millisecond)
	robotgo.KeyTap("enter")
	time.Sleep(t * time.Millisecond)
	fmt.Print(".")
}

func KeyDownEvent(keys []string){
	hook.Register(hook.KeyDown, keys, func(e hook.Event) {
		BottonLabel = "开始"
		fmt.Println(keys)
	})
	s := hook.Start()
	<-hook.Process(s)
}

func StartSettings() {
	defer func(){
		if err := recover(); err != nil {
			tags := []Tag{{"网友", make([]string, 0)}, {"弹幕", make([]string, 0)}}
			Settings = Setting{50, 10, 10, tags, []string{"control", "t"}, true, false, false}
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