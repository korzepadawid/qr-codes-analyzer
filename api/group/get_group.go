package group

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

func (h *groupHandler) getGroup(ctx *gin.Context) {
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

	arg := db.GetGroupByOwnerAndIDParams{
		Owner:   owner,
		GroupID: groupID,
	}
	group, err := h.store.GetGroupByOwnerAndID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Error(errors.ErrGroupNotFound)
			return
		}
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, group)
}
