#!/usr/bin/env python3
"""
指纹测试服务器
针对不同指纹格式的特点提供测试端点
"""

from http.server import HTTPServer, BaseHTTPRequestHandler
import hashlib
import base64
import struct

# 简单的 favicon 图标数据 (16x16 红色图标)
FAVICON_DATA = bytes([
    0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x10, 0x10,
    0x00, 0x00, 0x01, 0x00, 0x18, 0x00, 0x68, 0x03,
    0x00, 0x00, 0x16, 0x00, 0x00, 0x00
] + [0xFF, 0x00, 0x00] * 256 + [0x00] * 64)

def mmh3_hash(data):
    """计算 MMH3 hash (与 Shodan favicon hash 兼容)"""
    import struct
    
    def mmh3_32(key, seed=0):
        length = len(key)
        nblocks = length // 4
        h1 = seed
        c1 = 0xcc9e2d51
        c2 = 0x1b873593
        
        for i in range(nblocks):
            k1 = struct.unpack('<I', key[i*4:(i+1)*4])[0]
            k1 = (k1 * c1) & 0xffffffff
            k1 = ((k1 << 15) | (k1 >> 17)) & 0xffffffff
            k1 = (k1 * c2) & 0xffffffff
            h1 ^= k1
            h1 = ((h1 << 13) | (h1 >> 19)) & 0xffffffff
            h1 = ((h1 * 5) + 0xe6546b64) & 0xffffffff
        
        tail = key[nblocks * 4:]
        k1 = 0
        if len(tail) >= 3:
            k1 ^= tail[2] << 16
        if len(tail) >= 2:
            k1 ^= tail[1] << 8
        if len(tail) >= 1:
            k1 ^= tail[0]
            k1 = (k1 * c1) & 0xffffffff
            k1 = ((k1 << 15) | (k1 >> 17)) & 0xffffffff
            k1 = (k1 * c2) & 0xffffffff
            h1 ^= k1
        
        h1 ^= length
        h1 ^= h1 >> 16
        h1 = (h1 * 0x85ebca6b) & 0xffffffff
        h1 ^= h1 >> 13
        h1 = (h1 * 0xc2b2ae35) & 0xffffffff
        h1 ^= h1 >> 16
        
        # 转换为有符号整数
        if h1 >= 0x80000000:
            h1 -= 0x100000000
        return h1
    
    # Base64 编码后计算 hash
    b64_data = base64.b64encode(data)
    return str(mmh3_32(b64_data))

