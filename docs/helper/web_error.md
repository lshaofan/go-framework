# Helper Module: web_error.go

## 概述

`web_error.go` 文件定义了一个用于在应用程序中标准化错误信息的结构体 `ErrorModel`。它还提供了一个构造函数 `NewErrorModel` 来创建此模型的实例，并通过实现 `Error() string` 方法使 `*ErrorModel` 类型满足 Go 内置的 `error` 接口。

## 主要组件

### 1. `ErrorModel` 结构体

```go
package helper

//  示例: ServerError = NewErrorModel(500, "服务器错误", nil, http.StatusInternalServerError)

// ErrorModel 错误模型
type ErrorModel struct {
    Code       int         `json:"code"`
    Message    string      `json:"message"`
    Result     interface{} `json:"result"`    // 可用于携带额外错误上下文，如验证失败的字段详情
    HttpStatus int         `json:"httpStatus" swaggerignore:"true"` // 对应的HTTP状态码
}
```

*   **功能**: 提供一个标准的、结构化的方式来表示应用程序中发生的各种错误。
*   **字段**:
    *   `Code int `json:"code"`: 应用程序内部定义的错误码。这可以是一个通用的错误码（例如 `-1` 来自 `constants.go`），或者是一个更具体的业务逻辑错误码。
    *   `Message string `json:"message"`: 人类可读的错误描述信息。这个消息会直接或间接展示给用户或记录到日志中。
    *   `Result interface{} `json:"result"`: 一个灵活的字段，用于携带与错误相关的额外数据或上下文。例如，在参数验证失败时，这里可能包含导致错误的具体字段和原因。通常情况下，如果没有额外信息，此字段可以为 `nil`。
    *   `HttpStatus int `json:"httpStatus" swaggerignore:"true"`: 指示当此错误需要通过HTTP API返回给客户端时，应该使用的HTTP状态码（例如 `http.StatusBadRequest` (400), `http.StatusNotFound` (404), `http.StatusInternalServerError` (500)）。
        *   `swaggerignore:"true"` 标签: 这个标签通常用于告知 Swagger/OpenAPI 文档生成工具忽略此字段。这是因为 HTTP 状态码是响应元数据的一部分，通常不在响应体 JSON 中显式表示 `httpStatus` 字段。

### 2. `NewErrorModel(code int, message string, result interface{}, httpStatus int) *ErrorModel` 函数

*   **功能**: `ErrorModel` 结构体的构造函数。
*   **参数**: 接收错误码 (`code`)、错误消息 (`message`)、可选的额外结果数据 (`result`) 和相应的 HTTP 状态码 (`httpStatus`)。
*   **返回**: 返回一个初始化后的 `*ErrorModel` 指针。

### 3. `(e *ErrorModel) Error() string` 方法

```go
func (e *ErrorModel) Error() string {
    return e.Message
}
```

*   **功能**: 此方法使得 `*ErrorModel` 类型满足 Go 的标准 `error` 接口。
*   **行为**: 当一个 `*ErrorModel` 实例被当作一个 `error` 类型来处理时（例如，在 `fmt.Printf("Error: %v", err)` 中，或者在进行 `if err != nil` 判断后直接使用 `err.Error()`），此方法会被调用，并返回 `ErrorModel` 中的 `Message` 字段内容。

## 用法与影响

*   **标准化错误表示**: `ErrorModel` 为整个应用程序（包括服务层、Action 层、工具函数等）提供了一种统一的方式来创建和传递结构化的错误信息。
*   **丰富错误上下文**: 与 Go 内置的简单 `error` 字符串相比，`ErrorModel` 可以携带更多的上下文信息，如业务错误码、建议的 HTTP 状态码以及可选的详细数据，这对于调试、日志记录和向客户端提供有意义的错误反馈都非常有用。
*   **API 错误响应构建**:
    *   当 Action 层（如 `gin_action.go` 中的 `GinActionImpl.ThrowError`）需要向客户端报告错误时，它会接收一个 `*ErrorModel`。
    *   然后，Action 层会使用 `ErrorModel` 的 `Code`、`Message` 和 `Result` 字段来填充标准 `Response` 结构体（来自 `web_response.go`）的相应字段，并使用 `ErrorModel.HttpStatus` 字段来设置实际的 HTTP 响应状态码。
*   **错误转换与包装**:
    *   在 `web_request.go` 中的 `GetValidateErr` 方法会将参数验证库（`go-playground/validator`）产生的错误转换为 `*ErrorModel`。
    *   在 `web_response.go` 中的 `DefaultResult.SetError` 方法，如果接收到的是一个非 `*ErrorModel` 类型的普通 `error`，它会将其包装成一个新的 `*ErrorModel`，从而确保错误在内部流转时具有统一的结构。
*   **错误类型判断**: 代码中可以使用 `errors.As(err, &targetErrorModel)` 来安全地检查一个 `error` 变量是否是 `*ErrorModel` 类型，如果是，则可以访问其所有特定字段。

## 关联模块

*   **`web_response.go`**: `DefaultResult` 结构体中的 `Err` 字段是 `*ErrorModel` 类型。其 `SetError` 方法会创建或处理 `ErrorModel`。
*   **`web_request.go`**: `Request.GetValidateErr` 方法的主要职责就是将验证错误转换成 `*ErrorModel`。
*   **`gin_action.go`**: `GinActionImpl` 中的 `ThrowError` 方法直接消费 `*ErrorModel` 来构建 HTTP 错误响应。其 `Error` 和 `ThrowValidateError` 方法也会处理或间接触发 `ErrorModel` 的创建/使用。
*   **`constants.go`**: `ErrorModel` 的 `Code` 字段可能使用此文件中定义的 `ERROR` 常量（如 `-1`）作为通用错误代码。
*   **`interface.go`**: `DefaultResultInterface` 接口的 `GetError()` 方法声明返回 `*ErrorModel`。

## 总结

`web_error.go` 中定义的 `ErrorModel` 是该 `helper` 包中错误处理框架的核心。它通过提供一个结构化的错误表示，并使其兼容标准的 `error` 接口，极大地增强了应用程序处理和传递错误信息的能力。这不仅有助于开发和调试，也为向API客户端提供清晰、一致的错误反馈奠定了基础。 