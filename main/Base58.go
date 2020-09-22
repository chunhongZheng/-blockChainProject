package main

import (
	"bytes"
	"math/big"
)

////base58和base64一样是一种二进制转可视字符串的算法，主要用来转换大整数值。区别是，转换出来的字符串，去除了几个看起来会产生歧义的字符，
//// 如 0 (零), O (大写字母O), I (大写的字母i) and l (小写的字母L) ，和几个影响双击选择的字符，如/, +。
////
////结果字符集正好58个字符(包括9个数字，24个大写字母，25个小写字母)。
////base58的go实现   base58编码用于生成 比特币地址数据

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

////加密
func Base58Encode(input []byte) []byte {
	var result []byte                           //声明一个slice类型变量  切片
	x := big.NewInt(0).SetBytes(input)          //大数据类型
	base := big.NewInt(int64(len(b58Alphabet))) //相当于多少进制数
	zero := big.NewInt(0)
	mod := &big.Int{} //大整数的指针
	for x.Cmp(zero) != 0 {
		//x是商,mod是余数    10/2=5 ....0     5/2=2 ....1    2/2=1....0    1/2=0....1          1010
		x.DivMod(x, base, mod) //对x取余数
		result = append(result, b58Alphabet[mod.Int64()])
		//println(result)
	}
	//反转结果
	ReverseByte(result)
	//把这个字节前置的0 字节置换成第一个字符
	for _, b := range input {
		if b == 0x00 {

			result = append([]byte{b58Alphabet[0]}, result...)

		} else {
			break
		}
	}
	return result
}

//解密
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0 //统计前面为0的个数，方便截取解码
	for _, b := range input {
		if b != b58Alphabet[0] {
			break
		}
		zeroBytes++
	}
	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b) //找出b所在的索引位置，即是余数
		result.Mul(result, big.NewInt(int64(len(b58Alphabet))))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}
