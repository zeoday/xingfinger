// Package pkg 提供 xingfinger 的核心功能
// 本文件负责扫描结果的输出和保存
package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// saveResults 保存扫描结果到 JSON 文件
//
// 参数：
//   - filename: 输出文件路径
//   - results: 指纹识别结果切片
func saveResults(filename string, results []Result) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".json" {
		fmt.Println("[!] Only JSON format is supported")
		return
	}
	saveJSON(filename, results)
}

// saveJSON 将结果保存为 JSON 格式文件
// JSON 格式带有缩进，便于阅读和后续程序处理
//
// 参数：
//   - filename: 输出文件路径
//   - results: 指纹识别结果切片
func saveJSON(filename string, results []Result) {
	// 使用缩进格式化 JSON，提高可读性
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("[!] JSON error:", err)
		return
	}

	// 创建输出文件
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("[!] Create error:", err)
		return
	}
	defer f.Close()

	// 写入 JSON 数据
	f.Write(data)
}
