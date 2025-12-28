// Package pkg 提供 xingfinger 的核心功能
// 本文件实现 ARL YAML 格式指纹的解析和匹配
package pkg

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// ARLFingerprint ARL YAML 格式的指纹结构
type ARLFingerprint struct {
	Name string `yaml:"name"`
	Rule string `yaml:"rule"`
}

// ARLEngine ARL 指纹匹配引擎
type ARLEngine struct {
	fingerprints []ARLFingerprint
}

// ARLCondition 解析后的单个条件
type ARLCondition struct {
	Type    string // body, header, title, icon_hash
	Keyword string // 匹配的关键字
}

// NewARLEngine 创建 ARL 引擎
func NewARLEngine(filepath string) (*ARLEngine, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fingerprints []ARLFingerprint
	if err := yaml.Unmarshal(data, &fingerprints); err != nil {
		return nil, err
	}

	return &ARLEngine{fingerprints: fingerprints}, nil
}

// Match 匹配指纹，返回匹配到的 CMS 名称列表
func (e *ARLEngine) Match(body, header, title string, faviconHash string) []string {
	var matched []string
	seen := make(map[string]bool)

	for _, fp := range e.fingerprints {
		if fp.Rule == "" {
			continue
		}

		if e.matchRule(fp.Rule, body, header, title, faviconHash) {
			name := extractARLName(fp.Name)
			if !seen[name] {
				seen[name] = true
				matched = append(matched, name)
			}
		}
	}

	return matched
}

// matchRule 匹配单条规则
// 规则格式: body="xxx" && header="yyy" && title="zzz"
// 所有条件必须同时满足（AND 关系）
func (e *ARLEngine) matchRule(rule, body, header, title, faviconHash string) bool {
	conditions := parseARLConditions(rule)
	if len(conditions) == 0 {
		return false
	}

	// 所有条件都必须匹配
	for _, cond := range conditions {
		if !matchCondition(cond, body, header, title, faviconHash) {
			return false
		}
	}

	return true
}

// parseARLConditions 解析规则字符串为条件列表
func parseARLConditions(rule string) []ARLCondition {
	var conditions []ARLCondition

	// 正则匹配各类条件，支持转义引号
	bodyRe := regexp.MustCompile(`body="((?:[^"\\]|\\.)*)"`)
	headerRe := regexp.MustCompile(`header="((?:[^"\\]|\\.)*)"`)
	titleRe := regexp.MustCompile(`title="((?:[^"\\]|\\.)*)"`)
	iconHashRe := regexp.MustCompile(`icon_hash="([^"]+)"`)

	// 提取 body 条件
	for _, m := range bodyRe.FindAllStringSubmatch(rule, -1) {
		if len(m) > 1 && m[1] != "" {
			conditions = append(conditions, ARLCondition{
				Type:    "body",
				Keyword: unescapeARLString(m[1]),
			})
		}
	}

	// 提取 header 条件
	for _, m := range headerRe.FindAllStringSubmatch(rule, -1) {
		if len(m) > 1 && m[1] != "" {
			conditions = append(conditions, ARLCondition{
				Type:    "header",
				Keyword: unescapeARLString(m[1]),
			})
		}
	}

	// 提取 title 条件
	for _, m := range titleRe.FindAllStringSubmatch(rule, -1) {
		if len(m) > 1 && m[1] != "" {
			conditions = append(conditions, ARLCondition{
				Type:    "title",
				Keyword: unescapeARLString(m[1]),
			})
		}
	}

	// 提取 icon_hash 条件
	for _, m := range iconHashRe.FindAllStringSubmatch(rule, -1) {
		if len(m) > 1 && m[1] != "" {
			conditions = append(conditions, ARLCondition{
				Type:    "icon_hash",
				Keyword: m[1],
			})
		}
	}

	return conditions
}

// matchCondition 匹配单个条件
func matchCondition(cond ARLCondition, body, header, title, faviconHash string) bool {
	switch cond.Type {
	case "body":
		return strings.Contains(strings.ToLower(body), strings.ToLower(cond.Keyword))
	case "header":
		return strings.Contains(strings.ToLower(header), strings.ToLower(cond.Keyword))
	case "title":
		return strings.Contains(strings.ToLower(title), strings.ToLower(cond.Keyword))
	case "icon_hash":
		return faviconHash == cond.Keyword
	}
	return false
}

// unescapeARLString 处理转义字符
func unescapeARLString(s string) string {
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

// extractARLName 从 ARL name 中提取干净的名称
func extractARLName(name string) string {
	suffixes := []string{"_body", "_header", "_title", "_icon_hash"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(name, suffix)
		}
	}
	return name
}
