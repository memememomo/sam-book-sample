package utils

import (
	"crypto/rand"
	"fmt"
	"strconv"
)

func GenerateToken(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	l := length / 2
	if length%2 == 1 {
		l++
	}
	return fmt.Sprintf("%x", buf[0:l]), nil
}

func ParseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
