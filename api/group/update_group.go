package group

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
	"strings"
)

type updateGroupRequest struct {
	Title       string `json:"title,omitempty" binding:"max=255"`
	Description string `json:"description,omitempty" binding:"max=255"`
}

func (h *groupHandler) updateGroup(ctx *gin.Context) {
	var request updateGroupRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		return
	}

	groupID, err := GetGroupIDFromParams(ctx)

	if err != nil {
		ctx.Error(err)
		return
	}

	owner, err := auth.GetCurrentUserUsername(ctx)

	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.UpdateGroupTxParams{
		Title:       strings.TrimSpace(request.Title),
		Description: strings.TrimSpace(request.Description),
		Owner:       owner,
		ID:          groupID,
	}

	groupTx, err := h.store.UpdateGroupTx(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Error(errors.ErrGroupNotFound)
			return
		}
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, groupTx)
}
