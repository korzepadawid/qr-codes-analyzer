package group

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

func (h *groupHandler) deleteGroup(ctx *gin.Context) {
	groupID, ok := getGroupIDFromParams(ctx)

	if !ok {
		return
	}

	owner, ok := auth.GetCurrentUserUsername(ctx)

	if !ok {
		return
	}

	arg := db.DeleteGroupByOwnerAndIDParams{
		GroupID: groupID,
		Owner:   owner,
	}

	err := h.store.DeleteGroupByOwnerAndID(ctx, arg)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
