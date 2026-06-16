package dto

type TransferResponse struct {
	TransferID string `json:"transferId"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
