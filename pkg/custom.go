// Package pkg 提供 xingfinger 的核心功能
// 本文件实现自定义指纹文件加载功能
package pkg

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainreactors/fingers/resources"
)

// CustomFingerConfig 自定义指纹配置
type CustomFingerConfig struct {
	EHole       string // EHole 格式指纹文件路径
	Goby        string // Goby 格式指纹文件路径
	Wappalyzer  string // Wappalyzer 格式指纹文件路径
	Fingers     string // Fingers 原生格式指纹文件路径
	FingerPrint string // FingerPrintHub 格式指纹文件路径
	ARL         string // ARL YAML 格式指纹文件路径
	NoDefault   bool   // 禁用默认指纹
}

// LoadCustomFingerprints 加载自定义指纹文件
// 根据配置加载用户指定的指纹文件
// 自定义指纹默认与内置指纹叠加使用，除非指定 NoDefault
//
// 参数：
//   - config: 自定义指纹配置
//   - silent: 是否静默模式（不输出加载信息）
//
// 返回：
//   - error: 加载错误
func LoadCustomFingerprints(config *CustomFingerConfig, silent bool) error {
	// 如果指定了 NoDefault，清空所有内置指纹
	if config.NoDefault {
		resources.EholeData = []byte{}
		resources.GobyData = []byte{}
		resources.WappalyzerData = []byte{}
		resources.FingersHTTPData = []byte{}
		resources.FingerprinthubWebData = []byte{}
		if !silent {
			fmt.Println("[*] 已禁用默认指纹")
		}
	}

	if config.EHole != "" {
		data, err := loadFingerFile(config.EHole)
		if err != nil {
			return fmt.Errorf("加载 EHole 指纹失败: %v", err)
		}
		resources.EholeData = data
		if !silent {
			fmt.Printf("[*] 已加载自定义 EHole 指纹: %s\n", config.EHole)
		}
	}

	if config.Goby != "" {
		data, err := loadFingerFile(config.Goby)
		if err != nil {
			return fmt.Errorf("加载 Goby 指纹失败: %v", err)
		}
		resources.GobyData = data
		if !silent {
			fmt.Printf("[*] 已加载自定义 Goby 指纹: %s\n", config.Goby)
		}
	}

	if config.Wappalyzer != "" {
		data, err := loadFingerFile(config.Wappalyzer)
		if err != nil {
			return fmt.Errorf("加载 Wappalyzer 指纹失败: %v", err)
		}
		resources.WappalyzerData = data
		if !silent {
			fmt.Printf("[*] 已加载自定义 Wappalyzer 指纹: %s\n", config.Wappalyzer)
		}
	}

	if config.Fingers != "" {
		data, err := loadFingerFile(config.Fingers)
		if err != nil {
			return fmt.Errorf("加载 Fingers 指纹失败: %v", err)
		}
		resources.FingersHTTPData = data
		if !silent {
			fmt.Printf("[*] 已加载自定义 Fingers 指纹: %s\n", config.Fingers)
		}
	}

	if config.FingerPrint != "" {
		data, err := loadFingerFile(config.FingerPrint)
		if err != nil {
			return fmt.Errorf("加载 FingerPrintHub 指纹失败: %v", err)
		}
		resources.FingerprinthubWebData = data
		if !silent {
			fmt.Printf("[*] 已加载自定义 FingerPrintHub 指纹: %s\n", config.FingerPrint)
		}
	}

	// ARL 使用独立引擎，在 scanner.go 中初始化

	return nil
}

// loadFingerFile 加载指纹文件
// 支持 .json 和 .json.gz 格式
//
// 参数：
//   - path: 文件路径
//
// 返回：
//   - []byte: 文件内容（如果是 .json 会自动压缩为 gzip）
//   - error: 读取错误
func loadFingerFile(path string) ([]byte, error) {
	// 读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(path))

	// 如果是 .gz 文件，直接返回
	if ext == ".gz" {
		return data, nil
	}

	// 如果是 .json 文件，需要压缩为 gzip 格式
	// 因为 fingers 库期望的是 gzip 压缩的数据
	if ext == ".json" {
		return gzipCompress(data)
	}

	return data, nil
}

// gzipCompress 将数据压缩为 gzip 格式
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// gzipDecompress 解压 gzip 数据（用于调试）
func gzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
