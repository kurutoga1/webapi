package utils

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

// StringWithCharset は適当な文字を生成する。
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetRandomString は文字数を受け取り、ランダムな文字列を返す。
func GetRandomString(length int) string {
	return StringWithCharset(length, charset)
}

// GetNowTimeStringWithHyphen は現在時間の文字列を返す
// example -> 2021-11-09-21-51-25
func GetNowTimeStringWithHyphen() string {
	now := time.Now()

	s := now.String()

	n := s[0:10] + "-" + strings.Replace(s[11:19], ":", "-", -1)

	return n
}
