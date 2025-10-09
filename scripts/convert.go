package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func main() {
	fmt.Println("🚀 开始转换规则集...")

	// 转换单个源规则（保留兼容性）
	if err := convertChnlist(); err != nil {
		fmt.Printf("❌ chnlist 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnlist 转换成功")
	}

	if err := convertGfwlist(); err != nil {
		fmt.Printf("❌ gfwlist 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ gfwlist 转换成功")
	}

	if err := convertChnroute(); err != nil {
		fmt.Printf("❌ chnroute 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnroute 转换成功")
	}

	if err := convertChnroute6(); err != nil {
		fmt.Printf("❌ chnroute6 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnroute6 转换成功")
	}

	// === 合并多源规则（增强版）===
	fmt.Println("\n📦 开始合并多源规则...")

	// 合并 chnlist-all（7 个源）
	if err := convertChnlistAll(); err != nil {
		fmt.Printf("❌ chnlist-all 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnlist-all 转换成功")
	}

	// 合并 gfwlist-all（4 个源）
	if err := convertGfwlistAll(); err != nil {
		fmt.Printf("❌ gfwlist-all 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ gfwlist-all 转换成功")
	}

	// 合并 chnroute-all（6 个源）
	if err := convertChnrouteAll(); err != nil {
		fmt.Printf("❌ chnroute-all 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnroute-all 转换成功")
	}

	// 合并 chnroute6-all（3 个源）
	if err := convertChnroute6All(); err != nil {
		fmt.Printf("❌ chnroute6-all 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnroute6-all 转换成功")
	}

	// 使用 geoview 转换 Geosite/GeoIP（可选，需要 geoview 工具）
	// 如果 geoview 不可用，只会打印警告，不影响主流程
	convertGeoview()

	fmt.Println("🎉 所有规则集转换完成！")
}

// convertChnlist 转换 chnlist
func convertChnlist() error {
	domains, err := parseDnsmasqConf("source/chnlist.txt")
	if err != nil {
		return err
	}

	return writeSRS("compiled/chnlist.srs", domains, nil)
}

// convertChnlistAll 合并并转换 chnlist-all（7 个源）
func convertChnlistAll() error {
	allDomains := make(map[string]bool)

	// 7 个 chnlist 源文件（felixonmars 3个 + Loyalsoldier 3个 + ios_rule_script 1个）
	files := []struct {
		path   string
		parser func(string) ([]string, error)
	}{
		{"source/chnlist.txt", parseDnsmasqConf},
		{"source/chnlist-apple.txt", parseDnsmasqConf},
		{"source/chnlist-google.txt", parseDnsmasqConf},
		{"source/chnlist-loyalsoldier.txt", parseTextLines},        // Loyalsoldier/china-list
		{"source/chnlist-loyalsoldier-apple.txt", parseTextLines},  // Loyalsoldier/apple-cn
		{"source/chnlist-loyalsoldier-google.txt", parseTextLines}, // Loyalsoldier/google-cn
		{"source/chnlist-ios.txt", parseTextLines},                 // ios_rule_script/ChinaMax_Domain
	}

	for _, f := range files {
		domains, err := f.parser(f.path)
		if err != nil {
			fmt.Printf("⚠️  读取 %s 失败: %v\n", f.path, err)
			continue
		}
		for _, d := range domains {
			allDomains[d] = true
		}
	}

	// 转为切片
	domains := make([]string, 0, len(allDomains))
	for d := range allDomains {
		domains = append(domains, d)
	}

	fmt.Printf("📊 chnlist-all 合并后共 %d 条域名\n", len(domains))
	return writeSRS("compiled/chnlist-all.srs", domains, nil)
}

// convertGfwlist 转换 gfwlist
func convertGfwlist() error {
	domains, err := parseGfwlist("source/gfwlist.txt")
	if err != nil {
		return err
	}

	return writeSRS("compiled/gfwlist.srs", domains, nil)
}

// convertChnroute 转换 chnroute
func convertChnroute() error {
	cidrs, err := parseTextLines("source/chnroute.txt")
	if err != nil {
		return err
	}

	return writeSRS("compiled/chnroute.srs", nil, cidrs)
}

// convertChnroute6 转换 chnroute6
func convertChnroute6() error {
	cidrs, err := parseTextLines("source/chnroute6.txt")
	if err != nil {
		return err
	}

	return writeSRS("compiled/chnroute6.srs", nil, cidrs)
}

// ========== 多源合并函数 ==========

// convertGfwlistAll 合并多个 gfwlist 源（4 个源）
func convertGfwlistAll() error {
	allDomains := make(map[string]bool)

	// 4 个 gfwlist 源文件
	files := []string{
		"source/gfwlist-v2fly.txt",
		"source/gfwlist-loyalsoldier.txt",
		"source/gfwlist-loukky.txt",
		"source/gfwlist-original.txt",
	}

	for _, file := range files {
		domains, err := parseGfwlist(file)
		if err != nil {
			fmt.Printf("⚠️  读取 %s 失败: %v\n", file, err)
			continue
		}
		for _, d := range domains {
			allDomains[d] = true
		}
	}

	// 转为切片
	domains := make([]string, 0, len(allDomains))
	for d := range allDomains {
		domains = append(domains, d)
	}

	fmt.Printf("📊 gfwlist-all 合并后共 %d 条域名\n", len(domains))
	return writeSRS("compiled/gfwlist-all.srs", domains, nil)
}

