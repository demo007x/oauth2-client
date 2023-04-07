package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

// GenerateBaseAuthorization Generate Basic Authorization
func GenerateBaseAuthorization(name, password string) string {
	var buf = bytes.NewBufferString(fmt.Sprintf("%s:%s", name, password))
	var encode = base64.StdEncoding.EncodeToString(buf.Bytes())
	return "Basic " + encode
}

// GenerateBearAuthorization Generate Bearer Authorization
func GenerateBearAuthorization(accessToken string) string {
	buf := bytes.NewBuffer([]byte("Bearer "))

	if _, err := buf.Write([]byte(accessToken)); err != nil {
		fmt.Println(err.Error())
	}

	return buf.String()
}
