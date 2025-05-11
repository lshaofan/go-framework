# Helper Module: web.go (Generic Service Runner)

## 概述

`web.go` 文件（更准确地说，可以称之为服务运行器或应用生命周期管理器）提供了一个用于管理长时间运行的 Go 服务或应用程序的框架。它定义了一系列接口 (`Service`, `Context`, `Environment`, `Handler`) 和一个核心的 `Run` 函数，用于控制服务的初始化、启动、信号处理和优雅停止。这个模式对于构建需要妥善处理操作系统信号（如 `SIGINT`, `SIGTERM`）并能平滑关闭的后台守护进程或任何类型的服务都非常有用。

## 主要组件

### 1. 接口定义

*   **`Service` 接口**:
    ```go
    type Service interface {
        Init(Environment) error // 初始化，非阻塞
        Start() error           // 启动服务，非阻塞
        Stop() error            // 停止服务
    }
    ```
    *   **功能**: 任何希望通过此框架管理的服务都必须实现此核心接口。
    *   `Init(Environment) error`: 在服务启动前调用，用于执行初始化任务（如加载配置、连接数据库）。`Environment` 参数提供关于运行环境的信息。此方法必须是非阻塞的。
    *   `Start() error`: 在 `Init` 成功后调用，用于启动服务的主要逻辑（如启动一个 HTTP 服务器的监听 goroutine）。此方法也必须是非阻塞的。
    *   `Stop() error`: 当接收到停止信号或上下文被取消时调用，用于执行优雅关闭的逻辑（如关闭数据库连接、等待正在进行的任务完成）。

*   **`Environment` 接口**:
    ```go
    type Environment interface {
        IsWindowsService() bool // 判断是否作为 Windows 服务运行
    }
    ```
    *   **功能**: 提供关于服务运行环境的信息。
    *   `IsWindowsService() bool`: 报告程序是否作为 Windows 服务运行。文件内提供了一个简单的 `environment` 结构体实现，其此方法总是返回 `false`。

*   **`Context` 接口 (可选)**:
    ```go
    type Context interface {
        Context() context.Context // 返回一个与服务生命周期绑定的上下文
    }
    ```
    *   **功能**: 如果 `Service` 也实现了此接口，`Run` 函数会使用返回的 `context.Context` 的 `Done()` 通道作为额外的停止条件。

*   **`Handler` 接口 (可选)**:
    ```go
    type Handler interface {
        Handle(os.Signal) error // 处理接收到的操作系统信号
    }
    ```
    *   **功能**: 如果 `Service` 也实现了此接口，当 `Run` 函数捕获到指定的操作系统信号时，会调用此 `Handle` 方法。如果此方法返回 `ErrStop`，则会触发服务的 `Stop()`流程。

### 2. `environment` 结构体

```go
type environment struct{}

func (environment) IsWindowsService() bool {
    return false
}
```
*   **功能**: `Environment` 接口的一个简单默认实现。`IsWindowsService` 硬编码返回 `false`，表明完整的 Windows 服务特性可能需要更具体的实现或外部包支持。

### 3. `ErrStop` 变量

```go
var ErrStop = errors.New("stopping service")
```
*   **功能**: 一个哨兵错误。当 `Handler` 接口的 `Handle` 方法返回此错误时，`Run` 函数会知道应该调用服务的 `Stop()` 方法。

### 4. `signalNotify` 变量

```go
var signalNotify = signal.Notify
```
*   **功能**: 将标准库的 `signal.Notify` 函数赋值给一个变量，这样做主要是为了方便在单元测试中模拟 (mock) 信号通知的行为。

### 5. `Run(service Service, sig ...os.Signal) error` 函数

*   **功能**: 这是启动和管理服务生命周期的核心函数。它会阻塞执行，直到接收到指定的操作系统信号或服务的上下文被取消。
*   **执行流程**:
    1.  **初始化 (`service.Init`)**:
        *   创建一个 `environment` 实例。
        *   调用 `service.Init(env)`。如果失败，则返回错误。
    2.  **启动 (`service.Start`)**:
        *   调用 `service.Start()`。如果失败，则返回错误。
    3.  **信号监听设置**:
        *   如果没有在 `sig` 参数中指定信号，则默认监听 `syscall.SIGINT` (中断信号，通常由 Ctrl+C 触发) 和 `syscall.SIGTERM` (终止信号)。
        *   创建一个 `os.Signal` 通道 `signalChan`。
        *   使用 `signalNotify` (即 `signal.Notify`) 将指定的信号转发到 `signalChan`。
    4.  **上下文获取**:
        *   检查 `service` 是否实现了 `Context` 接口。
        *   如果是，则调用 `service.Context()` 获取上下文；否则，使用 `context.Background()`。
    5.  **主循环 (阻塞)**:
        *   进入一个无限循环，使用 `select` 语句监听两个事件：
            *   **接收到 OS 信号 (`case s := <-signalChan`)**:
                *   检查 `service` 是否实现了 `Handler` 接口。
                *   如果是，则调用 `service.Handle(s)`。若返回 `ErrStop`，则通过 `goto stop` 跳转到停止流程。若返回其他错误，`Run` 函数本身会忽略这个错误（由 `Handle` 方法自行处理）。
                *   如果否 (未实现 `Handler`)，则直接 `goto stop` (为了向后兼容)。
            *   **上下文完成 (`case <-ctx.Done()`)**:
                *   如果服务的上下文的 `Done()` 通道关闭，则 `goto stop`。
    6.  **停止流程 (`stop:` 标签)**:
        *   调用 `service.Stop()`，并将其返回的错误作为 `Run` 函数的最终结果。

