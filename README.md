# Eino-RAG

基于 [Eino](https://github.com/cloudwego/eino) 构建的 RAG (检索增强生成) 应用系统。

## 项目概述

本项目是一个基于 Golang 实现的 RAG (检索增强生成) 系统，主要用于文档检索和生成问答。系统使用 Redis 作为向量数据库，支持文档的加载、分割、索引和检索，并通过大语言模型生成回答。

## 核心功能

- 文档加载与分割：支持加载本地文档（如 Markdown 文件）并进行分割
- 向量索引：使用火山引擎的 Ark Embedding 模型生成文档向量表示
- 向量检索：基于 Redis 的向量检索能力
- 问答生成：通过 DeepSeek 大语言模型生成回答

## 系统架构

系统分为以下几个主要组件：

- **文档加载器**：负责加载文档
- **文本分割器**：将文档分割成较小的语义单元
- **Embedding 模型**：使用火山引擎的 Ark 模型生成向量表示
- **索引器**：将文档向量存储到 Redis 中
- **检索器**：从 Redis 中检索相关文档向量
- **大语言模型**：使用 DeepSeek 模型生成回答

## 快速开始

### 前置条件

- Go 1.19 或更高版本
- Redis 服务（支持向量搜索功能）
- 火山引擎 Ark API Key（用于 Embedding）
- DeepSeek API Key（用于大语言模型）

### 配置

1. 复制配置文件示例：

```bash
cp config.yaml.example config.yaml
```

2. 编辑 `config.yaml` 文件，填入您的 API Key 和其他配置：

```yaml
# RAG 配置
rag:
  dimension: 1536  # 向量嵌入维度

# 火山引擎embedding配置
embedder:
  api_key: ""  # 填入您的火山引擎 API Key
  model_id: "bge-large-zh"  # 嵌入模型ID

# 索引器配置
indexer:
  batch_size: 100  # 批量索引大小

# DeepSeek模型配置
DeekSeek:
  api_key: ""  # 填入您的 DeepSeek API Key
  model_id: "deepseek-chat"  # 大模型ID
  base_url: "https://api.deepseek.com/v1"  # API 基础URL

# Redis 配置
Redis:
  addr: "localhost:6379"  # Redis地址
  password: ""  # Redis密码，无密码则留空
```

### 运行示例

基本使用方法如下：

```go
// 初始化配置
if err := conf.Init(); err != nil {
    log.Fatal(err)
}

// 创建 RAG 引擎
ctx := context.Background()
ragCli, err := rag.NewEngine(ctx, "RAG_INDEX:", "RAG_PREFIX:")
if err != nil {
    log.Fatal(err)
}

// 加载文档
if err := ragCli.AddFile(ctx, "/path/to/your/document.md"); err != nil {
    log.Fatal(err)
}

// 查询文档并获取回答
rsp, err := ragCli.Query(ctx, "你的问题")
if err != nil {
    log.Fatal(err)
}

// 处理流式输出
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
```

## 项目结构

- `conf/`: 配置相关代码
- `rag/`: RAG 核心实现
  - `embedder.go`: Embedding 模型接口
  - `indexer.go`: 索引相关功能
  - `llm.go`: 大语言模型接口
  - `rag.go`: RAG 引擎主要实现
  - `retriever.go`: 检索相关功能
  - `splitter.go`: 文档分割功能
- `test/`: 测试文档和测试用例

## 贡献

欢迎提交 Issue 和 Pull Request

## 许可证

[Apache License 2.0](LICENSE) 