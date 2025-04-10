package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/hildam/eino-rag/conf"
	"github.com/hildam/eino-rag/rag"
)

func main() {
	// 初始化配置
	if err := conf.Init(); err != nil {
		log.Fatal(err)
	}

	// 创建 rag
	ctx := context.Background()
	ragCli, err := rag.NewEngine(ctx, "RAG_INDEX:", "RAG_PREFIX:")
	if err != nil {
		log.Fatal(err)
	}

	// 加载文档
	if err := ragCli.AddFile(ctx, "./test/mysql-1.md"); err != nil {
		log.Fatal(err)
	}

	// 检索
	var query string
	for {
		_, _ = fmt.Scan(&query)
		rsp, err := ragCli.Query(ctx, query)
		if err != nil {
			log.Fatal(err)
		}

		// 流式输出
		for {
			output, err := rsp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(output.Content)
		}
	}
}
