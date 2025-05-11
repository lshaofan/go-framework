# Helper Module: logger.go

## 概述

`logger.go` 文件主要定义了 `GetLoggerOutput` 函数。此函数使用 `github.com/lestrrat-go/file-rotatelogs` 库来创建和配置一个可自动轮转的日志文件写入器 (`io.Writer`)。这个函数是日志系统（如 `logger copy.go` 中定义的 `Logger`）设置日志输出目标（特别是对于文件日志及其轮转策略）的关键部分。

## 主要组件

### 1. `GetLoggerOutput(path, name string) (*rotatelogs.RotateLogs, error)` 函数

```go
package helper

import (
    "fmt"
    "time"

    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// GetLoggerOutput 获取 SetOutput
func GetLoggerOutput(path, name string) (*rotatelogs.RotateLogs, error) {
    // 获取当前时间
    now := time.Now()
    // 获取当前年月日
    year, month, day := now.Date()
    nowTime := fmt.Sprintf("%d-%d-%d", year, month, day) // e.g., "2023-10-27"

    // 构建日志文件基础路径
    // 例如: ./runtime/log/2023-10-27/http/access.log
    filePath := fmt.Sprintf("./runtime/log/%s/%s/%s.log", nowTime, path, name)

    return rotatelogs.New(
        // 轮转后的文件命名模式，例如: ./runtime/log/2023-10-27/http/access.log.20231027
        filePath+".%Y%m%d",
        // 创建一个符号链接，指向当前最新的日志文件
        rotatelogs.WithLinkName(filePath),
        // 每24小时轮转一次
        rotatelogs.WithRotationTime(24*time.Hour),
        // 最多保留7个轮转后的日志文件
        rotatelogs.WithRotationCount(7),
        // 当日志文件大小超过100MB时进行轮转
        rotatelogs.WithRotationSize(100*1024*1024),
    )
}
```

*   **功能**: 创建并返回一个 `*rotatelogs.RotateLogs` 对象，该对象实现了 `io.Writer` 接口，并能自动管理日志文件的轮转。
*   **参数**:
    *   `path string`: 用于构建日志文件路径的子目录名（例如，表示日志来源模块，如 "http", "grpc"）。
    *   `name string`: 用于构建日志文件路径的另一部分，通常是日志的种类或具体名称（例如，"access", "error", "debug"）。
*   **日志文件路径构造**:
    1.  获取当前日期 (`nowTime`，格式为 `YYYY-M-D`)。
    2.  日志文件的基础路径 `filePath` 被构造成 `./runtime/log/[当前日期]/[path参数]/[name参数].log`。
        *   这意味着每天的日志会存储在不同的日期命名的子目录下。例如，如果 `path="api"`，`name="errors"`，在 2023年10月27日，基础路径会是 `./runtime/log/2023-10-27/api/errors.log`。
*   **日志轮转配置 (`rotatelogs.New`)**:
    *   **轮转文件模式**: `filePath + ".%Y%m%d"`。这是实际写入的、带有日期后缀的轮转文件名。例如，上述 `errors.log` 在轮转后会变成 `errors.log.20231027`。
    *   **符号链接**: `rotatelogs.WithLinkName(filePath)` 在 `filePath`（即 `./runtime/log/[当前日期]/[path参数]/[name参数].log`）处创建一个符号链接，该链接始终指向当前正在写入的、带有日期后缀的实际日志文件。这提供了一个固定路径来访问最新的日志。
    *   **按时间轮转**: `rotatelogs.WithRotationTime(24 * time.Hour)` 设置日志每 24 小时轮转一次。
    *   **保留文件数**: `rotatelogs.WithRotationCount(7)` 设置在轮转时最多保留 7 个旧的日志文件。超过这个数量的旧文件会被删除。
    *   **按大小轮转**: `rotatelogs.WithRotationSize(100 * 1024 * 1024)` (100MB) 设置当日志文件大小达到 100MB 时也会触发轮转，即使未到 24 小时。
