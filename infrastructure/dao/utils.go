/*
 * 版权所有 (c) 2022 伊犁绿鸟网络科技团队。
 *  utils.go  utils.go 2022-11-30
 */

package dao

import (
	"github.com/lshaofan/go-framework/application/dto/response"
	"gorm.io/gorm"
	"strings"
)

type Util[T interface{}] struct {
	DB                *gorm.DB
	Model             *T
	PageRequestParams *PageRequest
}

func NewUtil[T interface{}](db *gorm.DB) *Util[T] {
	return &Util[T]{
		DB:                db,
		PageRequestParams: NewPageReq(),
	}
}

// GetOne 获取一条记录
func (u *Util[T]) GetOne(model *T) error {
	return u.DB.Model(u.Model).First(model).Error
}

// GetList 获取多条记录
func (u *Util[T]) GetList(request *PageRequest) (*response.PageList[T], error) {
	list := &response.PageList[T]{
		Data: make([]T, 0),
	}
	err := u.DB.Model(u.Model).Scopes(Paginate(request)).Find(&list.Data).Offset(-1).Count(&list.Total).Error
	list.Page = request.Page
	list.PageSize = request.PageSize
	if err != nil {
		return nil, err

	}
	return list, nil
}

// GetListWithData  获取多条记录 使用传入的data进行返回赋值 ， 第二个参数需要传入指针
func (u *Util[T]) GetListWithData(request *PageRequest, data interface{}) (*response.PageList[T], error) {
	list := &response.PageList[T]{}
	err := u.DB.Model(u.Model).Scopes(Paginate(request)).Find(data).Offset(-1).Count(&list.Total).Error
	list.Page = request.Page
	list.PageSize = request.PageSize
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAll 获取所有记录
func (u *Util[T]) GetAll() ([]T, error) {
	all := make([]T, 0)
	err := u.DB.Model(u.Model).Find(&all).Error
	if err != nil {
		return nil, err
	}
	return all, nil
}

// CreateOne 创建一条记录
func (u *Util[T]) CreateOne(model *T) error {
	return u.DB.Model(u.Model).Create(model).Error
}

// CreateMany 创建多条记录
func (u *Util[T]) CreateMany(model *[]T) error {
	return u.DB.Model(u.Model).Create(model).Error
}

// UpdateOne 更新一条记录
func (u *Util[T]) UpdateOne(model *T) error {
	return u.DB.Model(model).Updates(model).Error
}

// UpdateOneColumn 根据字段名更新单列
func (u *Util[T]) UpdateOneColumn(model *T, column ...string) error {
	return u.DB.Model(model).Select(column).Updates(model).Error

}

// UpdateMany 更新多条记录
func (u *Util[T]) UpdateMany(model *[]T) error {
	return u.DB.Model(model).Updates(model).Error
}

// DeleteOne 删除一条记录
func (u *Util[T]) DeleteOne(model *T) error {
	return u.DB.Model(u.Model).Delete(model).Error
}

// DeleteMany 删除多条记录
func (u *Util[T]) DeleteMany(model *[]T) error {
	return u.DB.Model(u.Model).Delete(model).Error
}

// SetDB 修改DB
func (u *Util[T]) SetDB(fn func(db *gorm.DB) *gorm.DB) {
	u.DB = fn(u.DB)
}

// GetDB 获取DB
func (u *Util[T]) GetDB() *gorm.DB {
	return u.DB
}

// PageRequest 分页请求的参数
type PageRequest struct {
	Page     int                    // 页码
	PageSize int                    // 每页数量
	Total    int64                  // total 总数
	Where    map[string]interface{} // 条件and 自行拼接
	OrWhere  map[string]interface{} // 条件or 自行拼接
	asc      string                 // 正序排序
	desc     string                 //倒序排序
}

// NewPageReq 初始化分页请求参数 默认第一页 每页10条
func NewPageReq() *PageRequest {
	return &PageRequest{
		Page:     1,
		PageSize: 10,
		Where:    make(map[string]interface{}),
		OrWhere:  make(map[string]interface{}),
		asc:      "",
		desc:     "",
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

		if p.asc != "" {
			// 判断p.asc是否有空格
			if strings.Contains(p.asc, " ") {
				// 有空格代表有多个排序字段
				order := strings.Split(p.asc, " ")
				var newOrder []string
				for k, v := range order {
					// 判断是否是最后一个
					if k == len(order)-1 {
						newOrder = append(newOrder, v+" asc")
					} else {
						newOrder = append(newOrder, v+" asc, ")
					}

				}
				db.Order(newOrder)
			} else {
				db.Order(p.asc + " asc")
			}

		}
		// 拼接倒序排序
		if p.desc != "" {
			if strings.Contains(p.desc, " ") {
				order := strings.Split(p.desc, " ")
				var newOrder []string
				for k, v := range order {

					// 判断是否是最后一个
					if k == len(order)-1 {
						newOrder = append(newOrder, v+" desc")
					} else {
						newOrder = append(newOrder, v+" desc, ")
					}

				}
				db.Order(newOrder)
			} else {
				db.Order(p.desc + " desc")
			}

		}
		offset := (p.Page - 1) * p.PageSize
		// 分页查询
		return db.Offset(offset).Limit(p.PageSize)
	}
}

// AscSort 正序排序 多个排序字段使用空格隔开
func (p *PageRequest) AscSort(field string) {
	p.asc = field
}

// DescSort 倒序排序 多个排序字段使用空格隔开
func (p *PageRequest) DescSort(field string) {
	p.desc = field
}
