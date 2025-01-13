package util

import (
	"math/rand"
	"time"
)

type (
	Random struct {
		charset Charset
	}

	Charset string
)

const (
	Alphanumeric Charset = Alphabetic + Numeric
	Alphabetic   Charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric      Charset = "0123456789"
	Hex          Charset = Numeric + "abcdef"
)

var (
	globalRandom = NewRandom()
)

func NewRandom() *Random {
	rand.Seed(time.Now().UnixNano())
	return &Random{
		charset: Alphanumeric,
	}
}

func (r *Random) SetCharset(charset Charset) {
	r.charset = charset
}

func (r *Random) String(length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = r.charset[rand.Int63()%int64(len(r.charset))]
	}
	return string(b)
}

func SetCharset(charset Charset) {
	globalRandom.SetCharset(charset)
}

func RandomString(length uint8) string {
	return globalRandom.String(length)
}

func RandNumStr(length uint8) string {
	r := NewRandom()
	r.SetCharset(Numeric)
	return r.String(length)
}
