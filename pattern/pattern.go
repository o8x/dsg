package pattern

import "regexp"

type Pattern struct {
	Origin        string         `json:"origin"`         // 规则原文
	Pattern       string         `json:"pattern"`        // 规则
	Exclude       bool           `json:"exclude"`        // 例外
	WildcardMatch bool           `json:"wildcard_match"` // 通配符
	RegMatch      *regexp.Regexp `json:"reg_match"`      // 正则
	PrefixMatch   bool           `json:"prefix_match"`   // 前置匹配
	SuffixMatch   bool           `json:"suffix_match"`   // 末尾匹配
	ProtoMatch    bool           `json:"proto_match"`    // 允许任何协议
	HTTP          bool           `json:"http"`           // 是否是 HTTP
	Custom        bool           `json:"custom"`         // 是否是自定义规则
}
