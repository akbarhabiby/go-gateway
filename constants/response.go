package constants

const (
	STATUS_SUCCESS = "00"
	STATUS_FAILED  = "99"
)

const (
	MESSAGE_SUCCESS = "Success"
)

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

type HealthCheckResponse struct {
	Message    string `json:"message"`
	ServerTime string `json:"serverTime"`
	Version    string `json:"version"`
}
