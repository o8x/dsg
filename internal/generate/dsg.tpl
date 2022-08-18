package dsg

// Code generated DO NOT EDIT

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/o8x/dsg/pattern"
)

var (
// REG DEFINE
)

type DSG struct {
	l        *sync.Mutex
	Index    map[string]int     `json:"index"`
	Patterns []*pattern.Pattern `json:"patterns"`
}

func New() *DSG {
	return &DSG{
		l:     &sync.Mutex{},
		Index: map[string]int{
            // INDEX
		},
		Patterns: []*pattern.Pattern{
            // PATTERNS
		},
	}
}

func (l DSG) Exist(pattern string) bool {
	l.l.Lock()
	defer l.l.Unlock()

	if strings.Contains(pattern, "://") {
		if seg := strings.Split(pattern, "://"); len(seg) > 1 {
			pattern = seg[1]
		}
	}

	_, ok := l.Index[utils.Sha1Sum(pattern)]
	return ok
}

func (l DSG) Each(fn func(*pattern.Pattern) bool) {
	l.l.Lock()
	defer l.l.Unlock()

	for _, it := range l.Patterns {
		if !fn(it) {
			return
		}
	}
}

func (l DSG) Match(s string) (*pattern.Pattern, bool) {
	l.l.Lock()
	defer l.l.Unlock()

	link, err := url.QueryUnescape(s)
	if err != nil {
		link = s
	}

	for _, it := range l.Patterns {
		// 暂时不处理正则和忽略
		if it.Exclude || it.RegMatch != nil {
			continue
		}

		if it.PrefixMatch && strings.HasPrefix(s, it.Pattern) {
			return it, true
		}

		if it.SuffixMatch && strings.HasSuffix(s, it.Pattern) {
			return it, true
		}

		if it.WildcardMatch {
			if match, _ := filepath.Match(it.Pattern, s); match {
				return it, true
			}
		}

		if it.HTTP && strings.Contains(link, it.Pattern) {
			return it, true
		}
	}

	return nil, false
}

func Sha1Sum(str string) string {
	sum := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", sum)
}
