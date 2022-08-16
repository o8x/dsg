package dsg

import (
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/o8x/dsg/internal/downloader"
	"github.com/o8x/dsg/internal/parser"
	"github.com/o8x/dsg/internal/utils"
	"github.com/o8x/dsg/pattern"
)

var (
	d *DSG
)

func Get() *DSG {
	return d
}

type DSG struct {
	l        *sync.Mutex
	Index    map[string]int     `json:"index"`
	Patterns []*pattern.Pattern `json:"patterns"`
}

func New() *DSG {
	return &DSG{
		l: &sync.Mutex{},
	}
}

func Load(url string) error {
	reader, err := downloader.DownAsReader(url)
	if err != nil {
		return err
	}

	LoadReader(reader)
	return nil
}

func LoadReader(reader io.Reader) {
	d = New()
	d.l.Lock()
	defer d.l.Unlock()

	d.Patterns = parser.ToPatterns(reader)
	d.Index = parser.ToIndexMap(d.Patterns)
}

func (l DSG) Exist(pattern string) bool {
	l.l.Lock()
	defer l.l.Unlock()

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
