#!/usr/bin/env python3
"""
指纹测试服务器
针对不同指纹格式的特点提供测试端点
"""

from http.server import HTTPServer, BaseHTTPRequestHandler
import json

class TestHandler(BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        print(f"[{self.address_string()}] {args[0]}")

    def do_GET(self):
        routes = {
            '/': self.index,
            '/ehole/shiro': self.ehole_shiro,
            '/ehole/swagger': self.ehole_swagger,
            '/ehole/spring': self.ehole_spring,
            '/ehole/tomcat': self.ehole_tomcat,
            '/wappalyzer/nginx': self.wappalyzer_nginx,
            '/wappalyzer/php': self.wappalyzer_php,
            '/goby/nacos': self.goby_nacos,
            '/goby/weblogic': self.goby_weblogic,
            '/fingers/apache': self.fingers_apache,
            '/fingerprinthub/thinkphp': self.fingerprinthub_thinkphp,
        }
        
        handler = routes.get(self.path, self.not_found)
        handler()

    def send_html(self, content, status=200, headers=None):
        self.send_response(status)
        self.send_header('Content-Type', 'text/html; charset=utf-8')
        if headers:
            for k, v in headers.items():
                self.send_header(k, v)
        self.end_headers()
        self.wfile.write(content.encode('utf-8'))

    def index(self):
        html = '''<!DOCTYPE html>
<html>
<head><title>指纹测试服务器</title></head>
<body>
<h1>指纹测试服务器</h1>
<h2>EHole 指纹测试</h2>
<ul>
    <li><a href="/ehole/shiro">Shiro (Header: rememberMe)</a></li>
    <li><a href="/ehole/swagger">Swagger UI (Body keyword)</a></li>
    <li><a href="/ehole/spring">Spring Boot (Body keyword)</a></li>
    <li><a href="/ehole/tomcat">Tomcat (Body keyword)</a></li>
</ul>
<h2>Wappalyzer 指纹测试</h2>
<ul>
    <li><a href="/wappalyzer/nginx">Nginx (Server header)</a></li>
    <li><a href="/wappalyzer/php">PHP (X-Powered-By header)</a></li>
</ul>
<h2>Goby 指纹测试</h2>
<ul>
    <li><a href="/goby/nacos">Nacos (Body keyword)</a></li>
    <li><a href="/goby/weblogic">WebLogic (Body keyword)</a></li>
</ul>
<h2>Fingers 指纹测试</h2>
<ul>
    <li><a href="/fingers/apache">Apache (Server header)</a></li>
</ul>
<h2>FingerPrintHub 指纹测试</h2>
<ul>
    <li><a href="/fingerprinthub/thinkphp">ThinkPHP (Header keyword)</a></li>
</ul>
</body>
</html>'''
        self.send_html(html)

    # EHole 指纹测试
    def ehole_shiro(self):
        """Shiro - 通过 header 中的 rememberMe 识别"""
        self.send_html('<h1>Shiro Test</h1>', headers={
            'Set-Cookie': 'rememberMe=deleteMe; Path=/'
        })

    def ehole_swagger(self):
        """Swagger UI - 通过 body 中的关键词识别"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Swagger UI</title></head>
<body>
<div id="swagger-ui"></div>
<script src="swagger-ui-bundle.js"></script>
</body>
</html>'''
        self.send_html(html)

    def ehole_spring(self):
        """Spring Boot - 通过 body 中的关键词识别"""
        html = '''{"servletContextInitParams":{},"profiles":["default"],"logback":"enabled"}'''
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(html.encode('utf-8'))

    def ehole_tomcat(self):
        """Tomcat - 通过 body 中的关键词识别"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Apache Tomcat</title></head>
<body>
<h1>Apache Tomcat</h1>
<a href="/manager/status">Server Status</a>
<a href="/manager/html">Manager App</a>
</body>
</html>'''
        self.send_html(html)

    # Wappalyzer 指纹测试
    def wappalyzer_nginx(self):
        """Nginx - 通过 Server header 识别"""
        self.send_html('<h1>Welcome to nginx!</h1>', headers={
            'Server': 'nginx/1.18.0'
        })

    def wappalyzer_php(self):
        """PHP - 通过 X-Powered-By header 识别"""
        self.send_html('<h1>PHP Test</h1>', headers={
            'X-Powered-By': 'PHP/7.4.3'
        })

    # Goby 指纹测试
    def goby_nacos(self):
        """Nacos - 通过 body 关键词识别"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Nacos</title></head>
<body>
<div class="nacos-console">Nacos Console</div>
</body>
</html>'''
        self.send_html(html)

    def goby_weblogic(self):
        """WebLogic - 通过 body 关键词识别"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Error 404--Not Found</title></head>
<body>
<h1>Error 404--Not Found</h1>
</body>
</html>'''
        self.send_html(html, status=404)

    # Fingers 指纹测试
    def fingers_apache(self):
        """Swagger UI - 通过 body 关键词识别 (Fingers格式)"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Swagger UI</title></head>
<body>
<div id="swagger-ui"></div>
<script src="swagger-ui.js"></script>
</body>
</html>'''
        self.send_html(html)

    # FingerPrintHub 指纹测试
    def fingerprinthub_thinkphp(self):
        """Example Domain - 通过 title 识别 (FingerPrintHub格式)"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Example Domain</title></head>
<body>
<h1>Example Domain</h1>
<p>This domain is for use in illustrative examples in documents.</p>
</body>
</html>'''
        self.send_html(html)

    def not_found(self):
        self.send_html('<h1>404 Not Found</h1>', status=404)

if __name__ == '__main__':
    port = 8888
    server = HTTPServer(('0.0.0.0', port), TestHandler)
    print(f'测试服务器启动在 http://localhost:{port}')
    server.serve_forever()
