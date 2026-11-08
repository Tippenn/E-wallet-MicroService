package request

type TransferRequest struct {
	ToUserID string  `json:"toUserID" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}
