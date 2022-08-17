package group

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

type createGroupRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=255"`
}

func (h *groupHandler) createGroup(ctx *gin.Context) {
	var request createGroupRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		return
	}

	owner, ok := auth.GetCurrentUserUsername(ctx)

	if !ok {
		return
	}

	arg := mapRequestToParams(owner, request)

	group, err := h.store.CreateGroup(ctx, arg)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, group)
}

func mapRequestToParams(owner string, request createGroupRequest) db.CreateGroupParams {
	return db.CreateGroupParams{
		Owner:       owner,
		Title:       request.Title,
		Description: request.Description,
	}
}
