package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// Clone 完整复制数据
func clone(a, b interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(a); err != nil {
		return err
	}
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}

type ReceiveMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func main() {
	data := &ReceiveMsg{
		Type: "login",
		Data: "Nf+bAwEBClZlcmlmeUluZm8B/5wAAQMBAklkAQQAAQdBY2NvdW50AQwAAQVUb2tlbgEMAAAAF/+cAf4DBgEEbWFzawEIWa/FVpjowpAA",
	}
	bs, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("bytes:%s\n", bs)
}
