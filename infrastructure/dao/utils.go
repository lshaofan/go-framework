package dao

import (
	"github.com/lshaofan/go-framework/application/dto/response"
	"gorm.io/gorm"
)

// PageRequest 分页请求的参数
type PageRequest struct {
	Page     int                    // 页码
	PageSize int                    // 每页数量
	Total    int64                  // total 总数
	Where    map[string]interface{} // 条件and 自行拼接
	OrWhere  map[string]interface{} // 条件or 自行拼接
	Asc      string                 // 正序排序
	Desc     string                 //倒序排序
}

func NewPageReq() *PageRequest {
	return &PageRequest{
		Page:     1,
		PageSize: 10,
		Where:    make(map[string]interface{}),
		OrWhere:  make(map[string]interface{}),
		Asc:      "",
		Desc:     "",
	}
}

// GetPageList 分页公共方法
func GetPageList[T any](page *PageRequest, model *gorm.DB, list *response.PageList[T]) (err error) {
	// 分页查询用户
	// 拼接where条件
	if page.Where != nil {
		for k, v := range page.Where {
			model = model.Where(k, v)
		}
	}
	model = model.Session(&gorm.Session{})
	err = model.Limit(page.PageSize).Offset((page.Page - 1) * page.PageSize).Count(&list.Total).Find(&list.Data).Error
	list.Page = page.Page
	list.PageSize = page.PageSize
	return
}

func Paginate(p *PageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Session(&gorm.Session{})

		if p.Page == 0 {
			p.Page = 1
		}
		switch {
		case p.PageSize > 100:
			p.PageSize = 100
		case p.PageSize <= 0:
			p.PageSize = 10
		}
		// 拼接where条件
		if p.Where != nil {
			for k, v := range p.Where {
				db.Where(k, v)
			}
		}
		// 拼接or条件
		if p.OrWhere != nil {
			for k, v := range p.OrWhere {
				db.Or(k, v)
			}
		}
		// 拼接正序排序
		if p.Asc != "" {
			db.Order(p.Asc)
		}
		// 拼接倒序排序
		if p.Desc != "" {
			db.Order(p.Desc)
		}
		offset := (p.Page - 1) * p.PageSize
		// 分页查询
		return db.Offset(offset).Limit(p.PageSize)

	}
}
