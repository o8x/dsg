DSG
====

[Adblock Plus filters](https://help.eyeo.com/en/adblockplus/how-to-write-filters) 解析工具

## 示例

```golang
package main

import (
	"fmt"
	"log"

	"github.com/o8x/dsg"
	"github.com/o8x/dsg/pattern"
)

func main() {
	if err := dsg.Load(""); err != nil {
		log.Fatalln(err.Error())
	}

	sni := "github.com"
	d := dsg.Get()

	// 索引查询
	if d.Exist(sni) {
		fmt.Println("hit index")
	}

	// 匹配
	if rule, ok := d.Match(sni); ok {
		fmt.Println("matched by", rule.Pattern)
	}

	// 遍历规则
	d.Each(func(p *pattern.Pattern) bool {
		fmt.Println("current pattern:", p.Pattern)
		return true
	})
}
```

## 安装模块

```shell
go get -d github.com/o8x/dsg
```

## 文件生成器

产物位于 dsg/dsg.go，文件结构类似 [dsg.go](dsg.go)

```shell
go run github.com/o8x/dsg/cmd/dsg generate --url https://example.com/rules.text [--dest dsgapp] 
```

## 匹配测试

产物位于 dsg/dsg.go，文件结构类似 [dsg.go](dsg.go)

```shell
go run github.com/o8x/dsg/cmd/dsg match --url https://example.com/rules.text  domain.com ftp://domain.com
```
