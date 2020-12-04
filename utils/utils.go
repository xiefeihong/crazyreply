package utils

import (
	"bufio"
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
	Tags []string `json:"tags"`
	EndKeys []string `json:"end_keys"`
	Random bool `json:"random"`
	Persion bool `json:"persion"`
}

var Texts = make([]*gtk.Entry, 0)
var Label string
var Settings Setting

func CarryReply(button *gtk.Button) {
	messages := getMessages()
	for i:= 0; i<Settings.EditNum; i++ {
		Texts[i].SetEditable(false)
	}
	if Settings.Random {
		for i:= 0; Label == "结束" &&  i < Settings.ReplyNum * Settings.EditNum; i++ {
			reply(messages[rand.Intn(Settings.EditNum)])
		}
	} else {
		for i:= 0 ;Label == "结束" && i < Settings.ReplyNum; i++ {
			for j:=0; j< Settings.EditNum; j++ {
				reply(messages[j])
			}
		}
	}
	for i:= 0; i<Settings.EditNum; i++ {
		Texts[i].SetEditable(true)
	}
	Label = "开始"
	button.SetLabel(Label)
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
	fmt.Printf("耗时：%d毫秒\n", (endTime - startTime)/1000000)
}

func getMessages() []string {
	messages := make([]string, Settings.EditNum)
	for i:= 0; i<Settings.EditNum; i++ {
		messages[i], _ = Texts[i].GetText()
	}
	return messages
}

func KeyEvent(keys []string){
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
		Label = "开始"
		fmt.Println("监测到按下了退出组合键")
	}
}

func StartSettings() {
	file, e := os.Open("config.json")
	if e != nil {
		fmt.Println("打开文件失败：", e)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024)
	for {
		n, e2 := reader.Read(buf)
		if e2 == io.EOF {
			break
		}
		e3 := json.Unmarshal(buf[:n], &Settings)
		if e3 != nil {
			fmt.Println("json解析出错： ", e3)
		}
	}
}