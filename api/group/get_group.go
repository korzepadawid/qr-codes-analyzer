package group

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

// todo: implement getting count of qr-codes in the given group
func (h *groupHandler) getGroup(ctx *gin.Context) {
	groupID, ok := getGroupIDFromParams(ctx)

	if !ok {
		return
	}

	owner, ok := auth.GetCurrentUserUsername(ctx)

	if !ok {
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
