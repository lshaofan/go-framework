# Helper Module: logger copy.go

## 概述

`logger copy.go` 文件提供了一个基于 `github.com/sirupsen/logrus` 库的日志记录器实现。它封装了 `logrus` 的配置，提供了结构化日志记录的方法。

**注意**: 文件名中的 "copy" 可能表明这是一个重复的或旧版本的文件 (例如，可能存在一个主要的 `logger.go`)。在使用前应确认其状态。

## 主要组件

### 1. `Logger` 结构体

```go
type Logger struct {
    logger *logrus.Logger
}
```

*   **功能**: 封装一个 `logrus.Logger` 实例，作为日志操作的句柄。

### 2. `NewLogger(args interface{}) *Logger` 函数

```go
func NewLogger(args interface{}) *Logger {
    l := logrus.New()
    l.SetLevel(logrus.InfoLevel) // 日志级别硬编码为 InfoLevel
    writer, err := GetLoggerOutput("http", "exception") // 依赖外部的 GetLoggerOutput 函数
    if err != nil {
        panic(err) // 初始化失败则 panic
    }
    l.SetOutput(writer)
    l.SetFormatter(&logrus.JSONFormatter{ // 日志格式硬编码为 JSON
        TimestampFormat: "2006-01-02 15:04:05",
    })
    return &Logger{logger: l}
}
```

*   **功能**: `Logger` 的构造函数，负责初始化和配置底层的 `logrus.Logger`。
*   **初始化步骤**:
    1.  创建一个新的 `logrus.Logger` 实例。
    2.  设置默认日志级别为 `logrus.InfoLevel`。
    3.  调用一个未在此文件中定义的 `GetLoggerOutput("http", "exception")` 函数来获取日志的输出 `io.Writer`。这个函数的具体实现会决定日志被写入到哪里（如文件、控制台等）。参数 `"http"` 和 `"exception"` 暗示可能用于区分不同来源或类型的日志。
    4.  如果 `GetLoggerOutput` 返回错误，程序将 `panic`。
    5.  设置 `logrus` Logger 的输出为从 `GetLoggerOutput` 获取的 `writer`。
    6.  设置日志格式化器为 `logrus.JSONFormatter`，并指定时间戳格式为 `"2006-01-02 15:04:05"`。
*   **参数 `args interface{}`**: 此参数在当前实现中未被使用。

### 3. `Logger` 的方法

*   **`AddErrorLog(fields map[string]interface{})`**:
    ```go
    func (l *Logger) AddErrorLog(fields map[string]interface{}) {
        l.logger.WithFields(fields).Error()
    }
    ```
    *   **功能**: 记录一条错误 (Error) 级别的日志。
    *   **参数**: `fields map[string]interface{}` 允许以键值对的形式添加结构化的日志信息。
    *   **注意**: 该方法直接调用 `Error()` 而没有传递主要错误消息字符串。Logrus 仍会记录 `fields` 中的内容。

*   **`AddInfoLog(fields map[string]interface{})`**:
    ```go
    func (l *Logger) AddInfoLog(fields map[string]interface{}) {
        l.logger.WithFields(fields).Info()
    }
    ```
    *   **功能**: 记录一条信息 (Info) 级别的日志。
    *   **参数**: `fields map[string]interface{}` 允许以键值对的形式添加结构化的日志信息。
    *   **注意**: 该方法直接调用 `Info()` 而没有传递主要信息字符串。Logrus 仍会记录 `fields` 中的内容。

## 用法示例

```go
// import helper "path/to/pkg/helper"

// 创建 Logger 实例 (args 在当前实现中未使用，可传 nil)
var appLogger *helper.Logger // 最好是单例或通过依赖注入获取

func InitializeLogger() {
    appLogger = helper.NewLogger(nil)
}

func DoSomething() {
    // ... some operation ...
    if err != nil {
        appLogger.AddErrorLog(map[string]interface{}{
            "module":    "my_module",
            "operation": "do_something",
            "error_msg": err.Error(), // 建议将错误文本也放入 fields
        })
        return
    }

    appLogger.AddInfoLog(map[string]interface{}{
        "module":    "my_module",
        "operation": "do_something",
        "status":    "success",
        "record_id": 12345,
    })
}
```

## 注意事项与潜在问题

*   **文件名**: `logger copy.go` 这个名称强烈暗示它可能是 `logger.go` 的一个副本或旧版本。需要确认哪个是当前项目实际使用的日志文件，以避免维护冗余代码。
*   **`GetLoggerOutput` 依赖**: 日志记录器的核心功能（日志输出位置）依赖于外部未定义的 `GetLoggerOutput` 函数。需要找到并理解这个函数的实现。
*   **初始化 Panic**: `NewLogger` 在 `GetLoggerOutput` 失败时会 `panic`。这可能导致程序意外终止。更健壮的做法是返回一个 `error`，让调用者决定如何处理初始化失败。
*   **硬编码配置**:
    *   日志级别硬编码为 `logrus.InfoLevel`。
    *   日志格式化器硬编码为 `logrus.JSONFormatter`。
    *   这些配置如果能通过参数或配置文件进行设置，会更加灵活。
*   **未使用的 `args` 参数**: `NewLogger` 的 `args interface{}` 参数当前未被使用，可以考虑移除或明确其用途。
*   **日志方法缺少消息参数**: `AddErrorLog` 和 `AddInfoLog` 方法没有直接的 `message string` 参数。虽然可以通过 `fields` 添加信息，但通常 `logrus` 的 `Error("message", ...)` 或 `Info("message", ...)` 风格更为常见和直观。可以考虑修改方法签名或在 `fields` 中约定一个标准的消息键。
*   **日志轮转和管理**: 该文件本身未涉及日志轮转、大小限制等实际部署中需要考虑的日志管理功能。这些可能由 `GetLoggerOutput` 或外部机制处理。

## 总结

`logger copy.go` 提供了一个使用 `logrus` 进行 JSON 格式结构化日志记录的基础封装。它的主要优点是简化了 `logrus` 的配置和使用。然而，其硬编码的配置、对外部函数的依赖以及文件名都提示在实际应用前需要进一步的审查和确认。 