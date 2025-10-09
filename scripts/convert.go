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

	// 转换 chnlist
	if err := convertChnlist(); err != nil {
		fmt.Printf("❌ chnlist 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnlist 转换成功")
	}

	// 合并 chnlist-all
	if err := convertChnlistAll(); err != nil {
		fmt.Printf("❌ chnlist-all 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnlist-all 转换成功")
	}

	// 转换 gfwlist
	if err := convertGfwlist(); err != nil {
		fmt.Printf("❌ gfwlist 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ gfwlist 转换成功")
	}

	// 转换 chnroute
	if err := convertChnroute(); err != nil {
		fmt.Printf("❌ chnroute 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnroute 转换成功")
	}

	// 转换 chnroute6
	if err := convertChnroute6(); err != nil {
		fmt.Printf("❌ chnroute6 转换失败: %v\n", err)
	} else {
		fmt.Println("✅ chnroute6 转换成功")
	}

	// 使用 geoview 转换 Geosite（可选，如果需要）
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

// convertChnlistAll 合并并转换 chnlist-all
func convertChnlistAll() error {
	allDomains := make(map[string]bool)

	// 读取三个文件
	files := []string{"source/chnlist.txt", "source/chnlist-apple.txt", "source/chnlist-google.txt"}
	for _, file := range files {
		domains, err := parseDnsmasqConf(file)
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

	fmt.Printf("📊 合并后共 %d 条域名\n", len(domains))
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

// parseGfwlist 解析 gfwlist（Base64 编码）
func parseGfwlist(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return nil, err
	}

	var domains []string
	lines := strings.Split(string(decoded), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "[") {
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
		} else if !strings.Contains(line, "/") && !strings.Contains(line, "*") {
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
