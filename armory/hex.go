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

func BytesToHex(data []byte) string {
	var sb strings.Builder

	for i, b := range data {
		// 每个字节格式化成两位十六进制
		sb.WriteString(fmt.Sprintf("%02x", b))

		sb.WriteByte(' ')

		// 每 16 个字节后加换行
		if (i+1)%16 == 0 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
