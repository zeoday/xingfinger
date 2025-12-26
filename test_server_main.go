//go:build ignore

package main

import (
	"fmt"
	"net/http"
)

func main() {
	// WordPress 测试页面 (EHole 指纹)
	http.HandleFunc("/wordpress", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>WordPress Site</title>
    <link rel="stylesheet" href="/wp-content/themes/default/style.css">
</head>
<body>
    <h1>WordPress Test</h1>
    <script src="/wp-includes/js/jquery.js"></script>
</body>
</html>`)
	})

	// Shiro 测试页面 (EHole 指纹 - header)
	http.HandleFunc("/shiro", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Set-Cookie", "rememberMe=deleteMe; Path=/")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Shiro App</title></head>
<body><h1>Shiro Test</h1></body>
</html>`)
	})

	// Spring Boot 测试页面 (EHole 指纹)
	http.HandleFunc("/spring", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"servletContextInitParams":{},"profiles":["default"],"logback":"enabled"}`)
	})

	// Swagger UI 测试页面 (EHole 指纹)
	http.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Swagger UI</title></head>
<body>
    <div id="swagger-ui"></div>
    <script src="swagger-ui-bundle.js"></script>
</body>
</html>`)
	})

	// Nginx 测试页面 (Wappalyzer 指纹)
	http.HandleFunc("/nginx", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Server", "nginx/1.18.0")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Welcome to nginx!</title></head>
<body>
<h1>Welcome to nginx!</h1>
</body>
</html>`)
	})

	// jQuery 测试页面 (Wappalyzer 指纹)
	http.HandleFunc("/jquery", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>jQuery Test</title>
    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
</head>
<body>
    <h1>jQuery Test Page</h1>
</body>
</html>`)
	})

	// PHP 测试页面 (Wappalyzer 指纹)
	http.HandleFunc("/php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Powered-By", "PHP/7.4.3")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>PHP Test</title></head>
<body><h1>PHP Test Page</h1></body>
</html>`)
	})

	// Tomcat 测试页面 (多种指纹)
	http.HandleFunc("/tomcat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Server", "Apache-Coyote/1.1")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Apache Tomcat</title></head>
<body>
    <h1>Apache Tomcat</h1>
    <p>If you're seeing this, you've successfully installed Tomcat.</p>
    <a href="/manager/status">Server Status</a>
    <a href="/manager/html">Manager App</a>
</body>
</html>`)
	})

	// Jenkins 测试页面 (多种指纹)
	http.HandleFunc("/jenkins", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Jenkins", "2.319.1")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Dashboard [Jenkins]</title>
    <link rel="icon" href="/static/favicon.ico">
</head>
<body>
    <div id="jenkins">
        <h1>Jenkins Dashboard</h1>
    </div>
</body>
</html>`)
	})

	// GitLab 测试页面 (多种指纹)
	http.HandleFunc("/gitlab", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>GitLab</title>
    <link rel="icon" href="/assets/gitlab_logo-icon.png">
</head>
<body>
    <img src="/assets/gitlab_logo-7ae504fe4f68fdebb3c2034e36621930cd36ea87924c11ff65dbcb8ed50dca58.png" alt="GitLab">
    <h1>GitLab</h1>
</body>
</html>`)
	})

	// 默认页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
    <title>Test Server</title>
</head>
<body>
    <h1>指纹测试服务器</h1>
    <ul>
        <li><a href="/wordpress">WordPress (EHole)</a></li>
        <li><a href="/shiro">Shiro (EHole - Header)</a></li>
        <li><a href="/spring">Spring Boot (EHole)</a></li>
        <li><a href="/swagger">Swagger UI (EHole)</a></li>
        <li><a href="/nginx">Nginx (Wappalyzer)</a></li>
        <li><a href="/jquery">jQuery (Wappalyzer)</a></li>
        <li><a href="/php">PHP (Wappalyzer)</a></li>
        <li><a href="/tomcat">Tomcat (多种指纹)</a></li>
        <li><a href="/jenkins">Jenkins (多种指纹)</a></li>
        <li><a href="/gitlab">GitLab (多种指纹)</a></li>
    </ul>
</body>
</html>`)
	})

	fmt.Println("测试服务器启动在 http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}
