package repository

type ILogger interface {
	// AddErrorLog 添加错误日志
	AddErrorLog(fields map[string]interface{})

	// AddInfoLog 添加信息日志
	AddInfoLog(fields map[string]interface{})
}
