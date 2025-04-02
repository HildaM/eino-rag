package main

import (
	"github.com/hildam/eino-rag/conf"
	"log"
)

func main() {
	// 初始化配置
	if err := conf.Init(); err != nil {
		log.Fatal(err)
	}

}
