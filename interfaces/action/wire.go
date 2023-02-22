//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package action

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitGinAction(c *gin.Context) *GinAction {
	wire.Build(NewGinAction)
	return &GinAction{}
}
