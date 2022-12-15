/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  logger.go  logger.go 2022-11-30
 */

package repository

type ILogger interface {
	// AddErrorLog 添加错误日志
	AddErrorLog(fields map[string]interface{})

	// AddInfoLog 添加信息日志
	AddInfoLog(fields map[string]interface{})
}
