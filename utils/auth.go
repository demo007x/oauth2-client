package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

// GenerateBaseAuth Generate Http Base Auth
func GenerateBaseAuth(name, password string) string {
	var buf = bytes.NewBufferString(fmt.Sprintf("%s:%s", name, password))
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
