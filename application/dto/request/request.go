/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  request.go  request.go 2022-11-30
 */

package request

// ListRequest 列表请求参数
type ListRequest struct {
	Page     int    `form:"page" json:"page" query:"page" binding:"omitempty,min=1" msg:"页码最小为1" `
	PageSize int    `form:"page_size" json:"page_size" query:"page_size" binding:"omitempty,min=1" msg:"每页数量最小为1" `
	Order    string `form:"order" json:"order" query:"order" msg:"排序" `
	Field    string `form:"field" json:"field" query:"field" msg:"排序字段" `
}

// ListResponse 列表响应参数
type ListResponse struct {
	Total int `json:"total"`
	Page  int `json:"page"`
}
