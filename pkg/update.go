// Package pkg 提供 xingfinger 的核心功能
// 本文件实现指纹库更新功能，从 GitHub 下载最新的指纹文件
package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chainreactors/fingers"
	"github.com/chainreactors/fingers/resources"
	"github.com/gookit/color"
)

// 指纹文件配置
// 映射引擎名称到对应的指纹文件名
var FingerConfigs = map[string]string{
	fingers.FingersEngine:     "fingers_http.json.gz",
	fingers.FingerPrintEngine: "fingerprinthub_web.json.gz",
	fingers.WappalyzerEngine:  "wappalyzer.json.gz",
	fingers.EHoleEngine:       "ehole.json.gz",
	fingers.GobyEngine:        "goby.json.gz",
}

// 指纹文件下载基础 URL
const baseURL = "https://raw.githubusercontent.com/chainreactors/fingers/master/resources/"

// 默认指纹存储目录
const DefaultFingerPath = "fingers"

// UpdateFingerprints 更新所有指纹库
// 从 GitHub 下载最新的指纹文件到本地目录
//
// 返回：
//   - error: 更新过程中的错误
func UpdateFingerprints() error {
	fmt.Println("[*] 开始更新指纹库...")

	// 确保指纹目录存在
	fingerPath := getFingerPath()
	if err := os.MkdirAll(fingerPath, 0755); err != nil {
		return fmt.Errorf("创建指纹目录失败: %v", err)
	}

	// 统计更新结果
	updated := 0
	failed := 0

	// 遍历所有指纹引擎进行更新
	for engineName, fileName := range FingerConfigs {
		ok, err := downloadFingerConfig(engineName, fileName, fingerPath)
		if err != nil {
			color.Red.Printf("[!] 更新 %s 失败: %v\n", engineName, err)
			failed++
			continue
		}
		if ok {
			updated++
		}
	}

	// 输出更新结果
	fmt.Println()
	if updated > 0 {
		color.Green.Printf("[+] 成功更新 %d 个指纹库\n", updated)
	}
	if failed > 0 {
		color.Red.Printf("[!] %d 个指纹库更新失败\n", failed)
	}
	if updated == 0 && failed == 0 {
		color.Green.Println("[+] 所有指纹库已是最新版本")
	}

	return nil
}

// LoadLocalFingerprints 加载本地指纹文件
// 如果本地存在更新的指纹文件，则替换内嵌的指纹数据
//
// 返回：
//   - error: 加载过程中的错误
func LoadLocalFingerprints() error {
	fingerPath := getFingerPath()

	for engineName, fileName := range FingerConfigs {
		filePath := filepath.Join(fingerPath, fileName)

		// 读取本地指纹文件
		content, err := os.ReadFile(filePath)
		if err != nil {
			// 文件不存在，使用内嵌指纹
			continue
		}

		// 计算本地文件的 MD5
		localMD5 := md5Hash(content)

		// 与内嵌指纹的校验和比较
		if checksum, ok := resources.CheckSum[engineName]; ok {
			if localMD5 != checksum {
				// 本地文件与内嵌不同，使用本地文件
				if err := replaceEmbedData(engineName, content); err != nil {
					return err
				}
				fmt.Printf("[*] 使用本地指纹: %s\n", fileName)
			}
		}
	}

	return nil
}

// downloadFingerConfig 下载单个指纹配置文件
//
// 参数：
//   - engineName: 引擎名称
//   - fileName: 文件名
//   - fingerPath: 保存路径
//
// 返回：
//   - bool: 是否有更新
//   - error: 下载错误
func downloadFingerConfig(engineName, fileName, fingerPath string) (bool, error) {
	url := baseURL + fileName
	filePath := filepath.Join(fingerPath, fileName)

	fmt.Printf("[*] 检查 %s ...\n", engineName)

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// 下载文件
	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("读取响应失败: %v", err)
	}

	// 计算下载内容的 MD5
	newMD5 := md5Hash(content)

	// 检查本地文件是否存在
	if existingContent, err := os.ReadFile(filePath); err == nil {
		// 文件存在，比较 MD5
		existingMD5 := md5Hash(existingContent)
		if newMD5 == existingMD5 {
			color.Gray.Printf("    %s 已是最新\n", fileName)
			return false, nil
		}
	}

	// 保存新文件
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return false, fmt.Errorf("保存文件失败: %v", err)
	}

	color.Green.Printf("    ✓ 已更新 %s\n", fileName)
	return true, nil
}

// replaceEmbedData 替换内嵌的指纹数据
//
// 参数：
//   - engineName: 引擎名称
//   - content: 新的指纹数据
//
// 返回：
//   - error: 替换错误
func replaceEmbedData(engineName string, content []byte) error {
	switch engineName {
	case fingers.FingersEngine:
		resources.FingersHTTPData = content
	case fingers.FingerPrintEngine:
		resources.FingerprinthubWebData = content
	case fingers.EHoleEngine:
		resources.EholeData = content
	case fingers.GobyEngine:
		resources.GobyData = content
	case fingers.WappalyzerEngine:
		resources.WappalyzerData = content
	default:
		return fmt.Errorf("未知的引擎名称: %s", engineName)
	}
	return nil
}

// getFingerPath 获取指纹存储路径
// 优先使用可执行文件所在目录下的 fingers 目录
//
// 返回：
//   - string: 指纹存储路径
func getFingerPath() string {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return DefaultFingerPath
	}

	// 使用可执行文件所在目录
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, DefaultFingerPath)
}

// md5Hash 计算数据的 MD5 哈希值
//
// 参数：
//   - data: 要计算哈希的数据
//
// 返回：
//   - string: MD5 哈希值（十六进制字符串）
func md5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
