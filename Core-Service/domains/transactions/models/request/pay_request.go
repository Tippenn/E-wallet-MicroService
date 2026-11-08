package request

type PayRequest struct {
	MerchantID  string  `json:"merchantID" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Description string  `json:"description" validate:"required"`
}