// convertChnrouteAll 合并多个 chnroute 源（6 个源）
func convertChnrouteAll() error {
	allCidrs := make(map[string]bool)

	// 6 个 chnroute 源文件
	files := []string{
		"source/chnroute-gaoyifan.txt",
		"source/chnroute-clang.txt",
		"source/chnroute-clang-cidr.txt",
		"source/chnroute-soffchen.txt",
		"source/chnroute-hackl0us.txt",
		"source/chnroute-ios.txt",
	}

	for _, file := range files {
		cidrs, err := parseTextLines(file)
		if err != nil {
			fmt.Printf("⚠️  读取 %s 失败: %v\n", file, err)
			continue
		}
		for _, c := range cidrs {
			allCidrs[c] = true
		}
	}

	// 转为切片
	cidrs := make([]string, 0, len(allCidrs))
	for c := range allCidrs {
		cidrs = append(cidrs, c)
	}

	fmt.Printf("📊 chnroute-all 合并后共 %d 条 IPv4 CIDR\n", len(cidrs))
	return writeSRS("compiled/chnroute-all.srs", nil, cidrs)
}

// convertChnroute6All 合并多个 chnroute6 源（3 个源）
func convertChnroute6All() error {
	allCidrs := make(map[string]bool)

	// 3 个 chnroute6 源文件
	files := []string{
		"source/chnroute6-gaoyifan.txt",
		"source/chnroute6-clang.txt",
		"source/chnroute6-ios.txt",
	}

	for _, file := range files {
		cidrs, err := parseTextLines(file)
		if err != nil {
			fmt.Printf("⚠️  读取 %s 失败: %v\n", file, err)
			continue
		}
		for _, c := range cidrs {
			allCidrs[c] = true
		}
	}

	// 转为切片
	cidrs := make([]string, 0, len(allCidrs))
	for c := range allCidrs {
		cidrs = append(cidrs, c)
	}

	fmt.Printf("📊 chnroute6-all 合并后共 %d 条 IPv6 CIDR\n", len(cidrs))
	return writeSRS("compiled/chnroute6-all.srs", nil, cidrs)
}

// convertGeoview 使用 geoview 转换 Geosite/GeoIP（可选）
func convertGeoview() {
	// Geosite tags
	geositeTags := []string{"cn", "google", "steam", "netflix", "disney", "openai", "github", "category-games"}
	for _, tag := range geositeTags {
		cmd := exec.Command("geoview", "-type", "geosite", "-action", "convert",
			"-input", "source/geosite.dat", "-list", tag,
			"-output", fmt.Sprintf("compiled/geosite-%s.srs", tag))
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  geosite-%s 转换失败: %v\n", tag, err)
		} else {
			fmt.Printf("✅ geosite-%s 转换成功\n", tag)
		}
	}

	// GeoIP tags
	geoipTags := []string{"cn", "us"}
	for _, tag := range geoipTags {
		cmd := exec.Command("geoview", "-type", "geoip", "-action", "convert",
			"-input", "source/geoip.dat", "-list", tag,
			"-output", fmt.Sprintf("compiled/geoip-%s.srs", tag))
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  geoip-%s 转换失败: %v\n", tag, err)
		} else {
			fmt.Printf("✅ geoip-%s 转换成功\n", tag)
		}
	}
}

// parseDnsmasqConf 解析 dnsmasq 配置文件
func parseDnsmasqConf(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 server=/domain/114.114.114.114 格式
		if strings.HasPrefix(line, "server=/") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 && parts[1] != "" {
				domains = append(domains, parts[1])
			}
		}
	}

	return domains, scanner.Err()
}

// parseGfwlist 解析 gfwlist（支持 Base64 编码和纯文本）
func parseGfwlist(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var textContent string

	// 尝试 Base64 解码（兼容旧格式）
	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err == nil && len(decoded) > 0 {
		// Base64 解码成功
		textContent = string(decoded)
	} else {
		// 使用原始文本（新格式 gfw.txt）
		textContent = string(content)
	}

	var domains []string
	lines := strings.Split(textContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "[") || strings.HasPrefix(line, "#") {
			continue
		}

		// 提取域名
		if strings.HasPrefix(line, "||") {
			domain := strings.TrimPrefix(line, "||")
			domain = strings.TrimSuffix(domain, "^")
			domain = strings.Split(domain, "/")[0]
			if domain != "" {
				domains = append(domains, domain)
			}
		} else if strings.HasPrefix(line, ".") {
			domains = append(domains, strings.TrimPrefix(line, "."))
		} else if strings.Contains(line, "://") {
			// 处理 URL 格式（如 http://example.com）
			parts := strings.Split(line, "://")
			if len(parts) > 1 {
				domain := strings.Split(parts[1], "/")[0]
				domain = strings.Split(domain, ":")[0]
				if domain != "" {
					domains = append(domains, domain)
				}
			}
		} else if !strings.Contains(line, "/") && !strings.Contains(line, "*") && !strings.Contains(line, "@") {
			// 纯域名
			domains = append(domains, line)
		}
	}

	return domains, nil
}

// parseTextLines 解析纯文本行（CIDR 列表）
func parseTextLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

// writeSRS 写入 SRS 文件
func writeSRS(outputPath string, domains []string, ipCidrs []string) error {
	var headlessRule option.DefaultHeadlessRule

	if len(domains) > 0 {
		headlessRule.DomainSuffix = domains
	}
	if len(ipCidrs) > 0 {
		headlessRule.IPCIDR = ipCidrs
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = srs.Write(outputFile, option.PlainRuleSet{
		Rules: []option.HeadlessRule{
			{
				Type:           C.RuleTypeDefault,
				DefaultOptions: headlessRule,
			},
		},
	}, 1)

	if err != nil {
		return err
	}

	// 获取文件大小
	stat, _ := outputFile.Stat()
	fmt.Printf("📦 %s: %.2f KB (%d 条规则)\n", filepath.Base(outputPath), float64(stat.Size())/1024, len(domains)+len(ipCidrs))

	return nil
}
