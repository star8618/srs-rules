package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type RuleSetMetadata struct {
	Name        string    `json:"name"`
	Tag         string    `json:"tag"`
	Type        string    `json:"type"`
	Format      string    `json:"format"`
	FileSize    int64     `json:"file_size"`
	RuleCount   int       `json:"rule_count,omitempty"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Metadata struct {
	Version   string            `json:"version"`
	UpdatedAt time.Time         `json:"updated_at"`
	RuleSets  []RuleSetMetadata `json:"rule_sets"`
}

func main() {
	ruleSets := []RuleSetMetadata{
		{
			Name:        "中国域名白名单",
			Tag:         "chnlist",
			Type:        "domain",
			Format:      "binary",
			Description: "中国域名白名单（通用，117K+ 条）",
		},
		{
			Name:        "中国域名白名单（合并版）",
			Tag:         "chnlist-all",
			Type:        "domain",
			Format:      "binary",
			Description: "合并 chnlist + Apple + Google 中国服务",
		},
		{
			Name:        "GFW 域名列表",
			Tag:         "gfwlist",
			Type:        "domain",
			Format:      "binary",
			Description: "被 GFW 屏蔽的域名列表（5.9K+ 条）",
		},
		{
			Name:        "中国 IPv4 路由",
			Tag:         "chnroute",
			Type:        "ip",
			Format:      "binary",
			Description: "中国 IPv4 地址段（4.2K+ 条）",
		},
		{
			Name:        "中国 IPv6 路由",
			Tag:         "chnroute6",
			Type:        "ip",
			Format:      "binary",
			Description: "中国 IPv6 地址段（1.5K+ 条）",
		},
	}

	// 获取文件大小
	for i := range ruleSets {
		srsPath := filepath.Join("compiled", ruleSets[i].Tag+".srs")
		if stat, err := os.Stat(srsPath); err == nil {
			ruleSets[i].FileSize = stat.Size()
			ruleSets[i].UpdatedAt = stat.ModTime()
		}

		// 设置下载 URL（需要替换为你的 GitHub 用户名和仓库名）
		ruleSets[i].URL = fmt.Sprintf(
			"https://github.com/YOUR_USERNAME/srs-rules/releases/latest/download/%s.srs",
			ruleSets[i].Tag,
		)
	}

	metadata := Metadata{
		Version:   "1.0",
		UpdatedAt: time.Now(),
		RuleSets:  ruleSets,
	}

	// 输出 JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		fmt.Fprintf(os.Stderr, "生成元数据失败: %v\n", err)
		os.Exit(1)
	}
}
