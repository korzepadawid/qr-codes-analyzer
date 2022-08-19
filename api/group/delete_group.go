package group

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

func (h *groupHandler) deleteGroup(ctx *gin.Context) {
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

	arg := db.DeleteGroupByOwnerAndIDParams{
		GroupID: groupID,
		Owner:   owner,
	}

	if err := h.store.DeleteGroupByOwnerAndID(ctx, arg); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
