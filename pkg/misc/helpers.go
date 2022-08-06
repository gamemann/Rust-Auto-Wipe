package misc

import (
	"bytes"
	"fmt"
)

func CreateKeyPairs(m map[string]interface{}) string {
	b := new(bytes.Buffer)

	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}

	return b.String()
}
