// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"regexp"
	"strings"

	"github.com/yinheli/mahonia"
	"golang.org/x/net/html/charset"
)

// decodeToUTF8 将内容转换为 UTF-8 编码
func decodeToUTF8(content, contentType string) string {
	encoding := detectEncoding(contentType)

	if enc := extractMetaCharset(content); enc != "" {
		encoding = enc
	}

	if enc := detectTitleEncoding(content); enc != "" && encoding == "utf-8" {
		encoding = enc
	}

	if encoding != "" && encoding != "utf-8" {
		return convertEncoding(content, encoding, "utf-8")
	}
	return content
}

// detectEncoding 检测编码
func detectEncoding(contentType string) string {
	contentType = strings.ToLower(contentType)
	switch {
	case strings.Contains(contentType, "gbk"),
		strings.Contains(contentType, "gb2312"),
		strings.Contains(contentType, "gb18030"),
		strings.Contains(contentType, "windows-1252"):
		return "gb18030"
	case strings.Contains(contentType, "big5"):
		return "big5"
	case strings.Contains(contentType, "utf-8"):
		return "utf-8"
	}
	return "gb18030"
}

// extractMetaCharset 从 meta 标签提取 charset
func extractMetaCharset(content string) string {
	re := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		return detectEncoding(match[1])
	}
	return ""
}

// detectTitleEncoding 通过 title 检测编码
func detectTitleEncoding(content string) string {
	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		_, enc, _ := charset.DetermineEncoding([]byte(match[1]), "")
		return detectEncoding(enc)
	}
	return ""
}

// convertEncoding 转换编码
func convertEncoding(src, from, to string) string {
	if from == to {
		return src
	}
	decoder := mahonia.NewDecoder(from)
	result := decoder.ConvertString(src)
	encoder := mahonia.NewDecoder(to)
	_, data, _ := encoder.Translate([]byte(result), true)
	return string(data)
}
