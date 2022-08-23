package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/group"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
)

type getQRCodesRequest struct {
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=10,max=20"`
}

func (h *qrCodeHandler) getQRCodes(ctx *gin.Context) {
	var queryParams getQRCodesRequest
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.Error(err)
		return
	}

	owner, err := auth.GetCurrentUserUsername(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	groupID, err := group.GetGroupIDFromParams(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	errCh := make(chan error)
	countCh := make(chan int64)
	resCh := make(chan []db.QrCode)

	go func() {
		arg := db.GetQRCodesCountByGroupAndOwnerParams{
			GroupID: groupID,
			Owner:   owner,
		}
		count, qErr := h.store.GetQRCodesCountByGroupAndOwner(ctx, arg)
		errCh <- qErr
		countCh <- count
	}()

	go func() {
		arg := db.GetQRCodesPageByGroupAndOwnerParams{
			Limit:   queryParams.PageSize,
			Offset:  common.GetPageOffset(queryParams.PageNumber, queryParams.PageSize),
			GroupID: groupID,
			Owner:   owner,
		}
		qrCodes, qErr := h.store.GetQRCodesPageByGroupAndOwner(ctx, arg)
		errCh <- qErr
		resCh <- qrCodes
	}()

	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			ctx.Error(err)
			return
		}
	}

	pageResponse := common.NewPageResponse(
		queryParams.PageNumber,
		queryParams.PageSize,
		common.GetLastPage(<-countCh, queryParams.PageSize),
		<-resCh,
	)

	ctx.JSON(http.StatusOK, pageResponse)
}
