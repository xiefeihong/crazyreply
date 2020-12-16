package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func Reverse(list []string){
	var temp string
	kl := len(list)
	for i := 0; i < kl / 2; i ++ {
		temp = list[i]
		list[i] = list[kl - i - 1]
		list[kl -i - 1] = temp
	}
}

func WriteStringToFile(filename string, s string){
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		var e2 error
		file, e2 = os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
		if e2 != nil {
			panic(e2)
		}
	}
	defer file.Close()
	file.WriteString(s)
}

func ReadBytesToFile(filename string) []byte {
	file, e := os.Open(filename)
	if e != nil {
		panic(e)
	}
	defer file.Close()
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
	return out
}