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
	UsagesCount    int64     `json:"usages_count,omitempty"`
	RedirectionURL string    `json:"redirection_url,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	QRCodeImageURL string    `json:"qr_code_image_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *qrCodeHandler) createQRCode(ctx *gin.Context) {
	var request createQRCodeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
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

	if err := h.groupExists(ctx, owner, groupID); err != nil {
		ctx.Error(err)
		return
	}

	keyUUID := uuid.NewString() // PK

	// generates QRCode to the redirect url
	qrCode, err := h.genereateQRCode(keyUUID)

	if err != nil {
		ctx.Error(err)
		return
	}

	storageKeyWithFileExtension, err := h.putQRCodeImageToFileStorage(ctx, keyUUID, qrCode)

	if err != nil {
		ctx.Error(err)
		return
	}

	createQRCodeArgs := h.newCreateQRCodeParams(owner, groupID, request, storageKeyWithFileExtension, keyUUID)

	QRCode, err := h.store.CreateQRCode(ctx, createQRCodeArgs)

	if err != nil {
		go h.deleteFile(ctx, storageKeyWithFileExtension)
		ctx.Error(err)
		return
	}

	// use any cache provider to cache uuid -> redirection-url
	ctx.JSON(http.StatusCreated, newCreateQRCodeResponse(QRCode))
}

func (h *qrCodeHandler) deleteFile(ctx *gin.Context, storageKeyWithFileExtension string) {
	if dErr := h.storage.DeleteFile(ctx, storageKeyWithFileExtension); dErr != nil {
		log.Println(dErr)
	}
}

func (h *qrCodeHandler) newCreateQRCodeParams(owner string, groupID int64, request createQRCodeRequest, storageKeyWithFileExtension string, keyUUID string) db.CreateQRCodeParams {
	createQRCodeArgs := db.CreateQRCodeParams{
		Owner:          owner,
		GroupID:        groupID,
		RedirectionUrl: strings.TrimSpace(request.URL),
		Title:          strings.TrimSpace(request.Title),
		Description:    strings.TrimSpace(request.Description),
		StorageUrl:     h.config.CDNAddress + storageKeyWithFileExtension,
		Uuid:           keyUUID,
	}
	return createQRCodeArgs
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

func (h qrCodeHandler) putQRCodeImageToFileStorage(ctx *gin.Context, uuid string, qrCode *[]byte) (string, error) {
	storageKeyWithFileExtension := uuid + storage.ImageExt

	putFileParams := storage.PutFileParams{
		Object:      qrCode,
		StorageKey:  storageKeyWithFileExtension,
		ContentType: storage.ImageMimeType,
	}

	err := h.storage.PutFile(ctx, putFileParams)

	if err != nil {
		return "", storage.ErrFailedToPutFile
	}

	return storageKeyWithFileExtension, nil
}

func (h qrCodeHandler) genereateQRCode(uuid string) (*[]byte, error) {
	qrCode, err := h.qrCodeEncoder.Encode(fmt.Sprintf("%s/encode-codes/%s/redirect", h.config.AppURL, uuid))

	if err != nil {
		return &[]byte{}, errors.ErrQRCodeGenerationFailed
	}

	return &qrCode, nil
}

func (h qrCodeHandler) groupExists(ctx *gin.Context, owner string, groupID int64) error {
	params := db.GetGroupByOwnerAndIDParams{
		Owner:   owner,
		GroupID: groupID,
	}

	_, err := h.store.GetGroupByOwnerAndID(ctx, params)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrGroupNotFound
		}
		return err
	}

	return nil
}
