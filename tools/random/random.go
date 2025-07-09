package random

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	mathrand "math/rand"
	"time"
)

const defaultLens = 30

var runes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func Str(lens int) string {
	if lens <= 0 {
		lens = defaultLens
	}
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	b := make([]byte, lens)
	for i := range b {
		b[i] = runes[r.Intn(len(runes))]
	}
	return string(b)
}

func Number(start, end int) int {
	if end < start {
		t := end
		end = start
		start = t
	}
	return mathrand.Intn(end-start) + start // (end-start)+start
}

// RandSlice 随机弹出切片中的元素
func RandSlice[T any](slice []T) (T, int) {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	if len(slice) == 0 {
		panic("empty slice")
	}
	randIndex := r.Intn(len(slice))
	return slice[randIndex], randIndex
}

// === 安全Salt生成功能（使用crypto/rand） ===

// GenerateSalt 生成指定长度的随机salt（Base64编码）
// 推荐用于密码加密，使用加密级别的随机数生成器
func GenerateSalt(byteLength int) (string, error) {
	if byteLength <= 0 {
		byteLength = 16 // 默认16字节
	}
	
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// GenerateSaltHex 生成指定长度的随机salt（Hex编码）
func GenerateSaltHex(byteLength int) (string, error) {
	if byteLength <= 0 {
		byteLength = 16 // 默认16字节
	}
	
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateSecureString 生成指定长度的安全随机字符串（自定义字符集）
// 使用crypto/rand，比Str()更安全，适合密码、token等场景
func GenerateSecureString(length int) (string, error) {
	if length <= 0 {
		length = defaultLens
	}
	
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	
	return string(result), nil
}

// GenerateURLSafeSalt 生成URL安全的随机salt
func GenerateURLSafeSalt(byteLength int) (string, error) {
	if byteLength <= 0 {
		byteLength = 16 // 默认16字节
	}
	
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
