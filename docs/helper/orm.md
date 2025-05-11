# Helper Module: orm.go

## 概述

`orm.go` 文件提供了一系列围绕 GORM (`gorm.io/gorm`) 构建的泛型工具和辅助函数，旨在简化常见的数据库操作，特别是 CRUD (创建、读取、更新、删除) 和分页列表查询。它引入了泛型结构体 `Util[T]`，以及用于分页请求的 `PageRequest` 结构体和相关的分页处理函数。

## 主要组件

### 1. `Util[T interface{}]` 结构体

```go
type Util[T interface{}] struct {
    DB                *gorm.DB
    Model             *T // GORM 模型类型实例的指针
    PageRequestParams *PageRequest
}
```

*   **功能**: 一个泛型的数据库操作工具类，封装了常用的 GORM 操作。
*   **泛型参数 `[T interface{}]`**: `T` 代表 GORM 模型（例如，`User`, `Product` 结构体）。
*   **字段**:
    *   `DB *gorm.DB`: GORM 数据库连接实例。
    *   `Model *T`: 指向模型类型 `T` 的实例。在 `Util` 的某些方法中，它被用于向 GORM 提供模型信息（例如确定表名）。**注意**: `Model` 字段在 `NewUtil` 中未初始化，其使用方式（如 `u.DB.Model(u.Model)`) 需要调用者注意或在特定方法中被覆盖。对于实例级别的操作（如 `UpdateOne`），GORM 通常从传递的参数 `model *T` 中推断模型。
    *   `PageRequestParams *PageRequest`: 默认的分页请求参数，通过 `NewPageReq()` 初始化。

*   **`NewUtil[T interface{}](db *gorm.DB) *Util[T]`**:
    *   `Util[T]` 的构造函数，初始化 `DB` 和 `PageRequestParams`。

### 2. `Util[T]` 的方法

封装了针对模型 `T` 的 GORM 数据库操作：

*   **读取**:
    *   `GetOne(model *T) error`: 获取单条记录。`DB.Model(u.Model).First(model)`。
    *   `GetList(request *PageRequest) (*PageList[T], error)`: 获取分页列表。使用 `Paginate(request)` Scope，并统计总数。返回 `*PageList[T]` (此结构体未在此文件定义，推测在 `dto.go` 或 `web_response.go`)。
    *   `GetListWithData(request *PageRequest, data interface{}) (*PageList[T], error)`: 功能类似 `GetList`，但允许传入一个 `data` (通常是 `*[]ModelType`) 来接收查询结果。
    *   `GetAll() ([]T, error)`: 获取模型 `T` 的所有记录。

*   **创建**:
    *   `CreateOne(model *T) error`: 创建单条记录。
    *   `CreateMany(model []T) error`: 批量创建多条记录。

*   **更新**:
    *   `UpdateOne(model *T) error`: 更新单条记录（通常基于主键）。`DB.Model(model).Updates(model)`。
    *   `UpdateOneColumn(model *T, column ...string) error`: 更新单条记录的指定列。`DB.Model(model).Select(column).Updates(model)`。
    *   `UpdateMany(model []T) error`: 批量更新多条记录。GORM 的批量更新行为依赖于模型中是否存在主键。

*   **删除**:
    *   `DeleteOne(model *T) error`: 删除单条记录（通常基于主键）。
    *   `DeleteMany(model []T) error`: 批量删除多条记录（通常基于主键列表）。

*   **数据库实例管理**:
    *   `SetDB(fn func(db *gorm.DB) *gorm.DB)`: 允许通过回调函数修改或替换内部的 `gorm.DB` 实例。
    *   `GetDB() *gorm.DB`: 返回当前的 `gorm.DB` 实例。

### 3. `PageRequest` 结构体

```go
type PageRequest struct {
    Page     int                    `json:"page"`
    PageSize int                    `json:"page_size"`
    Total    int64                  `json:"total"` // 注意: Total 通常是响应的一部分
    Where    map[string]interface{} // AND 条件, e.g., {"name = ?": "John", "age > ?": 20}
    OrWhere  map[string]interface{} // OR 条件
    asc      string                 // 升序排序字段 (多个用空格隔开)
    desc     string                 // 降序排序字段 (多个用空格隔开)
}
```

*   **功能**: 定义了分页和动态查询的请求参数。
*   **字段**:
    *   `Page`, `PageSize`: 分页参数。
    *   `Total`: 总记录数。**注意**: 此字段通常属于分页响应 (`PageList`) 而非请求，放在这里可能有些混淆。
    *   `Where`: 用于构建 `AND` 查询条件。键是查询片段 (如 `"name = ?"` 或 `"status"`), 值是对应的参数。
    *   `OrWhere`: 用于构建 `OR` 查询条件。
    *   `asc`, `desc`: 私有字段，用于存储排序信息，通过 `AscSort` 和 `DescSort` 方法设置。
*   **`NewPageReq() *PageRequest`**:
    *   `PageRequest` 的构造函数，提供默认值（`Page: 1`, `PageSize: 10`）。
*   **`AscSort(field string)` 和 `DescSort(field string)` 方法**:
    *   用于设置升序和降序排序的字段。支持通过空格分隔指定多个排序字段。

### 4. `Paginate(p *PageRequest) func(db *gorm.DB) *gorm.DB` 函数