class TestHandler(BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        print(f"[{self.address_string()}] {args[0]}")

    def do_GET(self):
        routes = {
            '/': self.index,
            # EHole 指纹测试
            '/ehole/shiro': self.ehole_shiro,
            '/ehole/swagger': self.ehole_swagger,
            '/ehole/spring': self.ehole_spring,
            '/ehole/tomcat': self.ehole_tomcat,
            # Wappalyzer 指纹测试
            '/wappalyzer/nginx': self.wappalyzer_nginx,
            '/wappalyzer/php': self.wappalyzer_php,
            # Goby 指纹测试
            '/goby/nacos': self.goby_nacos,
            '/goby/weblogic': self.goby_weblogic,
            # Fingers 指纹测试
            '/fingers/apache': self.fingers_apache,
            # FingerPrintHub 指纹测试
            '/fingerprinthub/thinkphp': self.fingerprinthub_thinkphp,
            # ARL 指纹测试 - 覆盖所有条件类型
            '/arl/body': self.arl_body,
            '/arl/header': self.arl_header,
            '/arl/title': self.arl_title,
            '/arl/icon_hash': self.arl_icon_hash,
            '/arl/combined': self.arl_combined,
            '/arl/favicon.ico': self.arl_favicon,
            '/favicon.ico': self.arl_favicon,
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
        # 计算 favicon hash 用于显示
        favicon_hash = mmh3_hash(FAVICON_DATA)
        
        html = f'''<!DOCTYPE html>
<html>
<head><title>指纹测试服务器</title></head>
<body>
<h1>指纹测试服务器</h1>

<h2>ARL 指纹测试 (覆盖所有条件类型)</h2>
<ul>
    <li><a href="/arl/body">body 条件测试</a> - body="thinkphp"</li>
    <li><a href="/arl/header">header 条件测试</a> - header="X-Powered-By: ThinkPHP"</li>
    <li><a href="/arl/title">title 条件测试</a> - title="若依管理系统"</li>
    <li><a href="/arl/icon_hash">icon_hash 条件测试</a> - icon_hash="{favicon_hash}"</li>
    <li><a href="/arl/combined">组合条件测试</a> - body && header && title</li>
</ul>

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

<p><strong>Favicon MMH3 Hash:</strong> {favicon_hash}</p>
</body>
</html>'''
        self.send_html(html)

    # ==================== ARL 指纹测试 ====================
    
    def arl_body(self):
        """ARL body 条件测试 - 匹配 ThinkPHP"""
        html = '''<!DOCTYPE html>
<html>
<head><title>ThinkPHP Test</title></head>
<body>
<h1>ThinkPHP Framework</h1>
<p>Powered by thinkphp</p>
<div class="thinkphp-logo">ThinkPHP V5.0</div>
</body>
</html>'''
        self.send_html(html)

    def arl_header(self):
        """ARL header 条件测试 - 匹配 ThinkPHP header"""
        html = '''<!DOCTYPE html>
<html>
<head><title>Header Test</title></head>
<body>
<h1>Header Fingerprint Test</h1>
</body>
</html>'''
        self.send_html(html, headers={
            'X-Powered-By': 'ThinkPHP',
            'Server': 'nginx/1.18.0'
        })

    def arl_title(self):
        """ARL title 条件测试 - 匹配若依管理系统"""
        html = '''<!DOCTYPE html>
<html>
<head><title>若依管理系统</title></head>
<body>
<h1>若依后台管理系统</h1>
<p>RuoYi Management System</p>
</body>
</html>'''
        self.send_html(html)

    def arl_icon_hash(self):
        """ARL icon_hash 条件测试"""
        favicon_hash = mmh3_hash(FAVICON_DATA)
        html = f'''<!DOCTYPE html>
<html>
<head>
<title>Icon Hash Test</title>
<link rel="icon" href="/arl/favicon.ico" type="image/x-icon">
</head>
<body>
<h1>Icon Hash Fingerprint Test</h1>
<p>Favicon MMH3 Hash: {favicon_hash}</p>
<p>访问 <a href="/arl/favicon.ico">/arl/favicon.ico</a> 获取图标</p>
</body>
</html>'''
        self.send_html(html)

    def arl_favicon(self):
        """返回 favicon 图标"""
        self.send_response(200)
        self.send_header('Content-Type', 'image/x-icon')
        self.send_header('Content-Length', str(len(FAVICON_DATA)))
        self.end_headers()
        self.wfile.write(FAVICON_DATA)

    def arl_combined(self):
        """ARL 组合条件测试 - body && header && title"""
        html = '''<!DOCTYPE html>
<html>
<head><title>若依管理系统</title></head>
<body>
<h1>若依后台管理系统</h1>
<p>Powered by RuoYi-Vue</p>
<div class="ruoyi-footer">Copyright © RuoYi</div>
</body>
</html>'''
        self.send_html(html, headers={
            'X-Powered-By': 'RuoYi',
            'Server': 'nginx'
        })

    # ==================== EHole 指纹测试 ====================
    
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

    # ==================== Wappalyzer 指纹测试 ====================
    
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

    # ==================== Goby 指纹测试 ====================
    
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

    # ==================== Fingers 指纹测试 ====================
    
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

    # ==================== FingerPrintHub 指纹测试 ====================
    
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
    port = 9999
    server = HTTPServer(('0.0.0.0', port), TestHandler)
    print(f'测试服务器启动在 http://localhost:{port}')
    print(f'Favicon MMH3 Hash: {mmh3_hash(FAVICON_DATA)}')
    server.serve_forever()
