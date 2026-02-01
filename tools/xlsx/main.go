package main

import (
	"flag"

	"github.com/hechh/library/toolkit"
	"github.com/hechh/library/toolkit/xlsx2code"
	"github.com/hechh/library/toolkit/xlsx2data"
	"github.com/hechh/library/toolkit/xlsx2proto"
	"github.com/hechh/library/util"
)

func main() {
	var action string
	var src, dst string
	var option, pkgname, desc string
	var pbimport string
	flag.StringVar(&action, "a", "proto", "proto,data,code")
	flag.StringVar(&src, "src", ".", "源目录")
	flag.StringVar(&dst, "dst", ".", "目的目录")
	flag.StringVar(&pbimport, "i", "", "import ?")
	flag.StringVar(&option, "o", "", "option go_package = ?")
	flag.StringVar(&pkgname, "p", "", "package ?")
	flag.StringVar(&desc, "d", "", "descriptor_set_out生成文件")
	flag.Parse()

	// 加载所有xlsx文件
	files, err := util.Glob(src, ".*\\.xlsx", true)
	if err != nil {
		panic(err)
	}

	switch action {
	case "proto":
		// 解析文件
		p := xlsx2proto.NewMsgParser()
		for _, filename := range files {
			if err := p.ParseFile(filename); err != nil {
				panic(err)
			}
		}
		// 生成文件
		if err := p.Gen(pkgname, option, dst); err != nil {
			panic(err)
		}
	case "data":
		if err := toolkit.Init(pkgname, desc); err != nil {
			panic(err)
		}
		// 解析 xlsx 文件
		parse := xlsx2data.NewMsgParser()
		for _, filename := range files {
			if err := parse.ParseFile(filename); err != nil {
				panic(err)
			}
		}
		// 生成 data 文件
		if err := parse.Gen(dst); err != nil {
			panic(err)
		}
	case "code":
		// 加载 proto文件
		parse, err := xlsx2code.NewMsgParser(pkgname, desc)
		if err != nil {
			panic(err)
		}
		// 解析 xlsx 文件
		for _, filename := range files {
			if err := parse.ParseFile(filename); err != nil {
				panic(err)
			}
		}
		// 生成 data 文件
		if err := parse.Gen(dst, pbimport); err != nil {
			panic(err)
		}
	}
}
