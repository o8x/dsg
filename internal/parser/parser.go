package parser

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/o8x/dsg/pattern"
)

type Index struct {
	Pattern *pattern.Pattern `json:"pattern"`
	Hash    string           `json:"hash"`
	Index   int              `json:"index"`
}

func ToList(reader io.Reader) []string {
	patterns := map[string]interface{}{}
	buf := bufio.NewReader(reader)
	for {
		l, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}

		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "!") {
			// 不读取白名单
			if strings.Contains(l, "Whitelist") {
				break
			}
			continue
		}

		if l == "" {
			continue
		}

		patterns[l] = nil
	}

	var res []string
	for k := range patterns {
		res = append(res, k)
	}

	return res
}

func ToPatterns(reader io.Reader) []*pattern.Pattern {
	var list []*pattern.Pattern
	set := map[string]interface{}{}
	for _, it := range ToList(reader) {
		if strings.HasPrefix(it, "@@") {
			continue
		}

		p := &pattern.Pattern{
			Origin: it,
		}
		if strings.Contains(it, "*") {
			p.WildcardMatch = true
		}

		// 点开头是通配符匹配
		if strings.HasPrefix(it, ".") {
			p.WildcardMatch = true
			// .google.com -> *.google.com
			it = fmt.Sprintf("*%s", it)
		}

		if strings.HasPrefix(it, `/`) {
			reg, err := regexp.Compile(it)
			if err != nil {
				continue
			}

			p.RegMatch = reg
		}

		if strings.HasPrefix(it, `||`) {
			p.ProtoMatch = true
			it = strings.TrimPrefix(it, "||")
		}

		// 自定义规则
		if strings.HasPrefix(it, `++`) {
			p.Custom = true
			it = strings.TrimPrefix(it, "++")
		}

		if !p.ProtoMatch && strings.HasPrefix(it, `|`) {
			p.PrefixMatch = true
			it = strings.TrimPrefix(it, "|")
		}

		if strings.HasSuffix(it, "|") {
			p.SuffixMatch = true
			it = strings.TrimSuffix(it, "|")
		}

		if strings.HasPrefix(it, `http://`) {
			p.HTTP = true

			it = strings.TrimPrefix(it, "http://")
			if unescape, err := url.QueryUnescape(it); err == nil {
				it = unescape
			}
		}

		// 去掉HTTPS头
		if strings.HasPrefix(it, `https://`) {
			it = strings.TrimPrefix(it, "https://")
		}

		// 忽略同名
		if _, ok := set[it]; ok {
			continue
		}

		set[it] = nil
		p.Pattern = it
		list = append(list, p)
	}

	return list
}

func HasIndex(it *pattern.Pattern) bool {
	// 没有任何规则的 pattern，生成 hash index 用于全匹配
	return it.RegMatch == nil &&
		// ProtoMatch 剥离掉协议之后可以做全匹配
		!it.WildcardMatch && !it.PrefixMatch && !it.SuffixMatch && !it.HTTP && !it.Custom
}

func ToIndex(patterns []*pattern.Pattern) []*Index {
	var res []*Index
	for _, it := range patterns {
		if HasIndex(it) {
			res = append(res, &Index{
				Hash: it.Pattern,
				// 暂时用不着
				Index:   0,
				Pattern: it,
			})
		}
	}

	return res
}

func ToIndexMap(patterns []*pattern.Pattern) map[string]int {
	var m = map[string]int{}
	for _, it := range ToIndex(patterns) {
		m[it.Hash] = it.Index
	}

	return m
}
