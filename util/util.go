package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

func Encrypt(plaintext string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(plaintext))
	return fmt.Sprintf("%x", m.Sum(nil))
}
