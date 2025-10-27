package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// hex 字符串（带空格） -> byte 数组
func HexStringToBytes(s string) ([]byte, error) {
	parts := strings.Fields(s)
	res := make([]byte, len(parts))
	for i, p := range parts {
		v, err := strconv.ParseUint(p, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析 hex %s 失败: %v", p, err)
		}
		res[i] = byte(v)
	}
	return res, nil
}

// byte 数组 -> 大写 hex，每16个换行
func BytesToHexString(data []byte) string {
	var sb strings.Builder
	for i, b := range data {
		sb.WriteString(fmt.Sprintf("%02X", b))
		if i != len(data)-1 {
			sb.WriteByte(' ')
		}
		if (i+1)%16 == 0 && i != len(data)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
