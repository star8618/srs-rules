# 🚀 SRS 自动化规则集仓库 - 完整部署指南

## 📋 目录

1. [快速开始](#快速开始)
2. [详细步骤](#详细步骤)
3. [验证部署](#验证部署)
4. [在 SingForge 中使用](#在-singforge-中使用)
5. [常见问题](#常见问题)
6. [高级配置](#高级配置)

---

## 🎯 快速开始（5 分钟）

### 前置要求

- ✅ GitHub 账号
- ✅ Git 已安装
- ✅ 基本的命令行操作

### 一键部署

```bash
# 1. 进入仓库目录
cd /Users/jacky/Desktop/sing-forger/srs-ruleset-automation

# 2. 运行部署脚本
./deploy.sh

# 3. 按照提示输入你的 GitHub 用户名
# 4. 脚本会自动配置所有文件
```

---

## 📖 详细步骤

### 步骤 1：准备 GitHub 仓库

#### 1.1 创建新仓库

1. 访问 https://github.com/new
2. 填写信息：
   - **Repository name**: `srs-rules`（或自定义名称）
   - **Description**: `Automated SRS rule sets for sing-box`
   - **Visibility**: 选择 **Public**（推荐）或 Private
   - ⚠️ **不要勾选** "Initialize this repository with..."
3. 点击 **Create repository**

#### 1.2 推送代码到 GitHub

```bash
# 添加远程仓库（替换 YOUR_USERNAME）
git remote add origin https://github.com/YOUR_USERNAME/srs-rules.git

# 推送到 GitHub
git branch -M main
git push -u origin main
```

**预期输出**：
```
Enumerating objects: 10, done.
...
To https://github.com/YOUR_USERNAME/srs-rules.git
 * [new branch]      main -> main
Branch 'main' set up to track remote branch 'main' from 'origin'.
```

---

### 步骤 2：启用 GitHub Actions

#### 2.1 启用 Actions 权限

1. 进入仓库页面：`https://github.com/YOUR_USERNAME/srs-rules`
2. 点击 **Settings** （⚙️ 设置）
3. 左侧菜单点击 **Actions** → **General**
4. 在 "Actions permissions" 下选择：
   - ✅ **Allow all actions and reusable workflows**
5. 在 "Workflow permissions" 下选择：
   - ✅ **Read and write permissions**
   - ✅ 勾选 **Allow GitHub Actions to create and approve pull requests**
6. 点击 **Save**

#### 2.2 启用工作流

1. 点击顶部的 **Actions** 标签
2. 如果看到提示，点击 **I understand my workflows, go ahead and enable them**

---

### 步骤 3：首次运行

#### 3.1 手动触发工作流

1. 在 **Actions** 页面
2. 左侧选择 **Update SRS Rule Sets** 工作流
3. 右侧点击 **Run workflow** 下拉按钮
4. 再次点击绿色的 **Run workflow** 按钮

#### 3.2 等待完成

- ⏱️ 预计时间：2-5 分钟
- 状态：
  - 🟡 黄色圆圈 = 运行中
  - ✅ 绿色对勾 = 成功
  - ❌ 红色叉号 = 失败（查看日志）

#### 3.3 查看运行日志（可选）

1. 点击正在运行的工作流
2. 点击 **update-rules** 任务
3. 展开步骤查看详细日志

---

### 步骤 4：验证 Release

#### 4.1 检查 Release

1. 返回仓库首页
2. 右侧点击 **Releases**
3. 应该能看到新创建的 Release（如 `v1`）
4. 展开 **Assets**，应该包含：
   - ✅ `chnlist.srs`
   - ✅ `chnlist-all.srs`
   - ✅ `gfwlist.srs`
   - ✅ `chnroute.srs`
   - ✅ `chnroute6.srs`
   - ✅ `metadata.json`

#### 4.2 测试下载

```bash
# 测试下载一个文件
wget https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnlist.srs

# 查看文件大小
ls -lh chnlist.srs
# 应该显示 ~525KB
```

---

## 🔗 在 SingForge 中使用

### 方法 1：修改 SingForge 源配置

编辑 `backend/services/ruleset/downloader.go`，添加你的仓库：

```go
// 在 RuleSources 中添加
"chnlist": {
    Name: "chnlist",
    Type: "domain",
    URL:  "https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnlist.srs",
    MirrorURLs: []string{
        "https://cdn.jsdelivr.net/gh/YOUR_USERNAME/srs-rules@latest/compiled/chnlist.srs",
    },
    FileName:    "chnlist.srs",
    Description: "中国域名白名单（自动更新）",
},
```

### 方法 2：直接使用（无需转换）

在 SingForge 中，将规则集类型改为 **remote**，直接指向你的 Release：

```json
{
  "rule_set": [
    {
      "tag": "chnlist",
      "type": "remote",
      "format": "binary",
      "url": "https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnlist.srs"
    }
  ]
}
```

**优点**：
- ✅ 无需本地转换
- ✅ 自动使用最新版本
- ✅ sing-box 自动缓存
- ✅ 节省本地存储

---

## ❓ 常见问题

### Q1: Actions 运行失败怎么办？

**A**: 检查错误日志：

1. 进入 Actions 页面
2. 点击失败的运行
3. 查看红色叉号的步骤
4. 展开查看错误信息

**常见错误**：
- `Permission denied` → 检查 Actions 权限设置
- `wget: command not found` → GitHub runner 问题（不太可能）
- `404 Not Found` → 源文件 URL 失效（等待源恢复）

### Q2: 如何更改更新频率？

**A**: 编辑 `.github/workflows/update-rules.yml` 第 5 行：

```yaml
schedule:
  # 每天 00:00 UTC (北京时间 08:00)
  - cron: '0 0 * * *'
  
  # 改为每 12 小时一次
  - cron: '0 */12 * * *'
  
  # 改为每周一次（周一 00:00）
  - cron: '0 0 * * 1'
```

### Q3: 可以添加自定义规则吗？

**A**: 可以！编辑 `scripts/convert.go`：

```go
// 添加新函数
func convertMyCustomRule() error {
    domains, err := parseTextLines("source/my-rule.txt")
    if err != nil {
        return err
    }
    return writeSRS("compiled/my-rule.srs", domains, nil)
}

// 在 main() 中调用
func main() {
    // ... 现有代码 ...
    
    // 添加自定义规则
    if err := convertMyCustomRule(); err != nil {
        fmt.Printf("❌ my-rule 转换失败: %v\n", err)
    }
}
```

然后在 workflow 中添加下载步骤。

### Q4: 如何查看规则集包含哪些域名？

**A**: 使用 `sing-box` 命令行工具：

```bash
# 查看规则内容
sing-box rule-set format -f binary compiled/chnlist.srs

# 或者使用 geoview
geoview -action list -input compiled/chnlist.srs
```

### Q5: 私有仓库可以用吗？

**A**: 可以，但需要额外配置：

1. Release 文件默认私有
2. 需要使用 GitHub Token 访问
3. sing-box 配置需要添加 `download_detour` 和 token

**不推荐**，建议使用 Public 仓库（规则文件本身是公开的）。

---

## 🔧 高级配置

### 配置 CDN 加速

使用 jsdelivr CDN 加速访问：

```
原始 URL:
https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/chnlist.srs

CDN URL:
https://cdn.jsdelivr.net/gh/YOUR_USERNAME/srs-rules@latest/compiled/chnlist.srs
```

**注意**：CDN 有缓存，更新后可能需要等待 1-24 小时。

### 添加 Webhook 通知

在 workflow 末尾添加通知步骤（Telegram/Discord/Email）：

```yaml
- name: Send notification
  if: success()
  run: |
    curl -X POST "https://api.telegram.org/bot${{ secrets.BOT_TOKEN }}/sendMessage" \
      -d chat_id="${{ secrets.CHAT_ID }}" \
      -d text="✅ SRS 规则集已更新"
```

### 性能优化

如果转换很慢，可以：

1. 移除不需要的规则集
2. 使用并行转换（Go goroutines）
3. 缓存源文件（避免重复下载）

---

## 📊 监控和统计

### 查看更新历史

```bash
# 查看所有 Release
https://github.com/YOUR_USERNAME/srs-rules/releases

# 查看提交历史
https://github.com/YOUR_USERNAME/srs-rules/commits/main
```

### 查看 Actions 运行统计

进入 Actions 页面，可以看到：
- ✅ 成功次数
- ❌ 失败次数
- ⏱️ 平均运行时间
- 📊 趋势图表

---

## 🎉 完成！

现在你有了一个：
- ✅ 自动更新的 SRS 规则集仓库
- ✅ 每天自动运行
- ✅ 自动创建 Release
- ✅ 可以在 sing-box 中直接使用

**下一步**：
1. 在 SingForge 中配置使用你的规则集
2. 监控 Actions 运行状态
3. 根据需要调整更新频率

---

## 📞 获取帮助

- 📖 查看 [README.md](README.md) 获取更多信息
- 🐛 提交 [Issue](https://github.com/YOUR_USERNAME/srs-rules/issues) 报告问题
- 💬 在 [Discussions](https://github.com/YOUR_USERNAME/srs-rules/discussions) 讨论交流

---

**祝你使用愉快！** 🎈

