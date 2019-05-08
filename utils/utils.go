package utils

import (
	"crypto/rand"
	"log"
	"runtime"
	"strings"
)

// PanicOnError : Prints the error & exits the program
func PanicOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}

// GetCurrentFuncName : Return current calling function name as string
func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
}

// GenerateCryptoRandomBytes : Generate n-sized random bytes array
func GenerateCryptoRandomBytes(n int) ([]byte, error) {
	r := make([]byte, n)

	_, err := rand.Read(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// IsStringIn : Check wether the string is present in array
func IsStringIn(s string, list []string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}
	return false
}

// CaseInsensitiveContains : Case insensitive contains
func CaseInsensitiveContains(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}
