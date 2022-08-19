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

func GetGroupIDFromParams(ctx *gin.Context) (int64, error) {
	param := ctx.Param(groupIDParamName)
	groupID, err := strconv.ParseInt(param, baseSystem, bitSize)

	if err != nil {
		return 0, errors.ErrInvalidParamFormat
	}

	return groupID, nil
}
