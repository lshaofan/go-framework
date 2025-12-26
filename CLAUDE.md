# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个 Go 通用 SDK/框架库（go-framework），为 Go 项目提供快捷工具方法和完整的 Web 开发支持。围绕 Gin Web 框架和 GORM ORM 构建，提供从请求处理、业务逻辑、数据访问到错误处理的完整解决方案。

## 常用命令

```bash
# 依赖管理
go mod tidy              # 整理依赖
go mod download          # 下载依赖

# 编译和测试
go build ./pkg/helper    # 编译核心包
go test ./...            # 运行所有测试
go vet ./...             # 代码静态检查
```

## 代码架构

### 目录结构

```
pkg/helper/              # 核心功能包（SDK 主体）
├── action.go            # 请求处理框架（Action 接口）
├── gin_action.go        # Gin 框架的 Action 实现
├── web_response.go      # API 响应标准结构
├── web_error.go         # 错误模型定义
├── web_request.go       # 请求绑定与验证
├── dto.go               # 数据传输对象基类
├── orm.go               # GORM 泛型工具类
├── repository.go        # 仓储接口定义
├── jwt.go               # JWT 认证工具
├── middlewares.go       # Gin 中间件
├── logger.go            # 日志记录（带轮转）
├── web.go               # 服务生命周期管理
├── interface.go         # 核心接口定义
├── constants.go         # 全局常量
└── utils.go             # 工具函数
docs/                    # 模块文档
GO_FRAMEWORK_DEVELOPMENT_GUIDE.md  # 完整开发规范指南
```

### 分层架构模式

框架采用 DDD 四层架构：

```
接口层 (Interfaces)      → API 路由、Action 处理器、中间件
应用层 (Application)     → DTO 定义、常量、业务异常
领域层 (Domain)          → Service 业务逻辑、仓储接口
基础设施层 (Infrastructure) → DAO 实现、数据库配置
```

### 请求处理流程

```
请求 → Action 层 (HandleRequest) → DTO 验证 → Service 层 → 响应处理 → Response
```

核心处理方式：
```go
helper.HandleRequest(c, req, serviceFunc, ctxFunc)
```

### 核心组件

**响应结构（三层）：**
- `DefaultResult` - Service 层内部返回
- `Response` - API 响应给客户端
- `PageList[T]` - 分页数据泛型

**错误模型：**
```go
ErrorModel {
    Code       int         // 业务错误码
    Message    string      // 错误描述
    Result     interface{} // 额外信息
    HttpStatus int         // HTTP 状态码
}
```

**泛型工具类：**
- `Util[T]` - GORM 通用数据库操作
- `BaseRepository[T]` - 仓储接口泛型定义
- `PageList[T]` - 分页数据泛型

### 关键接口

**Action 接口（action.go）：**
```go
type Action interface {
    Success(data any)
    Error(err any)
    ThrowError(err *ErrorModel)
    BindParam(param any) error  // 智能绑定（自动识别 URI/JSON/Query/Form）
}
```

**仓储接口（repository.go）：**
```go
type BaseRepository[T any] interface {
    GetOneById(id uint) (T, error)
    FindByField(field string, value string) (T, error)
    ExistsById(id uint) (bool, error)
    DeleteById(id uint) error
    Create(entity T) error
    Update(entity T) error
    GetListData(request *PageRequest) (*PageList[T], error)
}
```

## 开发规范

### DTO 定义模式

```go
// 创建请求 - 使用 json 标签
type CreateRequest struct {
    helper.BaseRequest
    Name string `json:"name" binding:"required,max=100"`
}

// 列表请求 - 使用 form 标签，嵌入 ListRequest
type ListRequest struct {
    helper.BaseRequest
    helper.ListRequest
    Name string `form:"name" binding:"omitempty"`
}

// 详情/删除请求 - 使用 uri 标签
type ShowRequest struct {
    helper.BaseRequest
    ID uint `uri:"id" binding:"required,gt=0"`
}

// 更新请求 - 混合 uri 和 json 标签
type UpdateRequest struct {
    helper.BaseRequest
    ID   uint   `uri:"id" binding:"required,gt=0"`
    Name string `json:"name" binding:"required"`
}
```

### Service 层规范

```go
func (s *Service) Create(in *dto.CreateRequest) *helper.DefaultResult {
    ret := helper.NewDefaultResult()
    ret.SetResponse(s.doCreate(in))  // SetResponse(data, err)
    return ret
}
```

### 错误码规范

- 10000-10999: 通用错误
- 11000-11999: 模块A错误
- 12000-12999: 模块B错误

```go
var ErrNotFound = helper.NewErrorModel(10100, "资源不存在", nil, http.StatusNotFound)
```

## 核心依赖

- `github.com/gin-gonic/gin v1.10.0` - Web 框架
- `gorm.io/gorm v1.26.1` - ORM 框架
- `github.com/golang-jwt/jwt/v5 v5.3.0` - JWT 认证
- `github.com/sirupsen/logrus v1.9.3` - 日志库
- `github.com/go-playground/validator/v10 v10.26.0` - 参数验证

## 文档资源

- `GO_FRAMEWORK_DEVELOPMENT_GUIDE.md` - 完整开发规范指南
- `docs/MIGRATION.md` - 框架升级迁移指南
- `docs/helper/*.md` - 各模块详细文档（13个）
