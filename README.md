# XingFinger 红队重点攻击系统指纹探测工具

![Author](https://img.shields.io/badge/Author-yyhuni-green)  ![language](https://img.shields.io/badge/language-Golang-green)

### 简介

```
  __  ___                _____ _                       
  \ \/ (_)___  ____ _   / ____(_)___  ____ ____  _____ 
   \  /| / _ \/ __ `/  / /_  / / __ \/ __ `/ _ \/ ___/ 
   /  \| |  __/ /_/ /  / __/ / / / / / /_/ /  __/ /     
  /_/\_\_|\___/\__, /  /_/   /_/_/ /_/\__, /\___/_/      
              /____/                 /____/   By:yyhuni
```

XingFinger 是一款对资产中重点系统指纹识别的工具，在红队作战中，信息收集是必不可少的环节。XingFinger 旨在帮助红队人员在信息收集期间能够快速从 C 段、大量杂乱的资产中精准定位到易被攻击的系统，从而实施进一步攻击。

### 安装

```bash
go install github.com/yyhuni/xingfinger@latest
```

或者从源码编译：

```bash
git clone https://github.com/yyhuni/xingfinger.git
cd xingfinger
go build -o xingfinger
```

### 使用

```bash
xingfinger -l url.txt              # 批量扫描
xingfinger -u https://example.com  # 单URL扫描
xingfinger -l url.txt -o result.json  # 输出JSON
xingfinger -l url.txt -t 50        # 设置线程数
xingfinger -l url.txt -p http://127.0.0.1:8080  # 设置代理
xingfinger -l url.txt --silent     # 安静模式，只输出命中结果
```

### 指纹编写

指纹文件为 `finger.json`，支持三种识别方式：

```json
{
    "cms": "系统名称",
    "method": "keyword|faviconhash|regular",
    "location": "body|header|title",
    "keyword": ["关键字"]
}
```

### License

MIT License
