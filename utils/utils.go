package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"
)

func Now() int {
	return int(time.Now().Unix())
}

func EncodeHexMD5(params string) string {
	sumString := md5.Sum([]byte(params))
	return hex.EncodeToString(sumString[:])
}

func SplitName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName) // split by whitespace

	if len(parts) == 0 {
		return "", ""
	}

	firstName = parts[0]

	if len(parts) == 1 {
		lastName = ""
	} else {
		lastName = strings.Join(parts[1:], " ")
	}

	return firstName, lastName
}
