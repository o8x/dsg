package parser

import (
	"bufio"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/o8x/dsg/internal/utils"
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

		if strings.HasPrefix(it, `/`) {
			reg, err := regexp.Compile(it)
			if err != nil {
				continue
			}

			p.RegMatch = reg
		}

		if strings.HasPrefix(it, `||`) {
			p.ProtoMatch = true
			it = strings.TrimLeft(it, "||")
		}

		if !p.ProtoMatch && strings.HasPrefix(it, `|`) {
			p.PrefixMatch = true
			it = strings.TrimLeft(it, "|")
		}

		// 点开头是通配符匹配，也认为是前缀匹配
		if strings.HasPrefix(it, ".") {
			p.PrefixMatch = true
			it = strings.TrimLeft(it, ".")
		}

		if strings.HasSuffix(it, "|") {
			p.SuffixMatch = true
			it = strings.TrimRight(it, "|")
		}

		if strings.HasPrefix(it, `http://`) {
			p.HTTP = true

			it = strings.TrimLeft(it, "http://")
			if unescape, err := url.QueryUnescape(it); err == nil {
				it = unescape
			}
		}

		// 去掉HTTPS头
		if strings.HasPrefix(it, `https://`) {
			it = strings.TrimLeft(it, "https://")
		}

		p.Pattern = it
		list = append(list, p)
	}

	return list
}

func ToIndex(patterns []*pattern.Pattern) []*Index {
	var res []*Index
	for i, it := range patterns {
		// 没有任何规则的 pattern，生成 hash index 用于全匹配
		if it.RegMatch == nil &&
			!it.WildcardMatch && !it.PrefixMatch && !it.SuffixMatch && !it.ProtoMatch && !it.HTTP {

			res = append(res, &Index{
				Hash:    utils.Sha1Sum(it.Pattern),
				Index:   i,
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
