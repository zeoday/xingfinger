// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Timeout 请求超时时间（秒）
var Timeout = 10

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Response HTTP 响应结构体
type Response struct {
	URL        string
	Body       string
	Header     map[string][]string
	Server     string
	StatusCode int
	Length     int
	Title      string
	JsURLs     []string
	FavHash    string
}

// userAgents 常用浏览器 User-Agent 列表
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/96.0.4664.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/97.0.4692.71 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/97.0.4692.71 Safari/537.36",
}

// randomUA 随机返回一个 User-Agent
func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// extractTitle 从 HTML 提取页面标题
func extractTitle(body string) string {
	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		title := strings.TrimSpace(match[1])
		title = strings.ReplaceAll(title, "\n", "")
		title = strings.ReplaceAll(title, "\r", "")
		title = strings.ReplaceAll(title, "\t", "")
		return title
	}
	return ""
}

// extractFavicon 提取并计算 Favicon hash
func extractFavicon(body, targetURL string) string {
	paths := extractRegex(`href="(.*?favicon....)"`, body)
	u, _ := url.Parse(targetURL)
	baseURL := u.Scheme + "://" + u.Host

	var faviconURL string
	if len(paths) > 0 {
		fav := paths[0][1]
		if strings.HasPrefix(fav, "//") {
			faviconURL = "http:" + fav
		} else if strings.HasPrefix(fav, "http") {
			faviconURL = fav
		} else {
			faviconURL = baseURL + "/" + fav
		}
	} else {
		faviconURL = baseURL + "/favicon.ico"
	}
	return calcFaviconHash(faviconURL)
}

// fetch 发送 HTTP 请求并解析响应
func fetch(task []string, proxy string) (*Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   time.Duration(Timeout) * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest("GET", task[0], nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "rememberMe", Value: "me"})
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", randomUA())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, _ := ioutil.ReadAll(resp.Body)

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	body := decodeToUTF8(string(rawBody), contentType)

	server := ""
	if s := resp.Header.Get("Server"); s != "" {
		server = s
	} else if p := resp.Header.Get("X-Powered-By"); p != "" {
		server = p
	}

	var jsURLs []string
	if task[1] == "0" {
		jsURLs = parseJSRedirect(body, task[0])
	}

	return &Response{
		URL:        task[0],
		Body:       body,
		Header:     resp.Header,
		Server:     server,
		StatusCode: resp.StatusCode,
		Length:     len(body),
		Title:      extractTitle(body),
		JsURLs:     jsURLs,
		FavHash:    extractFavicon(body, task[0]),
	}, nil
}