*   **功能**: 一个 GORM Scope 函数，用于构建动态的分页和条件查询。
*   **处理逻辑**:
    1.  创建新的 GORM Session (`db.Session(&gorm.Session{})`) 以避免条件污染。
    2.  **分页参数规范化**: 如果 `p.Page` 为 0，则设为 1。`p.PageSize` 被限制在 1 到 100 之间（默认为 10）。
    3.  **`WHERE` 条件**: 遍历 `p.Where` map，将条件应用到 `db` 查询。
    4.  **`OR` 条件**: 遍历 `p.OrWhere` map，将条件应用到 `db` 查询。
    5.  **排序**: 根据 `p.asc` 和 `p.desc` 字段的值，添加 `ORDER BY` 子句。支持多个字段排序（字段名通过空格分隔）。
    6.  **`OFFSET` 和 `LIMIT`**: 计算 `offset` 并应用 `db.Offset(offset).Limit(p.PageSize)`。
*   **用途**: 可通过 `db.Scopes(Paginate(pageRequest))` 应用到 GORM 查询链上。

### 5. `GetPageList[T any](page *PageRequest, model *gorm.DB, list *PageList[T]) (err error)` 函数

*   **功能**: 一个通用的分页列表获取函数。
*   **泛型参数 `[T any]`**: 用于 `list *PageList[T]` 参数的类型。
*   **参数**:
    *   `page *PageRequest`: 分页和条件参数。
    *   `model *gorm.DB`: 一个已经预设了 GORM 模型（例如通过 `DB.Model(&MyType{})`）的 `*gorm.DB` 实例。
    *   `list *PageList[T]`: 用于接收查询结果和分页信息的 `PageList` 实例。
*   **处理逻辑**:
    1.  应用 `page.Where` 中的 `AND` 条件。 **注意**: 此函数不直接处理 `page.OrWhere` 或 `page.asc`/`page.desc` 定义的排序。这些需要预先应用到传入的 `model` `*gorm.DB` 对象上。
    2.  执行 `Count(&list.Total)` 获取满足 `WHERE` 条件的总记录数。
    3.  执行 `Find(&list.Data)` 获取当前页的数据。
    4.  设置 `list.Page` 和 `list.PageSize`。
*   **与 `Util[T].GetList` 的区别**:
    *   `Util[T].GetList` 使用了更全面的 `Paginate` Scope，它能处理 `PageRequest` 中的 `Where`, `OrWhere` 以及 `asc`/`desc` 排序。
    *   `GetPageList` 是一个相对更基础的分页辅助函数，它仅直接应用 `Where` 条件，并将排序和更复杂的 `OR` 条件的处理留给调用者在传入 `model *gorm.DB` 前完成。

## 用法与影响

*   **简化 GORM 操作**: `Util[T]` 结构体及其方法为常见的数据库交互提供了便捷的封装，减少了重复代码。
*   **动态分页与查询**: `PageRequest` 结构体和 `Paginate` Scope 共同提供了一套强大的机制来实现动态的、可配置的数据列表查询，包括条件过滤、排序和分页。
*   **标准化**: `PageList[T]` (虽然未在此文件定义) 作为分页查询的标准返回结构，有助于统一接口。

## 注意事项与潜在问题

*   **`Util[T].Model *T` 字段**: `Util` 中的 `Model` 字段（一个指向 `T` 类型实例的指针）在 `NewUtil` 中未被初始化。其在 `DB.Model(u.Model)` 这样的调用中的使用方式需要开发者注意。GORM 通常通过 `DB.Model(&UserModel{})` （传递零值类型的指针）来确定表。如果 `u.Model` 未正确设置，某些 `Util` 方法可能无法按预期工作。
*   **`PageRequest.Total` 字段**: `Total` 字段通常是分页查询的结果，属于响应的一部分 (`PageList`)。将其包含在 `PageRequest` 中可能会引起混淆。
*   **`PageList[T]` 的定义缺失**: `PageList[T]` 是分页查询的关键返回类型，但其定义不在此文件中。它应该在 `dto.go` 或 `web_response.go` 中定义。
*   **SQL 注入风险 (`PageRequest.Where` / `OrWhere` 的键)**: `Paginate` 和 `GetPageList` 中处理 `Where` 条件时，如果 `map` 的键 (`k`) 直接由不受信任的外部输入构成SQL片段（而不是安全的占位符形式如 `"column = ?"`），则可能存在SQL注入风险。注释中提到的"自行拼接"需要调用者特别小心。通常，键应该是安全的列名或 GORM 支持的查询表达式，值是参数。
*   **排序字段的安全性**: `PageRequest` 的 `asc` 和 `desc` 字段用于指定排序的列名。如果这些列名直接来自用户输入而未经验证，也可能被用于恶意查询。应确保列名是有效的、预期的数据库列。
*   **`GetPageList` 功能的局限性**: 独立的 `GetPageList` 函数不像 `Paginate` Scope 那样全面，它不处理 `OrWhere` 和排序。使用时需明确其行为。
*   **GORM 事务**: 此文件中的工具未显式处理事务。事务管理通常需要在服务层或更高层级进行，或者在具体的 Repository 实现中封装。

## 总结

`orm.go` 文件提供了一套有用的 GORM 辅助工具，通过泛型 `Util[T]` 简化了 CRUD 操作，并通过 `PageRequest` 和 `Paginate` Scope 实现了强大的动态分页和过滤功能。这些工具可以显著提高开发效率并促进 GORM 使用的标准化。然而，开发者在使用时应注意其设计上的一些特性和潜在风险点，如 `Util[T].Model` 的使用方式和动态查询参数的安全性。 