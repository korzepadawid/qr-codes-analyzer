package ipapi

import "errors"

var (
	ErrRequestCreation       = errors.New("failed to create a new request")
	ErrExternalAPICallFailed = errors.New("failed to call an external api")
	ErrFailedJSONParse       = errors.New("failed to parse json response")
)

type IPAddr string

type IPDetails struct {
	Status  string  `json:"status"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"` // approximated
	Lon     float64 `json:"lon"` // approximated
	ISP     string  `json:"isp"` // internet service provider (ISP)
	Org     string  `json:"org"`
	AS      string  `json:"as"` // autonomous system (AS)
	Query   string  `json:"query"`
}

//Client the interface responsible for connecting
//with an external apis like ip-api.com, geo.ipify.com
type Client interface {

	//GetIPDetails gets details about requested ip address,
	//in case of failure returns an error
	GetIPDetails(addr IPAddr) (*IPDetails, error)
}
