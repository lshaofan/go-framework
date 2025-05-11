package helper

// BaseRepository 通用仓储接口定义
type BaseRepository[T any] interface {
	// 基础查询方法
	GetOneById(id uint) (T, error)
	FindByField(field string, value string) (T, error)

	// 存在性检查
	ExistsById(id uint) (bool, error)

	// 删除方法
	DeleteById(id uint) error

	// 新增方法
	Create(entity T) error

	// 更新方法
	Update(entity T) error                               // 更新所有字段
	UpdateWithColumns(entity T, columns ...string) error // 更新指定字段

	// 列表查询
	GetListData(request *PageRequest) (*PageList[T], error)
}
