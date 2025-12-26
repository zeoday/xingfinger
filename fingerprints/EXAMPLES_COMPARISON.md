# 指纹格式示例对比分析

本文档详细对比了 5 种指纹格式的示例文件，展示了每种格式的独特特性和差异。

## 📊 示例对比

所有示例都用来检测 **WordPress、Joomla、Drupal** 这三个常见的 CMS 系统。

### 检测特征

| CMS | 特征 | 位置 |
|-----|------|------|
| WordPress | wp-content、wp-includes、wp-admin | body、header |
| Joomla | Joomla!、com_、X-Powered-By: Joomla | body、header |
| Drupal | Drupal、X-Drupal-Cache | body、header |

---

## 1. EHole 格式示例

### 特点

- **最简洁** - 直接列出特征
- **支持多种方法** - keyword、regular、faviconhash
- **支持多个位置** - body、header、title
- **易于理解** - 初学者友好

### 示例代码

```json
{
  "fingerprint": [
    {
      "cms": "WordPress",
      "method": "keyword",
      "location": "body",
      "keyword": ["wp-content", "wp-includes"]
    },
    {
      "cms": "Joomla",
      "method": "regular",
      "location": "body",
      "keyword": ["Joomla!\\s+([\\d.]+)"]
    },
    {
      "cms": "Drupal",
      "method": "keyword",
      "location": "header",
      "keyword": ["X-Drupal-Cache"]
    }
  ]
}
```

### 关键特性

1. **keyword 方法** - 精确字符串匹配
   - WordPress 使用 keyword 匹配 "wp-content" 和 "wp-includes"
   - 两个关键词都需要匹配（AND 逻辑）

2. **regular 方法** - 正则表达式匹配
   - Joomla 使用 regular 提取版本号
   - 支持版本提取

3. **多个位置** - 检测不同位置
   - WordPress 检测 body
   - Drupal 检测 header

### 优点

- ✅ 格式简洁
- ✅ 易于学习
- ✅ 支持多种方法
- ✅ 支持多个位置

### 缺点

- ❌ 不支持逻辑组合
- ❌ 功能相对简单

---

## 2. Goby 格式示例

### 特点

- **支持逻辑组合** - AND、OR 逻辑
- **灵活的规则定义** - 复杂条件组合
- **JSON 数组格式** - 易于扩展
- **中等复杂度** - 适合中等场景

### 示例代码

```json
[
  {
    "name": "WordPress",
    "logic": "a|b|c",
    "rule": [
      {
        "label": "a",
        "feature": "wp-content",
        "is_equal": false
      },
      {
        "label": "b",
        "feature": "wp-includes",
        "is_equal": false
      },
      {
        "label": "c",
        "feature": "wp-admin",
        "is_equal": false
      }
    ]
  },
  {
    "name": "Joomla",
    "logic": "a&b",
    "rule": [
      {
        "label": "a",
        "feature": "Joomla",
        "is_equal": false
      },
      {
        "label": "b",
        "feature": "com_",
        "is_equal": false
      }
    ]
  }
]
```

### 关键特性

1. **OR 逻辑** - WordPress 使用 "a|b|c"
   - 匹配 wp-content 或 wp-includes 或 wp-admin 中的任意一个
   - 更灵活的匹配方式

2. **AND 逻辑** - Joomla 使用 "a&b"
   - 需要同时匹配 "Joomla" 和 "com_"
   - 更严格的匹配条件

3. **is_equal 字段** - 控制匹配方式
   - false：模糊匹配（包含）
   - true：精确匹配（相等）

### 优点

- ✅ 支持复杂逻辑
- ✅ 灵活的条件组合
- ✅ 易于维护

### 缺点

- ❌ 相对复杂
- ❌ 学习曲线陡峭

---

## 3. Wappalyzer 格式示例

### 特点

- **多种检测方式** - HTML、headers、scripts、cookies、meta
- **技术依赖关系** - implies 字段
- **JSON 对象格式** - 按技术名称组织
- **元数据丰富** - 包含图标、网站等信息

### 示例代码

```json
{
  "WordPress": {
    "cats": [1, 6],
    "headers": {
      "X-Powered-By": "WordPress"
    },
    "html": [
      "<link[^>]+href=\"[^\"]*wp-content/",
      "<script[^>]+src=\"[^\"]*wp-includes/"
    ],
    "scripts": [
      "/wp-includes/js/",
      "/wp-content/plugins/"
    ],
    "implies": "PHP",
    "icon": "WordPress.svg",
    "website": "https://wordpress.org"
  },
  "Joomla": {
    "cats": [1, 6],
    "headers": {
      "X-Powered-By": "Joomla"
    },
    "html": [
      "Joomla!",
      "com_"
    ],
    "meta": {
      "generator": "Joomla"
    },
    "implies": "PHP",
    "icon": "Joomla.svg",
    "website": "https://www.joomla.org"
  }
}
```

