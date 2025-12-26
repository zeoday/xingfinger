// Package finger 提供 Web 指纹识别核心功能
package finger

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/twmb/murmur3"
)

// mmh3Hash 计算 Murmur3 hash 值
func mmh3Hash(data []byte) string {
	h := murmur3.New32()
	h.Write(data)
	return fmt.Sprintf("%d", int32(h.Sum32()))
}

// base64Encode Base64 编码（每 76 字符换行，FOFA 格式）
func base64Encode(data []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(data)
	var buf bytes.Buffer
	for i, ch := range encoded {
		buf.WriteByte(byte(ch))
		if (i+1)%76 == 0 {
			buf.WriteByte('\n')
		}
	}
	buf.WriteByte('\n')
	return buf.Bytes()
}

// calcFaviconHash 获取 favicon 并计算 hash
func calcFaviconHash(url string) string {
	client := http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "0"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "0"
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "0"
	}

	return mmh3Hash(base64Encode(body))
}