## 用法示例 (概念性)

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "time"
    "your_project/pkg/helper" // 替换为你的 helper 包路径
)

// MyApplication 实现了 helper.Service, helper.Context, 和 helper.Handler
type MyApplication struct {
    httpSrv    *http.Server
    appCtx     context.Context
    cancelFunc context.CancelFunc
}

func (app *MyApplication) Init(env helper.Environment) error {
    fmt.Println("Initializing application...")
    // 配置加载、数据库连接等
    return nil
}

func (app *MyApplication) Start() error {
    fmt.Println("Starting application...")
    app.appCtx, app.cancelFunc = context.WithCancel(context.Background())

    // 启动一个 HTTP 服务器作为示例
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello from MyApplication!")
    })
    app.httpSrv = &http.Server{Addr: ":8080", Handler: mux}

    go func() {
        fmt.Println("HTTP server listening on :8080")
        if err := app.httpSrv.ListenAndServe(); err != http.ErrServerClosed {
            fmt.Printf("HTTP server error: %v\n", err)
            // 实际应用中可能需要一种方式通知主服务停止
            app.cancelFunc() // 例如，通过取消上下文来停止服务
        }
    }()
    return nil
}

func (app *MyApplication) Stop() error {
    fmt.Println("Stopping application...")
    if app.cancelFunc != nil {
        app.cancelFunc() // 确保上下文被取消
    }

    // 优雅关闭 HTTP 服务器
    if app.httpSrv != nil {
        shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer shutdownCancel()
        if err := app.httpSrv.Shutdown(shutdownCtx); err != nil {
            fmt.Printf("HTTP server shutdown error: %v\n", err)
            return err
        }
        fmt.Println("HTTP server stopped.")
    }
    // 其他清理工作，如关闭数据库连接
    fmt.Println("Application stopped gracefully.")
    return nil
}

// Context (可选)
func (app *MyApplication) Context() context.Context {
    return app.appCtx
}

// Handle (可选)
func (app *MyApplication) Handle(s os.Signal) error {
    fmt.Printf("Received OS signal: %s. Preparing to stop.\n", s)
    // 可以根据不同信号执行不同操作，这里统一返回 ErrStop
    return helper.ErrStop
}

func main() {
    myApp := &MyApplication{}
    fmt.Println("Running service via helper.Run...")
    if err := helper.Run(myApp, syscall.SIGINT, syscall.SIGTERM); err != nil {
        fmt.Fprintf(os.Stderr, "Service run failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Service has been shut down.")
}

```

## 注意事项与依赖

*   **非阻塞 `Init` 和 `Start`**: 实现 `Service` 接口时，`Init` 和 `Start` 方法必须是非阻塞的。长时间运行的任务（如启动 HTTP 监听）应在 `Start` 方法内部的 goroutine 中执行。
*   **优雅关闭**: `Stop` 方法是实现所有清理逻辑（关闭连接、保存状态、等待 goroutines 完成等）的关键。
*   **上下文管理**: `Context` 接口的实现允许服务通过 Go 的 `context` 包进行更细致的生命周期管理，这对于需要与应用其他部分（如图书馆、微服务客户端）进行协作取消的场景特别有用。
*   **信号处理的灵活性**: `Handler` 接口提供了对接收到的操作系统信号进行自定义处理的能力。如果服务未实现 `Handler`，则默认的 `SIGINT` 和 `SIGTERM` 信号会直接触发 `Stop()`。
*   **Windows 服务支持**: `Environment.IsWindowsService()` 方法和相关注释表明设计上考虑了作为 Windows 服务运行的可能性，但当前的 `environment` 实现并未提供实际功能。完整的 Windows 服务集成通常需要使用如 `golang.org/x/sys/windows/svc` 这样的特定包。
*   **错误处理**: `Run` 函数会传播来自 `Init`、`Start` 和 `Stop` 方法的错误。`Handler.Handle` 返回的错误（除了 `ErrStop`）不会被 `Run` 直接传播，应在 `Handle` 方法内部处理。
*   **命名**: 文件名 `web.go` 对于这个通用服务运行器框架来说可能不够精确，类似 `service_runner.go` 或 `app_lifecycle.go` 这样的名称可能更符合其功能。

## 总结

`web.go` 文件（尽管名为 "web"）提供了一个强大且通用的框架，用于构建能够响应操作系统信号并实现优雅关闭的 Go 服务。通过实现其定义的接口，开发者可以专注于服务的核心业务逻辑，同时利用该框架来处理复杂的应用生命周期管理。 