package group

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
	"strconv"
)

const (
	baseSystem       = 10
	bitSize          = 64
	groupIDParamName = "group_id"
)

// todo: implement getting count of qr-codes in the given group
func (h *groupHandler) getGroup(ctx *gin.Context) {
	param := ctx.Param(groupIDParamName)
	groupID, err := strconv.ParseInt(param, baseSystem, bitSize)

	if err != nil {
		ctx.Error(errors.ErrInvalidParamFormat)
		return
	}

	owner, err := auth.GetCurrentUserUsername(ctx)

	if err != nil {
		ctx.Error(errors.ErrFailedCurrentUserRetrieval)
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
