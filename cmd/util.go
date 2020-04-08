package cmd

import (
	"crypto/md5"
	"fmt"
)

func HashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}
