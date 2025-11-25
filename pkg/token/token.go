package token

import (
	"crypto/rand"
	"encoding/hex"
)

func CreateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
