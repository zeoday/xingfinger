// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"encoding/json"
	"io/ioutil"
)

// Fingerprints 指纹规则库结构体
type Fingerprints struct {
	Fingerprint []Fingerprint
}

// Fingerprint 单条指纹规则
type Fingerprint struct {
	Cms      string
	Method   string
	Location string
	Keyword  []string
}

// 全局指纹规则库实例
var fingerprints *Fingerprints

// LoadFingerprints 从文件加载指纹规则库
func LoadFingerprints(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Fingerprints
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	fingerprints = &config
	return nil
}

// GetFingerprints 获取全局指纹规则库实例
func GetFingerprints() *Fingerprints {
	return fingerprints
}
