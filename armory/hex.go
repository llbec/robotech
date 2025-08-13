package armory

import (
	"fmt"
	"strconv"
	"strings"
)

func HexToBytes(hex string) ([]byte, error) {
	// 移除所有空格
	parts := strings.Fields(hex)
	var bytesArr []byte

	for _, part := range parts {
		if len(part) != 2 {
			return nil, fmt.Errorf("invalid hex byte: %s", part)
		}
		b, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return nil, err
		}
		bytesArr = append(bytesArr, byte(b))
	}
	return bytesArr, nil
}
