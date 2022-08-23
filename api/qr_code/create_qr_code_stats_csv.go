package qr_code

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
	"strconv"
)

var (
	colNames = []string{"uuid", "title", "url", "ipv4", "isp", "as", "city", "country", "lat", "lon", "date"}
)

func (h *qrCodeHandler) createQRCodeStatsCSV(ctx *gin.Context) {
	uuid, err := getQRCodeUUID(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	owner, err := auth.GetCurrentUserUsername(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.GetQRCodeRedirectEntriesParams{
		Uuid:  uuid,
		Owner: owner,
	}
	entries, err := h.store.GetQRCodeRedirectEntries(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	statsFileCSV, err := convertEntriesToCSV(entries)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=stats-%s.csv", uuid))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", strconv.Itoa(len(statsFileCSV)))
	ctx.Writer.Write(statsFileCSV)
}

func convertEntriesToCSV(entries []db.GetQRCodeRedirectEntriesRow) ([]byte, error) {
	b := make([]byte, 0)
	buff := bytes.NewBuffer(b)
	writer := csv.NewWriter(buff)

	// write column names
	if err := writer.Write(colNames); err != nil {
		return make([]byte, 0), err
	}
	// write entries
	for _, e := range entries {
		r := convertEntryToRow(e)
		if err := writer.Write(r); err != nil {
			return make([]byte, 0), err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return make([]byte, 0), err
	}

	return buff.Bytes(), nil
}

func convertEntryToRow(entry db.GetQRCodeRedirectEntriesRow) []string {
	return []string{
		entry.Uuid,
		entry.Title,
		entry.Url,
		entry.Ipv4,
		entry.Isp,
		entry.AutonomousSys,
		entry.City,
		entry.Country,
		entry.Lat,
		entry.Lon,
		entry.Date.String(),
	}
}
