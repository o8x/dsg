package utils

import (
	"crypto/sha1"
	"fmt"
)

func Sha1Sum(str string) string {
	sum := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", sum)
}
