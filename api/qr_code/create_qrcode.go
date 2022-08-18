package qr_code

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/api/group"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
	"log"
	"net/http"
	"strings"
	"time"
)

// todo: add url validation
type createQRCodeRequest struct {
	URL         string `json:"url,omitempty" binding:"required"`
	Title       string `json:"title,omitempty" binding:"required,max=255"`
	Description string `json:"description,omitempty" binding:"max=255"`
}

type createQRCodeResponse struct {
	UUID           string    `json:"uuid,omitempty"`
	UsagesCount    int64     `json:"usages_count,omitempty" json:"usages_count,omitempty"`
	RedirectionURL string    `json:"redirection_url,omitempty" json:"redirection_url,omitempty"`
	Title          string    `json:"title,omitempty" json:"title,omitempty"`
	Description    string    `json:"description,omitempty" json:"description,omitempty"`
	QRCodeImageURL string    `json:"qr_code_image_url,omitempty" json:"qr_code_image_url,omitempty"`
	CreatedAt      time.Time `json:"created_at" json:"created_at"`
}

func (h qrCodeHandler) createQRCode(ctx *gin.Context) {
	var request createQRCodeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		return
	}

	owner, ok := auth.GetCurrentUserUsername(ctx)

	if !ok {
		return
	}

	groupID, ok := group.GetGroupIDFromParams(ctx)

	if !ok {
		return
	}

	ok = h.groupExists(ctx, owner, groupID)

	if !ok {
		return
	}

	keyUUID := uuid.NewString() // PK

	qrCode, ok := h.generateQRCodeToRedirectionProxy(ctx, keyUUID)

	if !ok {
		return
	}

	storageKeyWithFileExtension, ok := h.putQRCodeImageToFileStorage(ctx, keyUUID, qrCode)

	if !ok {
		return
	}

	// save QRCode to db
	createQRCodeArgs := db.CreateQRCodeParams{
		Owner:          owner,
		GroupID:        groupID,
		RedirectionUrl: strings.TrimSpace(request.URL),
		Title:          strings.TrimSpace(request.Title),
		Description:    strings.TrimSpace(request.Description),
		StorageUrl:     h.config.CDNAddress + storageKeyWithFileExtension,
		Uuid:           keyUUID,
	}

	QRCode, err := h.store.CreateQRCode(ctx, createQRCodeArgs)

	if err != nil {
		go func() {
			dErr := h.storage.DeleteFile(ctx, storageKeyWithFileExtension)
			if dErr != nil {
				log.Println(dErr)
			}
		}()
		ctx.Error(err)
		return
	}

	// use any cache provider to cache uuid -> redirection-url
	ctx.JSON(http.StatusCreated, newCreateQRCodeResponse(QRCode))
}

func newCreateQRCodeResponse(qrCode db.QrCode) createQRCodeResponse {
	return createQRCodeResponse{
		UUID:           qrCode.Uuid,
		UsagesCount:    qrCode.UsagesCount,
		RedirectionURL: qrCode.RedirectionUrl,
		Title:          qrCode.Title,
		Description:    qrCode.Description,
		QRCodeImageURL: qrCode.StorageUrl,
		CreatedAt:      qrCode.CreatedAt,
	}
}

func (h qrCodeHandler) putQRCodeImageToFileStorage(ctx *gin.Context, uuid string, qrCode *[]byte) (string, bool) {
	storageKeyWithFileExtension := uuid + storage.ImageExt

	putFileParams := storage.PutFileParams{
		Object:      qrCode,
		StorageKey:  storageKeyWithFileExtension,
		ContentType: storage.ImageMimeType,
	}

	err := h.storage.PutFile(ctx, putFileParams)

	if err != nil {
		ctx.Error(storage.ErrFailedToPutFile)
		return "", false
	}

	return storageKeyWithFileExtension, true
}

func (h qrCodeHandler) generateQRCodeToRedirectionProxy(ctx *gin.Context, uuid string) (*[]byte, bool) {
	qrCode, err := h.qrCodeEncoder.Encode(fmt.Sprintf("%s/qr-codes/%s/redirect", h.config.AppURL, uuid))

	if err != nil {
		ctx.Error(errors.ErrQRCodeGenerationFailed)
		return &[]byte{}, false
	}

	return &qrCode, true
}

func (h qrCodeHandler) groupExists(ctx *gin.Context, owner string, groupID int64) bool {
	params := db.GetGroupByOwnerAndIDParams{
		Owner:   owner,
		GroupID: groupID,
	}

	_, err := h.store.GetGroupByOwnerAndID(ctx, params)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Error(errors.ErrGroupNotFound)
			return false
		}
		ctx.Error(err)
		return false
	}

	return true
}
