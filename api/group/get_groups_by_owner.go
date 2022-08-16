package group

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

type getGroupsByOwnerQueryParams struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=20"`
}

func (h *groupHandler) getGroupsByOwner(ctx *gin.Context) {
	var queryParams getGroupsByOwnerQueryParams

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.Error(err)
		return
	}

	owner, err := auth.GetCurrentUserUsername(ctx)

	if err != nil {
		ctx.Error(errors.ErrFailedCurrentUserRetrieval)
		return
	}

	totalElementsChannel := make(chan int64)
	resultsChannel := make(chan []db.Group)
	errorsChannel := make(chan error)

	go func() {
		total, err2 := h.store.GetGroupsCountByOwner(ctx, owner)
		errorsChannel <- err2
		totalElementsChannel <- total
	}()

	go func() {
		arg := db.GetGroupsByOwnerParams{
			Limit:  queryParams.PageSize,
			Offset: common.GetPageOffset(queryParams.PageNumber, queryParams.PageSize),
			Owner:  owner,
		}
		groups, err2 := h.store.GetGroupsByOwner(ctx, arg)
		errorsChannel <- err2
		resultsChannel <- groups
	}()

	for i := 0; i < 2; i++ {
		if cErr := <-errorsChannel; cErr != nil {
			ctx.Error(cErr)
			return
		}
	}

	results := <-resultsChannel
	totalElements := <-totalElementsChannel

	response := common.NewPageResponse(
		queryParams.PageNumber,
		queryParams.PageSize,
		common.GetLastPage(totalElements, queryParams.PageSize),
		results,
	)

	ctx.JSON(http.StatusOK, response)
}
