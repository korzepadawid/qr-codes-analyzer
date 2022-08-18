package group

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"strconv"
)

const (
	baseSystem       = 10
	bitSize          = 64
	groupIDParamName = "group_id"
)

func GetGroupIDFromParams(ctx *gin.Context) (int64, bool) {
	param := ctx.Param(groupIDParamName)
	groupID, err := strconv.ParseInt(param, baseSystem, bitSize)

	if err != nil {
		ctx.Error(errors.ErrInvalidParamFormat)
		return 0, false
	}

	return groupID, true
}
