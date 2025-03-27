package main

import (
	"bufio"
	"bytes"
	"strings"
)

// parseProperties 将属性文件内容解析为 Map（忽略注释和空行）
func parseProperties(content []byte) (map[string]string, error) {
	props := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释（以 # 或 ! 开头）
		if len(line) == 0 || line[0] == '#' || line[0] == '!' {
			continue
		}

		// 分割键值对（支持 = 或 : 作为分隔符）
		sepIndex := strings.IndexAny(line, "=:")
		if sepIndex == -1 {
			continue // 忽略无效行
		}

		key := strings.TrimSpace(line[:sepIndex])
		value := strings.TrimSpace(line[sepIndex+1:])
		props[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return props, nil
}

func compareProperties(props1, props2 map[string]string) bool {
	for key, value1 := range props1 {
		value2, exists := props2[key]
		if !exists || value1 != value2 {
			return false
		}
	}
	return true
}
