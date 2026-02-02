package main

import (
	"flag"

	"github.com/hechh/library/toolkit/pb2redis"
	"github.com/hechh/library/toolkit/pbextend"
	"github.com/hechh/library/util"
)

func main() {
	var src, dst string
	var action, pbimport string
	flag.StringVar(&action, "a", "pb", "pb,redis")
	flag.StringVar(&src, "src", "", ".pb.go文件目录")
	flag.StringVar(&dst, "dst", "", "生成文件目录")
	flag.StringVar(&pbimport, "i", "", "import ?")
	flag.Parse()

	files, err := util.Glob(src, ".*\\.pb\\.go", true)
	if err != nil {
		panic(err)
	}

	switch action {
	case "pb":
		parser := &pbextend.Parser{}
		if err := util.ParseFiles(parser, files...); err != nil {
			panic(err)
		}
		// 生成文件
		if err := parser.Gen(dst); err != nil {
			panic(err)
		}
	case "redis":
		parse := &pb2redis.Parser{}
		if err := util.ParseFiles(parse, files...); err != nil {
			panic(err)
		}
		// 生成文件
		if err := parse.Gen(dst, pbimport); err != nil {
			panic(err)
		}
	}
}
