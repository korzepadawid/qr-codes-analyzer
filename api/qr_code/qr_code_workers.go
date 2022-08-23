package qr_code

import (
	"context"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/ipapi"
	"log"
)

type saveRedirectJob struct {
	UUID string
	IPv4 string
}

func (h *qrCodeHandler) saveRedirectWorker() {
	log.Printf("Registered save redirect worker")
	for job := range h.redirectionWorker {
		log.Printf("%v", job)
		c := ipapi.New()
		det, err := c.GetIPDetails("142.250.203.206")

		if err != nil {
			log.Println(err)
			return
		}

		arg := db.IncrementRedirectEntriesTxParams{
			UUID:      job.UUID,
			IPv4:      "142.250.203.206",
			IPDetails: det,
		}
		if err := h.store.IncrementRedirectEntriesTx(context.Background(), arg); err != nil {
			log.Printf("%v", err)
		}
	}
}
