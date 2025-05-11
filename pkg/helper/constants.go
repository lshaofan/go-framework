package helper

type ClientType string

const (
	ClientAdminType ClientType = "admin"
	// ClientOpenapiType openapi client type
	ClientOpenapiType ClientType = "openapi"
	// ClientWebType web client type
	ClientWebType ClientType = "web"
	// ClientAppType app client type
	ClientAppType ClientType = "app"
	// ClientH5Type h5 client type
	ClientH5Type ClientType = "h5"
	// ClientWechatMiniProgramType wechat mini program client type
	ClientWechatMiniProgramType ClientType = "wechat_mini_program"
)

const ClientHeaderKey = "X-Client-Type"

const (
	ERROR   = -1
	SUCCESS = 0

	CreateSuccess = "创建成功"
	UpdateSuccess = "更新成功"
	DeleteSuccess = "删除成功"
	GetSuccess    = "获取成功"
	OkSuccess     = "操作成功"
	Succeed       = "成功"
)
