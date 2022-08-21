package ipapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type clientIPAPI struct{}

func New() *clientIPAPI {
	return &clientIPAPI{}
}

func (c *clientIPAPI) GetIPDetails(addr IPAddr) (*IPDetails, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	requestURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,query,country,city,lat,lon,isp,org,as", addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)

	if err != nil {
		return nil, ErrRequestCreation
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, ErrExternalAPICallFailed
	}

	defer func(b io.ReadCloser) {
		err := b.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, ErrExternalAPICallFailed
	}

	var result IPDetails
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, ErrFailedJSONParse
	}

	return &result, nil
}
