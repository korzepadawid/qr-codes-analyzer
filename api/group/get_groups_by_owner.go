package group

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
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

	owner, ok := auth.GetCurrentUserUsername(ctx)

	if !ok {
		return
	}

	groupsChannel := make(chan []db.Group)
	totalElementsChannel := make(chan int64)
	errorsChannel := make(chan error)

	// getting total count of user's groups and requested group page in parallel
	go h.getGroupsCount(ctx, owner, errorsChannel, totalElementsChannel)
	go h.getGroupsPage(ctx, owner, queryParams, errorsChannel, groupsChannel)

	for i := 0; i < 2; i++ {
		if cErr := <-errorsChannel; cErr != nil {
			ctx.Error(cErr)
			return
		}
	}

	response := common.NewPageResponse(
		queryParams.PageNumber,
		queryParams.PageSize,
		common.GetLastPage(<-totalElementsChannel, queryParams.PageSize),
		<-groupsChannel,
	)

	ctx.JSON(http.StatusOK, response)
}

func (h *groupHandler) getGroupsCount(
	ctx *gin.Context,
	owner string,
	errorsChannel chan<- error,
	totalElementsChannel chan<- int64,
) {
	total, err := h.store.GetGroupsCountByOwner(ctx, owner)
	errorsChannel <- err
	totalElementsChannel <- total
}

func (h *groupHandler) getGroupsPage(
	ctx *gin.Context,
	owner string,
	queryParams getGroupsByOwnerQueryParams,
	errorsChannel chan<- error,
	groupsChannel chan<- []db.Group,
) {
	arg := db.GetGroupsByOwnerParams{
		Limit:  queryParams.PageSize,
		Offset: common.GetPageOffset(queryParams.PageNumber, queryParams.PageSize),
		Owner:  owner,
	}
	groups, err := h.store.GetGroupsByOwner(ctx, arg)
	errorsChannel <- err
	groupsChannel <- groups
}
