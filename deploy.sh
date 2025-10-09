#!/bin/bash

# SRS Rules 自动化仓库快速部署脚本

set -e

echo "🚀 SRS Rules 自动化仓库部署脚本"
echo "================================"
echo ""

# 检查是否在正确的目录
if [ ! -f ".github/workflows/update-rules.yml" ]; then
    echo "❌ 错误：请在仓库根目录下运行此脚本"
    exit 1
fi

# 获取 GitHub 用户名
echo "📝 请输入你的 GitHub 用户名："
read -r GITHUB_USERNAME

if [ -z "$GITHUB_USERNAME" ]; then
    echo "❌ 用户名不能为空"
    exit 1
fi

echo ""
echo "✅ 用户名：$GITHUB_USERNAME"
echo ""

# 替换 README.md 中的占位符
echo "🔧 更新 README.md..."
sed -i.bak "s/YOUR_USERNAME/$GITHUB_USERNAME/g" README.md && rm README.md.bak

# 替换 generate-metadata.go 中的占位符
echo "🔧 更新 scripts/generate-metadata.go..."
sed -i.bak "s/YOUR_USERNAME/$GITHUB_USERNAME/g" scripts/generate-metadata.go && rm scripts/generate-metadata.go.bak

echo ""
echo "✅ 配置文件已更新"
echo ""

# 初始化 Git
if [ ! -d ".git" ]; then
    echo "📦 初始化 Git 仓库..."
    git init
    git add .
    git commit -m "Initial commit: SRS Rules 自动化仓库"
    echo "✅ Git 仓库已初始化"
else
    echo "📦 Git 仓库已存在，跳过初始化"
fi

echo ""
echo "================================"
echo "🎉 配置完成！"
echo ""
echo "📋 接下来的步骤："
echo ""
echo "1️⃣  在 GitHub 创建新仓库："
echo "   • 访问：https://github.com/new"
echo "   • 仓库名称：srs-rules（或自定义）"
echo "   • 可见性：Public（推荐）或 Private"
echo "   • 不要勾选 'Initialize this repository...'"
echo ""
echo "2️⃣  推送到 GitHub："
echo "   git remote add origin https://github.com/$GITHUB_USERNAME/srs-rules.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "3️⃣  启用 GitHub Actions："
echo "   • 进入仓库的 Settings → Actions → General"
echo "   • 选择 'Allow all actions and reusable workflows'"
echo "   • 进入 Actions 页面，点击 'I understand...'"
echo ""
echo "4️⃣  手动触发首次运行："
echo "   • 进入 Actions 页面"
echo "   • 选择 'Update SRS Rule Sets' 工作流"
echo "   • 点击 'Run workflow' → 'Run workflow'"
echo "   • 等待 2-5 分钟完成"
echo ""
echo "5️⃣  使用规则集："
echo "   • 访问：https://github.com/$GITHUB_USERNAME/srs-rules/releases"
echo "   • 下载 .srs 文件或使用 URL："
echo "   https://github.com/$GITHUB_USERNAME/srs-rules/releases/latest/download/chnlist.srs"
echo ""
echo "================================"
echo "📖 完整文档：README.md"
echo "================================"

