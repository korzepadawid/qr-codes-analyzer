package group

import (
	"github.com/gin-gonic/gin"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
)

const (
	routerGroupPrefix = "/groups"
)

type groupHandler struct {
	store       db.Store
	middlewares []gin.HandlerFunc
}

func NewGroupHandler(store db.Store, middlewares ...gin.HandlerFunc) *groupHandler {
	return &groupHandler{
		store:       store,
		middlewares: middlewares,
	}
}

func (h *groupHandler) RegisterRoutes(router *gin.Engine) {
	r := router.Group(routerGroupPrefix)
	r.Use(h.middlewares...)
	r.POST("", h.createGroup)
	r.GET("", h.getGroupsByOwner)
	r.GET("/:group_id", h.getGroup)
	r.DELETE("/:group_id", h.deleteGroup)
}
