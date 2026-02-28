// Package validate
package validate

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// IsValidYAML 验证字符串是否为有效的YAML格式
func IsValidYAML(yamlStr string) bool {
	if strings.TrimSpace(yamlStr) == "" {
		return true // 空字符串被认为是有效的
	}

	var temp any
	err := yaml.Unmarshal([]byte(yamlStr), &temp)
	return err == nil
}

// ValidateYAML 验证YAML格式并返回错误信息
func ValidateYAML(yamlStr string) error {
	if strings.TrimSpace(yamlStr) == "" {
		return nil // 空字符串被认为是有效的
	}

	var temp any
	err := yaml.Unmarshal([]byte(yamlStr), &temp)
	return err
}

// ParseYAML 解析YAML字符串为interface{}
func ParseYAML(yamlStr string) (any, error) {
	var result any
	err := yaml.Unmarshal([]byte(yamlStr), &result)
	return result, err
}

// ToYAML 将interface{}转换为YAML字符串
func ToYAML(data any) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
