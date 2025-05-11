# Helper Module: web_response.go

## 概述

`web_response.go` 文件定义了用于构建和标准化 API 响应以及服务层内部数据传递的几种核心数据结构。主要包括：
*   `Response`: 通用的 API JSON 响应结构。
*   `PageList[T]`: 用于表示分页列表数据的泛型结构。
*   `DefaultResult`: 服务层向控制器/Action层返回结果时使用的标准包装结构，它实现了 `DefaultResultInterface` 接口（定义在 `interface.go`）。

## 主要组件

### 1. `Response` 结构体

```go
// Response  返回数据用于api接口
type Response struct {
    Code    int    `json:"code"`
    Result  any    `json:"result"`
    Message string `json:"message"`
}

// NewResponse 创建返回数据
func NewResponse(code int, message string, result any) *Response {
    return &Response{Code: code, Result: result, Message: message}
}
```

*   **功能**: 定义了API接口返回给客户端的标准化JSON响应格式。
*   **字段**:
    *   `Code int `json:"code"`: 应用程序级别的状态码。通常使用 `constants.go` 中定义的 `SUCCESS` (0) 或 `ERROR` (-1)，或其他自定义业务状态码。
    *   `Result any `json:"result"`: 实际的业务数据负载。成功时可能包含具体数据对象或列表；失败时可能为 `nil` 或包含更详细的错误信息结构。
    *   `Message string `json:"message"`: 人类可读的消息，用于描述操作结果（如 "操作成功"）或错误详情。
*   **`NewResponse(...)`**: `Response` 结构体的构造函数。

### 2. `PageList[T interface{}]` 结构体

```go
// PageList  分页数据
type PageList[T interface{}] struct {
    Total    int64 `json:"total"`
    Data     []T   `json:"data"`
    Page     int   `json:"page"`
    PageSize int   `json:"page_size"`
}

func NewPageList[T interface{}]() *PageList[T] {
    return &PageList[T]{}
}
```

*   **功能**: 一个泛型结构体，用于封装分页列表的查询结果。
*   **泛型参数 `[T interface{}]`**: `T` 代表列表中元素的具体类型。
*   **字段**:
    *   `Total int64 `json:"total"`: 满足查询条件的总记录数。
    *   `Data []T `json:"data"`: 当前页的实际数据列表。
    *   `Page int `json:"page"`: 当前返回的是第几页。
    *   `PageSize int `json:"page_size"`: 每页的条目数量。
*   **`NewPageList[T]()`**: `PageList[T]` 结构体的构造函数，返回一个空的 `PageList`。
*   **用途**: 当API接口返回列表数据时，`Response` 结构体的 `Result` 字段通常会是此 `PageList[T]` 类型的实例。例如，在 `orm.go` 中的 `Util[T].GetList` 方法就返回 `*PageList[T]`。

### 3. `DefaultResult` 结构体

```go
// DefaultResult 默认的返回数据结构,用于services处理完业务逻辑后返回给controller的数据结构
type DefaultResult struct {
    Err  *ErrorModel `json:"err"`
    Data any         `json:"data"`
}

// ... (NewDefaultResult, IsError, GetError, SetError, GetData, SetData, SetResponse 方法) ...
```

*   **功能**: 作为服务层方法返回给上层（如控制器或 Action 层）的标准数据包装器。它统一了服务层成功和失败时的返回方式，并实现了在 `interface.go` 中定义的 `DefaultResultInterface` 接口。
*   **字段**:
    *   `Err *ErrorModel `json:"err"`: 如果服务处理过程中发生错误，此字段会包含一个指向 `ErrorModel` 结构体（定义在 `web_error.go`）的指针。若无错误，则为 `nil`。
    *   `Data any `json:"data"`: 如果服务处理成功，此字段包含实际的业务数据。
*   **`NewDefaultResult()`**: `DefaultResult` 的构造函数，初始化时 `Err` 字段为 `nil`。
*   **核心方法 (实现 `DefaultResultInterface`)**:
    *   `IsError() bool`: 判断 `Err` 字段是否为 `nil`，从而确定服务调用是否出错。
    *   `GetError() *ErrorModel`: 返回内部的 `Err` 字段。
    *   `SetError(err error)`: 设置错误信息。
        *   如果传入的 `err` 本身就是 `*ErrorModel` 类型，则直接赋值给 `r.Err`。
        *   如果传入的是一个普通的 `error` 类型，则会使用 `NewErrorModel(-1, err.Error(), nil, http.StatusInternalServerError)` 将其包装成一个新的 `ErrorModel` 实例，并赋值给 `r.Err`。这确保了所有错误都以统一的 `ErrorModel` 结构存在。
    *   `GetData() any`: 返回内部的 `Data` 字段。
    *   `SetData(data any)`: 设置 `Data` 字段。
    *   `SetResponse(data any, err error)`: 一个便捷方法，同时调用 `SetData(data)` 和 `SetError(err)`。

## 用法与影响

*   **API 响应一致性**: `Response` 结构体确保了所有API端点都以统一的 `{"code": ..., "message": ..., "result": ...}` 格式返回JSON数据，极大地简化了客户端（前端、其他微服务）的对接和解析工作。
*   **清晰的分页数据结构**: `PageList[T]` 为分页查询结果提供了标准化的封装，包含了必要的分页元数据和当前页的数据列表。
*   **服务层与控制层解耦**: `DefaultResult` 结构体是服务层和控制层（Action层）之间的重要契约。
    *   服务层专注于业务逻辑，并将结果（成功数据或错误信息）封装在 `DefaultResult` 中返回。
    *   Action 层（如 `action.go` 中的 `BaseAction.HandleResult`）接收 `DefaultResult`，并根据其状态（`IsError()`）和内容 (`GetData()`, `GetError()`) 来构建最终的、面向客户端的 HTTP `Response`。
    *   这种分离使得服务层的逻辑不直接依赖于HTTP响应的具体格式。
*   **统一错误处理**: `DefaultResult.SetError` 的逻辑确保了无论服务层内部产生何种类型的 `error`，最终都会被统一转换为 `ErrorModel`，便于上层一致地处理和展示错误。

## 关联模块

*   **`interface.go`**: `DefaultResult` 结构体是 `DefaultResultInterface` 接口的实际实现者。
*   **`constants.go`**: `Response.Code` 字段的值通常来源于此文件中定义的 `SUCCESS` 和 `ERROR` 等常量。成功消息也可能来自此文件。
*   **`web_error.go`**: 定义了 `ErrorModel` 结构体和 `NewErrorModel` 函数，它们被 `DefaultResult` 用于封装和标准化错误信息。
*   **`action.go`**: `BaseAction` 中的 `HandleResult` 方法是 `DefaultResult` 的主要消费者，负责将其转换为最终的HTTP响应。
*   **`orm.go`**: 分页查询工具（如 `Util[T].GetList`）返回 `*PageList[T]`，这个 `PageList[T]` 实例之后会作为 `DefaultResult.Data` 的一部分，最终成为 `Response.Result` 的内容。
*   **`gin_action.go`**: `GinActionImpl` 会使用 `Response` 结构体的内容（通过 Action 接口方法间接获得）来序列化JSON并发送给客户端。

## 总结

`web_response.go` 对于构建结构清晰、易于维护的API至关重要。它通过提供标准化的响应结构 (`Response`, `PageList[T]`) 和服务层结果包装器 (`DefaultResult`)，确保了数据在应用不同层之间以及应用与外部客户端之间能够以一致和可预测的方式进行交换。 