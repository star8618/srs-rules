#!/bin/bash

# 本地测试脚本

set -e

echo "🧪 本地测试 SRS 转换"
echo "===================="
echo ""

# 创建目录
mkdir -p source compiled

echo "📥 下载测试源文件..."

# 下载 chnlist（小文件，测试用）
wget -q "https://cdn.jsdelivr.net/gh/felixonmars/dnsmasq-china-list@master/accelerated-domains.china.conf" -O source/chnlist.txt
echo "✅ chnlist.txt 已下载"

wget -q "https://cdn.jsdelivr.net/gh/felixonmars/dnsmasq-china-list@master/apple.china.conf" -O source/chnlist-apple.txt
echo "✅ chnlist-apple.txt 已下载"

wget -q "https://cdn.jsdelivr.net/gh/felixonmars/dnsmasq-china-list@master/google.china.conf" -O source/chnlist-google.txt
echo "✅ chnlist-google.txt 已下载"

wget -q "https://cdn.jsdelivr.net/gh/gfwlist/gfwlist@master/gfwlist.txt" -O source/gfwlist.txt
echo "✅ gfwlist.txt 已下载"

wget -q "https://cdn.jsdelivr.net/gh/misakaio/chnroutes2@master/chnroutes.txt" -O source/chnroute.txt
echo "✅ chnroute.txt 已下载"

wget -q "https://cdn.jsdelivr.net/gh/misakaio/chnroutes2@master/chnroutes6.txt" -O source/chnroute6.txt
echo "✅ chnroute6.txt 已下载"

echo ""
echo "🔄 开始转换..."
echo ""

# 运行转换脚本
go run scripts/convert.go

echo ""
echo "📊 生成元数据..."
go run scripts/generate-metadata.go > compiled/metadata.json

echo ""
echo "✅ 转换完成！"
echo ""
echo "📁 输出文件："
ls -lh compiled/

echo ""
echo "📄 元数据："
cat compiled/metadata.json

echo ""
echo "🎉 测试成功！"
echo ""
echo "💡 提示："
echo "   • 编译后的文件在 compiled/ 目录"
echo "   • 可以在 sing-box 配置中使用这些 .srs 文件"
echo "   • 推送到 GitHub 后会自动创建 Release"

