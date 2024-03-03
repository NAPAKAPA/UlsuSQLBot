package utils

import (
	"github.com/bytedance/sonic"
	"math/rand"
	"time"
)

func ToJsonBytes(input interface{}) []byte {
	marshal, err := sonic.Marshal(input)
	if err != nil {
		return []byte(``)
	}
	return marshal
}

func RRange(a, b int) int {
	rand.Seed(time.Now().UnixNano())
	n := a + rand.Intn(b-a+1)
	return n
}