### 关键特性

1. **多种检测方式**
   - headers：HTTP 响应头
   - html：HTML 内容（支持正则）
   - scripts：脚本路径
   - meta：Meta 标签

2. **技术依赖关系** - implies 字段
   - WordPress implies PHP
   - 自动推导相关技术

3. **元数据** - 丰富的信息
   - cats：分类 ID
   - icon：图标文件
   - website：官网地址

### 优点

- ✅ 检测方式多样
- ✅ 支持技术依赖
- ✅ 元数据丰富
- ✅ 易于集成

### 缺点

- ❌ 格式相对复杂
- ❌ 需要维护元数据

---

## 4. Fingers 格式示例

### 特点

- **功能完整** - 支持多种检测方式
- **灵活的规则** - 支持复杂的检测逻辑
- **JSON 数组格式** - 易于扩展
- **多个检测位置** - headers、html、scripts、cookies、meta

### 示例代码

```json
[
  {
    "name": "WordPress",
    "category": "CMS",
    "website": "https://wordpress.org",
    "headers": {
      "X-Powered-By": "WordPress"
    },
    "html": [
      "wp-content",
      "wp-includes"
    ],
    "scripts": [
      "/wp-includes/js/",
      "/wp-content/plugins/"
    ],
    "cookies": {
      "wordpress_logged_in": ""
    }
  },
  {
    "name": "Joomla",
    "category": "CMS",
    "website": "https://www.joomla.org",
    "headers": {
      "X-Powered-By": "Joomla"
    },
    "html": [
      "Joomla",
      "com_"
    ],
    "meta": {
      "generator": "Joomla"
    }
  }
]
```

### 关键特性

1. **多种检测方式**
   - headers：HTTP 响应头
   - html：HTML 内容
   - scripts：脚本路径
   - cookies：Cookie 检测
   - meta：Meta 标签

2. **灵活的规则**
   - 支持多个检测位置
   - 支持复杂的逻辑组合

3. **元数据**
   - category：分类
   - website：官网

### 优点

- ✅ 功能完整
- ✅ 检测方式多样
- ✅ 灵活性高
- ✅ 易于扩展

### 缺点

- ❌ 相对复杂
- ❌ 学习难度高

---

## 5. FingerPrintHub 格式示例

### 特点

- **最灵活和强大** - 基于 Nuclei 模板
- **多种 Matcher 类型** - word、regex、status-code、favicon 等
- **支持提取器** - 提取信息
- **支持条件逻辑** - AND、OR 逻辑
- **最高级功能** - 最复杂但最强大

### 示例代码

```json
[
  {
    "id": "wordpress-detect",
    "info": {
      "name": "WordPress",
      "author": "test",
      "tags": "detect,tech,wordpress,cms",
      "severity": "info",
      "metadata": {
        "product": "WordPress",
        "vendor": "WordPress"
      }
    },
    "http": [
      {
        "method": "GET",
        "path": ["{{BaseURL}}/"],
        "matchers": [
          {
            "type": "word",
            "words": ["wp-content", "wp-includes"],
            "case-insensitive": true
          }
        ]
      }
    ]
  },
  {
    "id": "joomla-detect",
    "info": {
      "name": "Joomla",
      "author": "test",
      "tags": "detect,tech,joomla,cms",
      "severity": "info",
      "metadata": {
        "product": "Joomla",
        "vendor": "Joomla"
      }
    },
    "http": [
      {
        "method": "GET",
        "path": ["{{BaseURL}}/"],
        "matchers": [
          {
            "type": "regex",
            "regex": ["Joomla!\\s+([\\d.]+)"],
            "case-insensitive": true
          }
        ]
      }
    ]
  },
  {
    "id": "drupal-detect",
    "info": {
      "name": "Drupal",
      "author": "test",
      "tags": "detect,tech,drupal,cms",
      "severity": "info",
      "metadata": {
        "product": "Drupal",
        "vendor": "Drupal"
      }
    },
    "http": [
      {
        "method": "GET",
        "path": ["{{BaseURL}}/"],
        "matchers": [
          {
            "type": "word",
            "words": ["X-Drupal-Cache"],
            "part": "header",
            "case-insensitive": true
          }
        ]
      }
    ]
  }
]
```

