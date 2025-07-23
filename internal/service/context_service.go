package service

import (
	"github.com/gin-gonic/gin"
)

type ContextLogic struct {
	Context *gin.Context
}

var ContextLogicInstance *ContextLogic

func init() {
	ContextLogicInstance = &ContextLogic{}
}

func (c *ContextLogic) SetContext(gc *gin.Context) {
	c.Context = gc
}
