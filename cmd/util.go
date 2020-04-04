package cmd

import (
	"crypto/md5"
	"fmt"
)

func HashPassword(password string) string {
	// always gives different value, so have to use md5
	//bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}
