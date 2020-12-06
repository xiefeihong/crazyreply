package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
	"github.com/gotk3/gotk3/gtk"
	"io"
	"math/rand"
	"os"
	"time"
)

type Setting struct {
	DateLimit int `json:"date_limit"`
	ReplyNum int `json:"reply_num"`
	EditNum int `json:"edit_num"`
	Tags map[string][]string `json:"tags"`
	EndKeys []string `json:"end_keys"`
	Random bool `json:"random"`
	Persion bool `json:"persion"`
}

var Texts = make(map[string][]*gtk.Entry, 0)
var PageLabel string
var BottonLabel string
var Settings Setting

func CarryReply(button *gtk.Button) {
	messages := Settings.Tags[PageLabel]
	msgLen := len(messages)
	for _, msg := range Texts[PageLabel] {
		msg.SetEditable(false)
	}
	if Settings.Random {
		for i:= 0; BottonLabel == "结束" &&  i < Settings.ReplyNum * msgLen; i++ {
			reply(messages[rand.Intn(msgLen)])
		}
	} else {
		for i:= 0 ; BottonLabel == "结束" && i < Settings.ReplyNum; i++ {
			for j:=0; j< msgLen; j++ {
				reply(messages[j])
			}
		}
	}
	for _, msg := range Texts[PageLabel] {
		msg.SetEditable(true)
	}
	BottonLabel = "开始"
	button.SetLabel(BottonLabel)
}

func reply(message string){
	var date int
	if Settings.Persion {
		date = rand.Intn(100) * Settings.DateLimit
	} else {
		date = Settings.DateLimit
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

func KeyEvent(keys []string){
	l := len(keys)
	var ok bool
	if l == 2 {
		ok = robotgo.AddEvents(keys[0], keys[1])
	} else if l == 3 {
		ok = robotgo.AddEvents(keys[0], keys[1], keys[2])
	} else {
		fmt.Println("不支持的组合键")
	}
	if ok {
		BottonLabel = "开始"
		fmt.Println("监测到按下了退出组合键")
	}
}

func StartSettings() {
	file, e := os.Open("config.json")
	if e != nil {
		tags := map[string][]string{"网友":make([]string, 0), "弹幕":make([]string, 0)}
		Settings = Setting{50, 10, 10, tags, []string{"ctrl", "t"}, true, false}
	} else {
		reader := bufio.NewReader(file)
		buf := make([]byte, 1024)
		out := make([]byte, 1024)
		for {
			n, e2 := reader.Read(buf)
			if e2 == io.EOF {
				break
			}
			out = append(out, buf[:n]...)
		}
		out = bytes.Trim(out,"\x00")
		e3 := json.Unmarshal(out, &Settings)
		if e3 != nil {
			panic(e3)
		}
		for label := range Settings.Tags {
			t := Settings.Tags[label]
			tag := make([]string, 0)
			copy(tag, t)
		}
	}
	defer file.Close()
}

func SettingToFile(){
	config, _ := json.Marshal(Settings)
	file, err := os.OpenFile("config.json", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		var e2 error
		file, e2 = os.OpenFile("config.json", os.O_WRONLY | os.O_CREATE, 0666)
		if e2 != nil {
			panic(e2)
		}
	}
	defer file.Close()
	var str bytes.Buffer
	json.Indent(&str, config,"", "\t")
	file.Write(str.Bytes())
}