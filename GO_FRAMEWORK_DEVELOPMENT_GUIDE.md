# go-framework 开发规范指南

本文档为 AI Agent 和开发者提供使用 `github.com/lshaofan/go-framework` 包进行开发的完整规范指南。

## 目录

- [1. 包结构概览](#1-包结构概览)
- [2. 核心类型定义](#2-核心类型定义)
- [3. 分层架构规范](#3-分层架构规范)
- [4. DTO 定义规范](#4-dto-定义规范)
- [5. Service 层规范](#5-service-层规范)
- [6. Action 层规范](#6-action-层规范)
- [7. Repository 层规范](#7-repository-层规范)
- [8. 错误处理规范](#8-错误处理规范)
- [9. 完整示例](#9-完整示例)
- [10. 代码生成模板](#10-代码生成模板)

---

## 1. 包结构概览

```
github.com/lshaofan/go-framework/
└── pkg/
    └── helper/
        ├── dto.go           # 请求基类 BaseRequest, IBaseRequest
        ├── web_response.go  # 响应结构 DefaultResult, PageList, Response
        ├── web_error.go     # 错误模型 ErrorModel
        ├── web_request.go   # 请求验证 ListRequest, Validate
        ├── action.go        # Action 处理 HandleRequest, BaseAction
        ├── gin_action.go    # Gin 实现 GinActionImpl
        ├── orm.go           # ORM 工具 Util, PageRequest, Paginate
        ├── repository.go    # 仓储接口 BaseRepository
        ├── jwt.go           # JWT 工具 JWTUtil, JWTResponse
        ├── interface.go     # 核心接口定义
        ├── constants.go     # 常量定义
        └── utils.go         # 工具函数
```

### 导入方式

```go
import "github.com/lshaofan/go-framework/pkg/helper"

// 如果项目中有同名 helper 包，使用别名
import web "github.com/lshaofan/go-framework/pkg/helper"
```

---

## 2. 核心类型定义

### 2.1 请求基类 (BaseRequest)

```go
// BaseRequest 所有请求 DTO 必须嵌入此结构
type BaseRequest struct {
    Ctx context.Context `json:"-" form:"-"`
}

// IBaseRequest 请求接口
type IBaseRequest interface {
    SetContext(ctx context.Context)
    GetContext() context.Context
}
```

### 2.2 列表请求 (ListRequest)

```go
// ListRequest 列表查询通用参数
type ListRequest struct {
    Page     int    `form:"page" json:"page" query:"page" binding:"required"`
    PageSize int    `form:"page_size" json:"page_size" query:"page_size" binding:"required"`
    Order    string `form:"order" json:"order" query:"order"`
    Field    string `form:"field" json:"field" query:"field"`
}
```

### 2.3 响应结构 (DefaultResult)

```go
// DefaultResult Service 层统一返回结构
type DefaultResult struct {
    Err  *ErrorModel `json:"err"`
    Data any         `json:"data"`
}

// 核心方法
func NewDefaultResult() *DefaultResult
func (r *DefaultResult) IsError() bool
func (r *DefaultResult) GetError() *ErrorModel
func (r *DefaultResult) SetError(err error)
func (r *DefaultResult) GetData() any
func (r *DefaultResult) SetData(data any)
func (r *DefaultResult) SetResponse(data any, err error)  // 推荐使用
```

### 2.4 分页结构 (PageList)

```go
// PageList 分页数据结构（泛型）
type PageList[T interface{}] struct {
    Total    int64 `json:"total"`
    Data     []T   `json:"data"`
    Page     int   `json:"page"`
    PageSize int   `json:"page_size"`
}
```

### 2.5 错误模型 (ErrorModel)

```go
// ErrorModel 标准错误结构
type ErrorModel struct {
    Code       int         `json:"code"`
    Message    string      `json:"message"`
    Result     interface{} `json:"result"`
    HttpStatus int         `json:"httpStatus" swaggerignore:"true"`
}

// 创建错误
func NewErrorModel(code int, message string, result interface{}, httpStatus int) *ErrorModel
```

### 2.6 API 响应 (Response)

```go
// Response API 统一响应格式
type Response struct {
    Code    int    `json:"code"`
    Result  any    `json:"result"`
    Message string `json:"message"`
}
```

---

## 3. 分层架构规范

项目采用 DDD 四层架构：

```
项目结构/
├── src/
│   ├── application/           # 应用层
│   │   ├── dto/              # 数据传输对象
│   │   ├── constants/        # 常量定义
│   │   └── exp/              # 错误定义
│   ├── domain/               # 领域层
│   │   ├── models/           # 领域模型
│   │   ├── services/         # 领域服务
│   │   ├── repository/       # 仓储接口
│   │   └── helper/           # 领域辅助工具
│   ├── infrastructure/       # 基础设施层
│   │   ├── dao/              # 数据访问对象
│   │   ├── config/           # 配置管理
│   │   └── utils/            # 工具函数
│   └── interfaces/           # 接口层
│       ├── api/              # 路由定义
│       ├── actions/          # 请求处理器
│       ├── middleware/       # 中间件
│       └── global/           # 全局配置
```

### 数据流向

```
Request → Action → Service → Repository → Database
                      ↓
                 DefaultResult
                      ↓
Response ← Action ← (HandleResult)
```

---

## 4. DTO 定义规范

### 4.1 基本规则

1. **必须嵌入 `helper.BaseRequest`**
2. **列表查询嵌入 `helper.ListRequest`**
3. **使用标准标签**: `json`, `form`, `uri`, `query`, `binding`
4. **文件命名**: `{模块}_{操作}.go`（如 `organization_create.go`）

### 4.2 请求 DTO 模板

#### 创建请求

```go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

// OrganizationCreateRequest 创建组织请求
type OrganizationCreateRequest struct {
    helper.BaseRequest
    Name        string `json:"name" binding:"required,max=100"`        // 名称（必填）
    Code        string `json:"code" binding:"required,max=50"`         // 编码（必填）
    Description string `json:"description" binding:"omitempty,max=500"` // 描述（可选）
    Status      uint   `json:"status" binding:"omitempty,oneof=1 2"`   // 状态（可选）
}

// OrganizationCreateResponse 创建组织响应
type OrganizationCreateResponse struct {
    OrganizationData
}
```

#### 列表请求

```go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

// OrganizationListRequest 组织列表请求
type OrganizationListRequest struct {
    helper.BaseRequest
    helper.ListRequest                                                    // 嵌入分页参数
    Name   string `form:"name" json:"name" query:"name" binding:"omitempty"`       // 名称（模糊查询）
    Status *uint  `form:"status" json:"status" query:"status" binding:"omitempty"` // 状态筛选
}

// OrganizationListResponse 组织列表响应
type OrganizationListResponse struct {
    Data     []OrganizationData `json:"data"`
    Total    int64              `json:"total"`
    Page     int                `json:"page"`
    PageSize int                `json:"page_size"`
}
```

#### 详情请求（URI 参数）

```go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

// OrganizationShowRequest 获取组织详情请求
type OrganizationShowRequest struct {
    helper.BaseRequest
    ID uint `uri:"id" binding:"required,gt=0"` // 组织ID（URI参数）
}

// OrganizationShowResponse 获取组织详情响应
type OrganizationShowResponse struct {
    OrganizationData
}
```

#### 更新请求（URI + Body 混合）

```go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

// OrganizationUpdateRequest 更新组织请求
type OrganizationUpdateRequest struct {
    helper.BaseRequest
    ID          uint   `uri:"id" binding:"required,gt=0"`              // 组织ID（URI）
    Name        string `json:"name" binding:"required,max=100"`        // 名称
    Code        string `json:"code" binding:"required,max=50"`         // 编码
    Description string `json:"description" binding:"omitempty,max=500"` // 描述
    Status      uint   `json:"status" binding:"omitempty,oneof=1 2"`   // 状态
}

// OrganizationUpdateResponse 更新组织响应
type OrganizationUpdateResponse struct {
    OrganizationData
}
```

#### 部分更新请求（Patch）

```go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

// OrganizationPatchRequest 部分更新组织请求
type OrganizationPatchRequest struct {
    helper.BaseRequest
    ID          uint    `uri:"id" binding:"required,gt=0"`                      // 组织ID
    Name        *string `json:"name" binding:"omitempty,max=100"`               // 名称（可选）
    Description *string `json:"description" binding:"omitempty,max=500"`        // 描述（可选）
    Status      *uint   `json:"status" binding:"omitempty,oneof=1 2"`           // 状态（可选）
}

// OrganizationPatchResponse 部分更新响应
type OrganizationPatchResponse struct {
    OrganizationData
}
```

#### 删除请求

```go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

// OrganizationDeleteRequest 删除组织请求
type OrganizationDeleteRequest struct {
    helper.BaseRequest
    ID uint `uri:"id" binding:"required,gt=0"` // 组织ID
}

// OrganizationDeleteResponse 删除组织响应
type OrganizationDeleteResponse struct{}
```

### 4.3 数据结构定义

```go
package dto

// OrganizationData 组织数据结构（用于响应）
type OrganizationData struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Code        string `json:"code"`
    Description string `json:"description,omitempty"`
    Status      uint   `json:"status,omitempty"`
    CreatedAt   string `json:"created_at,omitempty"`
    UpdatedAt   string `json:"updated_at,omitempty"`
}
```

### 4.4 参数绑定标签说明

| 标签 | 用途 | 示例 |
|------|------|------|
| `json` | JSON Body 参数 | `json:"name"` |
| `form` | Form 表单参数 | `form:"name"` |
| `query` | URL Query 参数 | `query:"page"` |
| `uri` | URL Path 参数 | `uri:"id"` |
| `binding` | 验证规则 | `binding:"required,max=100"` |

### 4.5 常用验证规则

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填 | `binding:"required"` |
| `omitempty` | 可选（空值跳过验证） | `binding:"omitempty"` |
| `max=N` | 最大长度/值 | `binding:"max=100"` |
| `min=N` | 最小长度/值 | `binding:"min=1"` |
| `gt=N` | 大于 | `binding:"gt=0"` |
| `gte=N` | 大于等于 | `binding:"gte=1"` |
| `oneof=...` | 枚举值 | `binding:"oneof=1 2 3"` |
| `email` | 邮箱格式 | `binding:"email"` |
| `url` | URL 格式 | `binding:"url"` |

---

## 5. Service 层规范

### 5.1 基本规则

1. **返回类型**: 统一返回 `*helper.DefaultResult`
2. **使用 `SetResponse`**: 推荐使用 `ret.SetResponse(data, err)` 设置响应
3. **文件命名**: `{模块}_{操作}.go`（如 `organization_create.go`）
4. **依赖注入**: 通过构造函数注入 Repository

### 5.2 Service 结构定义

```go
package services

import (
    "your-project/src/application/dto"
    "your-project/src/domain/repository"
    "your-project/src/infrastructure/dao"
)

type OrganizationService struct {
    repo repository.IOrganizationRepository
}

func NewOrganizationService() *OrganizationService {
    return &OrganizationService{
        repo: dao.NewOrganization(),
    }
}
```

### 5.3 CRUD 方法模板

#### Create 方法

```go
package services

import (
    "github.com/lshaofan/go-framework/pkg/helper"
    "your-project/src/application/dto"
    "your-project/src/application/exp"
    "your-project/src/domain/models"
)

// Create 创建组织
func (s *OrganizationService) Create(in *dto.OrganizationCreateRequest) *helper.DefaultResult {
    ret := helper.NewDefaultResult()
    ret.SetResponse(s.createData(in))
    return ret
}

// createData 内部创建逻辑
func (s *OrganizationService) createData(in *dto.OrganizationCreateRequest) (*dto.OrganizationCreateResponse, error) {
    // 1. 业务验证
    exist, err := s.repo.FindByField("name", in.Name)
    if err != nil {
        return nil, err
    }
    if exist.ID > 0 {
        return nil, exp.ErrOrganizationNameExists
    }

    // 2. 构建模型
    model := &models.Organization{
        Name:        in.Name,
        Code:        in.Code,
        Description: in.Description,
        Status:      in.Status,
    }

    // 3. 持久化
    if err := s.repo.Create(*model); err != nil {
        return nil, err
    }

    // 4. 返回响应
    return &dto.OrganizationCreateResponse{
        OrganizationData: dto.OrganizationData{
            ID:   model.ID,
            Name: model.Name,
            Code: model.Code,
        },
    }, nil
}
```

#### List 方法

```go
package services

import (
    "github.com/lshaofan/go-framework/pkg/helper"
    "your-project/src/application/dto"
)

// List 获取组织列表
func (s *OrganizationService) List(in *dto.OrganizationListRequest) *helper.DefaultResult {
    ret := helper.NewDefaultResult()

    // 1. 构建分页请求
    pageReq := helper.NewPageReq()
    pageReq.Page = in.Page
    pageReq.PageSize = in.PageSize

    // 2. 添加筛选条件
    if in.Name != "" {
        pageReq.Where["name"] = in.Name
    }
    if in.Status != nil {
        pageReq.Where["status"] = *in.Status
    }

    // 3. 查询数据
    list, err := s.repo.GetListData(pageReq)
    if err != nil {
        ret.SetError(err)
        return ret
    }

    // 4. 转换响应
    response := &dto.OrganizationListResponse{
        Total:    list.Total,
        Page:     list.Page,
        PageSize: list.PageSize,
        Data:     s.convertListToDTO(list.Data),
    }

    ret.SetData(response)
    return ret
}

// convertListToDTO 转换列表数据为 DTO
func (s *OrganizationService) convertListToDTO(list []models.Organization) []dto.OrganizationData {
    result := make([]dto.OrganizationData, 0, len(list))
    for _, item := range list {
        result = append(result, dto.OrganizationData{
            ID:        item.ID,
            Name:      item.Name,
            Code:      item.Code,
            Status:    item.Status,
            CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
        })
    }
    return result
}
```

#### Show 方法

```go
package services

import (
    "errors"

    "github.com/lshaofan/go-framework/pkg/helper"
    "gorm.io/gorm"
    "your-project/src/application/dto"
    "your-project/src/application/exp"
)

// Show 获取组织详情
func (s *OrganizationService) Show(in *dto.OrganizationShowRequest) *helper.DefaultResult {
    ret := helper.NewDefaultResult()
    ret.SetResponse(s.showData(in))
    return ret
}

func (s *OrganizationService) showData(in *dto.OrganizationShowRequest) (*dto.OrganizationShowResponse, error) {
    // 1. 查询数据
    model, err := s.repo.GetOneById(in.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, exp.ErrOrganizationNotFound
        }
        return nil, err
    }

    // 2. 返回响应
    return &dto.OrganizationShowResponse{
        OrganizationData: dto.OrganizationData{
            ID:          model.ID,
            Name:        model.Name,
            Code:        model.Code,
            Description: model.Description,
            Status:      model.Status,
            CreatedAt:   model.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:   model.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }, nil
}
```

#### Update 方法

```go
package services

import (
    "errors"

    "github.com/lshaofan/go-framework/pkg/helper"
    "gorm.io/gorm"
    "your-project/src/application/dto"
    "your-project/src/application/exp"
)

// Update 更新组织
func (s *OrganizationService) Update(in *dto.OrganizationUpdateRequest) *helper.DefaultResult {
    ret := helper.NewDefaultResult()
    ret.SetResponse(s.updateData(in))
    return ret
}

func (s *OrganizationService) updateData(in *dto.OrganizationUpdateRequest) (*dto.OrganizationUpdateResponse, error) {
    // 1. 检查记录是否存在
    model, err := s.repo.GetOneById(in.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, exp.ErrOrganizationNotFound
        }
        return nil, err
    }

    // 2. 检查名称唯一性（排除自身）
    if in.Name != model.Name {
        exist, _ := s.repo.FindByField("name", in.Name)
        if exist.ID > 0 && exist.ID != in.ID {
            return nil, exp.ErrOrganizationNameExists
        }
    }

    // 3. 更新字段
    model.Name = in.Name
    model.Code = in.Code
    model.Description = in.Description
    model.Status = in.Status

    // 4. 保存更新
    if err := s.repo.Update(model); err != nil {
        return nil, err
    }

    // 5. 返回响应
    return &dto.OrganizationUpdateResponse{
        OrganizationData: dto.OrganizationData{
            ID:   model.ID,
            Name: model.Name,
            Code: model.Code,
        },
    }, nil
}
```

#### Delete 方法

```go
package services

import (
    "errors"

    "github.com/lshaofan/go-framework/pkg/helper"
    "gorm.io/gorm"
    "your-project/src/application/dto"
    "your-project/src/application/exp"
)

// Delete 删除组织
func (s *OrganizationService) Delete(in *dto.OrganizationDeleteRequest) *helper.DefaultResult {
    ret := helper.NewDefaultResult()

    // 1. 检查记录是否存在
    _, err := s.repo.GetOneById(in.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            ret.SetError(exp.ErrOrganizationNotFound)
            return ret
        }
        ret.SetError(err)
        return ret
    }

    // 2. 执行删除
    if err := s.repo.DeleteById(in.ID); err != nil {
        ret.SetError(err)
        return ret
    }

    ret.SetData(&dto.OrganizationDeleteResponse{})
    return ret
}
```

---

## 6. Action 层规范

### 6.1 基本规则

1. **使用 `helper.HandleRequest`**: 统一请求处理入口
2. **Swagger 注解**: 必须添加完整的 API 文档注解
3. **文件命名**: `{模块}_actions.go`

### 6.2 HandleRequest 签名

```go
func HandleRequest(
    c *gin.Context,
    req IBaseRequest,
    serviceCall func(IBaseRequest) *DefaultResult,
    ctxFunc BaseCtxFunc,
)

// BaseCtxFunc 上下文处理函数类型
type BaseCtxFunc func(ctx context.Context, c *gin.Context) context.Context
```

### 6.3 Action 模板

```go
package actions

import (
    "github.com/gin-gonic/gin"
    "github.com/lshaofan/go-framework/pkg/helper"
    "your-project/src/application/dto"
    "your-project/src/domain/services"
)

// OrganizationCreate 创建组织
// @Summary      创建组织
// @Description  创建新的组织信息（需要认证）
// @Tags         组织管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.OrganizationCreateRequest true "组织信息"
// @Success      200 {object} helper.DefaultResult{data=dto.OrganizationCreateResponse} "创建成功"
// @Failure      400 {object} helper.DefaultResult "请求参数错误"
// @Failure      401 {object} helper.DefaultResult "未授权"
// @Failure      500 {object} helper.DefaultResult "服务器内部错误"
// @Router       /admin/organization [post]
func OrganizationCreate(c *gin.Context) {
    req := new(dto.OrganizationCreateRequest)
    helper.HandleRequest(c, req, func(r helper.IBaseRequest) *helper.DefaultResult {
        return services.NewOrganizationService().Create(r.(*dto.OrganizationCreateRequest))
    }, HandleCtxFunc)
}

// OrganizationList 获取组织列表
// @Summary      获取组织列表
// @Description  分页获取组织列表（需要认证）
// @Tags         组织管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Param        name query string false "组织名称（模糊搜索）"
// @Param        status query int false "状态(1=正常、2=停用)"
// @Success      200 {object} helper.DefaultResult{data=dto.OrganizationListResponse} "查询成功"
// @Failure      400 {object} helper.DefaultResult "请求参数错误"
// @Failure      401 {object} helper.DefaultResult "未授权"
// @Failure      500 {object} helper.DefaultResult "服务器内部错误"
// @Router       /admin/organization [get]
func OrganizationList(c *gin.Context) {
    req := new(dto.OrganizationListRequest)
    helper.HandleRequest(c, req, func(r helper.IBaseRequest) *helper.DefaultResult {
        return services.NewOrganizationService().List(r.(*dto.OrganizationListRequest))
    }, HandleCtxFunc)
}

// OrganizationShow 获取组织详情
// @Summary      获取组织详情
// @Description  根据ID获取组织详细信息（需要认证）
// @Tags         组织管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "组织ID"
// @Success      200 {object} helper.DefaultResult{data=dto.OrganizationShowResponse} "查询成功"
// @Failure      400 {object} helper.DefaultResult "请求参数错误"
// @Failure      401 {object} helper.DefaultResult "未授权"
// @Failure      404 {object} helper.DefaultResult "组织不存在"
// @Failure      500 {object} helper.DefaultResult "服务器内部错误"
// @Router       /admin/organization/{id} [get]
func OrganizationShow(c *gin.Context) {
    req := new(dto.OrganizationShowRequest)
    helper.HandleRequest(c, req, func(r helper.IBaseRequest) *helper.DefaultResult {
        return services.NewOrganizationService().Show(r.(*dto.OrganizationShowRequest))
    }, HandleCtxFunc)
}

// OrganizationUpdate 更新组织
// @Summary      更新组织
// @Description  更新组织信息（需要认证）
// @Tags         组织管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "组织ID"
// @Param        request body dto.OrganizationUpdateRequest true "组织信息"
// @Success      200 {object} helper.DefaultResult{data=dto.OrganizationUpdateResponse} "更新成功"
// @Failure      400 {object} helper.DefaultResult "请求参数错误"
// @Failure      401 {object} helper.DefaultResult "未授权"
// @Failure      404 {object} helper.DefaultResult "组织不存在"
// @Failure      500 {object} helper.DefaultResult "服务器内部错误"
// @Router       /admin/organization/{id} [put]
func OrganizationUpdate(c *gin.Context) {
    req := new(dto.OrganizationUpdateRequest)
    helper.HandleRequest(c, req, func(r helper.IBaseRequest) *helper.DefaultResult {
        return services.NewOrganizationService().Update(r.(*dto.OrganizationUpdateRequest))
    }, HandleCtxFunc)
}

// OrganizationDelete 删除组织
// @Summary      删除组织
// @Description  删除指定组织（需要认证）
// @Tags         组织管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "组织ID"
// @Success      200 {object} helper.DefaultResult "删除成功"
// @Failure      400 {object} helper.DefaultResult "请求参数错误"
// @Failure      401 {object} helper.DefaultResult "未授权"
// @Failure      404 {object} helper.DefaultResult "组织不存在"
// @Failure      500 {object} helper.DefaultResult "服务器内部错误"
// @Router       /admin/organization/{id} [delete]
func OrganizationDelete(c *gin.Context) {
    req := new(dto.OrganizationDeleteRequest)
    helper.HandleRequest(c, req, func(r helper.IBaseRequest) *helper.DefaultResult {
        return services.NewOrganizationService().Delete(r.(*dto.OrganizationDeleteRequest))
    }, HandleCtxFunc)
}
```

### 6.4 HandleCtxFunc 定义示例

```go
package actions

import (
    "context"
    "github.com/gin-gonic/gin"
)

// HandleCtxFunc 上下文处理函数
// 用于从 Gin Context 提取信息并设置到 context.Context
var HandleCtxFunc = func(ctx context.Context, c *gin.Context) context.Context {
    // 可以在这里添加用户信息、请求ID等到上下文
    // 例如: ctx = context.WithValue(ctx, "user_id", c.GetString("user_id"))
    return ctx
}
```

---

## 7. Repository 层规范

### 7.1 接口定义

```go
package repository

import "github.com/lshaofan/go-framework/pkg/helper"

// BaseRepository 基础仓储接口（泛型）
type BaseRepository[T any] interface {
    GetOneById(id uint) (T, error)
    FindByField(field string, value string) (T, error)
    DeleteById(id uint) error
    Create(entity T) error
    Update(entity T) error
    UpdateById(id uint, entity T) error
    UpdateWithColumns(entity T, columns ...string) error
    GetListData(request *helper.PageRequest) (*helper.PageList[T], error)
}

// IOrganizationRepository 组织仓储接口
type IOrganizationRepository interface {
    BaseRepository[models.Organization]

    // 自定义方法
    GetByCode(code string) (models.Organization, error)
    GetByName(name string) (models.Organization, error)
    UpdateStatus(id uint, status uint) error
}
```

### 7.2 DAO 实现

```go
package dao

import (
    "github.com/lshaofan/go-framework/pkg/helper"
    "gorm.io/gorm"
    "your-project/src/domain/models"
    "your-project/src/interfaces/global"
)

type Organization struct {
    *helper.Util[models.Organization]
}

func NewOrganization() *Organization {
    return &Organization{
        Util: helper.NewUtil[models.Organization](global.Config.GetDBConfig().GetGormDB()),
    }
}

// GetOneById 根据ID获取
func (o *Organization) GetOneById(id uint) (models.Organization, error) {
    model := models.Organization{}
    err := o.GetDB().Where("id = ?", id).First(&model).Error
    return model, err
}

// FindByField 根据字段查询
func (o *Organization) FindByField(field string, value string) (models.Organization, error) {
    model := models.Organization{}
    err := o.GetDB().Where(field+" = ?", value).First(&model).Error
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return model, err
    }
    return model, nil
}

// Create 创建记录
func (o *Organization) Create(entity models.Organization) error {
    return o.CreateOne(&entity)
}

// Update 更新记录
func (o *Organization) Update(entity models.Organization) error {
    return o.UpdateOne(&entity)
}

// UpdateById 根据ID更新
func (o *Organization) UpdateById(id uint, entity models.Organization) error {
    return o.GetDB().Model(&models.Organization{}).Where("id = ?", id).Updates(entity).Error
}

// UpdateWithColumns 更新指定字段
func (o *Organization) UpdateWithColumns(entity models.Organization, columns ...string) error {
    return o.UpdateOneColumn(&entity, columns...)
}

// DeleteById 根据ID删除
func (o *Organization) DeleteById(id uint) error {
    return o.GetDB().Delete(&models.Organization{}, id).Error
}

// GetListData 分页查询
func (o *Organization) GetListData(request *helper.PageRequest) (*helper.PageList[models.Organization], error) {
    db := o.GetDB().Model(&models.Organization{})

    // 处理筛选条件
    for key, value := range request.Where {
        if value == nil {
            continue
        }
        switch key {
        case "name":
            db = db.Where("name LIKE ?", "%"+value.(string)+"%")
        case "status":
            db = db.Where("status = ?", value)
        default:
            db = db.Where(key+" = ?", value)
        }
    }

    // 统计总数
    var total int64
    if err := db.Count(&total).Error; err != nil {
        return nil, err
    }

    // 查询数据
    var list []models.Organization
    if err := db.Offset((request.Page - 1) * request.PageSize).
        Limit(request.PageSize).
        Order("id DESC").
        Find(&list).Error; err != nil {
        return nil, err
    }

    return &helper.PageList[models.Organization]{
        Data:     list,
        Total:    total,
        Page:     request.Page,
        PageSize: request.PageSize,
    }, nil
}
```

---

## 8. 错误处理规范

### 8.1 错误定义

```go
package exp

import (
    "net/http"
    "github.com/lshaofan/go-framework/pkg/helper"
)

// 错误码规范：
// - 10000-10999: 通用错误
// - 11000-11999: 模块A错误
// - 12000-12999: 模块B错误
// ...

// Organization 错误码 (10100-10199)
var (
    ErrOrganizationNotFound = helper.NewErrorModel(
        10100,                    // 错误码
        "组织不存在",              // 错误消息
        nil,                      // 额外数据
        http.StatusNotFound,      // HTTP 状态码
    )

    ErrOrganizationNameExists = helper.NewErrorModel(
        10101,
        "组织名称已存在",
        nil,
        http.StatusConflict,
    )

    ErrOrganizationCodeExists = helper.NewErrorModel(
        10102,
        "组织编码已存在",
        nil,
        http.StatusConflict,
    )

    ErrOrganizationDeleteFailed = helper.NewErrorModel(
        10103,
        "组织删除失败",
        nil,
        http.StatusInternalServerError,
    )
)
```

### 8.2 HTTP 状态码使用规范

| HTTP 状态码 | 使用场景 |
|------------|----------|
| `200 OK` | 成功（查询、更新、删除） |
| `201 Created` | 创建成功 |
| `400 Bad Request` | 请求参数错误 |
| `401 Unauthorized` | 未认证 |
| `403 Forbidden` | 无权限 |
| `404 Not Found` | 资源不存在 |
| `409 Conflict` | 资源冲突（如名称重复） |
| `412 Precondition Failed` | 前置条件失败（参数验证） |
| `500 Internal Server Error` | 服务器内部错误 |

---

## 9. 完整示例

以下是一个完整的模块实现示例：

### 9.1 目录结构

```
src/
├── application/
│   ├── dto/
│   │   ├── product.go           # 数据结构
│   │   ├── product_create.go    # 创建 DTO
│   │   ├── product_list.go      # 列表 DTO
│   │   ├── product_show.go      # 详情 DTO
│   │   ├── product_update.go    # 更新 DTO
│   │   └── product_delete.go    # 删除 DTO
│   └── exp/
│       └── product_errors.go    # 错误定义
├── domain/
│   ├── models/
│   │   └── product.go           # 领域模型
│   ├── repository/
│   │   └── product.go           # 仓储接口
│   └── services/
│       ├── product.go           # Service 主文件
│       ├── product_create.go    # 创建逻辑
│       ├── product_list.go      # 列表逻辑
│       ├── product_show.go      # 详情逻辑
│       ├── product_update.go    # 更新逻辑
│       └── product_delete.go    # 删除逻辑
├── infrastructure/
│   └── dao/
│       └── product.go           # DAO 实现
└── interfaces/
    ├── api/
    │   └── product.go           # 路由注册
    └── actions/
        └── product_actions.go   # Action 处理
```

---

## 10. 代码生成模板

### 10.1 快速生成命令

使用以下模板快速生成模块代码：

```bash
# 模块名称（小写，如 product）
MODULE_NAME=product
# 模块名称（首字母大写，如 Product）
MODULE_NAME_UPPER=Product
```

### 10.2 DTO 模板

```go
// dto/${MODULE_NAME}.go
package dto

type ${MODULE_NAME_UPPER}Data struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at,omitempty"`
    UpdatedAt string `json:"updated_at,omitempty"`
}

// dto/${MODULE_NAME}_create.go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

type ${MODULE_NAME_UPPER}CreateRequest struct {
    helper.BaseRequest
    Name string `json:"name" binding:"required,max=100"`
}

type ${MODULE_NAME_UPPER}CreateResponse struct {
    ${MODULE_NAME_UPPER}Data
}

// dto/${MODULE_NAME}_list.go
package dto

import "github.com/lshaofan/go-framework/pkg/helper"

type ${MODULE_NAME_UPPER}ListRequest struct {
    helper.BaseRequest
    helper.ListRequest
    Name string `form:"name" json:"name" query:"name" binding:"omitempty"`
}

type ${MODULE_NAME_UPPER}ListResponse struct {
    Data     []${MODULE_NAME_UPPER}Data `json:"data"`
    Total    int64                      `json:"total"`
    Page     int                        `json:"page"`
    PageSize int                        `json:"page_size"`
}
```

### 10.3 Service 模板

```go
// services/${MODULE_NAME}.go
package services

import (
    "your-project/src/domain/repository"
    "your-project/src/infrastructure/dao"
)

type ${MODULE_NAME_UPPER}Service struct {
    repo repository.I${MODULE_NAME_UPPER}Repository
}

func New${MODULE_NAME_UPPER}Service() *${MODULE_NAME_UPPER}Service {
    return &${MODULE_NAME_UPPER}Service{
        repo: dao.New${MODULE_NAME_UPPER}(),
    }
}
```

### 10.4 Action 模板

```go
// actions/${MODULE_NAME}_actions.go
package actions

import (
    "github.com/gin-gonic/gin"
    "github.com/lshaofan/go-framework/pkg/helper"
    "your-project/src/application/dto"
    "your-project/src/domain/services"
)

// ${MODULE_NAME_UPPER}Create 创建${MODULE_NAME}
// @Summary      创建${MODULE_NAME}
// @Tags         ${MODULE_NAME}管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.${MODULE_NAME_UPPER}CreateRequest true "信息"
// @Success      200 {object} helper.DefaultResult{data=dto.${MODULE_NAME_UPPER}CreateResponse}
// @Router       /admin/${MODULE_NAME} [post]
func ${MODULE_NAME_UPPER}Create(c *gin.Context) {
    req := new(dto.${MODULE_NAME_UPPER}CreateRequest)
    helper.HandleRequest(c, req, func(r helper.IBaseRequest) *helper.DefaultResult {
        return services.New${MODULE_NAME_UPPER}Service().Create(r.(*dto.${MODULE_NAME_UPPER}CreateRequest))
    }, HandleCtxFunc)
}
```

---

## 附录：核心概念速查表

| 概念 | 类型 | 说明 |
|------|------|------|
| `BaseRequest` | struct | 所有请求 DTO 必须嵌入 |
| `ListRequest` | struct | 列表查询通用参数 |
| `IBaseRequest` | interface | 请求接口 |
| `DefaultResult` | struct | Service 统一返回结构 |
| `PageList[T]` | struct | 泛型分页结构 |
| `ErrorModel` | struct | 标准错误结构 |
| `Response` | struct | API 响应格式 |
| `HandleRequest` | func | 统一请求处理器 |
| `NewDefaultResult` | func | 创建 DefaultResult |
| `SetResponse` | method | 设置响应数据和错误 |
| `NewErrorModel` | func | 创建错误模型 |
| `Util[T]` | struct | 泛型 ORM 工具 |
| `PageRequest` | struct | 分页查询参数 |

---

## 版本信息

- **文档版本**: 1.0.0
- **go-framework 版本**: latest
- **Go 版本**: 1.23+
- **更新日期**: 2025-01

---

*本文档由 AI Agent 根据 go-framework 源码和项目最佳实践自动生成*
