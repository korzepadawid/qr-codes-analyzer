package util

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

func MarshallBody(requestBody interface{}) io.Reader {
	marshal, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(marshal)
}
