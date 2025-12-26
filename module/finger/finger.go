// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gookit/color"
	"github.com/yyhuni/xingfinger/module/finger/source"
	"github.com/yyhuni/xingfinger/module/queue"
)

// Result 扫描结果结构体
type Result struct {
	URL        string `json:"url"`
	CMS        string `json:"cms"`
	Server     string `json:"server"`
	StatusCode int    `json:"status_code"`
	Length     int    `json:"length"`
	Title      string `json:"title"`
}

// Scanner 指纹扫描器
type Scanner struct {
	queue        *queue.Queue
	wg           sync.WaitGroup
	mu           sync.Mutex // 保护结果切片的并发写入
	thread       int
	output       string
	proxy        string
	silent       bool
	allResults   []Result
	hitResults   []Result
	fingerprints *Packjson
}

// NewScanner 创建新的扫描器实例
func NewScanner(urls []string, thread int, output, proxy string, timeout int, silent bool) *Scanner {
	s := &Scanner{
		queue:      queue.NewQueue(),
		thread:     thread,
		output:     output,
		proxy:      proxy,
		silent:     silent,
		allResults: []Result{},
		hitResults: []Result{},
	}

	Timeout = timeout

	err := LoadWebfingerprint(source.GetExePath() + "/finger.json")
	if err != nil {
		fmt.Println("[!] Fingerprint file error")
		os.Exit(1)
	}
	s.fingerprints = GetWebfingerprint()

	for _, url := range urls {
		s.queue.Push([]string{url, "0"})
	}
	return s
}

// Run 启动扫描任务
func (s *Scanner) Run() {
	for i := 0; i < s.thread; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.scan()
		}()
	}

	s.wg.Wait()

	if !s.silent {
		color.RGBStyleFromString("244,211,49").Printf("\n[+] Scanned: %d, Matched: %d\n", len(s.allResults), len(s.hitResults))
	}

	if s.output != "" {
		saveResults(s.output, s.allResults)
	}
}

// scan 执行扫描逻辑
func (s *Scanner) scan() {
	for {
		item := s.queue.Pop()
		if item == nil {
			return
		}

		urlData, ok := item.([]string)
		if !ok {
			continue
		}

		resp, err := doHTTPRequest(urlData, s.proxy)
		if err != nil {
			urlData[0] = strings.ReplaceAll(urlData[0], "https://", "http://")
			resp, err = doHTTPRequest(urlData, s.proxy)
			if err != nil {
				continue
			}
		}

		for _, jurl := range resp.jsurl {
			if jurl != "" {
				s.queue.Push([]string{jurl, "1"})
			}
		}

		headers := toJSON(resp.header)
		var matched []string

		for _, fp := range s.fingerprints.Fingerprint {
			var target string
			switch fp.Location {
			case "body":
				target = resp.body
			case "header":
				target = headers
			case "title":
				target = resp.title
			default:
				continue
			}

			switch fp.Method {
			case "keyword":
				if matchKeyword(target, fp.Keyword) {
					matched = append(matched, fp.Cms)
				}
			case "faviconhash":
				if resp.favhash == fp.Keyword[0] {
					matched = append(matched, fp.Cms)
				}
			case "regular":
				if matchRegex(target, fp.Keyword) {
					matched = append(matched, fp.Cms)
				}
			}
		}

		matched = unique(matched)
		result := Result{
			URL:        resp.url,
			CMS:        strings.Join(matched, ","),
			Server:     resp.server,
			StatusCode: resp.statuscode,
			Length:     resp.length,
			Title:      resp.title,
		}

		// 加锁保护并发写入
		s.mu.Lock()
		s.allResults = append(s.allResults, result)
		if result.CMS != "" {
			s.hitResults = append(s.hitResults, result)
		}
		s.mu.Unlock()

		// httpx 风格输出
		if s.silent {
			if result.CMS != "" {
				fmt.Printf("%s [%s]\n", result.URL, result.CMS)
			}
		} else {
			// 格式: url [status] [length] [server] [title] [cms]
			var parts []string
			parts = append(parts, result.URL)
			parts = append(parts, fmt.Sprintf("[%d]", result.StatusCode))
			parts = append(parts, fmt.Sprintf("[%d]", result.Length))
			if result.Server != "None" && result.Server != "" {
				parts = append(parts, fmt.Sprintf("[%s]", result.Server))
			}
			if result.Title != "" {
				parts = append(parts, fmt.Sprintf("[%s]", result.Title))
			}
			if result.CMS != "" {
				parts = append(parts, fmt.Sprintf("[%s]", result.CMS))
			}

			line := strings.Join(parts, " ")
			if result.CMS != "" {
				color.RGBStyleFromString("237,64,35").Println(line)
			} else {
				fmt.Println(line)
			}
		}
	}
}

// toJSON 将 map 转换为 JSON 字符串
func toJSON(data map[string][]string) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// unique 去重
func unique(arr []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range arr {
		if v != "" && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