*   **返回值**:
    *   成功时返回配置好的 `*rotatelogs.RotateLogs` 对象和 `nil` 错误。
    *   如果 `rotatelogs.New` 初始化失败，则返回 `nil` 和相应的错误。

## 用法与影响

*   此函数通常在初始化日志系统时被调用（例如，在 `logger copy.go` 中的 `NewLogger` 函数内部）。
*   日志库（如 `logrus`）会将 `GetLoggerOutput` 返回的 `io.Writer` 作为其输出目标。
*   **日志结构**: 日志文件会存储在 `./runtime/log/` 目录下，并按 `[日期]/[path参数]/[name参数].log` 的结构组织。
*   **自动管理**: 实现了日志文件的自动轮转（基于时间和大小）和旧日志的自动清理。
*   **便捷访问**: 通过符号链接，可以方便地访问到最新的日志文件，而无需关心具体轮转后的文件名。

## 注意事项与潜在问题

*   **日志目录结构与符号链接**:
    *   由于 `filePath` 本身包含了当前日期 (`nowTime`)，`GetLoggerOutput` 每天首次被调用时（或者如果日志实例每天重新创建），传递给 `rotatelogs.New` 的基础路径都会改变。
    *   这意味着 `rotatelogs.WithLinkName(filePath)` 创建的符号链接实际上是位于当天的日期目录下的，例如 `./runtime/log/2023-10-27/http/access.log`。
    *   这种结构是可行的，但它意味着"固定名称的链接"只在当天的目录内是固定的。如果需要一个全局固定的链接（例如，一个始终指向最新 `http/access` 日志的链接，无论日期），那么传递给 `rotatelogs.New` 的 `logFile` 参数和 `WithLinkName` 的参数应该是一个不包含动态日期的稳定路径（例如 `./runtime/log/http/access.log`），并让 `rotatelogs` 的模式参数（如 `filePath + ".%Y%m%d"` 中的 `.%Y%m%d`）全权负责生成带日期的文件名。
*   **错误处理**: `GetLoggerOutput` 正确地将 `rotatelogs.New` 可能返回的错误向上传播。调用此函数的代码（如 `NewLogger`）应该妥善处理这个错误（目前在 `logger copy.go` 中是 `panic`）。
*   **文件系统权限**: 运行此应用的进程需要有在 `./runtime/log/` 目录下创建子目录和文件的权限。
*   **外部库依赖**: 依赖 `github.com/lestrrat-go/file-rotatelogs` 库。
*   **硬编码配置**: 轮转的时间、保留数量、大小限制以及基础路径 `./runtime/log/` 都是硬编码的。在更复杂的应用中，将这些配置参数化（例如通过配置文件或环境变量读取）会提供更大的灵活性。
*   **路径分隔符**: 代码中使用了 `/` 作为路径分隔符。虽然 Go 的 `os` 包在很多情况下能处理不同操作系统的路径差异，但在显式构造路径时，使用 `filepath.Join` 是确保跨平台兼容性的更稳妥做法（尽管在此场景下，日志通常是在部署的服务器上，操作系统类型是已知的）。

## 与 `logger copy.go` 的关系

在 `logger copy.go` 中，`NewLogger` 函数调用了 `GetLoggerOutput("http", "exception")`。这意味着由该 `NewLogger` 实例产生的日志：
*   将被写入到类似 `./runtime/log/[当前日期]/http/exception.log.YYYYMMDD` 的文件中。
*   会有一个符号链接 `./runtime/log/[当前日期]/http/exception.log` 指向当前的日志文件。

## 总结

`GetLoggerOutput` 函数为应用提供了一个健壮的日志文件管理方案，通过集成 `file-rotatelogs` 库实现了日志的按天、按大小轮转以及旧日志清理。这对于防止日志文件无限增长、耗尽磁盘空间以及保持日志的有序性至关重要。设计上的主要特点是按天创建新的日志子目录，并在其中进行轮转和链接。 