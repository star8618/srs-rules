# SRS Rules - 自动化规则集仓库

[![Update Rules](https://github.com/YOUR_USERNAME/srs-rules/actions/workflows/update-rules.yml/badge.svg)](https://github.com/YOUR_USERNAME/srs-rules/actions/workflows/update-rules.yml)
[![Latest Release](https://img.shields.io/github/v/release/YOUR_USERNAME/srs-rules)](https://github.com/YOUR_USERNAME/srs-rules/releases/latest)

🤖 **自动化维护的 sing-box SRS 规则集仓库**

每天自动从官方源下载并转换为 `.srs` 格式，无需手动维护！

---

## 📦 包含的规则集

| 规则集 | 类型 | 说明 | 规则数量 | 大小 |
|--------|------|------|----------|------|
| `chnlist.srs` | 域名 | 中国域名白名单（通用） | 117K+ | ~525KB |
| `chnlist-all.srs` | 域名 | 合并版中国域名（含 Apple/Google） | 118K+ | ~540KB |
| `gfwlist.srs` | 域名 | GFW 被墙域名列表 | 5.9K+ | ~39KB |
| `chnroute.srs` | IP | 中国 IPv4 地址段 | 4.2K+ | ~15KB |
| `chnroute6.srs` | IP | 中国 IPv6 地址段 | 1.5K+ | ~6KB |

---

## 🚀 使用方法

### 方法 1：直接引用 Release（推荐）

在你的 `sing-box` 配置中：

```json
{
  "route": {
    "rule_set": [
      {
        "tag": "chnlist",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnlist.srs",
        "download_detour": "direct"
      },
      {
        "tag": "gfwlist",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/gfwlist.srs",
        "download_detour": "proxy"
      },
      {
        "tag": "chnroute",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnroute.srs",
        "download_detour": "direct"
      }
    ],
    "rules": [
      {
        "rule_set": ["chnlist", "chnroute"],
        "outbound": "direct"
      },
      {
        "rule_set": ["gfwlist"],
        "outbound": "proxy"
      }
    ]
  }
}
```

### 方法 2：使用 jsdelivr CDN（更快）

```json
{
  "rule_set": [
    {
      "tag": "chnlist",
      "type": "remote",
      "format": "binary",
      "url": "https://cdn.jsdelivr.net/gh/YOUR_USERNAME/srs-rules@latest/compiled/chnlist.srs"
    }
  ]
}
```

### 方法 3：本地使用

```bash
# 克隆仓库
git clone https://github.com/YOUR_USERNAME/srs-rules.git

# 使用本地文件
{
  "rule_set": [
    {
      "tag": "chnlist",
      "type": "local",
      "format": "binary",
      "path": "/path/to/srs-rules/compiled/chnlist.srs"
    }
  ]
}
```

---

## 🔄 更新频率

- **自动更新**：每天北京时间 08:00 (UTC 00:00)
- **手动触发**：在 Actions 页面手动运行工作流
- **源推送时**：当有新提交推送到 main 分支时

---

## 📊 规则来源

| 规则集 | 来源 | 维护方 |
|--------|------|--------|
| chnlist | [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) | felixonmars |
| gfwlist | [gfwlist/gfwlist](https://github.com/gfwlist/gfwlist) | gfwlist |
| chnroute | [misakaio/chnroutes2](https://github.com/misakaio/chnroutes2) | misakaio |
| geosite | [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) | Loyalsoldier |
| geoip | [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) | Loyalsoldier |

---

## 🛠️ 如何部署到自己的仓库

### 1. Fork 或创建新仓库

```bash
# 克隆这个仓库模板
git clone https://github.com/YOUR_USERNAME/srs-rules.git
cd srs-rules

# 或者直接在 GitHub 上点击 "Use this template"
```

### 2. 替换仓库信息

编辑以下文件，将 `YOUR_USERNAME` 替换为你的 GitHub 用户名：

- `README.md` 
- `scripts/generate-metadata.go` (第54行)

### 3. 推送到 GitHub

```bash
git add .
git commit -m "Initial commit"
git push origin main
```

### 4. 启用 GitHub Actions

1. 进入仓库的 **Settings** → **Actions** → **General**
2. 确保 "Allow all actions and reusable workflows" 已启用
3. 进入 **Actions** 页面，点击 "I understand my workflows, go ahead and enable them"

### 5. 手动触发首次运行

1. 进入 **Actions** 页面
2. 选择 "Update SRS Rule Sets" 工作流
3. 点击 "Run workflow" → "Run workflow"
4. 等待约 2-5 分钟，完成后会自动创建 Release

### 6. 检查结果

- 进入 **Releases** 页面，应该能看到新创建的 Release
- 点击 Release，可以下载 `.srs` 文件
- 查看 `compiled/metadata.json` 获取规则集信息

---

## 📝 自定义配置

### 修改更新时间

编辑 `.github/workflows/update-rules.yml` 第 5 行：

```yaml
on:
  schedule:
    # 修改这里的 cron 表达式
    # 格式：分钟 小时 日 月 星期
    # 示例：每12小时更新一次
    - cron: '0 */12 * * *'
```

### 添加更多规则集

编辑 `scripts/convert.go`，添加新的转换函数：

```go
// 示例：添加 adblock 规则
func convertAdblock() error {
    domains, err := parseTextLines("source/adblock.txt")
    if err != nil {
        return err
    }
    return writeSRS("compiled/adblock.srs", domains, nil)
}
```

然后在 `main()` 函数中调用。

### 禁用 Geosite/GeoIP 转换

如果只需要文本规则，注释掉 `scripts/convert.go` 第 57 行：

```go
// convertGeoview()  // 注释掉这行
```

---

## 🔍 查看元数据

访问 `metadata.json` 获取所有规则集信息：

```bash
curl https://cdn.jsdelivr.net/gh/YOUR_USERNAME/srs-rules@latest/compiled/metadata.json
```

返回示例：

```json
{
  "version": "1.0",
  "updated_at": "2025-10-09T00:00:00Z",
  "rule_sets": [
    {
      "name": "中国域名白名单",
      "tag": "chnlist",
      "type": "domain",
      "format": "binary",
      "file_size": 537600,
      "url": "https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnlist.srs",
      "description": "中国域名白名单（通用，117K+ 条）",
      "updated_at": "2025-10-09T00:00:00Z"
    }
  ]
}
```

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

## 📄 许可证

MIT License

---

## 🙏 致谢

感谢以下项目和贡献者：

- [sing-box](https://github.com/SagerNet/sing-box) - sing-box 核心
- [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) - 中国域名列表
- [gfwlist/gfwlist](https://github.com/gfwlist/gfwlist) - GFW 列表
- [misakaio/chnroutes2](https://github.com/misakaio/chnroutes2) - 中国路由表
- [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) - Geosite/GeoIP 数据库
- [metacubex/geo](https://github.com/metacubex/geo) - geoview 工具

---

## 📮 联系方式

- 问题反馈：[Issues](https://github.com/YOUR_USERNAME/srs-rules/issues)
- 讨论交流：[Discussions](https://github.com/YOUR_USERNAME/srs-rules/discussions)

---

**⭐ 如果这个项目对你有帮助，请给个 Star！**

