package generate

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/o8x/dsg/internal/parser"
)

var (
	//go:embed dsg.tpl
	tpl []byte
)

func Template(reader io.Reader) string {
	var (
		patternBuilder strings.Builder
		regexpBuilder  strings.Builder
		indexBuilder   strings.Builder
	)

	patterns := parser.ToPatterns(reader)
	for i, it := range patterns {
		regName := "nil"
		if it.RegMatch != nil {
			regName = fmt.Sprintf("reg%d", i)
			regexpBuilder.WriteString("\t")
			regexpBuilder.WriteString(fmt.Sprintf("%s, _ = regexp.Compile(`%s`)", regName, it.Pattern))
			regexpBuilder.WriteString("\n")
			// 生成正则之后，将 Pattern 设置为空
			it.Pattern = ""
		}

		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("\t\t\t // %s\n", it.Origin))
		builder.WriteString("\t\t\t")
		builder.WriteString("{")
		builder.WriteString("Exclude: false, ")
		builder.WriteString(fmt.Sprintf(`Pattern: "%s", `, it.Pattern))
		builder.WriteString(fmt.Sprintf("WildcardMatch: %t, ", it.WildcardMatch))
		builder.WriteString(fmt.Sprintf("RegMatch: %s, ", regName))
		builder.WriteString(fmt.Sprintf("PrefixMatch: %t, ", it.PrefixMatch))
		builder.WriteString(fmt.Sprintf("SuffixMatch: %t, ", it.SuffixMatch))
		builder.WriteString(fmt.Sprintf("ProtoMatch: %t, ", it.ProtoMatch))
		builder.WriteString(fmt.Sprintf("HTTP: %t", it.HTTP))
		builder.WriteString("}")

		// 新增一行规则
		patternBuilder.WriteString(builder.String())
		patternBuilder.WriteString(",\n")
	}

	// 生成索引
	for _, ind := range parser.ToIndex(patterns) {
		builder := strings.Builder{}
		builder.WriteString("\t\t\t")
		builder.WriteString(fmt.Sprintf("// %s", ind.Pattern.Pattern))
		builder.WriteString("\n")
		builder.WriteString("\t\t\t")
		builder.WriteString(fmt.Sprintf(`"%s": %d`, ind.Hash, ind.Index))

		// 新增一行索引
		indexBuilder.WriteString(builder.String())
		indexBuilder.WriteString(",\n")
	}

	tpl := string(tpl)
	tpl = strings.ReplaceAll(tpl, "// PATTERNS", patternBuilder.String())
	tpl = strings.ReplaceAll(tpl, "// INDEX", indexBuilder.String())
	tpl = strings.ReplaceAll(tpl, "// REG DEFINE", regexpBuilder.String())

	return tpl
}
