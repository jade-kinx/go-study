package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

func ToBinaryString(obj any) (string, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, obj)
	if err != nil {
		return "", nil
	}

	strs := []string{}
	for _, b := range buf.Bytes() {
		strs = append(strs, fmt.Sprintf("%08b", b))
	}
	return strings.Join(strs, "_"), nil
}