### 关键特性

1. **多种 Matcher 类型**
   - word：字符串匹配
   - regex：正则表达式
   - status-code：HTTP 状态码
   - favicon：Favicon hash

2. **灵活的配置**
   - part：指定检测位置（body、header 等）
   - case-insensitive：大小写不敏感

3. **完整的元数据**
   - id：指纹 ID
   - info：详细信息
   - tags：标签
   - metadata：产品和厂商信息

### 优点

- ✅ 最灵活和强大
- ✅ 支持最多的功能
- ✅ 支持提取器
- ✅ 支持条件逻辑

### 缺点

- ❌ 最复杂
- ❌ 学习难度最高

---

## 📈 格式对比总结

### 复杂度对比

```
EHole      ████░░░░░░ 40%
Goby       ██████░░░░ 60%
Wappalyzer ██████░░░░ 60%
Fingers    ████████░░ 80%
FingerPrintHub ██████████ 100%
```

### 功能对比

| 功能 | EHole | Goby | Wappalyzer | Fingers | FingerPrintHub |
|------|-------|------|-----------|---------|----------------|
| 字符串匹配 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 正则匹配 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 逻辑组合 | ❌ | ✅ | ❌ | ✅ | ✅ |
| 多个位置 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 技术依赖 | ❌ | ❌ | ✅ | ❌ | ❌ |
| 提取器 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 条件逻辑 | ❌ | ✅ | ❌ | ✅ | ✅ |
| 元数据 | 少 | 少 | 多 | 中 | 多 |

---

## 🎯 选择建议

### 快速开始
```bash
# 使用 EHole 格式 - 最简单
./xingfinger -u https://example.com --ehole fingerprints/custom_ehole.json
```

### 中等复杂度
```bash
# 使用 Goby 格式 - 支持逻辑组合
./xingfinger -u https://example.com --goby fingerprints/custom_goby.json

# 或使用 Wappalyzer 格式 - 多种检测方式
./xingfinger -u https://example.com --wappalyzer fingerprints/custom_wappalyzer.json
```

### 复杂场景
```bash
# 使用 Fingers 格式 - 功能完整
./xingfinger -u https://example.com --fingers fingerprints/custom_fingers.json

# 或使用 FingerPrintHub 格式 - 最强大
./xingfinger -u https://example.com --fingerprinthub fingerprints/custom_fingerprinthub.json
```

---

## 📝 示例特点总结

### EHole 示例
- **特点**：简洁直接
- **方法**：keyword、regular
- **位置**：body、header
- **逻辑**：无（隐含 AND）

### Goby 示例
- **特点**：支持逻辑组合
- **逻辑**：a|b|c（OR）、a&b（AND）
- **灵活性**：高
- **复杂度**：中

### Wappalyzer 示例
- **特点**：多种检测方式
- **方式**：headers、html、scripts、meta
- **元数据**：丰富（图标、网站等）
- **依赖关系**：支持 implies

### Fingers 示例
- **特点**：功能完整
- **方式**：headers、html、scripts、cookies、meta
- **灵活性**：高
- **复杂度**：高

### FingerPrintHub 示例
- **特点**：最灵活和强大
- **Matcher 类型**：word、regex、status-code、favicon
- **灵活性**：最高
- **复杂度**：最高

---

## 🧪 测试验证

所有示例都已通过测试验证，可以直接使用。

### 测试命令

```bash
# 启动测试服务器
go run test_server_main.go &

# 测试 EHole 格式
./xingfinger -u http://localhost:8888 --ehole fingerprints/custom_ehole.json

# 测试 Goby 格式
./xingfinger -u http://localhost:8888 --goby fingerprints/custom_goby.json

# 测试 Wappalyzer 格式
./xingfinger -u http://localhost:8888 --wappalyzer fingerprints/custom_wappalyzer.json

# 测试 Fingers 格式
./xingfinger -u http://localhost:8888 --fingers fingerprints/custom_fingers.json

# 测试 FingerPrintHub 格式
./xingfinger -u http://localhost:8888 --fingerprinthub fingerprints/custom_fingerprinthub.json
```

---

## 总结

这 5 个示例文件展示了不同指纹格式的独特特性：

1. **EHole** - 最简洁，适合入门
2. **Goby** - 支持逻辑，适合中等复杂度
3. **Wappalyzer** - 多种方式，适合 Web 技术
4. **Fingers** - 功能完整，适合复杂场景
5. **FingerPrintHub** - 最强大，适合高级用户

选择合适的格式，根据你的需求创建指纹文件！
