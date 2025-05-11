# go-framework

一个为 Go 项目提供快捷工具方法的通用 SDK。

## 目录结构

```
go-framework/
├── cmd/            # 可选：命令行工具（如有）
├── pkg/            # 主工具包目录，按功能子包组织
│   ├── stringutil/ # 字符串处理工具
│   ├── fileutil/   # 文件操作工具
│   └── ...         # 未来可扩展更多工具子包
├── internal/       # 内部实现包，不对外暴露
├── examples/       # 用法示例
├── docs/           # 项目文档
├── test/           # 集成测试
├── scripts/        # 可选：自动化脚本
├── .gitignore      # Git 忽略文件
├── go.mod          # Go 依赖管理
├── go.sum          # Go 依赖校验
├── LICENSE         # 许可证
└── README.md       # 项目说明
```

## 快速开始

1. 安装依赖：
   ```sh
   go get github.com/lshaofan/go-framework/pkg/...
   ```
2. 按需引入工具包：
   ```go
   import "github.com/lshaofan/go-framework/pkg/stringutil"
   ```

## 贡献指南

- 每类工具建议单独一个子包，便于维护和按需引入。
- 示例代码请放在 `examples/` 目录。
- 详细文档请放在 `docs/` 目录。
- 内部实现细节请放在 `internal/` 目录。
