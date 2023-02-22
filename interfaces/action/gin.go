package action

import (
	"github.com/gin-gonic/gin"
	"github.com/lshaofan/go-framework/application/dto/request"
	"github.com/lshaofan/go-framework/application/dto/response"
)

type GinAction struct {
	*request.GinRequest
	*response.GinResponse
}

func NewGinAction(c *gin.Context) *GinAction {
	return &GinAction{
		request.NewGinRequest(c),
		response.NewGinResponse(c),
	}
}
