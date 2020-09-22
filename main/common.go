package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"time"
)

//生成随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

//字节反转，大小端转换方法
func ReverseByte(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
func Int32ToLittleEndianHex(num int32) []byte {
	buff := new(bytes.Buffer)
	//binary.LittleEndian 小端模式
	err := binary.Write(buff, binary.LittleEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func Int32ToBigEndianHex(num int32) []byte {
	buff := new(bytes.Buffer)
	//binary.BigEndian 大端模式
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
